package ptls_amqp

import (
	"errors"
	"log"

	"github.com/assembla/cony"
)

type Producer struct {
	config    Config
	exchange  string
	kind      string
	key       string
	publisher *cony.Publisher
}

func (p Producer) Publish(msg Publishing) error {
	if p.publisher == nil {
		return errors.New("amqp客户端不存在")
	}

	return p.publisher.Publish(msg.AmqpPublishing())
}

func NewProducer(cfg Config, exchange, kind, key string) (*Producer, error) {
	pdcer := Producer{
		config:   cfg,
		exchange: exchange,
		kind:     kind,
		key:      key,
	}

	url := cfg.String()

	cli := cony.NewClient(
		cony.URL(url),
		cony.Backoff(cony.DefaultBackoff),
	)

	// Declare the exchange we'll be using
	exc := cony.Exchange{
		Name:       exchange,
		Kind:       kind,
		AutoDelete: false,
	}
	cli.Declare([]cony.Declaration{
		cony.DeclareExchange(exc),
	})

	// Declare and register a publisher
	// with the cony client
	pbl := cony.NewPublisher(exc.Name, key)
	cli.Publish(pbl)
	pdcer.publisher = pbl

	// Client loop sends out declarations(exchanges, queues, bindings
	// etc) to the AMQP server. It also handles reconnecting.
	go func() {
		for cli.Loop() {
			select {
			case err := <-cli.Errors():
				log.Printf("Client error: %v\n", err)
			case blocked := <-cli.Blocking():
				log.Printf("Client is blocked %v\n", blocked)
			}
		}
	}()

	return &pdcer, nil
}
