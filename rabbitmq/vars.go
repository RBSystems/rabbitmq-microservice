package rabbitmq

type Message struct {
	Channel string `json:"channel"`
	Data    string `json:"data,omitempty"`
}

const LocalAPI = "local-api"
const TransmitAPI = "transmit-api"
const External = "external"
const NotifyStart = "notify-start"
const AckStart = "ack-start"
