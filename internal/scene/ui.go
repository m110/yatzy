package scene

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"golang.org/x/image/colornames"

	"github.com/m110/yatzy/internal/component"
	"github.com/m110/yatzy/internal/entity"
)

const (
	diceNumberUI = 5

	diePanelWidth  = 64*diceNumber + (1+diceNumber)*dieOffsetX
	diePanelHeight = 64 + 2*dieOffsetY

	tablePanelWidth  = diePanelWidth
	tablePanelHeight = 500

	dieOffsetX = 10.0
	dieOffsetY = 10.0
)

type CommandPublisher interface {
	PublishCommand(command any) error
}

type UI struct {
	dieIcons  []*entity.DieIcon
	table     *Table
	publisher CommandPublisher
}

func NewUI(publisher CommandPublisher) *UI {
	var dieIcons []*entity.DieIcon

	offsetX := dieOffsetX
	for i := 0; i < diceNumberUI; i++ {

		position := component.Position{
			X: offsetX,
			Y: dieOffsetY,
		}

		dieIcon := entity.NewDieIcon(entity.NewRandomDie(), position)
		dieIcons = append(dieIcons, dieIcon)

		offsetX += float64(dieIcon.Sprite().Image.Bounds().Max.X) + dieOffsetX
	}

	return &UI{
		dieIcons:  dieIcons,
		table:     NewTable(),
		publisher: publisher,
	}
}

var selectionKeys = map[ebiten.Key]int{
	ebiten.Key1: 0,
	ebiten.Key2: 1,
	ebiten.Key3: 2,
	ebiten.Key4: 3,
	ebiten.Key5: 4,
}

func (u *UI) Update() error {
	for _, d := range u.dieIcons {
		d.Update()
	}

	diceRolling := false
	for _, d := range u.dieIcons {
		if d.IsRolling() {
			diceRolling = true
			break
		}
	}

	if diceRolling {
		return nil
	}

	if u.rerolls < 2 {
		for k, v := range selectionKeys {
			if inpututil.IsKeyJustPressed(k) {
				u.dieIcons[v].OnClick()
			}
		}

		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			x, y := ebiten.CursorPosition()
			for _, d := range u.dieIcons {
				if int(d.Position().X) < x &&
					int(d.Position().Y) < y &&
					int(d.Position().X)+d.Sprite().Image.Bounds().Dx() > x &&
					int(d.Position().Y)+d.Sprite().Image.Bounds().Dy() > y {
					d.OnClick()
				}
			}
		}

		if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			u.RollSelectedDice()
		}
	} else {
		if !u.table.ShowingAvailablePoints {
			var dice []*entity.Die
			for _, d := range u.dieIcons {
				dice = append(dice, d.Die())
			}

			u.table.ShowAvailablePoints(dice)
		}

		err := u.table.Update()
		if err != nil {
			return err
		}

		if u.table.Full() {
			// TODO game over
			return nil
		}

		if u.table.Ready() {
			u.table.HideAvailablePoints()
			u.RollAllDice()
			u.rerolls = 0
		}
	}

	return nil
}

func (u *UI) rollDice() error {
	for _, d := range g.dieIcons {
		d.Roll()
	}
}

func (u *UI) rerollDice() error {
	for _, d := range g.dieIcons {
		if d.Selected() {
			d.Roll()
		}
	}
	g.rerolls++
}

func (u *UI) Draw(screen *ebiten.Image) {
	var drawers []component.Drawer
	for _, d := range u.dieIcons {
		drawers = append(drawers, d)
	}
	dicePanel := entity.NewPanel(component.Rect{
		Position: component.Position{
			X: 0,
			Y: 0,
		},
		Size: component.Size{
			Width:  diePanelWidth,
			Height: diePanelHeight,
		},
	}, colornames.Forestgreen, drawers)
	dicePanel.Draw(screen)

	tablePanel := ebiten.NewImage(tablePanelWidth, tablePanelHeight)
	tablePanel.Fill(colornames.Darkgreen)
	u.table.Draw(tablePanel)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(0, diePanelHeight)
	screen.DrawImage(tablePanel, op)
}

func (u *UI) WindowSize() (int, int) {
	return diePanelWidth, diePanelHeight + tablePanelHeight
}

func (u *UI) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}
