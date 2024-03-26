package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

var (
	ErrorLength         = errors.New("length")
	ErrorRegex          = errors.New("regex")
	ErrorMin            = errors.New("greater")
	ErrorMax            = errors.New("less")
	ErrorIn             = errors.New("lots of")
	ErrorExpectedStruct = errors.New("expected a struct")
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	var builder strings.Builder
	for _, err := range v {
		builder.WriteString("Field: " + err.Field + ", error: " + err.Err.Error() + "\n")
	}
	return builder.String()
}

// Validate - функция для валидации полей структуры, основываясь на тегах validate.
func Validate(v interface{}) error {
	validationErrors := make(ValidationErrors, 0)
	value := reflect.ValueOf(v)
	if value.Kind() != reflect.Struct {
		return fmt.Errorf("%w, received %s ", ErrorExpectedStruct, value.Kind())
	}
	for i := 0; i < value.NumField(); i++ {
		field := value.Type().Field(i)
		validateTag, ok := field.Tag.Lookup("validate")
		if !ok {
			continue
		}
		validationErrors = checkValue(validationErrors, field.Name, validateTag, value.Field(i))
	}
	if len(validationErrors) > 0 {
		return validationErrors
	}
	return nil
}

// checkLen - проверяет длину строки на соответствие указанному значению.
func checkLen(rv reflect.Value, ruleValue string) bool {
	if rv.Kind() == reflect.String {
		intValue, err := strconv.Atoi(ruleValue)
		if err != nil {
			return false
		}
		return rv.Len() == intValue
	}
	return false
}

// checkRegex - проверяет соответствие строки регулярному выражению.
func checkRegex(rv reflect.Value, ruleValue string) bool {
	if rv.Kind() == reflect.String {
		rx, err := regexp.Compile(ruleValue)
		if err != nil {
			return false
		}
		return rx.Match([]byte(rv.String()))
	}
	return false
}

// checkMin - проверяет, что числовое значение больше указанного минимума.
func checkMin(rv reflect.Value, ruleValue string) bool {
	if rv.Kind() == reflect.Int {
		intValue := int(rv.Int())
		min, err := strconv.Atoi(ruleValue)
		if err != nil {
			return false
		}
		return intValue > min
	}
	return false
}

// checkMax - проверяет, что числовое значение меньше указанного максимума.
func checkMax(rv reflect.Value, ruleValue string) bool {
	if rv.Kind() == reflect.Int {
		intValue := int(rv.Int())
		max, err := strconv.Atoi(ruleValue)
		if err != nil {
			return false
		}
		return intValue < max
	}
	return false
}

// checkIn - проверяет, содержится ли значение в указанном наборе.
func checkIn(rv reflect.Value, ruleValue string) bool {
	ins := strings.Split(ruleValue, ",")
	isValid := false

	switch rv.Kind() {
	case reflect.Int:
		intValue := int(rv.Int())

		for _, in := range ins {
			in, err := strconv.Atoi(in)
			if err != nil {
				continue
			}
			if in == intValue {
				isValid = true
			}
		}
	case reflect.String:
		strValue := rv.String()

		for _, in := range ins {
			if in == strValue {
				isValid = true
			}
		}
	case reflect.Array, reflect.Bool, reflect.Chan, reflect.Complex128, reflect.Complex64, reflect.Float32,
		reflect.Float64, reflect.Func, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int8, reflect.Interface,
		reflect.Invalid, reflect.Map, reflect.Ptr, reflect.Struct, reflect.Slice, reflect.Uint, reflect.Uint16,
		reflect.Uint32, reflect.Uint64, reflect.Uint8, reflect.Uintptr, reflect.UnsafePointer:
	}
	return isValid
}

// validateValue - валидирует значение на основе правил, указанных в теге validate.
func validateValue(validateTag string, rv reflect.Value) []error {
	rules := strings.Split(validateTag, "|")
	errs := make([]error, 0)

	for _, rule := range rules {
		r := strings.Split(rule, ":")
		if len(r) != 2 {
			continue
		}

		rType, rValue := r[0], r[1]
		var err error

		switch rType {
		case "len":
			if !checkLen(rv, rValue) {
				err = fmt.Errorf("%w must be equal %s", ErrorLength, rValue)
			}
		case "regexp":
			if !checkRegex(rv, rValue) {
				err = fmt.Errorf("must match %w %s", ErrorRegex, rValue)
			}
		case "min":
			if !checkMin(rv, rValue) {
				err = fmt.Errorf("must be %w than %s", ErrorMin, rValue)
			}
		case "max":
			if !checkMax(rv, rValue) {
				err = fmt.Errorf("must be %w than %s", ErrorMax, rValue)
			}
		case "in":
			if !checkIn(rv, rValue) {
				err = fmt.Errorf("must be %w %s", ErrorIn, rValue)
			}
		default:
			continue
		}

		if err != nil {
			errs = append(errs, err)
		}
	}
	return errs
}

// checkValue - обрабатывает значение поля и выполняет валидацию согласно правилам.
func checkValue(valErrs ValidationErrors, fName string, validateTag string, rv reflect.Value) ValidationErrors {
	var (
		errs       []error
		newValErrs = valErrs
	)
	switch rv.Kind() {
	case reflect.String:
		errs = validateValue(validateTag, rv)
	case reflect.Int:
		errs = validateValue(validateTag, rv)
	case reflect.Slice:
		for i := 0; i < rv.Len(); i++ {
			newValErrs = checkValue(newValErrs, fName, validateTag, rv.Index(i))
		}
	case reflect.Array, reflect.Bool, reflect.Chan, reflect.Complex128, reflect.Complex64, reflect.Float32,
		reflect.Float64, reflect.Func, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int8, reflect.Interface,
		reflect.Invalid, reflect.Map, reflect.Ptr, reflect.Struct, reflect.Uint, reflect.Uint16,
		reflect.Uint32, reflect.Uint64, reflect.Uint8, reflect.Uintptr, reflect.UnsafePointer:
	}

	if len(errs) > 0 {
		for _, err := range errs {
			newValErrs = append(newValErrs, ValidationError{fName, err})
		}
	}

	return newValErrs
}
