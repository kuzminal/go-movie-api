package validator

import "regexp"

var (
	EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\. [a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
)

type Validator struct {
	Errors map[string]string
}

// New создет новый Validator.
func New() *Validator {
	return &Validator{Errors: make(map[string]string)}
}

// Valid возвращает true если errors пустая мапа.
func (v *Validator) Valid() bool {
	return len(v.Errors) == 0
}

// AddError добавляет ошибку в мапу если такого ключа там еще нет.
func (v *Validator) AddError(key, message string) {
	if _, exists := v.Errors[key]; !exists {
		v.Errors[key] = message
	}
}

// Check добавляет ошибку в случае если проверка != 'ok'.
func (v *Validator) Check(ok bool, key, message string) {
	if !ok {
		v.AddError(key, message)
	}
}

// In возвращает true если надено указанное значение с писке строк.
func In(value string, list ...string) bool {
	for i := range list {
		if value == list[i] {
			return true
		}
	}
	return false
}

// Matches возвращает true если строка соответствует регулярке.
func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}

// Unique возвращает true если все строки в слайсе уникальны.
func Unique(values []string) bool {
	uniqueValues := make(map[string]bool)
	for _, value := range values {
		uniqueValues[value] = true
	}
	return len(values) == len(uniqueValues)
}
