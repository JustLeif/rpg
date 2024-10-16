package utils

import (
	"fmt"
	"io"
)

type Environment string

const (
	Production  Environment = "production"
	Development Environment = "development"
	Preview     Environment = "preview"
	UnitTest    Environment = "unit_test"
)

type Logger struct {
	Environment Environment
}

func NewLogger(Environment Environment) *Logger {
	return &Logger{
		Environment: Environment,
	}
}

/* Behaves like `fmt.Fprintf()`, but only in the development and preview environment. */
func (l Logger) DevLog(w io.Writer, format string, a ...interface{}) (n int, err error) {
	if l.Environment == Production {
		return 0, nil
	}
	str := append([]byte(format), []byte("\n")...)
	return fmt.Fprintf(w, string(str), a...)
}

/* Behaves like `fmt.Fprintf()`, but only in the production environment. */
func (l Logger) ProdLog(w io.Writer, format string, a ...interface{}) (n int, err error) {
	if l.Environment == Production {
		str := append([]byte(format), []byte("\n")...)
		return fmt.Fprintf(w, string(str), a...)
	} else {
		return 0, nil
	}
}

/* Behaves like `fmt.Printf()`. */
func (l Logger) Log(w io.Writer, format string, a ...interface{}) (n int, err error) {
	str := append([]byte(format), []byte("\n")...)
	return fmt.Fprintf(w, string(str), a...)
}
