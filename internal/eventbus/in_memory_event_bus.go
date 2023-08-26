package eventbus

import (
	"log"

	"github.com/m110/yatzy/internal/scene"
)

type EventSubscriber interface {
	HandleEvent(event any) error
}

type CommandSubscriber interface {
	HandleCommand(player scene.Player, cmd any) error
}

type EventBus struct {
	eventSubscribers []EventSubscriber
	cmdSubscribers   []CommandSubscriber

	eventsQueue []any
	cmdsQueue   []any
}

func NewEventBus() *EventBus {
	return &EventBus{}
}

func (e *EventBus) SubscribeToEvents(subscriber EventSubscriber) {
	e.eventSubscribers = append(e.eventSubscribers, subscriber)
}

func (e *EventBus) SubscribeToCommands(subscriber CommandSubscriber) {
	e.cmdSubscribers = append(e.cmdSubscribers, subscriber)
}

func (e *EventBus) PublishEvent(event any) {
	e.eventsQueue = append(e.eventsQueue, event)
}

func (e *EventBus) PublishCommand(cmd any) {
	e.cmdsQueue = append(e.cmdsQueue, cmd)
}

func (e *EventBus) Update() error {
	for _, event := range e.eventsQueue {
		for _, subscriber := range e.eventSubscribers {
			err := subscriber.HandleEvent(event)
			if err != nil {
				log.Println("error handling event:", err)
			}
		}
	}

	for _, cmd := range e.cmdsQueue {
		for _, subscriber := range e.cmdSubscribers {
			err := subscriber.HandleCommand(scene.Player{}, cmd)
			if err != nil {
				log.Println("error handling command:", err)
			}
		}
	}

	e.eventsQueue = nil
	e.cmdsQueue = nil

	return nil
}
