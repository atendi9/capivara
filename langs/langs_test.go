package langs

import (
	"testing"
)

func TestLangs(t *testing.T) {
	if PT_BR != "portuguese" {
		t.Fail()
	}
	if EN_US != "english" {
		t.Fail()
	}
}
