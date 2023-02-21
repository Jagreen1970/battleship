package battleship

type ScoreBoard struct {
	Scores []Player `json:"scores"`
}

func NewScoreBoard(_ string) *ScoreBoard {
	return &ScoreBoard{
		Scores: []Player{
			{
				Name:  "Player 1",
				Score: 1,
			},
			{
				Name:  "Player 2",
				Score: 2,
			},
		},
	}
}
