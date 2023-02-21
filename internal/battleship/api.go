package battleship

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
	if err != nil && errors.Is(err, ErrorNotFound) {
		player, err = A.db.CreatePlayer(playerName)
	}
	return player, nil
}

func (A *API) NewGame(player string) (*Game, error) {
	p, err := A.db.FindPlayerByName(player)
	if err != nil {
		return nil, err
	}

	game, err := A.db.CreateGame(NewGame(p))
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

func (A *API) ScoreBoard(playerName string) (*ScoreBoard, error) {
	return NewScoreBoard(playerName), nil
}
