package syncbreaker

import "testing"

func TestBreak(t *testing.T) {
	Break()
	if Break() {
		t.Error("expected semaphore to be gone after first Break")
	}
}
