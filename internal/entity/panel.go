package entity

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/m110/yatzy/internal/component"
)

type Panel struct {
	rect            component.Rect
	backgroundColor color.Color
	drawers         []component.Drawer
}

func NewPanel(rect component.Rect, backgroundColor color.Color, drawers []component.Drawer) Panel {
	return Panel{
		rect:            rect,
		backgroundColor: backgroundColor,
		drawers:         drawers,
	}
}

func (p Panel) Draw(screen *ebiten.Image) {
	panel := ebiten.NewImage(int(p.rect.Size.Width), int(p.rect.Size.Height))
	panel.Fill(p.backgroundColor)

	for _, d := range p.drawers {
		d.Draw(panel)
	}

	op := &ebiten.DrawImageOptions{}
	screen.DrawImage(panel, op)
}
