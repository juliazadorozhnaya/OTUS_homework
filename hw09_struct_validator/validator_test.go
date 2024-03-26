package hw09structvalidator

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type UserRole string

type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int      `validate:"min:18|max:50"`
		Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole `validate:"in:admin,stuff"`
		Phones []string `validate:"len:11"`
		Meta   json.RawMessage
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
	tests := []struct {
		in              interface{}
		expectedValErrs []error
	}{
		{
			in: User{
				ID:     "b75ece0b-85c1-4afb-b4bf-bb4512bf0fa8", // Корректная длина ID
				Name:   "testName",
				Age:    25,
				Email:  "mail@test.com",
				Role:   "admin",
				Phones: []string{"12345678901", "10987654321"},
				Meta:   nil,
			},
			expectedValErrs: []error{ErrorLength, ErrorMin, ErrorLength, ErrorLength},
		},
		{
			in: App{
				Version: "123456",
			},
			expectedValErrs: []error{ErrorLength},
		},
		{
			in:              Token{},
			expectedValErrs: []error{nil},
		},
		{
			in: Response{
				Code: 200,
				Body: "test",
			},
			expectedValErrs: []error{nil},
		},
	}
	for i, tc := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tc := tc
			t.Parallel()
			err := Validate(tc.in)
			var valErrs ValidationErrors
			if errors.As(err, &valErrs) {
				for i, err := range tc.expectedValErrs {
					require.ErrorIs(t, valErrs[i].Err, err, "Validation error should be like expected")
				}
			}
		})
	}

	t.Run("Handle non struct value", func(t *testing.T) {
		err := Validate(123)
		require.ErrorIs(t, err, ErrorExpectedStruct, "Throw nonStruct error")
	})

	t.Run("Handle nil value", func(t *testing.T) {
		err := Validate(nil)
		require.ErrorIs(t, err, ErrorExpectedStruct, "Throw nonStruct error")
	})
}
