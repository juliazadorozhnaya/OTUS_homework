package hw09structvalidator

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int      `validate:"min:18|max:50"`
		Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole `validate:"in:admin,stuff"`
		Phones []string `validate:"len:11"`
		meta   json.RawMessage
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}
)

func TestValidate(t *testing.T) {
	testsWithValidationErrors := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			in: App{
				Version: "some long name",
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Version",
					Err:   ErrFieldNotValid,
				},
			},
		},
		{
			in: Response{
				Code: 501,
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Code",
					Err:   ErrFieldNotValid,
				},
			},
		},
		{
			in: User{
				ID:     "01d0c798-f17d-481d-8a69-869870b5008b",
				Age:    16,
				Email:  "wrong_email",
				Role:   "adminn",
				Phones: []string{"891608307777"},
				meta:   json.RawMessage{},
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Age",
					Err:   ErrFieldNotValid,
				},
				ValidationError{
					Field: "Email",
					Err:   ErrFieldNotValid,
				},
				ValidationError{
					Field: "Role",
					Err:   ErrFieldNotValid,
				},
				ValidationError{
					Field: "Phones",
					Err:   ErrFieldNotValid,
				},
			},
		},
	}

	for i, tt := range testsWithValidationErrors {
		t.Run(fmt.Sprintf("testsWithValidationErrors: case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(tt.in)

			var validationErrors ValidationErrors

			require.ErrorAs(t, err, &validationErrors)

			var expectedValidationErrors ValidationErrors
			if errors.As(tt.expectedErr, &expectedValidationErrors) {
				for i, validationError := range validationErrors {
					require.Equal(t, validationError.Field, expectedValidationErrors[i].Field)
					require.ErrorIs(t, validationError.Err, expectedValidationErrors[i].Err)
				}
			}
		})
	}

	testsWithoutValidationErrors := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			in: App{
				Version: "name_",
			},
			expectedErr: nil,
		},
		{
			in: Response{
				Code: 500,
			},
			expectedErr: nil,
		},
		{
			in:          Token{},
			expectedErr: nil,
		},
		{
			in: User{
				ID:     "01d0c798-f17d-481d-8a69-869870b5008b",
				Age:    20,
				Email:  "some@email.com",
				Role:   "admin",
				Phones: []string{"89160830777"},
				meta:   json.RawMessage{},
			},
			expectedErr: nil,
		},
	}

	for i, tt := range testsWithoutValidationErrors {
		t.Run(fmt.Sprintf("testsWithoutValidationErrors: case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(tt.in)

			require.True(t, err == nil)
		})
	}
}
