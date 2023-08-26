package transport

type RollDice struct {
}

type RerollDice struct {
	Dice []int
}

type PassRoll struct {
}

type AssignScore struct {
	Type ScoreType
}
