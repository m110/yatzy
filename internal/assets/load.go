package assets

import (
	"image"
	_ "image/png"
	"log"
	"os"
	"path"

	"golang.org/x/image/font"

	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/font/opentype"
)

const dpi = 72

var (
	Dice = []*ebiten.Image{
		MustLoadSprite("dieWhite1.png"),
		MustLoadSprite("dieWhite2.png"),
		MustLoadSprite("dieWhite3.png"),
		MustLoadSprite("dieWhite4.png"),
		MustLoadSprite("dieWhite5.png"),
		MustLoadSprite("dieWhite6.png"),
	}

	NormalFont = MustLoadFont("Roboto-Regular.ttf")
)

func MustLoadFont(name string) font.Face {
	f, err := LoadFont(name)
	if err != nil {
		// TODO Used in tests where the relative path differs
		log.Println(err)
	}

	return f
}

func LoadFont(name string) (font.Face, error) {
	f, err := os.ReadFile(path.Join("assets", name))
	if err != nil {
		return nil, err
	}

	parsedFont, err := opentype.Parse(f)
	if err != nil {
		return nil, err
	}

	face, err := opentype.NewFace(parsedFont, &opentype.FaceOptions{
		Size:    16,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		return nil, err
	}

	return face, nil
}

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
