package scene

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"golang.org/x/image/colornames"

	"github.com/m110/yatzy/internal/entity"
)

const (
	diceNumber = 5
)

type Game struct {
	dice    []entity.Die
	table   *Table
	rerolls int
}

func NewGame() *Game {
	var dice []entity.Die

	for i := 0; i < diceNumber; i++ {
		dice = append(dice, entity.NewRandomDie())
	}

	return &Game{
		dice:    dice,
		table:   NewTable(),
		rerolls: 0,
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

const (
	diePanelWidth  = 64*diceNumber + (1+diceNumber)*dieOffsetX
	diePanelHeight = 64 + 2*dieOffsetY

	tablePanelWidth  = diePanelWidth
	tablePanelHeight = 500

	dieOffsetX = 10.0
	dieOffsetY = 10.0

	highlightBorderOffset = 10.0
	highlightBorderSize   = 20.0
)

func (g *Game) Draw(screen *ebiten.Image) {
	dicePanel := ebiten.NewImage(diePanelWidth, diePanelHeight)
	dicePanel.Fill(colornames.Forestgreen)

	offsetX := dieOffsetX
	offsetY := dieOffsetY
	for _, d := range g.dice {
		if d.Selected {
			selectionImage := ebiten.NewImage(d.Sprite.Image.Bounds().Dx(), highlightBorderSize)
			selectionImage.Fill(colornames.Blueviolet)
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(offsetX, offsetY+float64(d.Sprite.Image.Bounds().Dy())-highlightBorderOffset)
			dicePanel.DrawImage(selectionImage, op)
		}

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(offsetX, offsetY)
		dicePanel.DrawImage(d.Sprite.Image, op)
		offsetX += float64(d.Sprite.Image.Bounds().Max.X) + dieOffsetX
	}

	op := &ebiten.DrawImageOptions{}
	screen.DrawImage(dicePanel, op)

	tablePanel := ebiten.NewImage(tablePanelWidth, tablePanelHeight)
	tablePanel.Fill(colornames.Darkgreen)
	g.table.Draw(tablePanel)

	op = &ebiten.DrawImageOptions{}
	op.GeoM.Translate(0, diePanelHeight)
	screen.DrawImage(tablePanel, op)
}

func (g *Game) WindowSize() (int, int) {
	return diePanelWidth, diePanelHeight + tablePanelHeight
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}
