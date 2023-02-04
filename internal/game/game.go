package game

import "fmt"

type Database interface {
	FindPlayerByName(username string) (Player, error)
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

type FieldState string

const (
	FieldStateUnknown FieldState = "!"
	FieldStateEmpty   FieldState = " "
	FieldStatePin     FieldState = "O"
	FieldStateHit     FieldState = "X"
	FielStateMiss     FieldState = "-"
)

type Game struct {
	Boards  map[string]*Board `json:"boards"`
	History []Move            `json:"history"`
	Status  Status            `json:"status"`

	Player1      Player `json:"player_1"`
	Player2      Player `json:"player_2"`
	PlayerToMove string `json:"player_to_move"`
}

func NewGame(player1 Player) *Game {
	g := Game{
		Status:  StatusSetup,
		Player1: player1,
	}
	g.Boards = make(map[string]*Board)
	g.Boards[player1.Name] = NewBoard()
	g.History = make([]Move, 0)

	return &g
}

func (g *Game) Join(player2 Player) error {
	if len(g.Boards) > 1 {
		return fmt.Errorf("you are not allowed to join the game: %w", ErrorIllegal)
	}

	if g.Status != StatusSetup {
		return fmt.Errorf("you are not allowed to join the game: %w", ErrorInvalid)
	}

	if _, ok := g.Boards[player2.Name]; ok {
		return fmt.Errorf("seems you already joined the game, %s: %w", player2.Name, ErrorIllegal)
	}

	g.Boards[player2.Name] = NewBoard()
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

func (g *Game) Start(playerName string) error {
	if g.Status == StatusPlaying {
		return fmt.Errorf("you are already playing: %w", ErrorInvalid)
	}

	if g.Status != StatusSetup {
		return fmt.Errorf("you are not allowed to start the game: %w", ErrorInvalid)
	}

	if !g.allPinsPlaced() {
		return fmt.Errorf("can't start game - players not ready: %w", ErrorNotReady)
	}

	if _, ok := g.Boards[playerName]; !ok {
		return fmt.Errorf("you are not allowed to start the game: %w", ErrorIllegal)
	}

	g.PlayerToMove = playerName
	g.Status = StatusPlaying
	return nil
}

func (g *Game) allPinsPlaced() bool {
	for _, board := range g.Boards {
		if board.PinsAvailable > 0 {
			return false
		}
	}

	return true
}

func (g *Game) MakeMove(playerName string, x, y int) error {
	if err := g.checkGameStatusForPlayer(playerName); err != nil {
		return err
	}

	playerBoard, opponentBoard, err := g.getPlayerBoards(playerName)
	if err != nil {
		return err
	}

	if err := playerBoard.CanAttack(x, y); err != nil {
		return fmt.Errorf("you can't attack: %w", err)
	}

	if result, err := opponentBoard.Attack(x, y); err != nil {
		return fmt.Errorf("you can't attack: %w", err)
	} else {
		playerBoard.Track(result, x, y)
	}

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
