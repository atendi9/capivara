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

var translations = map[langs.Lang]map[string]string{
	langs.EN_US: {
		"fail_match":        "FAIL: Values do not match",
		"expected":          "Expected:",
		"got":               "Got:     ",
		"fail_true":         "FAIL: Expected true, got false",
		"succ_true":         "Value is true",
		"fail_false":        "FAIL: Expected false, got true",
		"succ_false":        "Value is false",
		"fail_err":          "FAIL: Unexpected error encountered",
		"err":               "Error:",
		"succ_err":          "No error returned",
		"fail_expected_err": "FAIL: Expected an error, got nil",
		"succ_expected_err": "Error encountered as expected",
		"fail_err_is":       "FAIL: Error does not match the expected target",
		"succ_err_is":       "Error matches target",
		"target":            "Target:  ",
		"fail_notnil":       "FAIL: Value should not be nil",
		"succ_notnil":       "Value is not nil",
		"fail_empty":        "FAIL: Expected value to be empty, but it was:",
		"succ_empty":        "Value is empty",
		"fail_notempty":     "FAIL: Expected value not to be empty",
		"succ_notempty":     "Value is not empty:",
		"fail_length":       "FAIL: Unexpected length",
		"succ_length":       "Length matches expected value",
		"msg":               "Message:",
	},
	langs.PT_BR: {
		"fail_match":        "FALHA: Os valores não correspondem",
		"expected":          "Esperado:",
		"got":               "Obtido:  ",
		"fail_true":         "FALHA: Esperava true, obteve false",
		"succ_true":         "O valor é true",
		"fail_false":        "FALHA: Esperava false, obteve true",
		"succ_false":        "O valor é false",
		"fail_err":          "FALHA: Erro inesperado encontrado",
		"err":               "Erro:",
		"succ_err":          "Nenhum erro retornado",
		"fail_expected_err": "FALHA: Esperava um erro, obteve nil",
		"succ_expected_err": "Erro encontrado conforme esperado",
		"fail_err_is":       "FALHA: O erro não corresponde ao alvo esperado",
		"succ_err_is":       "O erro corresponde ao alvo",
		"target":            "Alvo:    ",
		"fail_notnil":       "FALHA: O valor não deve ser nil",
		"succ_notnil":       "O valor não é nil",
		"fail_empty":        "FALHA: Esperava que o valor fosse vazio, mas foi:",
		"succ_empty":        "O valor está vazio",
		"fail_notempty":     "FALHA: Esperava que o valor não fosse vazio",
		"succ_notempty":     "O valor não está vazio:",
		"fail_length":       "FALHA: Tamanho inesperado",
		"succ_length":       "Tamanho corresponde ao esperado",
		"msg":               "Mensagem:",
	},
	langs.RU: {
		"fail_match":        "ОШИБКА: Значения не совпадают",
		"expected":          "Ожидалось:",
		"got":               "Получено: ",
		"fail_true":         "ОШИБКА: Ожидалось true, получено false",
		"succ_true":         "Значение равно true",
		"fail_false":        "ОШИБКА: Ожидалось false, получено true",
		"succ_false":        "Значение равно false",
		"fail_err":          "ОШИБКА: Обнаружена непредвиденная ошибка",
		"err":               "Ошибка:",
		"succ_err":          "Ошибок не возвращено",
		"fail_expected_err": "ОШИБКА: Ожидалась ошибка, получено nil",
		"succ_expected_err": "Ошибка возникла как и ожидалось",
		"fail_err_is":       "ОШИБКА: Ошибка не соответствует ожидаемой",
		"succ_err_is":       "Ошибка соответствует ожидаемой",
		"target":            "Цель:     ",
		"fail_notnil":       "ОШИБКА: Значение не должно быть nil",
		"succ_notnil":       "Значение не nil",
		"fail_empty":        "ОШИБКА: Ожидалось пустое значение, но получено:",
		"succ_empty":        "Значение пустое",
		"fail_notempty":     "ОШИБКА: Ожидалось непустое значение",
		"succ_notempty":     "Значение не пустое:",
		"fail_length":       "ОШИБКА: Неожиданная длина",
		"succ_length":       "Длина соответствует ожидаемой",
		"msg":               "Сообщение:",
	},
	langs.JAP: {
		"fail_match":        "失敗: 値が一致しません",
		"expected":          "期待値:",
		"got":               "実際:  ",
		"fail_true":         "失敗: true を期待しましたが false でした",
		"succ_true":         "値は true です",
		"fail_false":        "失敗: false を期待しましたが true でした",
		"succ_false":        "値は false です",
		"fail_err":          "失敗: 予期しないエラーが発生しました",
		"err":               "エラー:",
		"succ_err":          "エラーは返されませんでした",
		"fail_expected_err": "失敗: エラーを期待しましたが nil でした",
		"succ_expected_err": "期待通りエラーが発生しました",
		"fail_err_is":       "失敗: エラーが期待されたターゲットと一致しません",
		"succ_err_is":       "エラーがターゲットと一致しました",
		"target":            "ターゲット:",
		"fail_notnil":       "失敗: 値は nil ではないはずです",
		"succ_notnil":       "値は nil ではありません",
		"fail_empty":        "失敗: 空の値を期待しましたが、以下の通りでした:",
		"succ_empty":        "値は空です",
		"fail_notempty":     "失敗: 値が空ではないことを期待しました",
		"succ_notempty":     "値は空ではありません:",
		"fail_length":       "失敗: 予期しない長さです",
		"succ_length":       "長さが期待値と一致します",
		"msg":               "メッセージ:",
	},
	langs.CH: {
		"fail_match":        "失败：数值不匹配",
		"expected":          "预期值：",
		"got":               "实际值：",
		"fail_true":         "失败：预期为 true，实际为 false",
		"succ_true":         "数值为 true",
		"fail_false":        "失败：预期为 false，实际为 true",
		"succ_false":        "数值为 false",
		"fail_err":          "失败：遇到了意外错误",
		"err":               "错误：",
		"succ_err":          "未返回错误",
		"fail_expected_err": "失败：预期有错误，实际为 nil",
		"succ_expected_err": "如预期般遇到错误",
		"fail_err_is":       "失败：错误与预期目标不符",
		"succ_err_is":       "错误与目标相符",
		"target":            "目标：  ",
		"fail_notnil":       "失败：数值不应为 nil",
		"succ_notnil":       "数值不是 nil",
		"fail_empty":        "失败：预期值应为空，但实际为：",
		"succ_empty":        "数值为空",
		"fail_notempty":     "失败：预期值不应为空",
		"succ_notempty":     "数值不为空：",
		"fail_length":       "失败：长度不符合预期",
		"succ_length":       "长度与预期相符",
		"msg":               "信息：",
	},
}

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
	if langMap, ok := translations[a.Lang]; ok {
		if msg, exists := langMap[key]; exists {
			return msg
		}
	}
	return translations[langs.EN_US][key]
}
