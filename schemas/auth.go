package schemas

type EmailRequestSchema struct {
	Email				string				`json:"email"`
}

type RegisterResponseSchema struct {
	ResponseSchema
	Data			EmailRequestSchema		`json:"data"`
}