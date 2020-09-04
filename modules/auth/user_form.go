package auth

type SignUpForm struct {
	Nickname        string `form:"nickname" json:"nickname" binding:"required,min=4,max=20,alphanum"`
	Username        string `form:"username" json:"username" binding:"required,min=6,max=16,alphanum"`
	Email           string `form:"email" json:"email" binding:"required,email"`
	Password        string `form:"password" json:"password" binding:"required,min=6,max=16"`
	EmailVerifyCode string `form:"email_verify_code" json:"email_verify_code" binding:"required,len=6,numeric"`
}

type SignInForm struct {
	Username string `form:"username" json:"username" binding:"required,email"`
	Password string `form:"password" json:"password" binding:"required,min=6,max=16"`
}

type RequestActiveEmailForm struct {
	Email string `form:"email" json:"email" binding:"required,email"`
}
