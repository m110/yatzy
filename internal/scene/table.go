package scene

import (
	"log"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/m110/yatzy/internal/assets"
	"golang.org/x/image/colornames"

	"github.com/hajimehoshi/ebiten/v2"
)

type Table struct {
	Boxes []*Box
}

func NewTable() *Table {
	boxes := []*Box{
		// Upper section
		{Name: "Ones", Scoring: ones},
		{Name: "Twos", Scoring: twos},
		{Name: "Threes", Scoring: threes},
		{Name: "Fours", Scoring: fours},
		{Name: "Fives", Scoring: fives},
		{Name: "Sixes", Scoring: sixes},
		// Lower section
		{Name: "One Pair", Scoring: onePair},
		{Name: "Two Pairs", Scoring: twoPairs},
		{Name: "Three of a Kind", Scoring: threeOfAKind},
		{Name: "Four of a Kind", Scoring: fourOfAKind},
		{Name: "Small Straight", Scoring: smallStraight},
		{Name: "Large Straight", Scoring: largeStraight},
		{Name: "Full House", Scoring: fullHouse},
		{Name: "Chance", Scoring: chance},
		{Name: "Yatzy", Scoring: yatzy},
	}

	return &Table{
		Boxes: boxes,
	}
}

func (t *Table) Draw(screen *ebiten.Image) {
	offsetY := 20

	for _, b := range t.Boxes {
		text.Draw(screen, b.Name, assets.NormalFont, 10, offsetY, colornames.White)

		if b.Filled {
			text.Draw(screen, strconv.Itoa(int(b.Points)), assets.NormalFont, 100, offsetY, colornames.White)
		}

		offsetY += 20
	}
}

func (t *Table) UpperSectionBonus() uint {
	var s uint
	for _, b := range t.Boxes[:6] {
		if b.Filled {
			s += b.Points
		}
	}

	if s > 63 {
		return 50
	}

	return 0
}

type Box struct {
	Name    string
	Filled  bool
	Points  uint
	Scoring scoringFunc
}

func (b *Box) Fill(points uint) {
	if b.Filled {
		log.Fatal("Box already filled")
	}

	b.Filled = true
	b.Points = points
}
