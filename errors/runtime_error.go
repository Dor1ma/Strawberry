package errors

import (
	"github.com/Dor1ma/Strawberry/token"
)

type RuntimeError struct {
	s     string
	token token.Token
}

func (r *RuntimeError) Error() string {
	return r.s
}

// Класс для возврата ошибок в рантайме
func Error(token token.Token, s string) {
	panic(RuntimeError{token: token, s: s})
}
