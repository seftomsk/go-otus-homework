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
		strBuilder.WriteString(fmt.Sprintf("Field: %s: %s\n", err.Field, err.Err))
	}
	return strBuilder.String()
}

type ParseError struct {
	Msg string
	Err error
}

var (
	ErrLength = errors.New("invalid length")
	ErrMin    = errors.New("value less than min")
	ErrMax    = errors.New("value more than max")
	ErrIn     = errors.New("value is not included in the set")
	ErrRegexp = errors.New("value is not match with regexp")
)

func (e ParseError) Error() string {
	return e.Msg
}

var supportedTypesOfStruct = make(map[reflect.Kind]struct{})
var vErrors ValidationErrors

func init() {
	supportedTypesOfStruct[reflect.String] = struct{}{}
	supportedTypesOfStruct[reflect.Int] = struct{}{}
	supportedTypesOfStruct[reflect.Slice] = struct{}{}
}

func splitConstraint(list string, flag map[string]struct{}, field reflect.StructField) (string, string, error) {
	constraintWithValue := strings.Split(list, ":")
	if len(constraintWithValue) != 2 {
		return "", "", ParseError{
			Msg: fmt.Sprintf("invalid constraint for field %q, expected [name:value(s)], but received %v", field.Name, constraintWithValue),
		}
	}
	constraint := constraintWithValue[0]
	constraintV := constraintWithValue[1]

	if _, ok := flag[constraint]; ok {
		return "", "", ParseError{
			Msg: fmt.Sprintf("duplicate constraint %q for field %q", constraint, field.Name),
		}
	}
	return constraint, constraintV, nil
}

func lenValidator(fv, constraint, constraintV string, field reflect.StructField) (bool, error) {
	num, err := strconv.Atoi(constraintV)
	if err != nil {
		return false, ParseError{
			Msg: fmt.Sprintf("expected a number value for %q constraint for field %q, but received %v", constraint, field.Name, constraintV),
			Err: err,
		}
	}
	strLen := len(fv)
	if strLen != num {
		vError := ValidationError{
			Field: field.Name,
			Err:   fmt.Errorf("%w", ErrLength),
		}
		vErrors = append(vErrors, vError)
		return false, nil
	}
	return true, nil
}

func inStringValidator(fv, constraintV string, field reflect.StructField) bool {
	values := strings.Split(constraintV, ",")
	in := false
	for _, v := range values {
		if v == fv {
			in = true
		}
	}
	if !in {
		vError := ValidationError{
			Field: field.Name,
			Err:   fmt.Errorf("values are not included in the set %v: %w", values, ErrIn),
		}
		vErrors = append(vErrors, vError)
		return false
	}
	return true
}

func regexpValidator(fv, constraint, constraintV string, field reflect.StructField) (bool, error) {
	reg, err := regexp.Compile(constraintV)
	if err != nil {
		return false, ParseError{
			Msg: fmt.Sprintf("invalid regexp for %q constraint for field %q", constraint, field.Name),
			Err: err,
		}
	}
	if !reg.MatchString(fv) {
		vError := ValidationError{
			Field: field.Name,
			Err:   fmt.Errorf("values do not match to %v: %w", constraintV, ErrRegexp),
		}
		vErrors = append(vErrors, vError)
		return false, nil
	}
	return true, nil
}

func minValidator(fv int64, constraint, constraintV string, field reflect.StructField) (bool, error) {
	num, err := strconv.Atoi(constraintV)
	if err != nil {
		return false, ParseError{
			Msg: fmt.Sprintf("expected a number value for %q constraint for field %q, but received %v", constraint, field.Name, constraintV),
			Err: err,
		}
	}
	if fv < int64(num) {
		vError := ValidationError{
			Field: field.Name,
			Err:   fmt.Errorf("numbers must be more than %v: %w", num, ErrMin),
		}
		vErrors = append(vErrors, vError)
		return false, nil
	}
	return true, nil
}

func maxValidator(fv int64, constraint, constraintV string, field reflect.StructField) (bool, error) {
	num, err := strconv.Atoi(constraintV)
	if err != nil {
		return false, ParseError{
			Msg: fmt.Sprintf("expected a number value for %q constraint for field %q, but received %v", constraint, field.Name, constraintV),
			Err: err,
		}
	}
	if fv > int64(num) {
		vError := ValidationError{
			Field: field.Name,
			Err:   fmt.Errorf("numbers must be less than %v: %w", num, ErrMax),
		}
		vErrors = append(vErrors, vError)
		return false, nil
	}
	return true, nil
}

