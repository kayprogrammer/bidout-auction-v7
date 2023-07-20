package schemas

type ResponseSchema struct {
	Status		string		`json:"status"`
	Message		string		`json:"message"`
}

func (obj ResponseSchema) Init() ResponseSchema {
	if obj.Status == "" {
		obj.Status = "success"
	}
	return obj
}