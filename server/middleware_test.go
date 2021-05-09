package server

import (
	"testing"
)

func TestTelegram(t *testing.T) {
	err := SendTelegram("go test telegram")

	if err != nil {
		t.Errorf("handler returned unexpected body: got %v want nil",
			err)
	}
}
