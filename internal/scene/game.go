package scene

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/m110/yatzy/internal/entity"
	"golang.org/x/image/colornames"
)

const (
	diceNumber = 5
)

type Game struct {
	dice    []entity.Die
	table   *Table
	rerolls int
}

func NewGame() *Game {
	var dice []entity.Die

	for i := 0; i < diceNumber; i++ {
		dice = append(dice, entity.NewRandomDie())
	}

	return &Game{
		dice:  dice,
		table: NewTable(),
	}
}

var selectionKeys = map[ebiten.Key]int{
	ebiten.Key1: 0,
	ebiten.Key2: 1,
	ebiten.Key3: 2,
	ebiten.Key4: 3,
	ebiten.Key5: 4,
}

func (g *Game) Update() error {
	if g.rerolls < 2 {
		for k, v := range selectionKeys {
			if inpututil.IsKeyJustPressed(k) {
				g.dice[v].ToggleSelection()
			}
		}

		if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			g.RollSelectedDice()
		}
	}

	return nil
}

func (g *Game) RollSelectedDice() {
	for i, d := range g.dice {
		if d.Selected {
			g.dice[i].ClearSelection()
			g.dice[i].Roll()
		}
	}
	g.rerolls++

	if g.rerolls == 2 {
		g.table.ShowAvailablePoints(g.dice)
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	dicePanel := ebiten.NewImage(500, 100)
	dicePanel.Fill(colornames.Azure)

	offsetX := 10.0
	offsetY := 10.0
	for _, d := range g.dice {
		if d.Selected {
			selectionImage := ebiten.NewImage(d.Sprite.Image.Bounds().Dx()+6, d.Sprite.Image.Bounds().Dy()+6)
			selectionImage.Fill(colornames.Yellow)
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(offsetX-3, offsetY-3)
			dicePanel.DrawImage(selectionImage, op)
		}

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(offsetX, offsetY)
		dicePanel.DrawImage(d.Sprite.Image, op)
		offsetX += float64(d.Sprite.Image.Bounds().Max.X) + 10.0
	}

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(50, 50)
	screen.DrawImage(dicePanel, op)

	tablePanel := ebiten.NewImage(300, 500)
	tablePanel.Fill(colornames.Coral)
	g.table.Draw(tablePanel)

	op = &ebiten.DrawImageOptions{}
	op.GeoM.Translate(10, 200)
	screen.DrawImage(tablePanel, op)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}
