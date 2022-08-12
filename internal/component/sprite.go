package component

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Sprite struct {
	Image *ebiten.Image
}

type Drawer interface {
	Draw(image *ebiten.Image)
}
