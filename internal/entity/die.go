package entity

import (
	"math/rand"

	"github.com/m110/yatzy/internal/assets"

	"github.com/m110/yatzy/internal/component"
)

type Die struct {
	Value  uint
	Sprite component.Sprite
}

func NewRandomDie() Die {
	d := Die{}
	d.Roll()
	return d
}

func (d *Die) Roll() {
	value := 1 + rand.Intn(6)
	d.Value = uint(value)
	d.Sprite = component.Sprite{
		Image: assets.Dice[value-1],
	}
}
