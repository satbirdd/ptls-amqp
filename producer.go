package ptls_amqp

import (
	"errors"
	"log"

	"github.com/satbirdd/cony"
	"github.com/streadway/amqp"
)

type Producer struct {
	config    Config
	exchange  string
	kind      string
	key       string
	publisher *cony.Publisher
	confirmCh chan amqp.Confirmation
	returnCh  chan amqp.Return
}

func (p Producer) Publish(msg Publishing) error {
	if p.publisher == nil {
		return errors.New("amqp客户端不存在")
	}

	return p.publisher.Publish(msg.AmqpPublishing())
}

func NewProducer(cfg Config, exchange, kind, key string, confirmCh chan amqp.Confirmation, returnCh chan amqp.Return) (*Producer, error) {
	pdcer := Producer{
		config:    cfg,
		exchange:  exchange,
		kind:      kind,
		key:       key,
		confirmCh: confirmCh,
		returnCh:  returnCh,
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
		Durable:    true,
	}
	cli.Declare([]cony.Declaration{
		cony.DeclareExchange(exc),
	})

	declares := []cony.Declaration{}
	opts := []cony.PublisherOpt{}

	if confirmCh != nil {
		declares = append(declares, cony.DeclareConfirm(false))
		declares = append(declares, cony.DeclareNotifyPublish(confirmCh))
	}

	if returnCh != nil {
		opts = append(opts, cony.PublishingMandatory(true))
		declares = append(declares, cony.DeclareNotifyReturn(returnCh))
	}

	cli.PublishDeclare(declares)

	// Declare and register a publisher
	// with the cony client
	pbl := cony.NewPublisher(exc.Name, key, opts...)
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
