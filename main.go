package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/m110/yatzy/internal/game"
)

func main() {
	g, err := game.NewGame()
	if err != nil {
		log.Fatal(err)
	}

	ebiten.SetWindowSize(g.WindowSize())

	err = ebiten.RunGame(g)
	if err != nil {
		log.Fatal(err)
	}
}
