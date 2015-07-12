package main

import (
	"testing"
)

func TestGetNextQuestion(t *testing.T) {
	prompt, complete := GetNextQuestion(map[string]interface{}{
		"game":  "hanabi",
		"score": "25",
	})
	if prompt != "players" {
		t.Errorf("invalid prompt: \"%v\"", prompt)
	}
	if complete {
		t.Error("game should not be complete")
	}

	prompt, complete = GetNextQuestion(map[string]interface{}{
		"game":    "hanabi",
		"score":   "25",
		"players": "tom ryan luke",
	})
	if !complete {
		t.Error("game should be complete")
	}
}
