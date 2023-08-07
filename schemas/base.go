package schemas

type ResponseSchema struct {
	Status		string		`json:"status" example:"success"`
	Message		string		`json:"message" example:"Data fetched/created/updated/deleted"`
}

func (obj ResponseSchema) Init() ResponseSchema {
	if obj.Status == "" {
		obj.Status = "success"
	}
	return obj
}