func inIntValidator(fv int64, constraint, constraintV string, field reflect.StructField) (bool, error) {
	values := strings.Split(constraintV, ",")
	in := false
	var num int
	var err error
	for _, v := range values {
		num, err = strconv.Atoi(v)
		if err != nil {
			return false, ParseError{
				Msg: fmt.Sprintf("expected a number value for %q constraint for field %q, but received %v", constraint, field.Name, constraintV),
				Err: err,
			}
		}
		if fv == int64(num) {
			in = true
		}
	}
	if !in {
		vError := ValidationError{
			Field: field.Name,
			Err:   fmt.Errorf("numbers are not included in the set %v: %w", values, ErrIn),
		}
		vErrors = append(vErrors, vError)
		return false, nil
	}
	return true, nil
}

func Validate(v interface{}) error {
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

		constraintFlag := make(map[string]struct{})
		constraintLists := strings.Split(validate, "|")

		switch fieldValue.Kind() {
		case reflect.String:
			fv := fieldValue.String()
			for _, list := range constraintLists {
				constraint, constraintV, err := splitConstraint(list, constraintFlag, field)
				if err != nil {
					return err
				}
				constraintFlag[constraint] = struct{}{}
				switch constraint {
				case "len":
					_, err = lenValidator(fv, constraint, constraintV, field)
					if err != nil {
						return err
					}
				case "in":
					inStringValidator(fv, constraintV, field)
				case "regexp":
					_, err = regexpValidator(fv, constraint, constraintV, field)
					if err != nil {
						return err
					}
				}
			}
		case reflect.Int:
			fv := fieldValue.Int()
			for _, list := range constraintLists {
				constraint, constraintV, err := splitConstraint(list, constraintFlag, field)
				if err != nil {
					return err
				}
				constraintFlag[constraint] = struct{}{}
				switch constraint {
				case "min":
					_, err = minValidator(fv, constraint, constraintV, field)
					if err != nil {
						return err
					}
				case "max":
					_, err = maxValidator(fv, constraint, constraintV, field)
					if err != nil {
						return err
					}
				case "in":
					_, err = inIntValidator(fv, constraint, constraintV, field)
					if err != nil {
						return err
					}
				}
			}
		case reflect.Slice:
			switch fieldValue.Interface().(type) {
			case []int:
				sliceFv := fieldValue.Interface().([]int)
				for _, list := range constraintLists {
					constraint, constraintV, err := splitConstraint(list, constraintFlag, field)
					if err != nil {
						return err
					}
					constraintFlag[constraint] = struct{}{}
					switch constraint {
					case "min":
						for _, fv := range sliceFv {
							ok, err = minValidator(int64(fv), constraint, constraintV, field)
							if err != nil {
								return err
							}
							if !ok {
								break
							}
						}
					case "max":
						for _, fv := range sliceFv {
							ok, err = maxValidator(int64(fv), constraint, constraintV, field)
							if err != nil {
								return err
							}
							if !ok {
								break
							}
						}
					case "in":
						for _, fv := range sliceFv {
							ok, err = inIntValidator(int64(fv), constraint, constraintV, field)
							if err != nil {
								return err
							}
							if !ok {
								break
							}
						}
					}
				}
			case []string:
				sliceFv := fieldValue.Interface().([]string)
				for _, list := range constraintLists {
					constraint, constraintV, err := splitConstraint(list, constraintFlag, field)
					if err != nil {
						return err
					}
					constraintFlag[constraint] = struct{}{}
					switch constraint {
					case "len":
						for _, fv := range sliceFv {
							ok, err = lenValidator(fv, constraint, constraintV, field)
							if err != nil {
								return err
							}
							if !ok {
								break
							}
						}
					case "in":
						for _, fv := range sliceFv {
							ok = inStringValidator(fv, constraintV, field)
							if err != nil {
								return err
							}
							if !ok {
								break
							}
						}
					case "regexp":
						for _, fv := range sliceFv {
							ok, err = regexpValidator(fv, constraint, constraintV, field)
							if err != nil {
								return err
							}
							if !ok {
								break
							}
						}
					}
				}
			}
		}
	}
	if len(vErrors) != 0 {
		return vErrors
	}
	return nil
}
