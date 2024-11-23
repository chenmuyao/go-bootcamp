package web

import (
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"strconv"
	"text/template"
	"time"

	"github.com/chenmuyao/go-bootcamp/internal/domain"
	"github.com/chenmuyao/go-bootcamp/internal/service"
	"github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const (
	emailRegexPattern    = `^[a-zA-Z0-9._-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,6}$`
	passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,}$`
	bizLogin             = "login"
)

var JWTKey = []byte("xQUPmbb2TP9CUyFZkgOnV3JQdr22ZNBx")

type UserHandler struct {
	emailRegex    *regexp2.Regexp
	passwordRegex *regexp2.Regexp
	svc           *service.UserService
	codeSvc       *service.CodeService
}

type UserClaims struct {
	jwt.RegisteredClaims
	UID       int64
	UserAgent string
}

func NewUserHandler(svc *service.UserService, codeSvc *service.CodeService) *UserHandler {
	return &UserHandler{
		emailRegex:    regexp2.MustCompile(emailRegexPattern, regexp2.None),
		passwordRegex: regexp2.MustCompile(passwordRegexPattern, regexp2.None),
		svc:           svc,
		codeSvc:       codeSvc,
	}
}

func (h *UserHandler) RegisterRoutes(server *gin.Engine) {
	user := server.Group("/user/")
	user.POST("/signup", h.SignUp)
	// user.POST("/login", h.Login)
	user.POST("/login", h.LoginJWT)
	user.GET("/profile", h.Profile)
	user.GET("/profile/:id", h.Profile)
	user.POST("/edit", h.Edit)

	// SMS code login
	user.POST("/login_sms/code/send", h.SendSMSLoginCode)
	user.POST("/login_sms", h.LoginSMS)
}

func (h *UserHandler) SendSMSLoginCode(ctx *gin.Context) {
	type Req struct {
		Phone string `json:"phone"`
	}

	var req Req
	if err := ctx.Bind(&req); err != nil {
		slog.Error("bind error", "msg", err)
		return
	}

	if req.Phone == "" {
		ctx.JSON(http.StatusBadRequest, Result{
			Code: CodeUserSide,
			Msg:  "empty phone number",
		})
		return
	}
	tpl, err := template.New(bizLogin).
		Parse("Verification code for webook: {{.Code}}\nExpires in 10 min.\n[webook]")
	if err != nil {
		slog.Error("cannot parse sms template", "error", err)
		ctx.JSON(http.StatusInternalServerError, Result{
			Code: CodeServerSide,
			Msg:  "internal sever error",
		})
		return
	}
	err = h.codeSvc.Send(ctx, bizLogin, req.Phone, tpl)
	switch err {
	case nil:
		ctx.JSON(http.StatusOK, Result{
			Code: CodeOK,
			Msg:  "Sent successfully",
		})
		return
	case service.ErrCodeSendTooMany:
		ctx.JSON(http.StatusTooManyRequests, Result{
			Code: CodeUserSide,
			Msg:  "send too many",
		})
		return
	default:
		slog.Error("sms send code error", "error", err)
		ctx.JSON(http.StatusInternalServerError, InternalServerErrorResult)
		return
	}
}

func (h *UserHandler) LoginSMS(ctx *gin.Context) {
	type Req struct {
		Phone string `json:"phone"`
		Code  string `json:"code"`
	}

	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}

	if req.Phone == "" || req.Code == "" {
		ctx.JSON(http.StatusBadRequest, Result{
			Code: CodeUserSide,
			Msg:  "must have the phone number and the code",
		})
		return
	}

	ok, err := h.codeSvc.Verify(ctx, bizLogin, req.Phone, req.Code)
	if err != nil {
		slog.Error("verify", "msg", err)
		ctx.JSON(http.StatusInternalServerError, InternalServerErrorResult)
		return
	}
	if !ok {
		ctx.JSON(http.StatusBadRequest, Result{
			Code: CodeUserSide,
			Msg:  "Wrong code, please retry",
		})
		return
	}

	u, err := h.svc.FindOrCreate(ctx, req.Phone)
	if err != nil {
		slog.Error("find or create", "msg", err)
		ctx.JSON(http.StatusInternalServerError, InternalServerErrorResult)
		return
	}
	err = h.setJWTToken(ctx, u.ID)
	if err != nil {
		return // error message is set
	}
	ctx.JSON(http.StatusOK, "successful login")
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
		ctx.JSON(http.StatusInternalServerError, InternalServerErrorResult)
		return
	}
	if !isEmail {
		ctx.JSON(http.StatusBadRequest, Result{
			Code: CodeUserSide,
			Msg:  "not a valid email",
		})
		return
	}

	validPassword, err := h.passwordRegex.MatchString(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, InternalServerErrorResult)
		return
	}
	if !validPassword {
		ctx.JSON(http.StatusBadRequest, Result{
			Code: CodeUserSide,
			Msg:  "not a valid password",
		})
		return
	}

	if req.Password != req.ConfirmPassword {
		ctx.JSON(http.StatusBadRequest, Result{
			Code: CodeUserSide,
			Msg:  "2 passwords don't match",
		})
		return
	}
	u, err := h.svc.SignUp(ctx, domain.User{
		Email:    req.Email,
		Password: req.Password,
	})
	switch err {
	case nil:
		err = h.setJWTToken(ctx, u.ID)
		if err != nil {
			return // error message is set
		}
		ctx.JSON(http.StatusOK, Result{
			Code: CodeOK,
			Msg:  "signup success",
		})
	case service.ErrDuplicatedUser:
		ctx.JSON(http.StatusBadRequest, Result{
			Code: CodeUserSide,
			Msg:  err.Error(),
		})
	default:
		ctx.JSON(http.StatusInternalServerError, InternalServerErrorResult)
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
			MaxAge:   900, // 15min - expire time of the session (+ expire time of the userID entry in redis.)
			HttpOnly: true,
		})
		err = sess.Save()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, InternalServerErrorResult)
			return
		}
		ctx.JSON(http.StatusOK, Result{
			Code: CodeOK,
			Msg:  "successful login",
		})
	case service.ErrInvalidUserOrPassword:
		ctx.JSON(http.StatusBadRequest, Result{
			Code: CodeUserSide,
			Msg:  "wrong login or password",
		})
	default:
		ctx.JSON(http.StatusInternalServerError, InternalServerErrorResult)
	}
}

func (h *UserHandler) LoginJWT(ctx *gin.Context) {
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
		err = h.setJWTToken(ctx, u.ID)
		if err != nil {
			return // error message is set
		}
		ctx.JSON(http.StatusOK, "successful login")
	case service.ErrInvalidUserOrPassword:
		ctx.JSON(http.StatusBadRequest, Result{
			Code: CodeUserSide,
			Msg:  "wrong login or password",
		})
	default:
		ctx.JSON(http.StatusInternalServerError, InternalServerErrorResult)
	}
}

func (h *UserHandler) Profile(ctx *gin.Context) {
	var userID int64
	var err error
	id := ctx.Param("id")
	if id == "" {
		// get userID from session
		// userID, err = h.getUserIDFromSession(ctx)
		// if err != nil {
		// 	log.Println(err)
		// 	ctx.String(http.StatusInternalServerError, "system error")
		// 	return
		// }
		userID = h.getUserIDFromJWT(ctx)
	} else {
		if userID, err = strconv.ParseInt(id, 10, 64); err != nil {
			log.Println(err)
			ctx.JSON(http.StatusBadRequest, Result{
				Code: CodeUserSide,
				Msg:  fmt.Sprintf("unknown userID: %s", id),
			})
			return
		}
	}

	u, err := h.svc.GetProfile(ctx, userID)
	switch err {
	case nil:
		break
	case service.ErrInvalidUserID:
		ctx.JSON(http.StatusNotFound, Result{
			Code: CodeUserSide,
			Msg:  "unknown userID",
		})
		return
	default:
		log.Printf("failed to get user profile: %s\n", err.Error())
		ctx.JSON(http.StatusInternalServerError, InternalServerErrorResult)
		return
	}

	var birthday string
	if !u.Birthday.IsZero() {
		birthday = u.Birthday.Format("2006-01-02")
	}

	resp := struct {
		Email    string `json:"email"`
		Phone    string `json:"phone"`
		Name     string `json:"name"`
		Birthday string `json:"birthday"`
		Profile  string `json:"profile"`
	}{
		Email:    u.Email,
		Phone:    u.Phone,
		Name:     u.Name,
		Birthday: birthday,
		Profile:  u.Profile,
	}

	ctx.JSON(http.StatusOK, resp)
}

func (h *UserHandler) Edit(ctx *gin.Context) {
	type Req struct {
		Name     string `json:"name"`
		Birthday string `json:"birthday"`
		Profile  string `json:"profile"`
	}

	var req Req

	if err := ctx.Bind(&req); err != nil {
		log.Printf("Binding error: %s\n", err)
		return
	}

	// Get the userID from session
	// userID, err := h.getUserIDFromSession(ctx)
	// if err != nil {
	// 	log.Println(err)
	// 	ctx.String(http.StatusInternalServerError, "system error")
	// 	return
	// }
	userID := h.getUserIDFromJWT(ctx)

	// Update user profile
	// NOTE: if birthday is not set, set it to zero value. And it will be
	// ignored when getting the profile
	var birthday time.Time
	var err error
	if len(req.Birthday) != 0 {
		birthday, err = time.Parse("2006-01-02", req.Birthday)
		if err != nil {
			// NOTE: check should be done on the frontend. If we bypass the
			// frontend check, it must not be a normal user, and we don't care
			// about the error message.
			ctx.JSON(http.StatusInternalServerError, InternalServerErrorResult)
			return
		}
	}
	err = h.svc.EditProfile(ctx, &domain.User{
		ID:       userID,
		Name:     req.Name,
		Birthday: birthday,
		Profile:  req.Profile,
	})
	switch err {
	case nil:
		ctx.JSON(http.StatusOK, Result{
			Code: CodeOK,
			Msg:  "user profile update success",
		})
	case service.ErrInvalidUserID:
		ctx.JSON(http.StatusNotFound, Result{
			Code: CodeUserSide,
			Msg:  "unknown userID",
		})
	default:
		log.Printf("failed to update user profile: %s\n", err.Error())
		ctx.JSON(http.StatusInternalServerError, InternalServerErrorResult)
	}
}

func (h *UserHandler) getUserIDFromSession(ctx *gin.Context) (int64, error) {
	sess := sessions.Default(ctx)
	userID, ok := sess.Get("userID").(int64)
	if !ok {
		return 0, errors.New("failed to get userID from session")
	}
	return userID, nil
}

func (h *UserHandler) getUserIDFromJWT(ctx *gin.Context) int64 {
	uc := ctx.MustGet("user").(UserClaims)

	return uc.UID
}

func (h *UserHandler) setJWTToken(ctx *gin.Context, uid int64) error {
	uc := UserClaims{
		UID:       uid,
		UserAgent: ctx.GetHeader("User-Agent"),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 30)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, uc)
	tokenStr, err := token.SignedString(JWTKey)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, InternalServerErrorResult)
		return err
	}

	ctx.Header("x-jwt-token", tokenStr)
	return nil
}
