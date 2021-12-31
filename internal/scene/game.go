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
	dice  []entity.Die
	table *Table
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
	ebiten.Key6: 5,
}

func (g *Game) Update() error {
	for k, v := range selectionKeys {
		if inpututil.IsKeyJustPressed(k) {
			g.dice[v].ToggleSelection()
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.RollSelectedDice()
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
}

func (g *Game) Draw(screen *ebiten.Image) {
	dicePanel := ebiten.NewImage(500, 100)
	dicePanel.Fill(colornames.Azure)

	offsetX := 10.0
	for _, d := range g.dice {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(offsetX, 10)
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
