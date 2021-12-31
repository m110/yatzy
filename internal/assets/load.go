package assets

import (
	"image"
	_ "image/png"
	"log"
	"os"
	"path"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	Dice = []*ebiten.Image{
		MustLoadSprite("dieWhite1.png"),
		MustLoadSprite("dieWhite2.png"),
		MustLoadSprite("dieWhite3.png"),
		MustLoadSprite("dieWhite4.png"),
		MustLoadSprite("dieWhite5.png"),
		MustLoadSprite("dieWhite6.png"),
	}
)

func MustLoadSprite(name string) *ebiten.Image {
	i, err := LoadSprite(name)
	if err != nil {
		// TODO Used in tests where the relative path differs
		log.Println(err)
	}
	return i
}

func LoadSprite(name string) (*ebiten.Image, error) {
	f, err := os.OpenFile(path.Join("assets", name), os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}

	image, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}

	return ebiten.NewImageFromImage(image), nil
}
