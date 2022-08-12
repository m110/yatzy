package assets

import (
	"image"
	_ "image/png"
	"io/fs"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"

	"github.com/m110/yatzy/assets"
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
	b, err := fs.ReadFile(assets.FS, name)
	if err != nil {
		return nil, err
	}

	parsedFont, err := opentype.Parse(b)
	if err != nil {
		return nil, err
	}

	face, err := opentype.NewFace(parsedFont, &opentype.FaceOptions{
		Size:    24,
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
	f, err := assets.FS.Open(name)

	if err != nil {
		return nil, err
	}

	image, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}

	return ebiten.NewImageFromImage(image), nil
}
