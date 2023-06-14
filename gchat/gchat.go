package gchat

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type Message struct {
	Text string `json:"text"`
}

func SendAlert(n string, l string) error {
	webhookURL := os.Getenv("GCHAT_WEBHOOK_URL")

	message := Message{
		Text: "The application " + n + " have new release " + l,
	}

	payload, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("error to marsh JSON: %v", err)
	}

	resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("error to send a message: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("gchat notifications was failed: %v", resp.Status)
	}

	return nil
}
