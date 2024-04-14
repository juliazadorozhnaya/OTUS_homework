package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

// Предопределенные ошибки для различных типов нарушений валидации.
var (
	ErrorLength         = errors.New("length")
	ErrorRegex          = errors.New("regex")
	ErrorMin            = errors.New("greater")
	ErrorMax            = errors.New("less")
	ErrorIn             = errors.New("lots of")
	ErrorExpectedStruct = errors.New("expected a struct")

	// Кэш для регулярных выражений и целочисленных значений.
	regexCache = make(map[string]*regexp.Regexp)
	intCache   = make(map[string]int)
)

// ValidationError описывает ошибку валидации для конкретного поля.
type ValidationError struct {
	Field string
	Err   error
}

// ValidationErrors представляет собой список ошибок валидации.
type ValidationErrors []ValidationError

// Error реализует интерфейс ошибки для ValidationErrors.
func (v ValidationErrors) Error() string {
	var builder strings.Builder
	for _, err := range v {
		builder.WriteString("Field: " + err.Field + ", error: " + err.Err.Error() + "\n")
	}
	return builder.String()
}

// Validate выполняет валидацию структуры на основе тегов `validate`.
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
		validationErrors = checkValue(validationErrors, field.Name, validateTag, value.Field(i))
	}
	if len(validationErrors) > 0 {
		return validationErrors
	}
	return nil
}

// checkValue обрабатывает значение поля и выполняет валидацию согласно правилам.
func checkValue(valErrs ValidationErrors, fName string, validateTag string, rv reflect.Value) ValidationErrors {
	var (
		errs       []error
		newValErrs = valErrs
	)

	//nolint:exhaustive
	switch rv.Kind() {
	case reflect.String, reflect.Int:
		errs = validateValue(validateTag, rv)
	case reflect.Slice:
		for i := 0; i < rv.Len(); i++ {
			newValErrs = checkValue(newValErrs, fName, validateTag, rv.Index(i))
		}
	}
	if len(errs) > 0 {
		for _, err := range errs {
			newValErrs = append(newValErrs, ValidationError{fName, err})
		}
	}
	return newValErrs
}

// validateValue валидирует значение на основе правил, указанных в теге validate.
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
