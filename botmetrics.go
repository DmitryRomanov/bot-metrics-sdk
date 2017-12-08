package botmetrics

import (
	"encoding/json"
	"time"

	"github.com/rockneurotiko/gorequest"
)

const (
	URL                 = "https://api.bot-metrics.com/v1/messages"
	MessageTypeIncoming = "incoming"
	MessageTypeOutgoing = "outgoing"
	PlatformTelegram    = "telegram"
)

type Answer struct {
	Status string `json:"status"`
	Info   string `json:"info,omitempty"`
}

type sender struct {
	Token string `json:"token"`
}

type Envelope struct {
	Message Message `json:"message"`
}

type Message struct {
	Text        string      `json:"text"`
	MessageType string      `json:"message_type"`
	UserID      string      `json:"user_id"`
	Platform    string      `json:"platform"`
	CreatedAt   time.Time   `json:"created_at"`
	Metadata    interface{} `json:"metadata,omitempty"`
}

func createRequest(getp sender, payload Envelope) *gorequest.SuperAgent {
	return gorequest.New().
		Post(URL).
		DisableKeepAlives(true).
		CloseRequest(true).
		Query(getp).  // Get parameters
		Send(payload) // Post payload
}

type Botmetrics struct {
	Token string
}

func New(token string) Botmetrics {
	return Botmetrics{Token: token}
}

func (self Botmetrics) Track(message Message) (result Answer, err []error) {
	result = Answer{"failed", ""}
	request := createRequest(sender{self.Token}, Envelope{message})

	_, body, err := request.End()

	if err != nil {
		return
	}

	jerr := json.Unmarshal([]byte(body), &result)
	err = []error{jerr}

	return
}

func (self Botmetrics) TrackAsync(message Message, f func(Answer, []error)) {
	go func() {
		ans, err := self.Track(message)
		f(ans, err)
	}()
}
