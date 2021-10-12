package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	strBuilder := strings.Builder{}
	for _, err := range v {
		strBuilder.WriteString("Field \"" + err.Field + "\": " + err.Err.Error() + "\n")
	}
	return strBuilder.String()
}

type ParseError struct {
	Msg string
	Err error
}

func (e ParseError) Error() string {
	return e.Msg
}

var (
	ErrLength = errors.New("invalid length")
	ErrMin    = errors.New("number is less than min value")
	ErrMax    = errors.New("number is more than max value")
	ErrIn     = errors.New("value is not included in the set")
	ErrRegexp = errors.New("value is not matched with regexp")
)

var supportedTypesOfStruct = make(map[reflect.Kind]struct{})

func init() {
	supportedTypesOfStruct[reflect.String] = struct{}{}
	supportedTypesOfStruct[reflect.Int] = struct{}{}
	supportedTypesOfStruct[reflect.Slice] = struct{}{}
}

type validator struct {
	errors             ValidationErrors
	currentConstraints map[string]struct{}
	currentStructField reflect.StructField
}

func (v *validator) splitConstraint(list string) (string, string, error) {
	constraintWithValue := strings.Split(list, ":")
	if len(constraintWithValue) != 2 {
		return "", "", ParseError{
			Msg: fmt.Sprintf(
				"field %q - invalid constraint, expected [name:value(s)], but received %v",
				v.currentStructField.Name,
				constraintWithValue,
			),
		}
	}
	constraint := constraintWithValue[0]

	if _, ok := v.currentConstraints[constraint]; ok {
		return "", "", ParseError{
			Msg: fmt.Sprintf(
				"field %q - duplicate constraint %q",
				v.currentStructField.Name,
				constraint,
			),
		}
	}
	return constraint, constraintWithValue[1], nil
}

func (v *validator) validateString(val string, constraintLists []string) error {
	for _, list := range constraintLists {
		constraint, constraintV, err := v.splitConstraint(
			list,
		)
		if err != nil {
			return err
		}

		v.currentConstraints[constraint] = struct{}{}
		var ok bool

		switch constraint {
		case "len":
			ok, err = v.strLen(val, constraint, constraintV) //nolint:staticcheck
			if err != nil {
				return err
			}
			if !ok {
				break
			}
		case "in":
			ok = v.strIn(val, constraintV) //nolint:staticcheck
			if !ok {
				break
			}
		case "regexp":
			ok, err = v.strRegexp(val, constraint, constraintV) //nolint:staticcheck
			if err != nil {
				return err
			}
			if !ok {
				break
			}
		}
	}

	return nil
}

func (v *validator) strLen(val, constraint, constraintV string) (bool, error) {
	num, err := strconv.Atoi(constraintV)
	if err != nil {
		return false, ParseError{
			Msg: fmt.Sprintf(
				"field %q - expected a number value for %q constraint, but received %v",
				v.currentStructField.Name,
				constraint,
				constraintV,
			),
			Err: err,
		}
	}

	if len(val) != num {
		vError := ValidationError{
			Field: v.currentStructField.Name,
			Err:   fmt.Errorf("%w", ErrLength),
		}
		v.errors = append(v.errors, vError)
		return false, nil
	}

	return true, nil
}

func (v *validator) strIn(val, constraintV string) bool {
	values := strings.Split(constraintV, ",")
	in := false
	for _, item := range values {
		if item == val {
			in = true
		}
	}
	if !in {
		vError := ValidationError{
			Field: v.currentStructField.Name,
			Err:   fmt.Errorf("%w", ErrIn),
		}
		v.errors = append(v.errors, vError)
		return false
	}

	return true
}

func (v *validator) strRegexp(val, constraint, constraintV string) (bool, error) {
	reg, err := regexp.Compile(constraintV)
	if err != nil {
		return false, ParseError{
			Msg: fmt.Sprintf(
				"field %q - invalid regexp for %q constraint",
				v.currentStructField.Name,
				constraint,
			),
			Err: err,
		}
	}
	if !reg.MatchString(val) {
		vError := ValidationError{
			Field: v.currentStructField.Name,
			Err:   fmt.Errorf("%w", ErrRegexp),
		}
		v.errors = append(v.errors, vError)
		return false, nil
	}

	return true, nil
}

