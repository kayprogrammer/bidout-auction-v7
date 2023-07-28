package schemas

type EmailRequestSchema struct {
	Email				string				`json:"email" validate:"required,min=5,email" example:"johndoe@email.com"`
}

type VerifyEmailRequestSchema struct {
	Email				string				`json:"email" validate:"required,min=5,email" example:"johndoe@example.com"`
	Otp					int					`json:"otp" validate:"required" example:"123456"`
}


type RegisterResponseSchema struct {
	ResponseSchema
	Data			EmailRequestSchema		`json:"data"`
}