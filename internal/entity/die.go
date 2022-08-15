package entity

import (
	"log"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/colornames"

	"github.com/m110/yatzy/internal/assets"
	"github.com/m110/yatzy/internal/component"
	"github.com/m110/yatzy/internal/timer"
)

type Die struct {
	value uint
}

func MustNewDie(value uint) *Die {
	if value < 1 || value > 6 {
		log.Fatalf("invalid die value: %v", value)
	}

	return &Die{
		value: value,
	}
}

func NewRandomDie() *Die {
	d := &Die{}
	d.Roll()
	return d
}

func (d *Die) Roll() {
	value := 1 + rand.Intn(6)
	d.value = uint(value)
}

func (d *Die) Value() uint {
	return d.value
}

type DieIcon struct {
	die                  *Die
	selected             bool
	position             component.Position
	isRolling            bool
	rollingNextFaceTimer *timer.Timer
	rollingTimer         *timer.Timer
}

func NewDieIcon(die *Die, position component.Position) *DieIcon {
	return &DieIcon{
		die:      die,
		position: position,

		isRolling:            false,
		rollingNextFaceTimer: timer.NewTimer(5),
		rollingTimer:         timer.NewTimer(60),
	}
}

func (d *DieIcon) Selected() bool {
	return d.selected
}

func (d *DieIcon) Die() *Die {
	return d.die
}

func (d *DieIcon) IsRolling() bool {
	return d.isRolling
}

func (d *DieIcon) Update() {
	if d.isRolling {
		d.rollingNextFaceTimer.Update()
		if d.rollingNextFaceTimer.IsDone() {
			d.die.Roll()
			d.rollingNextFaceTimer.Reset()
		}

		d.rollingTimer.Update()
		if d.rollingTimer.IsDone() {
			d.isRolling = false
		}
	}
}

func (d *DieIcon) Sprite() component.Sprite {
	return component.Sprite{
		Image: assets.Dice[d.die.Value()-1],
	}
}

func (d *DieIcon) Position() component.Position {
	return d.position
}

func (d *DieIcon) Roll() {
	d.clearSelection()

	d.rollingNextFaceTimer.Reset()
	d.rollingTimer.Reset()
	d.isRolling = true
}

func (d *DieIcon) OnClick() {
	d.toggleSelection()
}

func (d *DieIcon) toggleSelection() {
	d.selected = !d.selected
}

func (d *DieIcon) clearSelection() {
	d.selected = false
}

const (
	highlightBorderOffset = 10.0
	highlightBorderSize   = 20.0
)

func (d *DieIcon) Draw(screen *ebiten.Image) {
	sprite := d.Sprite()

	if d.selected {
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
