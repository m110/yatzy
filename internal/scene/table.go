package scene

import (
	"image/color"
	"log"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/colornames"

	"github.com/m110/yatzy/internal/assets"
	"github.com/m110/yatzy/internal/entity"
)

type Table struct {
	Boxes                  []*Box
	ShowingAvailablePoints bool
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

func (t *Table) Update() error {
	var available []*Box
	for _, b := range t.Boxes {
		if !b.Filled {
			available = append(available, b)
		}
	}

	if len(available) > 1 {
		if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
			for i, b := range available {
				if b.Selected {
					var next int
					if i == len(available)-1 {
						next = 0
					} else {
						next = i + 1
					}
					b.Selected = false
					available[next].Selected = true
					break
				}
			}
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) {
			for i, b := range available {
				if b.Selected {
					var next int
					if i == 0 {
						next = len(available) - 1
					} else {
						next = i - 1
					}
					b.Selected = false
					available[next].Selected = true
					break
				}
			}
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		for _, b := range t.Boxes {
			if b.Selected {
				b.Fill()
				break
			}
		}
	}
	return nil
}

func (t *Table) Ready() bool {
	for _, b := range t.Boxes {
		if b.Selected {
			return false
		}
	}

	return true
}

func (t *Table) Full() bool {
	for _, b := range t.Boxes {
		if !b.Filled {
			return false
		}
	}

	return true
}

const (
	tableStartOffsetY = 30
	tableOffsetY      = 30
	tableOffsetX      = 10

	tableScoreOffsetX = 200
)

func (t *Table) Draw(screen *ebiten.Image) {
	offsetY := tableStartOffsetY

	for _, b := range t.Boxes {
		var c color.RGBA
		if b.Selected {
			c = colornames.Yellow
		} else {
			c = colornames.White
		}
		text.Draw(screen, b.Name, assets.NormalFont, tableOffsetX, offsetY, c)

		if b.Filled {
			text.Draw(screen, strconv.Itoa(int(b.Points)), assets.NormalFont, tableScoreOffsetX, offsetY, colornames.White)
		} else if t.ShowingAvailablePoints {
			var c color.RGBA
			if b.AvailablePoints == 0 {
				c = colornames.Gray
			} else {
				c = colornames.Yellow
			}
			text.Draw(screen, strconv.Itoa(int(b.AvailablePoints)), assets.NormalFont, tableScoreOffsetX, offsetY, c)
		}

		offsetY += tableOffsetY
	}
}

func (t *Table) ShowAvailablePoints(dice []entity.Die) {
	for i, b := range t.Boxes {
		if b.Filled {
			continue
		}

		t.Boxes[i].AvailablePoints = b.Scoring(dice)
	}

	t.ShowingAvailablePoints = true

	for i, b := range t.Boxes {
		if !b.Filled {
			t.Boxes[i].Selected = true
			break
		}
	}
}

func (t *Table) HideAvailablePoints() {
	t.ShowingAvailablePoints = false
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
	Name            string
	Filled          bool
	Points          uint
	AvailablePoints uint
	Scoring         scoringFunc
	Selected        bool
}

func (b *Box) Fill() {
	if b.Filled {
		log.Fatal("Box already filled")
	}

	b.Filled = true
	b.Points = b.AvailablePoints
	b.AvailablePoints = 0
	b.Selected = false
}
