package hw09structvalidator

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
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
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		//{App{Version: "1234"}, ErrLength},
		//{App{Version: "12345"}, nil},
		{App{Version: "123456"}, ValidationErrors{ValidationError{Field: "Version", Err: ErrLength}}},
	}

	for i, tt := range tests {
		tt := tt
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			t.Parallel()

			err := Validate(tt.in)
			//if tt.expectedErr == nil {
			//	require.NoError(t, err)
			//} else {
			//	require.EqualError(t, err, tt.expectedErr.Error())
			//}

			var vErrors ValidationErrors
			require.ErrorAs(t, err, &vErrors)
			if errors.As(err, &vErrors) {
				for _, r := range vErrors {
					require.ErrorIs(t, r.Err, tt.expectedErr)
				}
			}
		})
	}
}
