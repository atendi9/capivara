// Package assert provides high-performance assertion functions for Go tests.
//
// Instead of using reflection (reflect), this package takes advantage of Generics
// ([T comparable]) to perform compile-time comparisons, ensuring zero heap allocations
// (zero-allocation) and maximum performance.
package assert

import (
	"errors"
	"fmt"
	"testing"

	"github.com/atendi9/capivara/langs"
)

const (
	iconPass  = "✅"
	iconFail  = "❌"
	iconGot   = "🔍"
	iconWant  = "🎯"
	iconError = "🔥"
	iconMsg   = "💬"
)

// translations is the internal dictionary that maps error messages
// according to the configured language.
var translations = map[langs.Lang]map[string]string{
    langs.EN_US: {
        "fail_match":        "FAIL: Values do not match",
        "expected":          "Expected:",
        "got":               "Got:     ",
        "succ_match":        "SUCCESS",
        "fail_true":         "FAIL: Expected true, got false",
        "succ_true":         "SUCCESS: Value is true",
        "fail_false":        "FAIL: Expected false, got true",
        "succ_false":        "SUCCESS: Value is false",
        "fail_err":          "FAIL: Unexpected error encountered",
        "err":               "Error:",
        "succ_err":          "SUCCESS: No error returned",
        "fail_expected_err": "FAIL: Expected an error, got nil",
        "succ_expected_err": "SUCCESS: Error encountered as expected",
        "fail_err_is":       "FAIL: Error does not match the expected target",
        "succ_err_is":       "SUCCESS: Error matches target",
        "target":            "Target:  ",
        "fail_notnil":       "FAIL: Value should not be nil",
        "succ_notnil":       "SUCCESS: Value is not nil",
        "fail_empty":        "FAIL: Expected value to be empty, but it was:",
        "succ_empty":        "SUCCESS: Value is empty",
        "fail_notempty":     "FAIL: Expected value not to be empty",
        "succ_notempty":     "SUCCESS: Value is not empty:",
        "fail_length":       "FAIL: Unexpected length",
        "succ_length":       "SUCCESS: Length matches expected value",
        "msg":               "Message:",
    },
    langs.PT_BR: {
        "fail_match":        "FALHA: Os valores não correspondem",
        "expected":          "Esperado:",
        "got":               "Obtido:  ",
        "succ_match":        "SUCESSO",
        "fail_true":         "FALHA: Esperava true, obteve false",
        "succ_true":         "SUCESSO: O valor é true",
        "fail_false":        "FALHA: Esperava false, obteve true",
        "succ_false":        "SUCESSO: O valor é false",
        "fail_err":          "FALHA: Erro inesperado encontrado",
        "err":               "Erro:",
        "succ_err":          "SUCESSO: Nenhum erro retornado",
        "fail_expected_err": "FALHA: Esperava um erro, obteve nil",
        "succ_expected_err": "SUCESSO: Erro encontrado conforme esperado",
        "fail_err_is":       "FALHA: O erro não corresponde ao alvo esperado",
        "succ_err_is":       "SUCESSO: O erro corresponde ao alvo",
        "target":            "Alvo:    ",
        "fail_notnil":       "FALHA: O valor não deve ser nil",
        "succ_notnil":       "SUCESSO: O valor não é nil",
        "fail_empty":        "FALHA: Esperava que o valor fosse vazio, mas foi:",
        "succ_empty":        "SUCESSO: O valor está vazio",
        "fail_notempty":     "FALHA: Esperava que o valor não fosse vazio",
        "succ_notempty":     "SUCESSO: O valor não está vazio:",
        "fail_length":       "FALHA: Tamanho inesperado",
        "succ_length":       "SUCESSO: Tamanho corresponde ao esperado",
        "msg":               "Mensagem:",
    },
    langs.RU: {
        "fail_match":        "ОШИБКА: Значения не совпадают",
        "expected":          "Ожидалось:",
        "got":               "Получено: ",
        "succ_match":        "УСПЕХ",
        "fail_true":         "ОШИБКА: Ожидалось true, получено false",
        "succ_true":         "УСПЕХ: Значение равно true",
        "fail_false":        "ОШИБКА: Ожидалось false, получено true",
        "succ_false":        "УСПЕХ: Значение равно false",
        "fail_err":          "ОШИБКА: Обнаружена непредвиденная ошибка",
        "err":               "Ошибка:",
        "succ_err":          "УСПЕХ: Ошибок не возвращено",
        "fail_expected_err": "ОШИБКА: Ожидалась ошибка, получено nil",
        "succ_expected_err": "УСПЕХ: Ошибка возникла как и ожидалось",
        "fail_err_is":       "ОШИБКА: Ошибка не соответствует ожидаемой",
        "succ_err_is":       "УСПЕХ: Ошибка соответствует ожидаемой",
        "target":            "Цель:     ",
        "fail_notnil":       "ОШИБКА: Значение не должно быть nil",
        "succ_notnil":       "УСПЕХ: Значение не nil",
        "fail_empty":        "ОШИБКА: Ожидалось пустое значение, но получено:",
        "succ_empty":        "УСПЕХ: Значение пустое",
        "fail_notempty":     "ОШИБКА: Ожидалось непустое значение",
        "succ_notempty":     "УСПЕХ: Значение не пустое:",
        "fail_length":       "ОШИБКА: Неожиданная длина",
        "succ_length":       "УСПЕХ: Длина соответствует ожидаемой",
        "msg":               "Сообщение:",
    },
    langs.JAP: {
        "fail_match":        "失敗: 値が一致しません",
        "expected":          "期待値:",
        "got":               "実際:  ",
        "succ_match":        "成功",
        "fail_true":         "失敗: true を期待しましたが false でした",
        "succ_true":         "成功: 値は true です",
        "fail_false":        "失敗: false を期待しましたが true でした",
        "succ_false":        "成功: 値は false です",
        "fail_err":          "失敗: 予期しないエラーが発生しました",
        "err":               "エラー:",
        "succ_err":          "成功: エラーは返されませんでした",
        "fail_expected_err": "失敗: エラーを期待しましたが nil でした",
        "succ_expected_err": "成功: 期待通りエラーが発生しました",
        "fail_err_is":       "失敗: エラーが期待されたターゲットと一致しません",
        "succ_err_is":       "成功: エラーがターゲットと一致しました",
        "target":            "ターゲット:",
        "fail_notnil":       "失敗: 値は nil ではないはずです",
        "succ_notnil":       "成功: 値は nil ではありません",
        "fail_empty":        "失敗: 空の値を期待しましたが、以下の通りでした:",
        "succ_empty":        "成功: 値は空です",
        "fail_notempty":     "失敗: 値が空ではないことを期待しました",
        "succ_notempty":     "成功: 値は空ではありません:",
        "fail_length":       "失敗: 予期しない長さです",
        "succ_length":       "成功: 長さが期待値と一致します",
        "msg":               "メッセージ:",
    },
    langs.CH: {
        "fail_match":        "失败：数值不匹配",
        "expected":          "预期值：",
        "got":               "实际值：",
        "succ_match":        "成功",
        "fail_true":         "失败：预期为 true，实际为 false",
        "succ_true":         "成功：数值为 true",
        "fail_false":        "失败：预期为 false，实际为 true",
        "succ_false":        "成功：数值为 false",
        "fail_err":          "失败：遇到了意外错误",
        "err":               "错误：",
        "succ_err":          "成功：未返回错误",
        "fail_expected_err": "失败：预期有错误，实际为 nil",
        "succ_expected_err": "成功：如预期般遇到错误",
        "fail_err_is":       "失败：错误与预期目标不符",
        "succ_err_is":       "成功：错误与目标相符",
        "target":            "目标：  ",
        "fail_notnil":       "失败：数值不应为 nil",
        "succ_notnil":       "成功：数值不是 nil",
        "fail_empty":        "失败：预期值应为空，但实际为：",
        "succ_empty":        "成功：数值为空",
        "fail_notempty":     "失败：预期值不应为空",
        "succ_notempty":     "成功：数值不为空：",
        "fail_length":       "失败：长度不符合预期",
        "succ_length":       "成功：长度与预期相符",
        "msg":               "信息：",
    },
}

