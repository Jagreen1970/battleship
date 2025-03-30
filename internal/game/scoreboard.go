package game

type ScoreBoard struct {
	Scores []Player `json:"scores"`
}

// DummyScoreBoard is a dummy scoreboard
var DummyScoreBoard = &ScoreBoard{
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

// NewScoreBoard creates a new scoreboard. For now, it is a dummy with two players
func NewScoreBoard(_ string) *ScoreBoard {
	return DummyScoreBoard
}
