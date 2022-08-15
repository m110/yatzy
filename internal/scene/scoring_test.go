package scene

import (
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/m110/yatzy/internal/entity"
)

func TestSixes(t *testing.T) {
	testCases := []struct {
		Dice          []uint
		ExpectedScore uint
	}{
		{
			Dice:          []uint{1, 2, 3, 4, 5},
			ExpectedScore: 0,
		},
		{
			Dice:          []uint{1, 2, 3, 4, 6},
			ExpectedScore: 6,
		},
		{
			Dice:          []uint{1, 2, 3, 6, 6},
			ExpectedScore: 12,
		},
		{
			Dice:          []uint{5, 5, 5, 5, 5},
			ExpectedScore: 0,
		},
		{
			Dice:          []uint{6, 6, 6, 6, 6},
			ExpectedScore: 30,
		},
	}

	for _, tc := range testCases {
		t.Run(testName(tc.Dice), func(t *testing.T) {
			dice := diceFromValues(tc.Dice)
			score := sixes(dice)
			assert.Equal(t, tc.ExpectedScore, score)
		})
	}
}

func TestOnePair(t *testing.T) {
	testCases := []struct {
		Dice          []uint
		ExpectedScore uint
	}{
		{
			Dice:          []uint{1, 2, 3, 4, 5},
			ExpectedScore: 0,
		},
		{
			Dice:          []uint{1, 1, 3, 4, 5},
			ExpectedScore: 2,
		},
		{
			Dice:          []uint{1, 2, 2, 4, 5},
			ExpectedScore: 4,
		},
		{
			Dice:          []uint{1, 1, 2, 2, 5},
			ExpectedScore: 4,
		},
		{
			Dice:          []uint{5, 5, 5, 5, 5},
			ExpectedScore: 10,
		},
	}

	for _, tc := range testCases {
		t.Run(testName(tc.Dice), func(t *testing.T) {
			dice := diceFromValues(tc.Dice)
			score := onePair(dice)
			assert.Equal(t, tc.ExpectedScore, score)
		})
	}
}

func TestTwoPairs(t *testing.T) {
	testCases := []struct {
		Dice          []uint
		ExpectedScore uint
	}{
		{
			Dice:          []uint{1, 2, 3, 4, 5},
			ExpectedScore: 0,
		},
		{
			Dice:          []uint{1, 1, 3, 4, 5},
			ExpectedScore: 0,
		},
		{
			Dice:          []uint{1, 1, 2, 2, 5},
			ExpectedScore: 6,
		},
		{
			Dice:          []uint{5, 5, 5, 5, 5},
			ExpectedScore: 0,
		},
		{
			Dice:          []uint{6, 5, 5, 6, 5},
			ExpectedScore: 22,
		},
	}

	for _, tc := range testCases {
		t.Run(testName(tc.Dice), func(t *testing.T) {
			dice := diceFromValues(tc.Dice)
			score := twoPairs(dice)
			assert.Equal(t, tc.ExpectedScore, score)
		})
	}
}

func TestThreeOfAKind(t *testing.T) {
	testCases := []struct {
		Dice          []uint
		ExpectedScore uint
	}{
		{
			Dice:          []uint{1, 2, 3, 4, 5},
			ExpectedScore: 0,
		},
		{
			Dice:          []uint{1, 1, 3, 4, 5},
			ExpectedScore: 0,
		},
		{
			Dice:          []uint{1, 1, 2, 2, 5},
			ExpectedScore: 0,
		},
		{
			Dice:          []uint{5, 5, 1, 1, 5},
			ExpectedScore: 15,
		},
		{
			Dice:          []uint{5, 5, 5, 5, 5},
			ExpectedScore: 15,
		},
	}

	for _, tc := range testCases {
		t.Run(testName(tc.Dice), func(t *testing.T) {
			dice := diceFromValues(tc.Dice)
			score := threeOfAKind(dice)
			assert.Equal(t, tc.ExpectedScore, score)
		})
	}
}

func TestFourOfAKind(t *testing.T) {
	testCases := []struct {
		Dice          []uint
		ExpectedScore uint
	}{
		{
			Dice:          []uint{1, 2, 3, 4, 5},
			ExpectedScore: 0,
		},
		{
			Dice:          []uint{1, 1, 3, 4, 5},
			ExpectedScore: 0,
		},
		{
			Dice:          []uint{1, 1, 2, 2, 5},
			ExpectedScore: 0,
		},
		{
			Dice:          []uint{5, 5, 1, 1, 5},
			ExpectedScore: 0,
		},
		{
			Dice:          []uint{1, 5, 5, 5, 5},
			ExpectedScore: 20,
		},
		{
			Dice:          []uint{5, 5, 5, 5, 5},
			ExpectedScore: 20,
		},
	}

	for _, tc := range testCases {
		t.Run(testName(tc.Dice), func(t *testing.T) {
			dice := diceFromValues(tc.Dice)
			score := fourOfAKind(dice)
			assert.Equal(t, tc.ExpectedScore, score)
		})
	}
}

