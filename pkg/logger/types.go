package logger

// {{{ Consts

// }}}
// {{{ Global Varirables

// }}}
// {{{ Interface

type Logger interface {
	Debug(msg string, args ...Field)
	Info(msg string, args ...Field)
	Warn(msg string, args ...Field)
	Error(msg string, args ...Field)
	With(args ...Field) Logger
}

// }}}
// {{{ Struct

// }}}
// {{{ Other structs

type Field struct {
	Key   string
	Value any
}

// }}}
// {{{ Struct Methods

// }}}
// {{{ Private functions

// }}}
// {{{ Package functions

func Error(err error) Field {
	return Field{
		Key:   "err",
		Value: err,
	}
}

// }}}
