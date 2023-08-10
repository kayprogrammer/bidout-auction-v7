package schemas

// REQUEST BODY SCHEMAS


// RESPONSE BODY SCHEMAS
type ProfileResponseDataSchema struct {
	FirstName				string				`json:"first_name"`
	LastName				string				`json:"last_name"`
	Avatar					*string				`json:"avatar"`
}

type ProfileResponseSchema struct {
	ResponseSchema
	Data					ProfileResponseDataSchema			`json:"data"`
}