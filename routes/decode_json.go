package routes

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/kayprogrammer/bidout-auction-v7/utils"
)


func DecodeJSONBody(c *fiber.Ctx, dst interface{}) (int, *utils.ErrorResponse) {
	var errData *utils.ErrorResponse
	code := 200
	if c.Get("Content-Type") != "application/json" {
		errD := utils.ErrorResponse{Message: "Content-Type header is not application/json"}.Init()
		errData = &errD
		return code, errData
	}

	dec := json.NewDecoder(bytes.NewReader(c.Body()))
	dec.DisallowUnknownFields()

	err := dec.Decode(&dst)

	msg := "Invalid Entry"
	fieldErrors := make(map[string]string)
	status_code := 422
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		errStr := err.Error()
		switch {
			case errors.As(err, &syntaxError):
				msg = fmt.Sprintf(
					"Request body contains badly-formed JSON (at position %d)",
					syntaxError.Offset,
				)

			case errors.Is(err, io.ErrUnexpectedEOF):
				status_code = http.StatusBadRequest
				msg="Request body contains badly-formed JSON"

			case errors.As(err, &unmarshalTypeError):
				fieldName := unmarshalTypeError.Field
				fieldErrors[fieldName] = "Invalid format"
			case strings.HasPrefix(errStr, "json: unknown field "):
				fieldName := strings.TrimPrefix(errStr, "json: unknown field ")
				fieldErrors[fieldName] = "Unknown field"
			case errors.Is(err, io.EOF):
				status_code = http.StatusBadRequest
				msg = "Request body must not be empty"

			case errStr == "http: request body too large":
				status_code = http.StatusRequestEntityTooLarge
				msg = "Request body must not be larger than 1MB"

			default:
				status_code = 400
				msg = "Invalid request"
		}
		errData := utils.ErrorResponse{Message: msg}.Init()
		if len(fieldErrors) > 0 {
			errData.Data = &fieldErrors
		}
		code = status_code
		return code, &errData
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		errData := utils.ErrorResponse{Message: "Request body must only contain a single JSON object"}.Init()
		return 400, &errData
	}
	return code, nil
}