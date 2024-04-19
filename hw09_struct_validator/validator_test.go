package hw09structvalidator

import (
	"encoding/json"
	"testing"
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
		name    string
		input   interface{}
		wantErr bool
	}{
		{
			name: "Valid User",
			input: User{
				ID:     "12345678-1234-1234-1234-123456789012",
				Name:   "John Doe",
				Age:    25,
				Email:  "john@example.com",
				Role:   "admin",
				Phones: []string{"12345678901", "10987654321"},
			},
			wantErr: false,
		},
		{
			name: "Invalid User Email",
			input: User{
				ID:     "12345678-1234-1234-1234-123456789012",
				Name:   "John Doe",
				Age:    25,
				Email:  "johnexample.com",
				Role:   "admin",
				Phones: []string{"12345678901", "10987654321"},
			},
			wantErr: true,
		},
		{
			name: "Invalid User Age",
			input: User{
				ID:     "12345678-1234-1234-1234-123456789012",
				Name:   "John Doe",
				Age:    17, // Below minimum age
				Email:  "john@example.com",
				Role:   "admin",
				Phones: []string{"12345678901", "10987654321"},
			},
			wantErr: true,
		},
		{
			name: "Valid App Version",
			input: App{
				Version: "1.0.0",
			},
			wantErr: false,
		},
		{
			name: "Invalid Response Code",
			input: Response{
				Code: 403, // Not in the list of valid codes
				Body: "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Validate(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate(%v) gotErr = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
		})
	}
}
