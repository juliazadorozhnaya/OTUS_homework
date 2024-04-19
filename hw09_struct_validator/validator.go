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
	ErrorLength           = errors.New("length")
	ErrorRegex            = errors.New("regex")
	ErrorMin              = errors.New("greater")
	ErrorMax              = errors.New("less")
	ErrorIn               = errors.New("lots of")
	ErrorExpectedStruct   = errors.New("expected a struct")
	ErrorInvalidValidator = errors.New("invalid validator for the type")
	regexCache            = make(map[string]*regexp.Regexp)
	intCache              = make(map[string]int)
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

func Validate(v interface{}) error {
	validationErrors := make(ValidationErrors, 0)
	value := reflect.ValueOf(v)
	if value.Kind() != reflect.Struct {
		return fmt.Errorf("%w, received %s", ErrorExpectedStruct, value.Kind())
	}
	for i := 0; i < value.NumField(); i++ {
		field := value.Type().Field(i)
		validateTag, ok := field.Tag.Lookup("validate")
		if !ok {
			continue
		}
		err := checkValue(&validationErrors, field.Name, validateTag, value.Field(i))
		if err != nil {
			return err
		}
	}
	if len(validationErrors) > 0 {
		return validationErrors
	}
	return nil
}

func checkValue(valErrs *ValidationErrors, fName string, validateTag string, rv reflect.Value) error {
	var errs []error
	//nolint:exhaustive
	switch rv.Kind() {
	case reflect.String, reflect.Int:
		errs = validateValue(validateTag, rv)
	case reflect.Slice:
		for i := 0; i < rv.Len(); i++ {
			err := checkValue(valErrs, fName, validateTag, rv.Index(i))
			if err != nil {
				return err
			}
		}
	}
	if len(errs) > 0 {
		for _, err := range errs {
			*valErrs = append(*valErrs, ValidationError{fName, err})
		}
	}
	return nil
}

func isTypeValidForRule(ruleType string, rv reflect.Value) bool {
	switch ruleType {
	case "len", "regexp", "in":
		return rv.Kind() == reflect.String
	case "min", "max":
		return rv.Kind() == reflect.Int
	default:
		return true
	}
}

func validateValue(validateTag string, rv reflect.Value) []error {
	rules := strings.Split(validateTag, "|")
	var errs []error
	for _, rule := range rules {
		r := strings.Split(rule, ":")
		if len(r) != 2 {
			continue
		}
		rType, rValue := r[0], r[1]
		if !isTypeValidForRule(rType, rv) {
			return []error{fmt.Errorf("%w: %s rule cannot be applied to %s type", ErrorInvalidValidator, rType, rv.Kind())}
		}
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

// getIntFromCache извлекает целочисленное значение из кэша или добавляет его при отсутствии.
func getIntFromCache(str string) (int, error) {
	if val, exists := intCache[str]; exists {
		return val, nil
	}
	val, err := strconv.Atoi(str)
	if err != nil {
		return 0, err
	}
	intCache[str] = val
	return val, nil
}

// checkLen проверяет длину строки на соответствие указанному значению.
func checkLen(rv reflect.Value, ruleValue string) bool {
	if rv.Kind() == reflect.String {
		intValue, err := getIntFromCache(ruleValue)
		if err != nil {
			return false
		}
		return rv.Len() == intValue
	}
	return false
}

// checkRegex проверяет соответствие строки регулярному выражению.
func checkRegex(rv reflect.Value, ruleValue string) bool {
	if rv.Kind() == reflect.String {
		rx, exists := regexCache[ruleValue]
		if !exists {
			var err error
			rx, err = regexp.Compile(ruleValue)
			if err != nil {
				return false
			}
			regexCache[ruleValue] = rx
		}
		return rx.MatchString(rv.String())
	}
	return false
}

// checkMin проверяет, что числовое значение больше указанного минимума.
func checkMin(rv reflect.Value, ruleValue string) bool {
	if rv.Kind() == reflect.Int {
		intValue := int(rv.Int())
		min, err := getIntFromCache(ruleValue)
		if err != nil {
			return false
		}
		return intValue > min
	}
	return false
}

// checkMax проверяет, что числовое значение меньше указанного максимума.
func checkMax(rv reflect.Value, ruleValue string) bool {
	if rv.Kind() == reflect.Int {
		intValue := int(rv.Int())
		max, err := getIntFromCache(ruleValue)
		if err != nil {
			return false
		}
		return intValue < max
	}
	return false
}

// checkIn проверяет, содержится ли значение в указанном наборе.
func checkIn(rv reflect.Value, ruleValue string) bool {
	ins := strings.Split(ruleValue, ",")
	isValid := false

	//nolint:exhaustive
	switch rv.Kind() {
	case reflect.Int:
		intValue := int(rv.Int())
		for _, in := range ins {
			inInt, err := getIntFromCache(in)
			if err != nil {
				continue
			}
			if inInt == intValue {
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
	}
	return isValid
}