// Assert holds the state of the current test context, including
// the preferred language for messages and the underlying test interface.
type Assert struct {
	Lang langs.Lang
	t    testing.TB // Kept private to prevent unwanted mutations
}

// New initializes a new assertion context instance.
// The testing.TB interface allows it to be used with both *testing.T, *testing.B and *testing.F.
//
// Usage example:
//
//	func TestSomething(t *testing.T) {
//	    a := assert.New(langs.EN_US, t)
//	    assert.Equal(a, "golang", "golang", "strings should be equal")
//	}
func New(lang langs.Lang, t testing.TB) *Assert {
	return &Assert{
		Lang: lang,
		t:    t,
	}
}

// Equal verifies if two values are deeply equal.
// It uses Generics [T comparable] to ensure fast compile-time checking,
// rejecting invalid types (like slices or maps) without causing a panic.
func Equal[T comparable](a *Assert, expected, actual T, msg ...string) {
	a.t.Helper()

	if expected != actual {
		m := formatMessage(a, msg)
		a.t.Errorf("\n%s %s%s\n\t%s %s %#v\n\t%s %s %#v\n",
			iconFail, translate(a, "fail_match"), m, iconWant, translate(a, "expected"), expected, iconGot, translate(a, "got"), actual)
	} else {
		a.t.Logf("%s %s: %#v == %#v", iconPass, translate(a, "succ_match"), expected, actual)
	}
}

