package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	err := strings.Builder{}
	for _, v := range v {
		err.WriteString(v.Err.Error() + ", ")
	}
	return err.String()
}

var (
	ErrNotAStruct    = errors.New("not a struct")
	ErrConvertToInt  = errors.New("could not convert to int")
	ErrCompileRegexp = errors.New("could not compile regexp")

	ErrValidation           = errors.New("could not validate field")
	ErrValidateStringLen    = fmt.Errorf("%w: string not correct length", ErrValidation)
	ErrValidateStringRegexp = fmt.Errorf("%w: string not correct regexp", ErrValidation)
	ErrValidateStringIn     = fmt.Errorf("%w: string is not in array", ErrValidation)
	ErrValidateIntMin       = fmt.Errorf("%w: integer less than minimum", ErrValidation)
	ErrValidateIntMax       = fmt.Errorf("%w: integer more than maximum", ErrValidation)
	ErrValidateIntIn        = fmt.Errorf("%w: integer is not in array", ErrValidation)
)

func Validate(v interface{}) error {
	reflVal := reflect.ValueOf(v)
	reflType := reflect.TypeOf(v)

	var errs ValidationErrors

	if reflVal.Kind() != reflect.Struct {
		return ErrNotAStruct
	}

	for i := 0; i < reflVal.NumField(); i++ {
		fieldName := reflVal.Type().Field(i).Name
		value := reflVal.Field(i)
		tag := reflType.Field(i).Tag.Get("validate")

		if tag == "" {
			continue
		}

		parts := strings.Split(tag, "|")
		for _, part := range parts {
			err := validateStruct(value, part)
			if err != nil {
				if errors.Is(err, ErrValidation) {
					errs = append(errs, ValidationError{
						Field: fieldName,
						Err:   err,
					})
				} else {
					return err
				}
			}
		}
	}

	if len(errs) > 0 {
		return errs
	}

	return nil
}

func validateStruct(field reflect.Value, tag string) error {
	parts := strings.Split(tag, ":")
	tagKey, tagValue := parts[0], parts[1]

	switch field.Kind() { //nolint:exhaustive
	case reflect.String:
		return validateString(field, tagKey, tagValue)
	case reflect.Int:
		return validateInt(field, tagKey, tagValue)
	case reflect.Slice:
		return validateSlice(field, tag)
	}
	return nil
}

func validateString(field reflect.Value, tagKey, tagValue string) error {
	switch tagKey {
	case "len":
		val, err := strconv.Atoi(tagValue)
		if err != nil {
			return ErrConvertToInt
		}
		if len(field.String()) != val {
			return ErrValidateStringLen
		}
		return nil
	case "regexp":
		reg, err := regexp.Compile(tagValue)
		if err != nil {
			return ErrCompileRegexp
		}
		if !reg.MatchString(field.String()) {
			return ErrValidateStringRegexp
		}
		return nil
	case "in":
		arrayVals := strings.Split(tagValue, ",")
		for _, val := range arrayVals {
			if field.String() == val {
				return nil
			}
		}
		return ErrValidateStringIn
	}
	return nil
}

func validateInt(field reflect.Value, tagKey, tagValue string) error {
	switch tagKey {
	case "min":
		val, err := strconv.Atoi(tagValue)
		if err != nil {
			return ErrConvertToInt
		}
		if int(field.Int()) < val {
			return ErrValidateIntMin
		}
		return nil
	case "max":
		val, err := strconv.Atoi(tagValue)
		if err != nil {
			return ErrConvertToInt
		}
		if int(field.Int()) > val {
			return ErrValidateIntMax
		}
		return nil
	case "in":
		arrayVals := strings.Split(tagValue, ",")
		for _, val := range arrayVals {
			val, err := strconv.Atoi(val)
			if err != nil {
				return ErrConvertToInt
			}
			if int(field.Int()) == val {
				return nil
			}
		}
		return ErrValidateIntIn
	}
	return nil
}

func validateSlice(field reflect.Value, tag string) error {
	for i := 0; i < field.Len(); i++ {
		err := validateStruct(field.Index(i), tag)
		if err != nil {
			return err
		}
	}

	return nil
}