func TestSmallStraight(t *testing.T) {
	testCases := []struct {
		Dice          []uint
		ExpectedScore uint
	}{
		{
			Dice:          []uint{1, 2, 3, 4, 5},
			ExpectedScore: 15,
		},
		{
			Dice:          []uint{1, 1, 3, 4, 5},
			ExpectedScore: 0,
		},
		{
			Dice:          []uint{2, 3, 4, 5, 6},
			ExpectedScore: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(testName(tc.Dice), func(t *testing.T) {
			dice := diceFromValues(tc.Dice)
			score := smallStraight(dice)
			assert.Equal(t, tc.ExpectedScore, score)
		})
	}
}

func TestLargeStraight(t *testing.T) {
	testCases := []struct {
		Dice          []uint
		ExpectedScore uint
	}{
		{
			Dice:          []uint{1, 2, 3, 4, 5},
			ExpectedScore: 0,
		},
		{
			Dice:          []uint{2, 3, 4, 5, 6},
			ExpectedScore: 20,
		},
		{
			Dice:          []uint{2, 3, 4, 5, 5},
			ExpectedScore: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(testName(tc.Dice), func(t *testing.T) {
			dice := diceFromValues(tc.Dice)
			score := largeStraight(dice)
			assert.Equal(t, tc.ExpectedScore, score)
		})
	}
}

func TestFullHouse(t *testing.T) {
	testCases := []struct {
		Dice          []uint
		ExpectedScore uint
	}{
		{
			Dice:          []uint{1, 2, 3, 4, 5},
			ExpectedScore: 0,
		},
		{
			Dice:          []uint{2, 2, 3, 3, 6},
			ExpectedScore: 0,
		},
		{
			Dice:          []uint{2, 2, 3, 3, 3},
			ExpectedScore: 13,
		},
		{
			Dice:          []uint{2, 3, 3, 3, 3},
			ExpectedScore: 0,
		},
		{
			Dice:          []uint{3, 3, 3, 3, 3},
			ExpectedScore: 0,
		},
		{
			Dice:          []uint{5, 5, 6, 6, 6},
			ExpectedScore: 28,
		},
	}

	for _, tc := range testCases {
		t.Run(testName(tc.Dice), func(t *testing.T) {
			dice := diceFromValues(tc.Dice)
			score := fullHouse(dice)
			assert.Equal(t, tc.ExpectedScore, score)
		})
	}
}

func TestChance(t *testing.T) {
	testCases := []struct {
		Dice          []uint
		ExpectedScore uint
	}{
		{
			Dice:          []uint{1, 2, 3, 4, 5},
			ExpectedScore: 15,
		},
		{
			Dice:          []uint{2, 2, 5, 5, 5},
			ExpectedScore: 19,
		},
		{
			Dice:          []uint{6, 6, 6, 6, 6},
			ExpectedScore: 30,
		},
	}

	for _, tc := range testCases {
		t.Run(testName(tc.Dice), func(t *testing.T) {
			dice := diceFromValues(tc.Dice)
			score := chance(dice)
			assert.Equal(t, tc.ExpectedScore, score)
		})
	}
}

func TestYatzy(t *testing.T) {
	testCases := []struct {
		Dice          []uint
		ExpectedScore uint
	}{
		{
			Dice:          []uint{1, 2, 3, 4, 5},
			ExpectedScore: 0,
		},
		{
			Dice:          []uint{2, 5, 5, 5, 5},
			ExpectedScore: 0,
		},
		{
			Dice:          []uint{1, 1, 1, 1, 1},
			ExpectedScore: 50,
		},
		{
			Dice:          []uint{6, 6, 6, 6, 6},
			ExpectedScore: 50,
		},
	}

	for _, tc := range testCases {
		t.Run(testName(tc.Dice), func(t *testing.T) {
			dice := diceFromValues(tc.Dice)
			score := yatzy(dice)
			assert.Equal(t, tc.ExpectedScore, score)
		})
	}
}

func diceFromValues(values []uint) []*entity.Die {
	var dice []*entity.Die

	for _, v := range values {
		dice = append(dice, entity.MustNewDie(v))
	}

	return dice
}

func testName(values []uint) string {
	var s []string
	for _, v := range values {
		s = append(s, strconv.Itoa(int(v)))
	}
	return strings.Join(s, "_")
}