// True verifies if the provided boolean value is true.
func True(a *Assert, actual bool, msg ...string) {
	a.t.Helper()

	if !actual {
		m := formatMessage(a, msg)
		a.t.Errorf("\n%s %s%s\n", iconFail, translate(a, "fail_true"), m)
	} else {
		a.t.Logf("%s %s", iconPass, translate(a, "succ_true"))
	}
}

// False verifies if the provided boolean value is false.
func False(a *Assert, actual bool, msg ...string) {
	a.t.Helper()

	if actual {
		m := formatMessage(a, msg)
		a.t.Errorf("\n%s %s%s\n", iconFail, translate(a, "fail_false"), m)
	} else {
		a.t.Logf("%s %s", iconPass, translate(a, "succ_false"))
	}
}

// NoError fails the test if the error interface is not nil.
func NoError(a *Assert, err error, msg ...string) {
	a.t.Helper()

	if err != nil {
		m := formatMessage(a, msg)
		a.t.Errorf("\n%s %s%s\n\t%s %s %v\n", iconError, translate(a, "fail_err"), m, iconGot, translate(a, "err"), err)
	} else {
		a.t.Logf("%s %s", iconPass, translate(a, "succ_err"))
	}
}

// Error fails the test if the error interface is nil.
func Error(a *Assert, err error, msg ...string) {
	a.t.Helper()

	if err == nil {
		m := formatMessage(a, msg)
		a.t.Errorf("\n%s %s%s\n", iconFail, translate(a, "fail_expected_err"), m)
	} else {
		a.t.Logf("%s %s: %v", iconPass, translate(a, "succ_expected_err"), err)
	}
}

// ErrorIs verifies if the provided error matches the target error using [errors.Is].
func ErrorIs(a *Assert, err, target error, msg ...string) {
	a.t.Helper()

	if !errors.Is(err, target) {
		m := formatMessage(a, msg)
		a.t.Errorf("\n%s %s%s\n\t%s %s %v\n\t%s %s %v\n",
			iconFail, translate(a, "fail_err_is"), m, iconWant, translate(a, "target"), target, iconGot, translate(a, "err"), err)
	} else {
		a.t.Logf("%s %s", iconPass, translate(a, "succ_err_is"))
	}
}

