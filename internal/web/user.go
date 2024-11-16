package web

import (
	"log"
	"net/http"
	"time"

	"github.com/chenmuyao/go-bootcamp/internal/domain"
	"github.com/chenmuyao/go-bootcamp/internal/service"
	"github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

const (
	emailRegexPattern    = `^[a-zA-Z0-9._-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,6}$`
	passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,}$`
)

type UserHandler struct {
	emailRegex    *regexp2.Regexp
	passwordRegex *regexp2.Regexp
	svc           *service.UserService
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{
		emailRegex:    regexp2.MustCompile(emailRegexPattern, regexp2.None),
		passwordRegex: regexp2.MustCompile(passwordRegexPattern, regexp2.None),
		svc:           svc,
	}
}

func (h *UserHandler) RegisterRoutes(server *gin.Engine) {
	user := server.Group("/user/")
	user.POST("/signup", h.SignUp)
	user.POST("/login", h.Login)
	user.GET("/profile", h.Profile)
	user.POST("/edit", h.Edit)
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
		ctx.String(http.StatusOK, "2 passwords don't match")
		return
	}
	err = h.svc.SignUp(ctx, domain.User{
		Email:    req.Email,
		Password: req.Password,
	})
	switch err {
	case nil:
		ctx.String(http.StatusOK, "signup success")
	case service.ErrDuplicatedEmail:
		ctx.String(http.StatusBadRequest, err.Error())
	default:
		ctx.String(http.StatusInternalServerError, "system error")
	}
}

func (h *UserHandler) Login(ctx *gin.Context) {
	type Req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var req Req

	if err := ctx.Bind(&req); err != nil {
		return
	}

	// NOTE: No need to check, because if it's not valid, we won't get
	// anything from the DB anyway.

	u, err := h.svc.Login(ctx, req.Email, req.Password)
	switch err {
	case nil:
		sess := sessions.Default(ctx)
		sess.Set("userID", u.ID)
		sess.Options(sessions.Options{
			MaxAge:   900, // 15min
			HttpOnly: true,
		})
		err = sess.Save()
		if err != nil {
			ctx.String(http.StatusInternalServerError, "system error")
			return
		}
		ctx.String(http.StatusOK, "successful login")
	case service.ErrInvalidUserOrPassword:
		ctx.String(http.StatusBadRequest, "wrong login or password")
	default:
		ctx.String(http.StatusInternalServerError, "system error")
	}
}

func (h *UserHandler) Profile(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, "Get User profile")
}

func (h *UserHandler) Edit(ctx *gin.Context) {
	type Req struct {
		Name     string `json:"name"     binding:"max=50,omitempty"`
		Birthday string `json:"birthday" binding:"date,omitempty"`
		Profile  string `json:"profile"  binding:"max=2000,omitempty"`
	}

	var req Req

	if err := ctx.Bind(&req); err != nil {
		log.Printf("Binding error: %s\n", err)
		return
	}

	// Get the userID from session
	sess := sessions.Default(ctx)
	userID, ok := sess.Get("userID").(int64)
	if !ok {
		log.Printf("failed to get userID from session")
		ctx.String(http.StatusInternalServerError, "system error")
	}

	// Update user profile
	// Ignore error because it is already checked.
	birthday, _ := time.Parse("2006-01-02", req.Birthday)
	err := h.svc.EditProfile(ctx, &domain.User{
		ID:       userID,
		Name:     req.Name,
		Birthday: birthday,
		Profile:  req.Profile,
	})
	switch err {
	case nil:
		ctx.String(http.StatusOK, "user profile update success")
	case service.ErrInvalidUserID:
		ctx.String(http.StatusBadRequest, "unknown userID")
	default:
		log.Printf("failed to update user profile: %s\n", err.Error())
		ctx.String(http.StatusInternalServerError, "system error")
	}
}
