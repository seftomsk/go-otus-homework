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
		ID        string `json:"id" validate:"len:36"`
		Name      string
		Age       int      `validate:"min:18|max:50"`
		Email     string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role      UserRole `validate:"in:admin,stuff"`
		Phones    []string `validate:"len:11"`
		CityCodes []int    `validate:"min:2|max:3"`
		meta      json.RawMessage
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

	Duplicate struct {
		Name string `validate:"len:2|len:3"`
	}

	InvalidConstraint struct {
		Name string `validate:"len:2:4"`
	}

	InvalidConstraintArguments struct {
		Name string `validate:"len:sdf"`
	}

	InvalidConstraintRegexp struct {
		Name string `validate:"regexp:\\\\\\\\\\"`
	}

	PrivateFields struct {
		name string `validate:"len:5"`
		age  int    `validate:"min:18|max:50"`
	}

	NotSupportedFields struct {
		Name struct{} `validate:"len:5"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{App{Version: "1234"}, ValidationErrors{ValidationError{"Version", ErrLength}}},
		{App{Version: "12345"}, nil},
		{App{Version: "123456"}, ValidationErrors{ValidationError{Field: "Version", Err: ErrLength}}},

		{Response{Code: 199, Body: ""}, ValidationErrors{ValidationError{Field: "Code", Err: ErrIn}}},
		{Response{Code: 200, Body: ""}, nil},
		{Response{Code: 201, Body: ""}, ValidationErrors{ValidationError{Field: "Code", Err: ErrIn}}},

		{Response{Code: 403, Body: ""}, ValidationErrors{ValidationError{Field: "Code", Err: ErrIn}}},
		{Response{Code: 404, Body: ""}, nil},
		{Response{Code: 405, Body: ""}, ValidationErrors{ValidationError{Field: "Code", Err: ErrIn}}},

		{Response{Code: 499, Body: ""}, ValidationErrors{ValidationError{Field: "Code", Err: ErrIn}}},
		{Response{Code: 500, Body: ""}, nil},
		{Response{Code: 501, Body: ""}, ValidationErrors{ValidationError{Field: "Code", Err: ErrIn}}},

		{User{
			ID:        "",
			Name:      "",
			Age:       0,
			Email:     "",
			Role:      "",
			Phones:    []string{""},
			CityCodes: []int{0},
			meta:      json.RawMessage(""),
		}, ValidationErrors([]ValidationError{
			{Field: "ID", Err: ErrLength},
			{Field: "Age", Err: ErrMin},
			{Field: "Email", Err: ErrRegexp},
			{Field: "Role", Err: ErrIn},
			{Field: "Phones", Err: ErrLength},
			{Field: "CityCodes", Err: ErrMin},
		})},
		{User{
			ID:     "123456789123456789123456789123456789",
			Name:   "",
			Age:    19,
			Email:  "test@test.com",
			Role:   "admin",
			Phones: []string{"12345678912"},
		}, nil},
	}

	for i, tt := range tests {
		tt := tt
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			t.Parallel()

			err := Validate(tt.in)

			if tt.expectedErr == nil {
				require.NoError(t, err)
			} else {
				require.ErrorAs(t, err, &ValidationErrors{})
				require.EqualError(t, err, tt.expectedErr.Error())
			}
		})
	}
}

func TestDuplicateConstraint(t *testing.T) {
	err := Validate(Duplicate{Name: ""})
	require.ErrorAs(t, err, &ParseError{})
}

func TestNotAStruct(t *testing.T) {
	err := Validate(2)
	require.ErrorAs(t, err, &ParseError{})
}

func TestInvalidConstraintSignature(t *testing.T) {
	err := Validate(InvalidConstraint{Name: ""})
	require.ErrorAs(t, err, &ParseError{})
}

func TestInvalidConstraintArguments(t *testing.T) {
	err := Validate(InvalidConstraintArguments{Name: ""})
	require.ErrorAs(t, err, &ParseError{})
}

func TestInvalidConstraintRegexp(t *testing.T) {
	err := Validate(InvalidConstraintRegexp{Name: ""})
	require.ErrorAs(t, err, &ParseError{})
}

func TestPrivateFields(t *testing.T) {
	err := Validate(PrivateFields{name: "", age: 30})
	require.NoError(t, err)
}

func TestNotSupportedFields(t *testing.T) {
	err := Validate(NotSupportedFields{Name: struct{}{}})
	require.NoError(t, err)
	err = Validate(Token{})
	require.NoError(t, err)
}
