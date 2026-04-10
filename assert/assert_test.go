package assert

import (
	"errors"
	"fmt"
	"testing"

	"github.com/atendi9/capivara/langs"
)

type spyTB struct {
	testing.TB
	errorfCalled bool
	errorfMsg    string
	logfCalled   bool
}

func (s *spyTB) Errorf(format string, args ...any) {
	s.errorfCalled = true
	s.errorfMsg = fmt.Sprintf(format, args...)
}

func (s *spyTB) Logf(format string, args ...any) {
	s.logfCalled = true
}

func (s *spyTB) Helper() {}

func TestAssertions_Success(t *testing.T) {
	a := New(langs.EN_US, t)

	t.Run("Equal", func(t *testing.T) {
		Equal(a, 10, 10, "Inteiros devem ser iguais")
		Equal(a, "golang", "golang")
	})

	t.Run("True", func(t *testing.T) {
		True(a, true)
	})

	t.Run("NoError", func(t *testing.T) {
		NoError(a, nil)
	})

	t.Run("NotNil", func(t *testing.T) {
		NotNil(a, "não sou nulo")
	})

	t.Run("Empty", func(t *testing.T) {
		Empty(a, 0)
		Empty(a, "")
		Empty(a, false)
	})

	t.Run("NotEmpty", func(t *testing.T) {
		NotEmpty(a, 42)
		NotEmpty(a, "texto")
	})
}

func TestAssertions_Failures(t *testing.T) {
	newMockAssert := func(lang langs.Lang) (*Assert, *spyTB) {
		mockT := &spyTB{TB: t}
		return New(lang, mockT), mockT
	}

	t.Run("Equal fail", func(t *testing.T) {
		a, mockT := newMockAssert(langs.EN_US)
		Equal(a, 10, 20)
		if !mockT.errorfCalled {
			t.Fatal("Esperava que Equal chamasse Errorf")
		}
	})

	t.Run("True fail", func(t *testing.T) {
		a, mockT := newMockAssert(langs.EN_US)
		True(a, false)
		if !mockT.errorfCalled {
			t.Fatal("Esperava que True chamasse Errorf")
		}
	})

	t.Run("NoError fail", func(t *testing.T) {
		a, mockT := newMockAssert(langs.EN_US)
		NoError(a, errors.New("erro forçado"))
		if !mockT.errorfCalled {
			t.Fatal("Esperava que NoError chamasse Errorf")
		}
	})

	t.Run("NotNil fail", func(t *testing.T) {
		a, mockT := newMockAssert(langs.EN_US)
		NotNil(a, nil)
		if !mockT.errorfCalled {
			t.Fatal("Esperava que NotNil chamasse Errorf")
		}
	})

	t.Run("Empty fail", func(t *testing.T) {
		a, mockT := newMockAssert(langs.EN_US)
		Empty(a, 42)
		if !mockT.errorfCalled {
			t.Fatal("Esperava que Empty chamasse Errorf")
		}
	})

	t.Run("NotEmpty fail", func(t *testing.T) {
		a, mockT := newMockAssert(langs.EN_US)
		NotEmpty(a, "")
		if !mockT.errorfCalled {
			t.Fatal("Esperava que NotEmpty chamasse Errorf")
		}
	})
}
