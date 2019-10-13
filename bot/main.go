package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Response events.APIGatewayProxyResponse
type Request events.APIGatewayProxyRequest

type RequestBody struct {
	Message Message `json:"message"`
}
type From struct {
	ID           int    `json:"id"`
	IsBot        bool   `json:"is_bot"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Username     string `json:"username"`
	LanguageCode string `json:"language_code"`
}
type Chat struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
	Type      string `json:"type"`
}
type Entities struct {
	Offset int    `json:"offset"`
	Length int    `json:"length"`
	Type   string `json:"type"`
}
type Message struct {
	MessageID int        `json:"message_id"`
	From      From       `json:"from"`
	Chat      Chat       `json:"chat"`
	Date      int        `json:"date"`
	Text      string     `json:"text"`
	Entities  []Entities `json:"entities"`
}

func Handler(request Request) {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_API_TOKEN"))
	if err != nil {
		log.Fatal(err)
	}

	var body RequestBody
	if err := json.Unmarshal([]byte(request.Body), &body); err != nil {
		log.Fatal(err)
	}

	var msg tgbotapi.MessageConfig
	if body.Message.Text == "/yesterday" {
		msg = tgbotapi.NewMessage(int64(body.Message.Chat.ID), GetYesterdayScores())
	} else if strings.Contains(body.Message.Text, "Hello, bot!") {
		msg = tgbotapi.NewMessage(
			int64(body.Message.Chat.ID),
			fmt.Sprintf("Hi there, %s!", body.Message.From.FirstName),
		)
		msg.ReplyToMessageID = body.Message.MessageID
	}

	if _, err := bot.Send(msg); err != nil {
		log.Fatal(err)
	}
}

func main() {
	lambda.Start(Handler)
}
