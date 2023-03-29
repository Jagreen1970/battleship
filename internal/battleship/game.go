package battleship

import (
	"fmt"
)

type Database interface {
	CreatePlayer(playerName string) (*Player, error)
	FindPlayerByName(username string) (*Player, error)

	QueryGames(page int, count int) ([]*Game, error)
	CreateGame(game *Game) (*Game, error)
	FindGameByID(id string) (*Game, error)
	UpdateGame(g *Game) (*Game, error)
}

type Move struct {
	Player string `json:"player"`
	Hit    bool   `json:"hit"`
	X      int    `json:"x"`
	Y      int    `json:"y"`
}

type Status int

const (
	StatusSetup Status = iota
	StatusPlaying
	StatusWon
	StatusLost
)

type FieldState byte
type FieldRow [10]FieldState

const (
	FieldStateUnknown FieldState = '!'
	FieldStateEmpty   FieldState = ' '
	FieldStatePin     FieldState = 'O'
	FieldStateHit     FieldState = 'X'
	FielStateMiss     FieldState = '-'
)

type Game struct {
	ID      string            `json:"_id,omitempty" bson:"_id,omitempty"`
	Boards  map[string]*Board `json:"boards" bson:"boards"`
	History []Move            `json:"history" bson:"history"`
	Status  Status            `json:"status" bson:"status"`

	Player1      *Player `json:"player_1" bson:"player1"`
	Player2      *Player `json:"player_2" bson:"player2"`
	PlayerToMove string  `json:"player_to_move" bson:"player_to_move"`
}

func NewGame(player1 *Player) *Game {
	g := Game{
		Status:  StatusSetup,
		Player1: player1,
		Player2: &Player{
			Name: "nobody",
		},
	}

	g.InitBoards()
	g.InitHistory()

	return &g
}

func (g *Game) InitBoards() {
	g.Boards = make(map[string]*Board)
	g.Boards[g.Player1.Name] = NewBoard(g.Player1.Name, g.Player2.Name)
	g.Boards[g.Player2.Name] = NewBoard(g.Player2.Name, g.Player1.Name)
}

func (g *Game) InitHistory() {
	g.History = make([]Move, 0)
}

func (g *Game) Join(player2 *Player) error {
	if len(g.Boards) > 1 {
		return fmt.Errorf("you are not allowed to join the game: %w", ErrorIllegal)
	}

	if g.Status != StatusSetup {
		return fmt.Errorf("you are not allowed to join the game: %w", ErrorInvalid)
	}

	if _, ok := g.Boards[player2.Name]; ok {
		return fmt.Errorf("seems you already joined the game, %s: %w", player2.Name, ErrorIllegal)
	}

	g.Boards[player2.Name] = NewBoard(player2.Name, g.Player1.Name)
	g.Player2 = player2
	return nil
}

func (g *Game) PlacePin(playerName string, x int, y int) error {
	if g.Status != StatusSetup {
		return fmt.Errorf("you are not allowed to set a pin: %w", ErrorIllegal)
	}

	var err error
	if board, ok := g.Boards[playerName]; ok {
		err = board.PlacePin(x, y)
	} else {
		return fmt.Errorf("you are not allowed to place a pin: %w", ErrorIllegal)
	}

	return err
}

func (g *Game) RecoverPin(playerName string, x int, y int) error {
	if g.Status != StatusSetup {
		return fmt.Errorf("you are not allowed to recover a pin: %w", ErrorIllegal)
	}

	var err error
	if board, ok := g.Boards[playerName]; ok {
		err = board.RecoverPin(x, y)
	} else {
		return fmt.Errorf("you are not allowed to recover a pin: %w", ErrorIllegal)
	}

	return err
}

func (g *Game) Start(playerName string) error {
	err := g.CanStart(playerName)
	if err != nil {
		return err
	}

	g.PlayerToMove = playerName
	g.Status = StatusPlaying
	return nil
}

func (g *Game) CanStart(playerName string) error {
	if g.Status == StatusPlaying {
		return fmt.Errorf("you are already playing: %w", ErrorInvalid)
	}

	err := g.ValidSetup()
	if err != nil {
		return err
	}

	if !g.allPinsPlaced() {
		return fmt.Errorf("can't start game - players not ready: %w", ErrorNotReady)
	}

	if _, ok := g.Boards[playerName]; !ok {
		return fmt.Errorf("you are not allowed to start the game: %w", ErrorIllegal)
	}

	return nil
}

// ValidSetup checks if all placed pins are valid
func (g *Game) ValidSetup() error {
	if g.Status != StatusSetup {
		return fmt.Errorf("error checking setup - status: %w", ErrorInvalid)
	}

	for player, board := range g.Boards {
		if err := board.ValidSetup(); err != nil {
			return fmt.Errorf("error checking setup player %s: %w", player, ErrorInvalid)
		}
	}

	return nil
}

func (g *Game) MakeMove(move Move) error {
	playerName := move.Player
	if err := g.checkGameStatusForPlayer(playerName); err != nil {
		return err
	}

	playerBoard, opponentBoard, err := g.getPlayerBoards(playerName)
	if err != nil {
		return err
	}

	if err := playerBoard.CanAttack(move.X, move.Y); err != nil {
		return fmt.Errorf("you can't attack: %w", err)
	}

	if result, err := opponentBoard.Attack(move.X, move.Y); err != nil {
		return fmt.Errorf("you can't attack: %w", err)
	} else {
		playerBoard.Track(result, move.X, move.Y)
	}

	g.History = append(g.History, move)
	g.UpdateGameState()
	if g.Status == StatusPlaying {
		g.cyclePlayerToMove(playerName)
	}

	return nil
}

func (g *Game) UpdateGameState() {
	if g.Status != StatusPlaying {
		return
	}

	if g.Boards[g.Player1.Name].Lost() {
		g.Status = StatusLost
	}

	if g.Boards[g.Player2.Name].Lost() {
		g.Status = StatusWon
	}
}

func (g *Game) allPinsPlaced() bool {
	for _, board := range g.Boards {
		if board.PinsAvailable > 0 {
			return false
		}
	}

	return true
}

func (g *Game) getPlayerBoards(playerName string) (*Board, *Board, error) {
	playerBoard, ok := g.Boards[playerName]
	if !ok {
		return nil, nil, fmt.Errorf("could not read board of player %q : %w", playerName, ErrorNotFound)
	}

	opponentBoard, ok := g.Boards[g.opponent(playerName)]
	if !ok {
		return nil, nil, fmt.Errorf("could not read board for opponent of player %q (name: %q) : %w", playerName, g.opponent(playerName), ErrorNotFound)
	}

	return playerBoard, opponentBoard, nil
}

func (g *Game) checkGameStatusForPlayer(playerName string) error {
	if g.Status != StatusPlaying {
		return fmt.Errorf("you are not allowed to make a move: %w", ErrorNotReady)
	}

	if playerName != g.PlayerToMove {
		return fmt.Errorf("it's not your turn, %s (%w)", playerName, ErrorIllegal)
	}

	return nil
}

func (g *Game) opponent(playerName string) string {
	for name := range g.Boards {
		if name != playerName {
			return name
		}
	}
	return ""
}

func (g *Game) cyclePlayerToMove(playerName string) {
	g.PlayerToMove = g.opponent(playerName)
}
