package web

import (
	"net/http"

	"github.com/dlclark/regexp2"
	"github.com/gin-gonic/gin"
)

const (
	emailRegexPattern    = `^[a-zA-Z0-9._-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,6}$`
	passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,}$`
)

type UserHandler struct {
	emailRegex    *regexp2.Regexp
	passwordRegex *regexp2.Regexp
}

func NewUserHandler() *UserHandler {
	return &UserHandler{
		emailRegex:    regexp2.MustCompile(emailRegexPattern, regexp2.None),
		passwordRegex: regexp2.MustCompile(passwordRegexPattern, regexp2.None),
	}
}

func (h *UserHandler) RegisterRoutes(server *gin.Engine) {
	user := server.Group("/user/")
	user.POST("/signup", h.SignUp)
	user.POST("/login", h.Login)
	user.GET("/profile", h.Profile)
	user.POST("/profile", h.Edit)
}

func (h *UserHandler) SignUp(ctx *gin.Context) {
	type SignUpReq struct {
		Email           string `json:"email"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirm_password"`
	}

	// Get request
	var req SignUpReq

	if err := ctx.Bind(&req); err != nil {
		return
	}

	// Check request
	isEmail, err := h.emailRegex.MatchString(req.Email)
	if err != nil {
		ctx.String(http.StatusInternalServerError, "internal server error")
		return
	}
	if !isEmail {
		ctx.String(http.StatusBadRequest, "not a valid email")
		return
	}

	validPassword, err := h.passwordRegex.MatchString(req.Password)
	if err != nil {
		ctx.String(http.StatusInternalServerError, "internal server error")
		return
	}
	if !validPassword {
		ctx.String(http.StatusBadRequest, "not a valid password")
		return
	}

	if req.Password != req.ConfirmPassword {
		ctx.String(http.StatusBadRequest, "2 passwords don't match")
		return
	}
}

func (h *UserHandler) Login(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, "User Login")
}

func (h *UserHandler) Profile(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, "Get User profile")
}

func (h *UserHandler) Edit(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, "Edit User profile")
}