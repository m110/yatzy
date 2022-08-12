package entity

import (
	"log"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/colornames"

	"github.com/m110/yatzy/internal/assets"
	"github.com/m110/yatzy/internal/component"
)

type Die struct {
	Value    uint
	Selected bool
}

func MustNewDie(value uint) *Die {
	if value < 1 || value > 6 {
		log.Fatalf("invalid die value: %v", value)
	}

	return &Die{
		Value: value,
	}
}

func NewRandomDie() *Die {
	d := &Die{}
	d.Roll()
	return d
}

func (d *Die) Roll() {
	value := 1 + rand.Intn(6)
	d.Value = uint(value)
}

func (d *Die) ToggleSelection() {
	d.Selected = !d.Selected
}

func (d *Die) ClearSelection() {
	d.Selected = false
}

type DieIcon struct {
	die      *Die
	position component.Position
}

func NewDieIcon(die *Die, position component.Position) DieIcon {
	return DieIcon{
		die:      die,
		position: position,
	}
}

func (d DieIcon) Sprite() component.Sprite {
	return component.Sprite{
		Image: assets.Dice[d.die.Value-1],
	}
}

func (d DieIcon) Die() *Die {
	return d.die
}

func (d DieIcon) Position() component.Position {
	return d.position
}

const (
	highlightBorderOffset = 10.0
	highlightBorderSize   = 20.0
)

func (d DieIcon) Draw(screen *ebiten.Image) {
	sprite := d.Sprite()

	if d.die.Selected {
		selectionImage := ebiten.NewImage(sprite.Image.Bounds().Dx(), highlightBorderSize)
		selectionImage.Fill(colornames.Blueviolet)
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(d.position.X, d.position.Y+float64(sprite.Image.Bounds().Dy())-highlightBorderOffset)
		screen.DrawImage(selectionImage, op)
	}

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(d.position.X, d.position.Y)
	screen.DrawImage(sprite.Image, op)
}
