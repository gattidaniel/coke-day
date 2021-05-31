package validators

import (
	"errors"
	"fmt"
	"gopkg.in/go-playground/validator.v9"
	"strconv"
	"strings"
)

var (
	CokeDomain  = "coke.com.us" // TODO: Avoid hardcoding this here
	PepsiDomain = "pepsi.com.us"
	validDomain = map[string]string{
		CokeDomain:  "C",
		PepsiDomain: "P",
	}
)

func CreateValidator() validator.Validate {
	v := validator.New()
	addEmailCustomValidation(v)
	addRoomCustomValidation(v)
	return *v
}

func addEmailCustomValidation(v *validator.Validate) {
	_ = v.RegisterValidation("emailCustom", func(fl validator.FieldLevel) bool {
		flParts := strings.Split(fl.Field().String(), "@")
		if len(flParts) < 2 {
			return false
		}
		domain := flParts[1]
		_, isValid := validDomain[domain]
		return isValid
	})
}

func addRoomCustomValidation(v *validator.Validate) {
	_ = v.RegisterValidation("roomCustom", func(fl validator.FieldLevel) bool {
		roomName := strings.ToLower(fl.Field().String())
		if !strings.HasPrefix(roomName, "c") && !strings.HasPrefix(roomName, "p") {
			return false
		}
		roomNumber := strings.TrimLeft(strings.TrimLeft(roomName, "c"), "p")
		number, err := strconv.Atoi(roomNumber)
		if err != nil {
			return false
		}
		if number < 1 || number > 10 {
			return false
		}
		return true
	})
}

func ParseValidationErrors(err error) error {
	var errorList []string
	for _, ef := range err.(validator.ValidationErrors) {
		field := ef.Field()
		tag := ef.Tag()
		errorList = append(errorList, fmt.Sprintf("field: %s, tag: %s", field, tag))
	}
	return errors.New(strings.Join(errorList, "\n"))
}
