package scene

import (
	"sort"

	"github.com/m110/yatzy/internal/entity"
)

type appliesFunc func([]entity.Die) uint

func ones(dice []entity.Die) uint {
	return sumOfDice(dice, 1)
}

func twos(dice []entity.Die) uint {
	return sumOfDice(dice, 2)
}

func threes(dice []entity.Die) uint {
	return sumOfDice(dice, 3)
}

func fours(dice []entity.Die) uint {
	return sumOfDice(dice, 4)
}

func fives(dice []entity.Die) uint {
	return sumOfDice(dice, 5)
}

func sixes(dice []entity.Die) uint {
	return sumOfDice(dice, 6)
}

func sumOfDice(dice []entity.Die, value uint) uint {
	var score uint

	for _, d := range dice {
		if d.Value == value {
			score += d.Value
		}
	}

	return score
}

func onePair(dice []entity.Die) uint {
	pairs := map[uint]uint{}

	for _, d := range dice {
		v, ok := pairs[d.Value]
		if ok {
			pairs[d.Value] = v + 1
		} else {
			pairs[d.Value] = 1
		}
	}

	var score uint

	for k, v := range pairs {
		if v >= 2 {
			s := k * 2
			if s > score {
				score = s
			}
		}
	}

	return score
}

func twoPairs(dice []entity.Die) uint {
	pairs := map[uint]uint{}

	for _, d := range dice {
		v, ok := pairs[d.Value]
		if ok {
			pairs[d.Value] = v + 1
		} else {
			pairs[d.Value] = 1
		}
	}

	var scores []int

	for k, v := range pairs {
		if v >= 2 {
			scores = append(scores, int(k)*2)
		}
	}

	if len(scores) < 2 {
		return 0
	}

	sort.Ints(scores)

	return uint(scores[len(scores)-1] + scores[len(scores)-2])
}

func threeOfAKind(dice []entity.Die) uint {
	var score uint

	return score
}

func fourOfAKind(dice []entity.Die) uint {
	var score uint

	return score
}

func smallStraight(dice []entity.Die) uint {
	var score uint

	return score
}

func largeStraight(dice []entity.Die) uint {
	var score uint

	return score
}

func fullHouse(dice []entity.Die) uint {
	var score uint

	return score
}

func chance(dice []entity.Die) uint {
	var score uint

	for _, d := range dice {
		score += d.Value
	}

	return score
}

func yatzy(dice []entity.Die) uint {
	var score uint

	return score
}
