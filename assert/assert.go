// Package assert provides high-performance, generic-based assertion functions for testing.
// It bypasses default testing prefixes and provides localized output with emojis.
package assert

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"testing"

	"github.com/atendi9/capivara/langs"
)

const (
	// IconPass represents the visual indicator for a successfully passed assertion.
	IconPass = "✅"
	// IconFail represents the visual indicator for a failed assertion.
	IconFail = "❌"
	// IconGot represents the visual indicator for the received actual value.
	IconGot = "🔍"
	// IconWant represents the visual indicator for the expected value.
	IconWant = "🎯"
	// IconError represents the visual indicator for an unexpected error.
	IconError = "🔥"
	// IconMsg represents the visual indicator for custom assertion messages.
	IconMsg = "💬"
	// BaseIndent defines the primary indentation level for the main assertion text.
	BaseIndent = "    "
	// ArrowIndent defines the secondary indentation level for metadata like file paths.
	ArrowIndent = "      ↳ "
	// SubIndent defines the tertiary indentation level for failure details (expected vs got).
	SubIndent = "        "
)

var (
	projectRoot string
	rootOnce    sync.Once
)

// getProjectRoot traverses upwards from the current working directory to find the directory containing go.mod.
// It caches the result to maintain high performance across multiple assertions.
func getProjectRoot() string {
	rootOnce.Do(func() {
		dir, err := os.Getwd()
		if err != nil {
			return
		}
		for {
			if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
				projectRoot = dir
				return
			}
			parent := filepath.Dir(dir)
			if parent == dir {
				break
			}
			dir = parent
		}
		// Fallback to the current working directory if go.mod isn't found
		projectRoot, _ = os.Getwd()
	})
	return projectRoot
}

// Assert holds the active state and localization preferences for the current test context.
// By embedding testing.TB, *Assert inherently implements the testing.TB interface.
type Assert struct {
	Lang langs.Lang
	testing.TB
}

// New initializes and returns a new Assert instance mapped to the provided testing context and language.
// If the provided context is already an *Assert, it re-wraps the underlying testing.TB to update the language.
func New(lang langs.Lang, t testing.TB) *Assert {
	if a, ok := t.(*Assert); ok {
		return &Assert{
			Lang: lang,
			TB:   a.TB,
		}
	}
	return &Assert{
		Lang: lang,
		TB:   t,
	}
}

// resolveAssert retrieves the *Assert instance from the testing.TB interface or creates a default one.
func resolveAssert(t testing.TB) *Assert {
	if a, ok := t.(*Assert); ok {
		return a
	}
	// Default to EN_US when standard testing.TB is provided.
	return New(langs.EN_US, t)
}

// printOut formats and outputs the assertion results, marking the test as failed if needed.
func printOut(a *Assert, passed bool, mainText string, msg []string) {
	a.Helper()

	_, file, line, ok := runtime.Caller(2)
	shortFile := "unknown"
	if ok {
		shortFile = file
		rootDir := getProjectRoot()
		if rootDir != "" {
			if rel, err := filepath.Rel(rootDir, file); err == nil {
				shortFile = "./" + filepath.ToSlash(rel)
			}
		}
	}

	m := ""
	if len(msg) > 0 {
		m = fmt.Sprintf("\n%s%s %s %s", ArrowIndent, IconMsg, translate(a, "msg"), msg[0])
	}

	if !passed {
		a.Fail()
	}

	fmt.Printf("\n%s%s\n%s📄 %s:%d%s\n", BaseIndent, mainText, ArrowIndent, shortFile, line, m)
}

// Equal verifies that the expected value is deeply and strictly equal to the actual value.
func Equal[T comparable](t testing.TB, expected, actual T, msg ...string) {
	t.Helper()
	a := resolveAssert(t)

	if expected != actual {
		text := fmt.Sprintf("%s %s\n%s%s %s %#v\n%s%s %s %#v",
			IconFail, translate(a, "fail_match"), SubIndent, IconWant, translate(a, "expected"), expected, SubIndent, IconGot, translate(a, "got"), actual)
		printOut(a, false, text, msg)
	} else {
		text := fmt.Sprintf("%s %#v == %#v", IconPass, expected, actual)
		printOut(a, true, text, msg)
	}
}

// True verifies that the provided boolean expression evaluates to true.
func True(t testing.TB, actual bool, msg ...string) {
	t.Helper()
	a := resolveAssert(t)

	if !actual {
		text := fmt.Sprintf("%s %s", IconFail, translate(a, "fail_true"))
		printOut(a, false, text, msg)
	} else {
		text := fmt.Sprintf("%s %s", IconPass, translate(a, "succ_true"))
		printOut(a, true, text, msg)
	}
}

// False verifies that the provided boolean expression evaluates to false.
func False(t testing.TB, actual bool, msg ...string) {
	t.Helper()
	a := resolveAssert(t)

	if actual {
		text := fmt.Sprintf("%s %s", IconFail, translate(a, "fail_false"))
		printOut(a, false, text, msg)
	} else {
		text := fmt.Sprintf("%s %s", IconPass, translate(a, "succ_false"))
		printOut(a, true, text, msg)
	}
}

// NoError verifies that the provided error interface is nil, indicating no error occurred.
func NoError(t testing.TB, err error, msg ...string) {
	t.Helper()
	a := resolveAssert(t)

	if err != nil {
		text := fmt.Sprintf("%s %s\n%s%s %s %v", IconError, translate(a, "fail_err"), SubIndent, IconGot, translate(a, "err"), err)
		printOut(a, false, text, msg)
	} else {
		text := fmt.Sprintf("%s %s", IconPass, translate(a, "succ_err"))
		printOut(a, true, text, msg)
	}
}

