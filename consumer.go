package ptls_amqp

import (
	"fmt"

	"github.com/assembla/cony"
)

type Comsumer struct {
	config   Config
	exchange string
	kind     string
	key      string
	queue    string
	comsumer *cony.Consumer
}

func NewComsumer(cfg Config, exchange, kind, key, queue string, msgCh chan<- Delivery) (*Comsumer, error) {
	comsumer := Comsumer{
		config:   cfg,
		exchange: exchange,
		kind:     kind,
		key:      key,
		queue:    queue,
	}

	url := cfg.String()

	cli := cony.NewClient(
		cony.URL(url),
		cony.Backoff(cony.DefaultBackoff),
	)

	que := &cony.Queue{
		AutoDelete: false,
		Name:       queue,
	}

	// Declare the exchange we'll be using
	exc := cony.Exchange{
		Name:       exchange,
		Kind:       kind,
		AutoDelete: false,
	}

	bnd := cony.Binding{
		Queue:    que,
		Exchange: exc,
		Key:      key,
	}
	cli.Declare([]cony.Declaration{
		cony.DeclareQueue(que),
		cony.DeclareExchange(exc),
		cony.DeclareBinding(bnd),
	})

	// Declare and register a publisher
	// with the cony client
	cns := cony.NewConsumer(
		que,
	)
	comsumer.comsumer = cns
	cli.Consume(cns)

	// Client loop sends out declarations(exchanges, queues, bindings
	// etc) to the AMQP server. It also handles reconnecting.
	go func() {
		for cli.Loop() {
			select {
			case msg := <-cns.Deliveries():
				// log.Printf("Received body: %q\n", msg.Body)
				msgCh <- Delivery(msg)

				// log.Printf("Received body: %q\n", msg.Body)
			case err := <-cns.Errors():
				fmt.Printf("Consumer error: %v\n", err)
			case err := <-cli.Errors():
				fmt.Printf("Client error: %v\n", err)
			}
		}
	}()

	return &comsumer, nil
}
