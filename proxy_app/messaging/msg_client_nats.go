package messaging

import (
	"github.com/nats-io/nats.go"
	"github.com/unibrightio/proxy-api/logger"
)

type IMessagingClient interface {
	// used to send an OffchainProcessMessage
	// message - message payload
	// recipient - address of the recipient (i.e. NATS server url) taken from workgroup
	// token - token used to authenticate (i.e. NATS server token) taken from the workgroup
	SendMessage(message []byte, recipient string, token string)

	// used to receive messages sent by other participants to our nats server
	// serverUrl - local server url
	// token - local server token
	// topic - listening topic
	// onMessageReceived - callback function
	Subscribe(serverUrl string, token string, topic string, onMessageReceived func(string, *nats.Msg))
}

type NatsMessagingClient struct {
}

func (client *NatsMessagingClient) SendMessage(message []byte, recipient string, token string) {
	// https://docs.nats.io/developing-with-nats/security/token
	nc, err := nats.Connect("nats://" + token + "@" + recipient)

	if err != nil {
		logger.Errorf("Error while trying to connect to Nats: %v, message: %s, recipient: %s, token: %s", err, message, recipient, token)
		return
	}

	defer nc.Close()

	// TODO: https://docs.nats.io/developing-with-nats/sending/replyto
	err = nc.Publish("baseledger", message)

	if err != nil {
		logger.Errorf("Error while trying to send NATS message: %v, message: %s, recipient: %s, token: %s", err, message, recipient, token)
	}
}

func (client *NatsMessagingClient) Subscribe(serverUrl string, token string, topic string, onMessageReceived func(string, *nats.Msg)) {
	// https://docs.nats.io/developing-with-nats/security/token
	nc, err := nats.Connect("nats://" + token + "@" + serverUrl)

	if err != nil {
		logger.Errorf("Error while trying to connect to local Nats: %v", err)
		return
	}

	nc.Subscribe(topic, func(m *nats.Msg) {
		onMessageReceived(string("TODO: m.Sender"), m)
	})
}
