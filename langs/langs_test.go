package langs

import (
	"testing"
)

func TestLangs(t *testing.T) {
	if string(PT_BR) != "portuguese" {
		t.Fail()
	}

	if string(EN_US) != "english" {
		t.Fail()
	}

	if string(RU) != "russian" {
		t.Fail()
	}

	if string(JAP) != "japanese" {
		t.Fail()
	}

	if string(CH) != "chinese" {
		t.Fail()
	}
}
