package main

import (
	"os"

	"github.com/streadway/amqp"
)

func main() {
	// Get the connection string from the environment variable
	url := os.Getenv("AMQP_URL")

	//If it doesn't exist, use the default connection string.

	if url == "" {
		//Don't do this in production, this is for testing purposes only.
		url = "amqp://guest:guest@localhost:5672"
	}

	// Connect to the rabbitMQ instance
	connection, err := amqp.Dial(url)

	if err != nil {
		panic("could not establish connection with RabbitMQ:" + err.Error())
	}

	channel, err := connection.Channel()

	if err != nil {
		panic("could not open RabbitMQ channel:" + err.Error())
	}

	err = channel.ExchangeDeclare("events-di3", "direct", true, false, false, false, nil)

	if err != nil {
		panic(err)
	}

	// We create a queue named Test
	_, err = channel.QueueDeclare("test-di3", true, false, false, false, nil)

	if err != nil {
		panic("error declaring the queue: " + err.Error())
	}

	// We bind the queue to the exchange to send and receive data from the queue
	err = channel.QueueBind("test-di3", "test", "events-di3", false, nil)

	if err != nil {
		panic("error binding to the queue: " + err.Error())
	}

	// We create a message to be sent to the queue.
	// It has to be an instance of the aqmp publishing struct
	message := amqp.Publishing{
		Body: []byte("Hello World"),
	}

	// We publish the message to the exahange we created earlier
	err = channel.Publish("events-di3", "test", false, false, message)

	if err != nil {
		panic("error publishing a message to the queue:" + err.Error())
	}
}