func (v *validator) validateInt(val int64, constraintLists []string) error {
	for _, list := range constraintLists {
		constraint, constraintV, err := v.splitConstraint(
			list,
		)
		if err != nil {
			return err
		}

		v.currentConstraints[constraint] = struct{}{}
		var ok bool

		switch constraint {
		case "min":
			ok, err = v.intMin(val, constraint, constraintV) //nolint:staticcheck
			if err != nil {
				return err
			}
			if !ok {
				break
			}
		case "max":
			ok, err = v.intMax(val, constraint, constraintV) //nolint:staticcheck
			if err != nil {
				return err
			}
			if !ok {
				break
			}
		case "in":
			ok, err = v.intIn(val, constraint, constraintV) //nolint:staticcheck
			if err != nil {
				return err
			}
			if !ok {
				break
			}
		}
	}

	return nil
}

func (v *validator) intMin(val int64, constraint, constraintV string) (bool, error) {
	num, err := strconv.Atoi(constraintV)
	if err != nil {
		return false, ParseError{
			Msg: fmt.Sprintf(
				"field %q - expected a number value for %q constraint, but received %v",
				v.currentStructField.Name,
				constraint,
				constraintV,
			),
			Err: err,
		}
	}

	if val < int64(num) {
		vError := ValidationError{
			Field: v.currentStructField.Name,
			Err:   fmt.Errorf("%w", ErrMin),
		}
		v.errors = append(v.errors, vError)
		return false, nil
	}

	return true, nil
}

func (v *validator) intMax(val int64, constraint, constraintV string) (bool, error) {
	num, err := strconv.Atoi(constraintV)
	if err != nil {
		return false, ParseError{
			Msg: fmt.Sprintf(
				"field %q - expected a number value for %q constraint, but received %v",
				v.currentStructField.Name,
				constraint,
				constraintV,
			),
			Err: err,
		}
	}

	if val > int64(num) {
		vError := ValidationError{
			Field: v.currentStructField.Name,
			Err:   fmt.Errorf("%w", ErrMax),
		}
		v.errors = append(v.errors, vError)
		return false, nil
	}

	return true, nil
}

func (v *validator) intIn(val int64, constraint, constraintV string) (bool, error) {
	values := strings.Split(constraintV, ",")
	in := false
	var num int
	var err error
	for _, item := range values {
		num, err = strconv.Atoi(item)
		if err != nil {
			return false, ParseError{
				Msg: fmt.Sprintf(
					"field %q - expected a number value for %q constraint, but received %v",
					v.currentStructField.Name,
					constraint,
					constraintV,
				),
				Err: err,
			}
		}
		if val == int64(num) {
			in = true
		}
	}

	if !in {
		vError := ValidationError{
			Field: v.currentStructField.Name,
			Err:   fmt.Errorf("%w", ErrIn),
		}
		v.errors = append(v.errors, vError)
		return false, nil
	}

	return true, nil
}

func Validate(v interface{}) error { //nolint:gocognit
	worker := validator{
		errors: make(ValidationErrors, 0),
	}
	rValue := reflect.ValueOf(v)
	if rValue.Kind() != reflect.Struct {
		return ParseError{
			Msg: fmt.Sprintf("expected a struct, but received %T", v),
		}
	}
	rType := rValue.Type()

	for i := 0; i < rType.NumField(); i++ {
		fieldValue := rValue.Field(i)
		// Only public field
		if !fieldValue.CanInterface() {
			continue
		}
		// Is support field?
		if _, ok := supportedTypesOfStruct[fieldValue.Kind()]; !ok {
			continue
		}
		field := rType.Field(i)
		tags := field.Tag
		// Are tags empty, and they have "validate" field?
		if tags == "" {
			continue
		}
		validate, ok := tags.Lookup("validate")
		if !ok {
			continue
		}

		worker.currentConstraints = make(map[string]struct{})
		worker.currentStructField = field
		constraintLists := strings.Split(validate, "|")

		switch fieldValue.Kind() {
		case reflect.String:
			if err := worker.validateString(fieldValue.String(), constraintLists); err != nil {
				return err
			}
		case reflect.Int:
			if err := worker.validateInt(fieldValue.Int(), constraintLists); err != nil {
				return err
			}
		case reflect.Slice:
			switch fieldValue.Interface().(type) {
			case []int:
				sliceValues := fieldValue.Interface().([]int)
				for _, val := range sliceValues {
					if err := worker.validateInt(int64(val), constraintLists); err != nil {
						return err
					}
				}
			case []string:
				sliceValues := fieldValue.Interface().([]string)
				for _, val := range sliceValues {
					if err := worker.validateString(val, constraintLists); err != nil {
						return err
					}
				}
			}
		}
	}
	if len(worker.errors) != 0 {
		return worker.errors
	}

	return nil
}
