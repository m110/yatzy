package scene

import (
	"sort"

	"github.com/m110/yatzy/internal/entity"
)

type scoringFunc func([]*entity.Die) uint

func ones(dice []*entity.Die) uint {
	return sumOfDice(dice, 1)
}

func twos(dice []*entity.Die) uint {
	return sumOfDice(dice, 2)
}

func threes(dice []*entity.Die) uint {
	return sumOfDice(dice, 3)
}

func fours(dice []*entity.Die) uint {
	return sumOfDice(dice, 4)
}

func fives(dice []*entity.Die) uint {
	return sumOfDice(dice, 5)
}

func sixes(dice []*entity.Die) uint {
	return sumOfDice(dice, 6)
}

func sumOfDice(dice []*entity.Die, value uint) uint {
	var score uint

	for _, d := range dice {
		if d.Value == value {
			score += d.Value
		}
	}

	return score
}

func onePair(dice []*entity.Die) uint {
	var score uint

	counts := diceCounts(dice)

	for k, v := range counts {
		if v >= 2 {
			s := k * 2
			if s > score {
				score = s
			}
		}
	}

	return score
}

func twoPairs(dice []*entity.Die) uint {
	counts := diceCounts(dice)

	var scores []int

	for k, v := range counts {
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

func threeOfAKind(dice []*entity.Die) uint {
	return atLeastCount(dice, 3)
}

func fourOfAKind(dice []*entity.Die) uint {
	return atLeastCount(dice, 4)
}

func smallStraight(dice []*entity.Die) uint {
	counts := diceCounts(dice)
	if counts[1] > 0 &&
		counts[2] > 0 &&
		counts[3] > 0 &&
		counts[4] > 0 &&
		counts[5] > 0 {
		return 15
	}

	return 0
}

func largeStraight(dice []*entity.Die) uint {
	counts := diceCounts(dice)
	if counts[2] > 0 &&
		counts[3] > 0 &&
		counts[4] > 0 &&
		counts[5] > 0 &&
		counts[6] > 0 {
		return 20
	}

	return 0
}

func fullHouse(dice []*entity.Die) uint {
	var pairScore uint
	var threeOfAKindScore uint

	counts := diceCounts(dice)

	for k, v := range counts {
		if v == 2 {
			s := k * 2
			if s > pairScore {
				pairScore = s
			}
		} else if v == 3 {
			s := k * 3
			if s > threeOfAKindScore {
				threeOfAKindScore = s
			}
		}
	}

	if pairScore > 0 && threeOfAKindScore > 0 {
		return pairScore + threeOfAKindScore
	}

	return 0
}

func chance(dice []*entity.Die) uint {
	var score uint

	for _, d := range dice {
		score += d.Value
	}

	return score
}

func yatzy(dice []*entity.Die) uint {
	if atLeastCount(dice, 5) > 0 {
		return 50
	}

	return 0
}

func atLeastCount(dice []*entity.Die, c uint) uint {
	counts := diceCounts(dice)
	for k, v := range counts {
		if v >= c {
			return k * c
		}
	}

	return 0
}

func diceCounts(dice []*entity.Die) map[uint]uint {
	counts := map[uint]uint{}

	for _, d := range dice {
		v, ok := counts[d.Value]
		if ok {
			counts[d.Value] = v + 1
		} else {
			counts[d.Value] = 1
		}
	}

	return counts
}
