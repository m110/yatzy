package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/m110/yatzy/internal/eventbus"

	"github.com/m110/yatzy/internal/scene"
)

type Game struct {
	eventBus *eventbus.EventBus
	client   *scene.Client
	server   *scene.Server
}

func NewGame() (*Game, error) {
	eventBus := eventbus.NewEventBus()
	server, err := scene.NewServer(eventBus, 1)
	if err != nil {
		return nil, err
	}

	client := scene.NewClient(eventBus)

	eventBus.SubscribeToEvents(client)
	eventBus.SubscribeToCommands(server)

	return &Game{
		eventBus: eventBus,
		server:   server,
		client:   client,
	}, nil
}

func (g *Game) Update() error {
	err := g.client.Update()
	if err != nil {
		return err
	}

	err = g.eventBus.Update()
	if err != nil {
		return err
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.client.Draw(screen)
}

func (g *Game) WindowSize() (int, int) {
	return scene.DiePanelWidth, scene.DiePanelHeight + scene.TablePanelHeight
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}
