package scene

import (
	"log"
	"math/rand"

	"github.com/google/uuid"

	"github.com/m110/yatzy/internal/entity"
	"github.com/m110/yatzy/internal/transport"
)

const (
	diceNumber = 5
)

type State int

const (
	WaitingForRoll State = iota
	WaitingForReroll
	WaitingForScore
)

type Player struct {
	ID      string
	Dice    []*entity.Die
	Scores  map[transport.ScoreType]*uint
	Rerolls int
}

func NewPlayer() *Player {
	p := &Player{
		ID:      uuid.NewString(),
		Dice:    make([]*entity.Die, diceNumber),
		Rerolls: 0,
	}

	for i := 0; i < diceNumber; i++ {
		p.Dice[i] = entity.NewRandomDie()
	}

	return p
}

type EventPublisher interface {
	Publish(event any) error
}

type Game struct {
	players       []*Player
	currentPlayer int
	state         State

	publisher EventPublisher
}

func NewGame(publisher EventPublisher, players int) (*Game, error) {
	g := &Game{
		players:       make([]*Player, players),
		currentPlayer: rand.Intn(players),
		state:         WaitingForRoll,

		publisher: publisher,
	}

	for i := 0; i < players; i++ {
		g.players[i] = NewPlayer()
	}

	eventPlayers := make([]transport.Player, players)
	for i, p := range g.players {
		eventPlayers[i] = transport.Player{
			ID: p.ID,
		}
	}
	err := publisher.Publish(transport.GameStarted{
		Players: eventPlayers,
	})
	if err != nil {
		return nil, err
	}

	err = publisher.Publish(transport.TurnStarted{
		Player: transport.Player{
			ID: g.players[g.currentPlayer].ID,
		},
	})
	if err != nil {
		return nil, err
	}

	return g, nil
}

func (g *Game) HandleCommand(player Player, command any) error {
	if g.players[g.currentPlayer].ID != player.ID {
		log.Println("invalid player")
		return nil
	}

	switch c := command.(type) {
	case transport.RollDice:
		return g.handleRollDice(c)
	case transport.RerollDice:
		return g.handleRerollDice(c)
	case transport.PassRoll:
		return g.handlePassRoll(c)
	case transport.AssignScore:
		return g.handleAssignScore(c)
	}

	return nil
}

func (g *Game) handleRollDice(command transport.RollDice) error {
	if g.state != WaitingForRoll {
		log.Println("not waiting for roll")
		return nil
	}

	for _, d := range g.players[g.currentPlayer].Dice {
		d.Roll()
	}

	diceRolled := make([]transport.Die, len(g.players[g.currentPlayer].Dice))
	for i, d := range g.players[g.currentPlayer].Dice {
		diceRolled[i] = transport.Die{
			ID:    i,
			Value: d.Value(),
		}
	}
	err := g.publisher.Publish(transport.DiceRolled{
		Dice: diceRolled,
	})
	if err != nil {
		return err
	}

	g.state = WaitingForReroll
	return nil
}

func (g *Game) handleRerollDice(command transport.RerollDice) error {
	if g.state != WaitingForReroll {
		log.Println("not waiting for reroll")
		return nil
	}

	if g.players[g.currentPlayer].Rerolls >= 2 {
		log.Println("no more rerolls")
		return nil
	}

	for _, i := range command.Dice {
		g.players[g.currentPlayer].Dice[i].Roll()
	}

	diceRolled := make([]transport.Die, len(command.Dice))
	for i, d := range command.Dice {
		diceRolled[i] = transport.Die{
			ID:    d,
			Value: g.players[g.currentPlayer].Dice[d].Value(),
		}
	}
	err := g.publisher.Publish(transport.DiceRolled{
		Dice: diceRolled,
	})
	if err != nil {
		return err
	}

	g.players[g.currentPlayer].Rerolls++

	if g.players[g.currentPlayer].Rerolls >= 2 {
		err := g.publishScoresReady()
		if err != nil {
			return err
		}
		g.state = WaitingForScore
	}

	return nil
}

func (g *Game) handlePassRoll(command transport.PassRoll) error {
	if g.state != WaitingForReroll {
		log.Println("not waiting for reroll (pass)")
		return nil
	}

	err := g.publishScoresReady()
	if err != nil {
		return err
	}

	g.state = WaitingForScore
	return nil
}

func (g *Game) publishScoresReady() error {
	var scores []transport.Score
	for t, f := range scoringFunctions {
		scores = append(scores, transport.Score{
			Type:  t,
			Value: f(g.players[g.currentPlayer].Dice),
		})
	}

	err := g.publisher.Publish(transport.ScoresReady{
		Scores: scores,
	})
	if err != nil {
		return err
	}

	return nil
}

func (g *Game) handleAssignScore(command transport.AssignScore) error {
	if g.state != WaitingForScore {
		log.Println("not waiting for score")
		return nil
	}

	if g.players[g.currentPlayer].Scores[command.Type] != nil {
		log.Println("score already assigned")
		return nil
	}

	scoreFun := scoringFunctions[command.Type]
	score := scoreFun(g.players[g.currentPlayer].Dice)
	g.players[g.currentPlayer].Scores[command.Type] = &score

	err := g.publisher.Publish(transport.ScoreAssigned{
		Player: transport.Player{
			ID: g.players[g.currentPlayer].ID,
		},
		Score: transport.Score{
			Type:  command.Type,
			Value: score,
		},
	})

	g.state = WaitingForRoll
	g.currentPlayer = (g.currentPlayer + 1) % len(g.players)

	err = g.publisher.Publish(transport.TurnStarted{
		Player: transport.Player{
			ID: g.players[g.currentPlayer].ID,
		},
	})
	if err != nil {
		return err
	}

	return nil
}

var scoringFunctions = map[transport.ScoreType]scoringFunc{
	transport.ScoreTypeOnes:          ones,
	transport.ScoreTypeTwos:          twos,
	transport.ScoreTypeThrees:        threes,
	transport.ScoreTypeFours:         fours,
	transport.ScoreTypeFives:         fives,
	transport.ScoreTypeSixes:         sixes,
	transport.ScoreTypeThreeOfAKind:  threeOfAKind,
	transport.ScoreTypeFourOfAKind:   fourOfAKind,
	transport.ScoreTypeFullHouse:     fullHouse,
	transport.ScoreTypeSmallStraight: smallStraight,
	transport.ScoreTypeLargeStraight: largeStraight,
	transport.ScoreTypeYatzy:         yatzy,
	transport.ScoreTypeChance:        chance,
}
