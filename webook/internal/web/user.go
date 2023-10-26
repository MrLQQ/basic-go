package web

import (
	"basic-go/webook/internal/domain"
	"basic-go/webook/internal/service"
	"errors"
	"fmt"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v5"
	"net/http"
	"time"
)

const emailRegexPatterm = "^[a-zA-Z0-9]+([-_.][a-zA-Z0-9]+)*@[a-zA-Z0-9]+([-_.][a-zA-Z0-9]+)*\\.[a-z]{2,}$"
const passwordRegexPatterm = "^(?=.*\\d)(?=.*[A-z])[\\da-zA-Z]{1,15}$"
const nicknameRegexPatterm = "^[\\u4e00-\\u9fa5_a-zA-Z0-9_]{4,10}$"
const aboutMeRegexPatterm = "^.{0,100}$"

type UserHandler struct {
	emailRexRxp    *regexp.Regexp
	passwordRexExp *regexp.Regexp
	nicknameRexExp *regexp.Regexp
	aboutMeRexExp  *regexp.Regexp
	svc            *service.UserService
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{
		emailRexRxp:    regexp.MustCompile(emailRegexPatterm, regexp.None),
		passwordRexExp: regexp.MustCompile(passwordRegexPatterm, regexp.None),
		nicknameRexExp: regexp.MustCompile(nicknameRegexPatterm, regexp.None),
		aboutMeRexExp:  regexp.MustCompile(aboutMeRegexPatterm, regexp.None),
		svc:            svc,
	}
}

func (h *UserHandler) RegisterRoutes(server *gin.Engine) {

	ug := server.Group("/users")
	// 相当于/users/signup
	ug.POST("/signup", h.SignUp)
	// 相当于/users/login
	//ug.POST("/login", h.Login)
	ug.POST("/login", h.LoginJWT)
	// 相当于/users/edit
	ug.POST("/edit", h.Edit)
	// 相当于/users/profile
	ug.GET("/profile", h.Profile)
}

func (h *UserHandler) SignUp(ctx *gin.Context) {
	type SignupReq struct {
		Email           string `json:"email"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirmPassword"`
	}
	var req SignupReq
	if err := ctx.Bind(&req); err != nil {
		return
	}
	isEmail, err := h.emailRexRxp.MatchString(req.Email)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误……")
		return
	}
	if !isEmail {
		ctx.String(http.StatusOK, "邮箱格式非法")
		return
	}

	isPassword, err := h.passwordRexExp.MatchString(req.Password)

	if err != nil {
		ctx.String(http.StatusOK, "系统错误……")
		return
	}
	if !isPassword {
		ctx.String(http.StatusOK, "密码格式非法：密码必须要同时包含字母和数字，位数1~15")
		return
	}
	if req.Password != req.ConfirmPassword {
		ctx.String(http.StatusOK, "两次密码不一致")
		return
	}

	err = h.svc.Signup(ctx, domain.User{
		Email:    req.Email,
		Password: req.Password,
	})
	// 要判定邮箱冲突
	switch {
	case err == nil:
		ctx.String(http.StatusOK, "注册成功")
	case errors.Is(err, service.ErrDuplicateEmail):
		ctx.String(http.StatusOK, "邮箱冲突,请换一个")
	default:
		ctx.String(http.StatusOK, "系统错误")
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
	u, err := h.svc.Login(ctx, req.Email, req.Password)
	switch {
	case err == nil:
		uc := UserClaims{
			Uid:       u.Id,
			UserAgent: ctx.GetHeader("User-Agent"),
			RegisteredClaims: jwt.RegisteredClaims{
				// 30分钟过期
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 30)),
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS512, uc)
		tokenStr, err := token.SignedString(JWTkey)
		if err != nil {
			ctx.String(http.StatusOK, "系统错误")
		}
		ctx.Header("x-jwt-token", tokenStr)
		ctx.String(http.StatusOK, "登录成功")
	case errors.Is(err, service.ErrInvalidUserOrPassword):
		ctx.String(http.StatusOK, "用户名或密码错误")
	default:
		ctx.String(http.StatusOK, "系统错误")
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
	u, err := h.svc.Login(ctx, req.Email, req.Password)
	switch {
	case err == nil:

		sess := sessions.Default(ctx)
		sess.Set("userId", u.Id)
		sess.Options(sessions.Options{
			// 有效期 15分钟
			MaxAge: 900,
		})
		err := sess.Save()
		if err != nil {
			ctx.String(http.StatusOK, "系统错误")
			return
		}
		ctx.String(http.StatusOK, "登录成功")
	case errors.Is(err, service.ErrInvalidUserOrPassword):
		ctx.String(http.StatusOK, "用户名或密码错误")
	default:
		ctx.String(http.StatusOK, "系统错误")
	}
}

func (h *UserHandler) Edit(ctx *gin.Context) {
	us := ctx.MustGet("user").(UserClaims)
	type EditReq struct {
		Nickname string `json:"nickname"`
		Birthday string `json:"birthday"`
		AboutMe  string `json:"aboutMe"`
	}
	var req EditReq
	if err := ctx.Bind(&req); err != nil {
		return
	}

	isNickName, err := h.nicknameRexExp.MatchString(req.Nickname)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误……")
		return
	}
	if !isNickName {
		ctx.String(http.StatusOK, "昵称只能由汉字、字母、数字、下划线组成，长度4~10位")
		return
	}

	isAboutMe, err := h.aboutMeRexExp.MatchString(req.AboutMe)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误……")
		return
	}
	if !isAboutMe {
		ctx.String(http.StatusOK, "个人介绍不能超过100个字符")
		return
	}
	//sess := sessions.Default(ctx)
	//userID := sess.Get("userId")
	userID := us.Uid
	if userID == 0 {
		// 中断，不要往后执行，也就是不要执行后面的业务逻辑
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	user := fmt.Sprintf("%v", userID)
	err = h.svc.Edit(ctx, domain.UserProfile{
		User_id:  user,
		Nickname: req.Nickname,
		Birthday: req.Birthday,
		About_me: req.AboutMe,
	})
	switch {
	case err == nil:
		ctx.String(http.StatusOK, "更新成功")
	default:
		ctx.String(http.StatusOK, "系统错误")
	}

}

func (h *UserHandler) Profile(ctx *gin.Context) {
	us := ctx.MustGet("user").(UserClaims)
	type Profile struct {
		Nickname string `json:"nickname"`
		Birthday string `json:"birthday"`
		AboutMe  string `json:"aboutMe"`
	}
	//sess := sessions.Default(ctx)
	//userID := sess.Get("userId")
	userID := us.Uid
	if userID == 0 {
		// 中断，不要往后执行，也就是不要执行后面的业务逻辑
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	user := fmt.Sprintf("%v", userID)
	profile, err := h.svc.Profile(ctx, domain.UserProfile{
		User_id: user,
	})
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
	}
	userProfile := Profile{
		Nickname: profile.Nickname,
		Birthday: profile.Birthday,
		AboutMe:  profile.About_me,
	}
	ctx.JSON(http.StatusOK, userProfile)
}

var JWTkey = []byte("jBxoQWRS5L9vYr$mYq5U9d5BRPHfSBAe")

type UserClaims struct {
	jwt.RegisteredClaims
	Uid       int64
	UserAgent string
}
