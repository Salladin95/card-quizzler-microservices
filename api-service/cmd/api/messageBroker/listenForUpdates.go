package messageBroker

import (
	"fmt"
	"github.com/Salladin95/card-quizzler-microservices/api-service/cmd/api/constants"
	"github.com/Salladin95/rmqtools"
)

func (mb *messageBroker) ListenForUpdates(topics []string, mh func(routingKey string, payload []byte)) {
	consumer, err := rmqtools.NewConsumer(
		mb.rabbitConn,
		constants.AmqpExchange,
		constants.AmqpQueue,
	)
	if err != nil {
		fmt.Println(err)
	}
	err = consumer.Listen(topics, mh)
	if err != nil {
		fmt.Println(err)
	}
}
