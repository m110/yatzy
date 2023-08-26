package transport

type Player struct {
	ID string
}

type Die struct {
	ID    int
	Value uint
}

type ScoreType int

const (
	ScoreTypeOnes ScoreType = iota
	ScoreTypeTwos
	ScoreTypeThrees
	ScoreTypeFours
	ScoreTypeFives
	ScoreTypeSixes
	ScoreTypeOnePair
	ScoreTypeTwoPairs
	ScoreTypeThreeOfAKind
	ScoreTypeFourOfAKind
	ScoreTypeSmallStraight
	ScoreTypeLargeStraight
	ScoreTypeFullHouse
	ScoreTypeChance
	ScoreTypeYatzy
)

type Score struct {
	Type  ScoreType
	Value uint
}
