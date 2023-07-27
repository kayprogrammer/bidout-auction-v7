package utils

type ErrorResponse struct {
	Status				string				`json:"status"`
	Message				string				`json:"message"`
	Data				*map[string]string	`json:"data,omitempty"`
}

func (obj ErrorResponse) Init() ErrorResponse {
	if obj.Status == "" {
		obj.Status = "failure"
	}
	return obj
}