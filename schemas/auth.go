package schemas

// REQUEST BODY SCHEMAS
type EmailRequestSchema struct {
	Email				string				`json:"email" validate:"required,min=5,email" example:"johndoe@email.com"`
}

type VerifyEmailRequestSchema struct {
	Email				string				`json:"email" validate:"required,min=5,email" example:"johndoe@example.com"`
	Otp					int					`json:"otp" validate:"required" example:"123456"`
}

type SetNewPasswordSchema struct {
	VerifyEmailRequestSchema
	Password			string				`json:"password" validate:"required,min=8,max=50" example:"newstrongpassword"`
}

type LoginSchema struct {
	Email				string				`json:"email" validate:"required,email" example:"johndoe@email.com"`
	Password			string				`json:"password" validate:"required" example:"password"`
}

type RefreshTokenSchema struct {
	Refresh			string					`json:"refresh" validate:"required" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InNpbXBsZWlkIiwiZXhwIjoxMjU3ODk0MzAwfQ.Ys_jP70xdxch32hFECfJQuvpvU5_IiTIN2pJJv68EqQ"`
}

// RESPONSE BODY SCHEMAS
type RegisterResponseSchema struct {
	ResponseSchema
	Data			EmailRequestSchema		`json:"data"`
}

type TokensResponseSchema struct {
	Access			string					`json:"access"`
	Refresh			string					`json:"refresh"`
}

type LoginResponseSchema struct {
	ResponseSchema
	Data			TokensResponseSchema		`json:"data"`
}
