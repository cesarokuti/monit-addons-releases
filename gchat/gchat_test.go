package gchat

import "testing"

func TestGchat(t *testing.T) {
	err := SendAlert("gchat", "3.2.1")
	if err != nil {
		t.Errorf("gchat not working %v", err)
	}
}
