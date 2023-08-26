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
	PublishEvent(event any)
}

type Server struct {
	players       []*Player
	currentPlayer int
	state         State

	publisher EventPublisher
}

func NewServer(publisher EventPublisher, players int) (*Server, error) {
	g := &Server{
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
	publisher.PublishEvent(transport.GameStarted{
		Players: eventPlayers,
	})

	publisher.PublishEvent(transport.TurnStarted{
		Player: transport.Player{
			ID: g.players[g.currentPlayer].ID,
		},
	})

	return g, nil
}

func (s *Server) HandleCommand(player Player, command any) error {
	if s.players[s.currentPlayer].ID != player.ID {
		log.Println("invalid player")
		return nil
	}

	switch c := command.(type) {
	case transport.RollDice:
		return s.handleRollDice(c)
	case transport.RerollDice:
		return s.handleRerollDice(c)
	case transport.PassRoll:
		return s.handlePassRoll(c)
	case transport.AssignScore:
		return s.handleAssignScore(c)
	}

	return nil
}

func (s *Server) handleRollDice(command transport.RollDice) error {
	if s.state != WaitingForRoll {
		log.Println("not waiting for roll")
		return nil
	}

	for _, d := range s.players[s.currentPlayer].Dice {
		d.Roll()
	}

	diceRolled := make([]transport.Die, len(s.players[s.currentPlayer].Dice))
	for i, d := range s.players[s.currentPlayer].Dice {
		diceRolled[i] = transport.Die{
			ID:    i,
			Value: d.Value(),
		}
	}
	s.publisher.PublishEvent(transport.DiceRolled{
		Dice: diceRolled,
	})

	s.state = WaitingForReroll
	return nil
}

func (s *Server) handleRerollDice(command transport.RerollDice) error {
	if s.state != WaitingForReroll {
		log.Println("not waiting for reroll")
		return nil
	}

	if s.players[s.currentPlayer].Rerolls >= 2 {
		log.Println("no more rerolls")
		return nil
	}

	for _, i := range command.Dice {
		s.players[s.currentPlayer].Dice[i].Roll()
	}

	diceRolled := make([]transport.Die, len(command.Dice))
	for i, d := range command.Dice {
		diceRolled[i] = transport.Die{
			ID:    d,
			Value: s.players[s.currentPlayer].Dice[d].Value(),
		}
	}
	s.publisher.PublishEvent(transport.DiceRolled{
		Dice: diceRolled,
	})

	s.players[s.currentPlayer].Rerolls++

	if s.players[s.currentPlayer].Rerolls >= 2 {
		err := s.publishScoresReady()
		if err != nil {
			return err
		}
		s.state = WaitingForScore
	}

	return nil
}

func (s *Server) handlePassRoll(command transport.PassRoll) error {
	if s.state != WaitingForReroll {
		log.Println("not waiting for reroll (pass)")
		return nil
	}

	err := s.publishScoresReady()
	if err != nil {
		return err
	}

	s.state = WaitingForScore
	return nil
}

func (s *Server) publishScoresReady() error {
	var scores []transport.Score
	for t, f := range scoringFunctions {
		scores = append(scores, transport.Score{
			Type:  t,
			Value: f(s.players[s.currentPlayer].Dice),
		})
	}

	s.publisher.PublishEvent(transport.ScoresReady{
		Scores: scores,
	})

	return nil
}

func (s *Server) handleAssignScore(command transport.AssignScore) error {
	if s.state != WaitingForScore {
		log.Println("not waiting for score")
		return nil
	}

	if s.players[s.currentPlayer].Scores[command.Type] != nil {
		log.Println("score already assigned")
		return nil
	}

	scoreFun := scoringFunctions[command.Type]
	score := scoreFun(s.players[s.currentPlayer].Dice)
	s.players[s.currentPlayer].Scores[command.Type] = &score

	s.publisher.PublishEvent(transport.ScoreAssigned{
		Player: transport.Player{
			ID: s.players[s.currentPlayer].ID,
		},
		Score: transport.Score{
			Type:  command.Type,
			Value: score,
		},
	})

	s.state = WaitingForRoll
	s.currentPlayer = (s.currentPlayer + 1) % len(s.players)

	s.publisher.PublishEvent(transport.TurnStarted{
		Player: transport.Player{
			ID: s.players[s.currentPlayer].ID,
		},
	})

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
