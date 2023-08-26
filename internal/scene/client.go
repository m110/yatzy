package scene

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/m110/yatzy/internal/component"
	"github.com/m110/yatzy/internal/entity"
	"github.com/m110/yatzy/internal/transport"
	"golang.org/x/image/colornames"
)

const (
	diceNumberUI = 5

	DiePanelWidth  = 64*diceNumber + (1+diceNumber)*dieOffsetX
	DiePanelHeight = 64 + 2*dieOffsetY

	TablePanelWidth  = DiePanelWidth
	TablePanelHeight = 500

	dieOffsetX = 10.0
	dieOffsetY = 10.0
)

type CommandPublisher interface {
	PublishCommand(command any)
}

type Client struct {
	dieIcons  []*entity.DieIcon
	table     *Table
	publisher CommandPublisher
}

func NewClient(publisher CommandPublisher) *Client {
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

	return &Client{
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

func (c *Client) HandleEvent(event any) error {
	switch event := event.(type) {
	case transport.DiceRolled:
		_ = event
	}

	return nil
}

func (c *Client) Update() error {
	for _, d := range c.dieIcons {
		d.Update()
	}

	diceRolling := false
	for _, d := range c.dieIcons {
		if d.IsRolling() {
			diceRolling = true
			break
		}
	}

	if diceRolling {
		return nil
	}

	if false {
		for k, v := range selectionKeys {
			if inpututil.IsKeyJustPressed(k) {
				c.dieIcons[v].OnClick()
			}
		}

		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			x, y := ebiten.CursorPosition()
			for _, d := range c.dieIcons {
				if int(d.Position().X) < x &&
					int(d.Position().Y) < y &&
					int(d.Position().X)+d.Sprite().Image.Bounds().Dx() > x &&
					int(d.Position().Y)+d.Sprite().Image.Bounds().Dy() > y {
					d.OnClick()
				}
			}
		}

		if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			//c.RollSelectedDice()
		}
	} else {
		if !c.table.ShowingAvailablePoints {
			var dice []*entity.Die
			for _, d := range c.dieIcons {
				dice = append(dice, d.Die())
			}

			c.table.ShowAvailablePoints(dice)
		}

		err := c.table.Update()
		if err != nil {
			return err
		}

		if c.table.Full() {
			// TODO game over
			return nil
		}

		if c.table.Ready() {
			c.table.HideAvailablePoints()
			//c.RollAllDice()
			//c.rerolls = 0
		}
	}

	return nil
}

func (c *Client) rollDice() error {
	/*
		for _, d := range g.dieIcons {
			d.Roll()
		}

	*/

	return nil
}

func (c *Client) rerollDice() error {
	/*
		for _, d := range g.dieIcons {
			if d.Selected() {
				d.Roll()
			}
		}
		g.rerolls++

	*/
	return nil
}

func (c *Client) Draw(screen *ebiten.Image) {
	var drawers []component.Drawer
	for _, d := range c.dieIcons {
		drawers = append(drawers, d)
	}
	dicePanel := entity.NewPanel(component.Rect{
		Position: component.Position{
			X: 0,
			Y: 0,
		},
		Size: component.Size{
			Width:  DiePanelWidth,
			Height: DiePanelHeight,
		},
	}, colornames.Forestgreen, drawers)
	dicePanel.Draw(screen)

	tablePanel := ebiten.NewImage(TablePanelWidth, TablePanelHeight)
	tablePanel.Fill(colornames.Darkgreen)
	c.table.Draw(tablePanel)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(0, DiePanelHeight)
	screen.DrawImage(tablePanel, op)
}
