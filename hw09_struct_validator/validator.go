package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

// ValidationError wraps field-specific errors.
type ValidationError struct {
	Field string
	Err   error
}

// Error implements the error interface for ValidationError.
func (v ValidationError) Error() string {
	return fmt.Sprintf("%s: %v", v.Field, v.Err)
}

// ValidationErrors is a slice of ValidationError, implementing the error interface.
type ValidationErrors []ValidationError

// Error concatenates the error messages of the underlying ValidationError slice.
func (v ValidationErrors) Error() string {
	errMsgs := make([]string, 0, len(v))
	for _, ve := range v {
		errMsgs = append(errMsgs, ve.Error())
	}
	return strings.Join(errMsgs, "; ")
}

var ErrFieldNotValid = errors.New("validation error")

// Validator is an interface that different types of validators will implement.
type Validator interface {
	Validate(field string, tag string, value reflect.Value) (bool, ValidationErrors)
}

// StringValidator checks string values for various constraints.
type StringValidator struct{}

// Validate implements validation checks for strings based on the provided tag.
func (v StringValidator) Validate(field, tag string, value reflect.Value) (bool, ValidationErrors) {
	var errs ValidationErrors
	strVal := value.String()

	for _, part := range strings.Split(tag, "|") {
		tagParts := strings.SplitN(part, ":", 2)
		if len(tagParts) != 2 {
			continue
		}

		key, tagValue := tagParts[0], tagParts[1]
		switch key {
		case "len":
			requiredLen, err := strconv.Atoi(tagValue)
			if err != nil {
				errs = append(errs, ValidationError{field, fmt.Errorf("invalid length value: %w", err)})
				continue
			}
			if len(strVal) != requiredLen {
				errs = append(errs, ValidationError{field, fmt.Errorf("must be %d characters long, "+
					"but got %d: %w", requiredLen, len(strVal), ErrFieldNotValid)})
			}
		case "regexp":
			re, err := regexp.Compile(tagValue)
			if err != nil {
				errs = append(errs, ValidationError{field, fmt.Errorf("invalid regexp: %w", err)})
				continue
			}
			if !re.MatchString(strVal) {
				errs = append(errs, ValidationError{field, fmt.Errorf("must match the pattern %s: %w",
					tagValue, ErrFieldNotValid)})
			}
		case "in":
			allowedValues := strings.Split(tagValue, ",")
			found := false
			for _, allowedValue := range allowedValues {
				if strVal == allowedValue {
					found = true
					break
				}
			}
			if !found {
				errs = append(errs, ValidationError{field, fmt.Errorf("must be one of [%s]: %w",
					tagValue, ErrFieldNotValid)})
			}
		}
	}

	if len(errs) > 0 {
		return false, errs
	}
	return true, nil
}

// NumberValidator checks int values for various constraints.
type NumberValidator struct{}

// Validate checks integer fields for minimum, maximum, and inclusion constraints based on the provided tag
//
//nolint:gocognit
func (v NumberValidator) Validate(field, tag string, value reflect.Value) (bool, ValidationErrors) {
	var errs ValidationErrors
	intVal := int(value.Int())

	for _, part := range strings.Split(tag, "|") {
		tagParts := strings.SplitN(part, ":", 2)
		if len(tagParts) != 2 {
			continue
		}

		key, tagValue := tagParts[0], tagParts[1]
		switch key {
		case "min":
			min, err := strconv.Atoi(tagValue)
			if err != nil {
				errs = append(errs, ValidationError{field, fmt.Errorf("invalid min value: %w", err)})
				continue
			}
			if intVal < min {
				errs = append(errs, ValidationError{field, fmt.Errorf("must be greater than or equal "+
					"to %d: %w", min, ErrFieldNotValid)})
			}
		case "max":
			max, err := strconv.Atoi(tagValue)
			if err != nil {
				errs = append(errs, ValidationError{field, fmt.Errorf("invalid max value: %w", err)})
				continue
			}
			if intVal > max {
				errs = append(errs, ValidationError{field, fmt.Errorf("must be less than or equal "+
					"to %d: %w", max, ErrFieldNotValid)})
			}
		case "in":
			inValuesStr := strings.Split(tagValue, ",")
			inValues := make([]int, len(inValuesStr))
			valid := false
			for i, strVal := range inValuesStr {
				inVal, err := strconv.Atoi(strVal)
				if err != nil {
					errs = append(errs, ValidationError{field, fmt.Errorf("invalid value in 'in' "+
						"constraint: %w", err)})
					continue
				}
				inValues[i] = inVal
				if intVal == inVal {
					valid = true
				}
			}
			if !valid {
				errs = append(errs, ValidationError{field, fmt.Errorf("value must be one of [%s]: %w",
					tagValue, ErrFieldNotValid)})
			}
		}
	}

	if len(errs) > 0 {
		return false, errs
	}
	return true, nil
}

// SliceValidator checks slice values by validating each element in the slice.
type SliceValidator struct{}

// Validate implements validation for each element in a slice based on the provided tag.
func (v SliceValidator) Validate(field string, tag string, value reflect.Value) (bool, ValidationErrors) {
	validationErrors := make(ValidationErrors, 0)

	for i := 0; i < value.Len(); i++ {
		value := value.Index(i)
		isValid, err := validateTag(tag, field, value)
		if !isValid {
			var ve ValidationErrors
			if errors.As(err, &ve) {
				validationErrors = append(validationErrors, ve...)
			}
		}
	}

	if len(validationErrors) > 0 {
		return false, validationErrors
	}

	return true, nil
}

// ValidatorFactory returns a Validator based on the field kind.
//
//nolint:exhaustive
func ValidatorFactory(kind reflect.Kind) (Validator, error) {
	switch kind {
	case reflect.Slice:
		return SliceValidator{}, nil
	case reflect.String:
		return StringValidator{}, nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return NumberValidator{}, nil
	default:
		return nil, fmt.Errorf("unsupported kind: %s", kind)
	}
}

// validateTag dispatches the validation based on the field's kind (type) using a validator factory.
func validateTag(tag, fieldName string, fieldValue reflect.Value) (bool, error) {
	validator, err := ValidatorFactory(fieldValue.Kind())
	if err != nil {
		return false, err
	}

	isValid, validationErrors := validator.Validate(fieldName, tag, fieldValue)
	if !isValid {
		return false, validationErrors
	}

	return true, nil
}

// Validate performs validation on a struct based on its field tags.
func Validate(v interface{}) error {
	var validationErrors ValidationErrors

	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Struct {
		return fmt.Errorf("expected a struct, received %T", v)
	}

	for i := 0; i < rv.NumField(); i++ {
		fieldValue := rv.Field(i)
		fieldType := rv.Type().Field(i)
		tag := fieldType.Tag.Get("validate")

		if tag == "" || tag == "-" {
			continue
		}

		if fieldValue.CanInterface() {
			isValid, err := validateTag(tag, fieldType.Name, fieldValue)
			if !isValid {
				var ve ValidationErrors
				if errors.As(err, &ve) {
					validationErrors = append(validationErrors, ve...)
				}
			}
		}
	}

	if len(validationErrors) > 0 {
		return validationErrors
	}

	return nil
}
