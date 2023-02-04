package game

type API struct {
	db Database
}

func (A *API) NewGame(player string) (Game, error) {
	//TODO implement me
	panic("implement me")
}

func (A *API) UpdateGame(g Game) error {
	//TODO implement me
	panic("implement me")
}

func NewApi(db Database) *API {
	return &API{
		db: db,
	}
}

func (A *API) Games(page int, count int) ([]Game, error) {
	//TODO implement me
	panic("implement me")
}

func (A *API) GetGame(id string) (Game, error) {
	//TODO implement me
	panic("implement me")
}

func (A *API) ScoreBoard(playerName string) (*ScoreBoard, error) {
	return NewScoreBoard(), nil
}
