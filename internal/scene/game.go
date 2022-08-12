package scene

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/m110/yatzy/internal/component"
	"golang.org/x/image/colornames"

	"github.com/m110/yatzy/internal/entity"
)

const (
	diceNumber = 5

	diePanelWidth  = 64*diceNumber + (1+diceNumber)*dieOffsetX
	diePanelHeight = 64 + 2*dieOffsetY

	tablePanelWidth  = diePanelWidth
	tablePanelHeight = 500

	dieOffsetX = 10.0
	dieOffsetY = 10.0
)

type Game struct {
	dice     []*entity.Die
	dieIcons []entity.DieIcon
	table    *Table
	rerolls  int
}

func NewGame() *Game {
	var dice []*entity.Die
	var dieIcons []entity.DieIcon

	offsetX := dieOffsetX
	for i := 0; i < diceNumber; i++ {
		die := entity.NewRandomDie()
		dice = append(dice, die)

		position := component.Position{
			X: offsetX,
			Y: dieOffsetY,
		}

		dieIcon := entity.NewDieIcon(die, position)
		dieIcons = append(dieIcons, dieIcon)

		offsetX += float64(dieIcon.Sprite().Image.Bounds().Max.X) + dieOffsetX
	}

	return &Game{
		dice:     dice,
		dieIcons: dieIcons,
		table:    NewTable(),
		rerolls:  0,
	}
}

var selectionKeys = map[ebiten.Key]int{
	ebiten.Key1: 0,
	ebiten.Key2: 1,
	ebiten.Key3: 2,
	ebiten.Key4: 3,
	ebiten.Key5: 4,
}

func (g *Game) Update() error {
	if g.rerolls < 2 {
		for k, v := range selectionKeys {
			if inpututil.IsKeyJustPressed(k) {
				g.dice[v].ToggleSelection()
			}
		}

		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			x, y := ebiten.CursorPosition()
			for _, d := range g.dieIcons {
				if int(d.Position().X) < x &&
					int(d.Position().Y) < y &&
					int(d.Position().X)+d.Sprite().Image.Bounds().Dx() > x &&
					int(d.Position().Y)+d.Sprite().Image.Bounds().Dy() > y {
					d.Die().ToggleSelection()
				}
			}
		}

		if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			g.RollSelectedDice()
		}
	} else {
		err := g.table.Update()
		if err != nil {
			return err
		}

		if g.table.Full() {
			// TODO game over
			return nil
		}

		if g.table.Ready() {
			g.table.HideAvailablePoints()
			g.RollAllDice()
			g.rerolls = 0
		}
	}

	return nil
}

func (g *Game) RollAllDice() {
	for i := range g.dice {
		g.dice[i].ClearSelection()
		g.dice[i].Roll()
	}
}

func (g *Game) RollSelectedDice() {
	for i, d := range g.dice {
		if d.Selected {
			g.dice[i].ClearSelection()
			g.dice[i].Roll()
		}
	}
	g.rerolls++

	if g.rerolls == 2 {
		g.table.ShowAvailablePoints(g.dice)
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	var drawers []component.Drawer
	for _, d := range g.dieIcons {
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
	g.table.Draw(tablePanel)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(0, diePanelHeight)
	screen.DrawImage(tablePanel, op)
}

func (g *Game) WindowSize() (int, int) {
	return diePanelWidth, diePanelHeight + tablePanelHeight
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}
