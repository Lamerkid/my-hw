package hw09structvalidator

import (
	"encoding/json"
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
		Age    int             `validate:"min:18|max:50"`
		Email  string          `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole        `validate:"in:admin,stuff"`
		Phones []string        `validate:"len:11"`
		meta   json.RawMessage //nolint:unused
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
		in          interface{}
		expectedErr error
	}{
		{
			in: User{
				ID:     "1",
				Name:   "2",
				Age:    3,
				Email:  "4",
				Role:   "5",
				Phones: []string{"6", "7"},
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "ID",
					Err:   ErrValidateStringLen,
				},
				ValidationError{
					Field: "Age",
					Err:   ErrValidateIntMin,
				},
				ValidationError{
					Field: "Email",
					Err:   ErrValidateStringRegexp,
				},
				ValidationError{
					Field: "Role",
					Err:   ErrValidateStringIn,
				},
				ValidationError{
					Field: "Phones",
					Err:   ErrValidateStringLen,
				},
			},
		},
		{
			in: User{
				ID:     "466c223c-d3fb-4297-81f9-14bc10049ca8",
				Name:   "name",
				Age:    33,
				Email:  "abc@cba.xyz",
				Role:   "admin",
				Phones: []string{"01234567891", "19876543210"},
			},
			expectedErr: nil,
		},
		{
			in: App{
				Version: "newerThanYestarday",
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Version",
					Err:   ErrValidateStringLen,
				},
			},
		},
		{
			in: App{
				Version: "12345",
			},
			expectedErr: nil,
		},
		{
			in: Token{
				Header:  nil,
				Payload: nil,
			},
			expectedErr: nil,
		},
		{
			in: Token{
				Header:    []byte("application: text/json"),
				Payload:   []byte("Somedata"),
				Signature: []byte("signature"),
			},
			expectedErr: nil,
		},
		{
			in: Response{
				Code: 429,
				Body: "",
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Code",
					Err:   ErrValidateIntIn,
				},
			},
		},
		{
			in: Response{
				Code: 200,
			},
			expectedErr: nil,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(tt.in)
			require.Equal(t, tt.expectedErr, err)
		})
	}
}
