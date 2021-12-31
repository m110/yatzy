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
	dice []entity.Die
}

func NewGame() *Game {
	var dice []entity.Die

	for i := 0; i < diceNumber; i++ {
		dice = append(dice, entity.NewRandomDie())
	}

	return &Game{
		dice: dice,
	}
}

func (g *Game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.RollDice()
	}

	return nil
}

func (g *Game) RollDice() {
	for i := range g.dice {
		g.dice[i].Roll()
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
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}
