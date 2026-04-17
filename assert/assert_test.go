package assert

import (
	"errors"
	"fmt"
	"testing"

	"github.com/atendi9/capivara/langs"
)

type spyTB struct {
	testing.TB
	failed bool
}

func (s *spyTB) Fail() {
	s.failed = true
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

	t.Run("False", func(t *testing.T) {
		False(a, false)
	})

	t.Run("NoError", func(t *testing.T) {
		NoError(a, nil)
	})

	t.Run("Error", func(t *testing.T) {
		Error(a, errors.New("erro esperado"))
	})

	t.Run("ErrorIs", func(t *testing.T) {
		targetErr := errors.New("erro alvo")
		wrappedErr := fmt.Errorf("embrulhando o erro: %w", targetErr)
		ErrorIs(a, wrappedErr, targetErr)
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

	t.Run("LengthSlice", func(t *testing.T) {
		slice := []int{1, 2, 3}
		LengthSlice(a, 3, slice)
	})

	t.Run("LengthMap", func(t *testing.T) {
		m := map[string]int{"a": 1, "b": 2}
		LengthMap(a, 2, m)
	})

	t.Run("LengthString", func(t *testing.T) {
		str := "golang"
		LengthString(a, 6, str)
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
		if !mockT.failed {
			t.Fatal("Expected Equal to call Fail")
		}
	})

	t.Run("True fail", func(t *testing.T) {
		a, mockT := newMockAssert(langs.EN_US)
		True(a, false)
		if !mockT.failed {
			t.Fatal("Expected True to call Fail")
		}
	})

	t.Run("False fail", func(t *testing.T) {
		a, mockT := newMockAssert(langs.EN_US)
		False(a, true)
		if !mockT.failed {
			t.Fatal("Expected False to call Fail")
		}
	})

	t.Run("NoError fail", func(t *testing.T) {
		a, mockT := newMockAssert(langs.EN_US)
		NoError(a, errors.New("erro forçado"))
		if !mockT.failed {
			t.Fatal("Expected NoError to call Fail")
		}
	})

	t.Run("Error fail", func(t *testing.T) {
		a, mockT := newMockAssert(langs.EN_US)
		Error(a, nil)
		if !mockT.failed {
			t.Fatal("Expected Error to call Fail")
		}
	})

	t.Run("ErrorIs fail", func(t *testing.T) {
		a, mockT := newMockAssert(langs.EN_US)
		targetErr := errors.New("erro alvo")
		otherErr := errors.New("erro diferente")
		ErrorIs(a, otherErr, targetErr)
		if !mockT.failed {
			t.Fatal("Expected ErrorIs to call Fail")
		}
	})

	t.Run("NotNil fail", func(t *testing.T) {
		a, mockT := newMockAssert(langs.EN_US)
		NotNil(a, nil)
		if !mockT.failed {
			t.Fatal("Expected NotNil to call Fail")
		}
	})

	t.Run("Empty fail", func(t *testing.T) {
		a, mockT := newMockAssert(langs.EN_US)
		Empty(a, 42)
		if !mockT.failed {
			t.Fatal("Expected Empty to call Fail")
		}
	})

	t.Run("NotEmpty fail", func(t *testing.T) {
		a, mockT := newMockAssert(langs.EN_US)
		NotEmpty(a, "")
		if !mockT.failed {
			t.Fatal("Expected NotEmpty to call Fail")
		}
	})

	t.Run("LengthSlice fail", func(t *testing.T) {
		a, mockT := newMockAssert(langs.EN_US)
		slice := []int{1, 2, 3}
		LengthSlice(a, 5, slice)
		if !mockT.failed {
			t.Fatal("Expected LengthSlice to call Fail")
		}
	})

	t.Run("LengthMap fail", func(t *testing.T) {
		a, mockT := newMockAssert(langs.EN_US)
		m := map[string]int{"a": 1, "b": 2}
		LengthMap(a, 5, m)
		if !mockT.failed {
			t.Fatal("Expected LengthMap to call Fail")
		}
	})

	t.Run("LengthString fail", func(t *testing.T) {
		a, mockT := newMockAssert(langs.EN_US)
		str := "golang"
		LengthString(a, 10, str)
		if !mockT.failed {
			t.Fatal("Expected LengthString to call Fail")
		}
	})
}
