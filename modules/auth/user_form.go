package auth

type SignUpForm struct {
	Nickname        string `form:"nickname" json:"nickname" binding:"required,nickname"`
	Username        string `form:"username" json:"username" binding:"required,username"`
	Email           string `form:"email" json:"email" binding:"required,email"`
	Password        string `form:"password" json:"password" binding:"required,password"`
	EmailVerifyCode string `form:"email_verify_code" json:"email_verify_code" binding:"required,len=6,numeric"`
}

type SignInForm struct {
	Username string `form:"username" json:"username" binding:"required,email"`
	Password string `form:"password" json:"password" binding:"required,password"`
}

type RequestActiveEmailForm struct {
	Email string `form:"email" json:"email" binding:"required,email"`
}
