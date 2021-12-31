package main

import (
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/m110/yatzy/internal/scene"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	ebiten.SetWindowSize(800, 600)

	err := ebiten.RunGame(scene.NewGame())
	if err != nil {
		log.Fatal(err)
	}
}
