package ginx

import "net/http"

// TODO: Use a better error code design

const (
	CodeOK         = 2
	CodeRedirect   = 3
	CodeUserSide   = 4
	CodeUnauth     = 401
	CodeServerSide = 5
)

type Result struct {
	Data any    `json:"data"`
	Msg  string `json:"message"`
	Code int    `json:"code"`
}

var InternalServerErrorResult = Result{
	Code: CodeServerSide,
	Msg:  "internal server error",
}

var UnauthorizedResult = Result{
	Code: CodeUnauth,
	Msg:  "unauthorized request",
}

func ResultToStatus(res Result) int {
	switch res.Code {
	case 2:
		return http.StatusOK
	case 3:
		return http.StatusSeeOther
	case 4:
		return http.StatusBadRequest
	case 401:
		return http.StatusUnauthorized
	case 5:
		return http.StatusInternalServerError
	default:
		return http.StatusOK
	}
}