// Error verifies that the provided error interface is not nil, indicating an error occurred.
func Error(t testing.TB, err error, msg ...string) {
	t.Helper()
	a := resolveAssert(t)

	if err == nil {
		text := fmt.Sprintf("%s %s", IconFail, translate(a, "fail_expected_err"))
		printOut(a, false, text, msg)
	} else {
		text := fmt.Sprintf("%s %s: %v", IconPass, translate(a, "succ_expected_err"), err)
		printOut(a, true, text, msg)
	}
}

// ErrorIs verifies that the provided error matches the specified target error within its wrap chain.
func ErrorIs(t testing.TB, err, target error, msg ...string) {
	t.Helper()
	a := resolveAssert(t)

	if !errors.Is(err, target) {
		text := fmt.Sprintf("%s %s\n%s%s %s %v\n%s%s %s %v",
			IconFail, translate(a, "fail_err_is"), SubIndent, IconWant, translate(a, "target"), target, SubIndent, IconGot, translate(a, "err"), err)
		printOut(a, false, text, msg)
	} else {
		text := fmt.Sprintf("%s %s", IconPass, translate(a, "succ_err_is"))
		printOut(a, true, text, msg)
	}
}

// NotNil strictly verifies that the actual interface or value is not nil.
func NotNil(t testing.TB, actual any, msg ...string) {
	t.Helper()
	a := resolveAssert(t)

	if actual == nil {
		text := fmt.Sprintf("%s %s", IconFail, translate(a, "fail_notnil"))
		printOut(a, false, text, msg)
	} else {
		text := fmt.Sprintf("%s %s", IconPass, translate(a, "succ_notnil"))
		printOut(a, true, text, msg)
	}
}

// Empty verifies that the actual value is equivalent to the zero-value of its specific generic type.
func Empty[T comparable](t testing.TB, actual T, msg ...string) {
	t.Helper()
	a := resolveAssert(t)
	var zero T

	if actual != zero {
		text := fmt.Sprintf("%s %s %#v", IconFail, translate(a, "fail_empty"), actual)
		printOut(a, false, text, msg)
	} else {
		text := fmt.Sprintf("%s %s", IconPass, translate(a, "succ_empty"))
		printOut(a, true, text, msg)
	}
}

// NotEmpty verifies that the actual value differs from the zero-value of its specific generic type.
func NotEmpty[T comparable](t testing.TB, actual T, msg ...string) {
	t.Helper()
	a := resolveAssert(t)
	var zero T

	if actual == zero {
		text := fmt.Sprintf("%s %s", IconFail, translate(a, "fail_notempty"))
		printOut(a, false, text, msg)
	} else {
		text := fmt.Sprintf("%s %s %#v", IconPass, translate(a, "succ_notempty"), actual)
		printOut(a, true, text, msg)
	}
}

// LengthSlice verifies that the number of elements in the provided slice matches the expected count.
func LengthSlice[S ~[]E, E any](t testing.TB, expected int, actual S, msg ...string) {
	t.Helper()
	a := resolveAssert(t)
	actualLen := len(actual)

	if actualLen != expected {
		text := fmt.Sprintf("%s %s\n%s%s %s %d\n%s%s %s %d",
			IconFail, translate(a, "fail_length"), SubIndent, IconWant, translate(a, "expected"), expected, SubIndent, IconGot, translate(a, "got"), actualLen)
		printOut(a, false, text, msg)
	} else {
		text := fmt.Sprintf("%s %s: len() == %d", IconPass, translate(a, "succ_length"), expected)
		printOut(a, true, text, msg)
	}
}

// LengthMap verifies that the number of key-value pairs in the provided map matches the expected count.
func LengthMap[M ~map[K]V, K comparable, V any](t testing.TB, expected int, actual M, msg ...string) {
	t.Helper()
	a := resolveAssert(t)
	actualLen := len(actual)

	if actualLen != expected {
		text := fmt.Sprintf("%s %s\n%s%s %s %d\n%s%s %s %d",
			IconFail, translate(a, "fail_length"), SubIndent, IconWant, translate(a, "expected"), expected, SubIndent, IconGot, translate(a, "got"), actualLen)
		printOut(a, false, text, msg)
	} else {
		text := fmt.Sprintf("%s %s: len() == %d", IconPass, translate(a, "succ_length"), expected)
		printOut(a, true, text, msg)
	}
}

// LengthString verifies that the character length of the provided string matches the expected count.
func LengthString[S ~string](t testing.TB, expected int, actual S, msg ...string) {
	t.Helper()
	a := resolveAssert(t)
	actualLen := len(actual)

	if actualLen != expected {
		text := fmt.Sprintf("%s %s\n%s%s %s %d\n%s%s %s %d",
			IconFail, translate(a, "fail_length"), SubIndent, IconWant, translate(a, "expected"), expected, SubIndent, IconGot, translate(a, "got"), actualLen)
		printOut(a, false, text, msg)
	} else {
		text := fmt.Sprintf("%s %s: len() == %d", IconPass, translate(a, "succ_length"), expected)
		printOut(a, true, text, msg)
	}
}

// translate retrieves the localized message for the given key based on the Assert language setting.
func translate(a *Assert, key string) string {
    translations := langs.Translate(langs.PkgAssert)
	if langMap, ok := translations[a.Lang]; ok {
		if msg, exists := langMap[key]; exists {
			return msg
		}
	}
	return translations[langs.EN_US][key]
}