// NotNil strictly verifies if an interface is not nil.
// Note: Due to the absence of reflect, it does not inspect "typed nils", focusing on absolute performance.
func NotNil(a *Assert, actual any, msg ...string) {
	a.t.Helper()

	if actual == nil {
		m := formatMessage(a, msg)
		a.t.Errorf("\n%s %s%s\n", iconFail, translate(a, "fail_notnil"), m)
	} else {
		a.t.Logf("%s %s", iconPass, translate(a, "succ_notnil"))
	}
}

// Empty verifies if the provided value matches the "zero-value" of its type
// (e.g., 0 for integers, "" for strings, false for booleans).
func Empty[T comparable](a *Assert, actual T, msg ...string) {
	a.t.Helper()
	var zero T // Implicitly instantiates the native zero-value for type T

	if actual != zero {
		m := formatMessage(a, msg)
		a.t.Errorf("\n%s %s %#v%s\n", iconFail, translate(a, "fail_empty"), actual, m)
	} else {
		a.t.Logf("%s %s", iconPass, translate(a, "succ_empty"))
	}
}

// NotEmpty verifies if the provided value is different from the "zero-value" of its respective type.
func NotEmpty[T comparable](a *Assert, actual T, msg ...string) {
	a.t.Helper()
	var zero T

	if actual == zero {
		m := formatMessage(a, msg)
		a.t.Errorf("\n%s %s%s\n", iconFail, translate(a, "fail_notempty"), m)
	} else {
		a.t.Logf("%s %s %#v", iconPass, translate(a, "succ_notempty"), actual)
	}
}

// LengthSlice verifies if the length of the provided slice matches the expected length.
// By using type constraints (~[]E), it enforces compile-time checks and avoids reflection.
func LengthSlice[S ~[]E, E any](a *Assert, expected int, actual S, msg ...string) {
	a.t.Helper()

	actualLen := len(actual)
	if actualLen != expected {
		m := formatMessage(a, msg)
		a.t.Errorf("\n%s %s%s\n\t%s %s %d\n\t%s %s %d\n",
			iconFail, translate(a, "fail_length"), m, iconWant, translate(a, "expected"), expected, iconGot, translate(a, "got"), actualLen)
	} else {
		a.t.Logf("%s %s: len() == %d", iconPass, translate(a, "succ_length"), expected)
	}
}

// LengthMap verifies if the length of the provided map matches the expected length.
// By using type constraints (~map[K]V), it enforces compile-time checks and avoids reflection.
func LengthMap[M ~map[K]V, K comparable, V any](a *Assert, expected int, actual M, msg ...string) {
	a.t.Helper()

	actualLen := len(actual)
	if actualLen != expected {
		m := formatMessage(a, msg)
		a.t.Errorf("\n%s %s%s\n\t%s %s %d\n\t%s %s %d\n",
			iconFail, translate(a, "fail_length"), m, iconWant, translate(a, "expected"), expected, iconGot, translate(a, "got"), actualLen)
	} else {
		a.t.Logf("%s %s: len() == %d", iconPass, translate(a, "succ_length"), expected)
	}
}

// LengthString verifies if the length of the provided string matches the expected length.
func LengthString[S ~string](a *Assert, expected int, actual S, msg ...string) {
	a.t.Helper()

	actualLen := len(actual)
	if actualLen != expected {
		m := formatMessage(a, msg)
		a.t.Errorf("\n%s %s%s\n\t%s %s %d\n\t%s %s %d\n",
			iconFail, translate(a, "fail_length"), m, iconWant, translate(a, "expected"), expected, iconGot, translate(a, "got"), actualLen)
	} else {
		a.t.Logf("%s %s: len() == %d", iconPass, translate(a, "succ_length"), expected)
	}
}

// translate is an internal helper function that fetches the localized text.
func translate(a *Assert, key string) string {
	if langMap, ok := translations[a.Lang]; ok {
		if msg, exists := langMap[key]; exists {
			return msg
		}
	}
	// Safety fallback to English
	return translations[langs.EN_US][key]
}

// formatMessage is an internal helper function to format the custom additional message.
func formatMessage(a *Assert, msg []string) string {
	if len(msg) > 0 {
		return fmt.Sprintf("\n\t%s %s %s", iconMsg, translate(a, "msg"), msg[0])
	}
	return ""
}
