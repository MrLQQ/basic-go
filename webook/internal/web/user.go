package web

import (
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-gonic/gin"
	"net/http"
)

const emailRegexPatterm = "^[a-zA-Z0-9]+([-_.][a-zA-Z0-9]+)*@[a-zA-Z0-9]+([-_.][a-zA-Z0-9]+)*\\.[a-z]{2,}$"
const passwordRegexPatterm = "^(?=.*\\d)(?=.*[A-z])[\\da-zA-Z]{1,9}$"

type UserHandler struct {
	emailRexRxp    *regexp.Regexp
	passwordRexExp *regexp.Regexp
}

func NewUserHandler() *UserHandler {
	return &UserHandler{
		emailRexRxp:    regexp.MustCompile(emailRegexPatterm, regexp.None),
		passwordRexExp: regexp.MustCompile(passwordRegexPatterm, regexp.None),
	}
}

func (h *UserHandler) RegisterRoutes(server *gin.Engine) {

	ug := server.Group("/users")
	// 相当于/users/signup
	ug.POST("/signup", h.SignUp)
	// 相当于/users/login
	ug.POST("/login", h.Login)
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
		ctx.String(http.StatusOK, "密码格式非法：密码必须要同时包含字母和数字，位数1~9")
		return
	}
	if req.Password != req.ConfirmPassword {
		ctx.String(http.StatusOK, "两次密码不一致")
		return
	}
	ctx.String(http.StatusOK, "hello,Signup注册成功")
	return
}

func (h *UserHandler) Login(ctx *gin.Context) {

}

func (h *UserHandler) Edit(ctx *gin.Context) {

}

func (h *UserHandler) Profile(ctx *gin.Context) {

}
