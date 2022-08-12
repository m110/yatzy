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
	game := scene.NewGame()
	ebiten.SetWindowSize(game.WindowSize())

	err := ebiten.RunGame(game)
	if err != nil {
		log.Fatal(err)
	}
}
