package entity

import (
	"log"
	"math/rand"

	"github.com/m110/yatzy/internal/assets"

	"github.com/m110/yatzy/internal/component"
)

type Die struct {
	Value  uint
	Sprite component.Sprite
}

func MustNewDie(value uint) Die {
	if value < 1 || value > 6 {
		log.Fatalf("invalid die value: %v", value)
	}

	return Die{
		Value: value,
		Sprite: component.Sprite{
			Image: assets.Dice[value-1],
		},
	}
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
