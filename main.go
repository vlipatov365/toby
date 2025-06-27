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

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		userInput := update.Message.Text
		reply, err := queryOllama(userInput)
		if err != nil {
			reply = "Ошибка при обращении к Ollama: " + err.Error()
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
		bot.Send(msg)
	}
}
