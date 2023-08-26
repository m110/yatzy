package transport

type GameStarted struct {
	Players []Player
}

type TurnStarted struct {
	Player Player
}

type DiceRolled struct {
	Dice []Die
}

type ScoresReady struct {
	Scores []Score
}

type ScoreAssigned struct {
	Player Player
	Score  Score
}
