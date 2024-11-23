package web

// TODO: Use a better error code design

const (
	CodeOK         = 2
	CodeRedirect   = 3
	CodeUserSide   = 4
	CodeServerSide = 5
)

type Result struct {
	Code int    `json:"code"`
	Msg  string `json:"message"`
	Data any    `json:"data"`
}

var InternalServerErrorResult = Result{
	Code: CodeServerSide,
	Msg:  "internal server error",
}
