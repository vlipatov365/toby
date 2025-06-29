package main

import (
	"bytes"
	"encoding/json"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"net/http"
	"os"
)

type OllamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

type OllamaResponse struct {
	Response string `json:"response"`
}

func queryOllama(prompt string) (string, error) {
	reqBody := OllamaRequest{
		Model:  os.Getenv("OLLAMA_MODEL"),
		Prompt: prompt,
		Stream: false,
	}

	jsonData, _ := json.Marshal(reqBody)

	resp, err := http.Post("http://host.docker.internal:11434/api/generate", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	var res OllamaResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return "", err
	}
	return res.Response, nil
}

func main() {
	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	commands := []tgbotapi.BotCommand{
		{
			Command:     "model_list",
			Description: "Выбери модель",
		},
	}

	cfg := tgbotapi.NewSetMyCommands(commands...)
	scope := tgbotapi.NewBotCommandScopeDefault()
	cfg.Scope = &scope

	if _, err := bot.Request(cfg); err != nil {
		log.Fatalf("Request error: %v", err)
	}

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		replyKeyboard := tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton("Выбери модель"),
				tgbotapi.NewKeyboardButton("Выбери персону"),
			),
		)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Выберите опцию ниже:")
		msg.ReplyMarkup = replyKeyboard
		bot.Send(msg)
	}
}
