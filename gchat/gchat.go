package gchat

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

type Message struct {
	Text string `json:"text"`
}

func SendAlert(n string, l string) string {
	webhookURL := os.Getenv("GCHAT_WEBHOOK_URL")

	message := Message{
		Text: "The application " + n + " have new release " + l,
	}

	payload, err := json.Marshal(message)
	if err != nil {
		log.Fatalf("Error to marsh JSON: %v", err)
	}

	resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		fmt.Printf("Error to send a message: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Gchat notifications was failed: %v", resp.Status)
	}

	return "Notification was sended to gchat"
}
