package logger

import "fmt"

type LogError struct {
	Msg    string
	Fields []Field
}

func (l *LogError) Error() string {
	return fmt.Sprintf("%s:%+v", l.Msg, l.Fields)
}

func LError(msg string, fields ...Field) *LogError {
	return &LogError{
		Msg:    msg,
		Fields: fields,
	}
}
