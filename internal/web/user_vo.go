package web

type UserSignUpReq struct {
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}

type UserLoginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserLoginSMSReq struct {
	Phone string `json:"phone"`
	Code  string `json:"code"`
}

type UserSMSCodeReq struct {
	Phone string `json:"phone"`
}

type UserEditReq struct {
	Name     string `json:"name"`
	Birthday string `json:"birthday"`
	Profile  string `json:"profile"`
}
