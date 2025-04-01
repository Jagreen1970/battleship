package game

import "errors"

type API struct {
	db Database
}

func NewApi(db Database) *API {
	return &API{
		db: db,
	}
}

func (A *API) GetPlayer(playerName string) (*Player, error) {
	player, err := A.db.FindPlayerByName(playerName)
	if err != nil {
		return nil, err
	}

	return player, nil
}

func (A *API) NewPlayer(playerName string) (*Player, error) {
	player, err := A.db.FindPlayerByName(playerName)
	if err == nil {
		return player, err
	}
	if errors.Is(err, ErrorNotFound) {
		player, err = A.db.CreatePlayer(playerName)
	}
	return player, err
}

func (A *API) NewGame(player string, name string) (*Game, error) {
	p, err := A.db.FindPlayerByName(player)
	if err != nil {
		return nil, err
	}

	game, err := A.db.CreateGame(NewGame(p, name))
	if err != nil {
		return nil, err
	}

	return game, nil
}

func (A *API) UpdateGame(g *Game) (*Game, error) {
	return A.db.UpdateGame(g)
}

func (A *API) Games(page int, count int) ([]*Game, error) {
	games, err := A.db.QueryGames(page, count)
	if err != nil {
		return nil, err
	}

	return games, nil
}

func (A *API) GetGame(id string) (*Game, error) {
	return A.db.FindGameByID(id)
}

func (A *API) GetGameByName(name string) (*Game, error) {
	return A.db.FindGameByName(name)
}

func (A *API) ScoreBoard(playerName string) (*ScoreBoard, error) {
	return NewScoreBoard(playerName), nil
}

// DeleteGame deletes a game by ID
func (A *API) DeleteGame(id string) error {
	return A.db.DeleteGame(id)
}

// DeleteAllGames deletes all games
func (A *API) DeleteAllGames() (int, error) {
	return A.db.DeleteAllGames()
}
