package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

var (
	ErrValidationNImpl  = errors.New("validation: not implemented")
	ErrValidationParse  = errors.New("validation: parse \"cmd:arg\" failure")
	ErrValidationMin    = errors.New("validation: min failure")
	ErrValidationMax    = errors.New("validation: max failure")
	ErrValidationLen    = errors.New("validation: len failure")
	ErrValidationIn     = errors.New("validation: in failure")
	ErrValidationRegExp = errors.New("validation: regexp failure")
	ErrVarlidationType  = errors.New("validation: wrong type")
)

type ValidationError struct {
	Field string
	Err   error
}

func (v ValidationError) Error() string {
	return fmt.Sprintf("Field[%s] %s", v.Field, v.Err)
}

type ValidationErrors []ValidationError

func (v ValidationErrors) As(target interface{}) bool {
	if len(v) == 0 || target == nil {
		return false
	}

	if t, ok := target.(*ValidationError); ok {
		for _, i := range v {
			if t.Field == i.Field {
				if errors.Is(t.Err, i.Err) {
					return true
				}
			}
		}
	}

	return false
}

func (v ValidationErrors) Is(target error) bool {
	if len(v) == 0 && target == nil {
		return true
	}

	for _, i := range v {
		if errors.Is(i.Err, target) {
			return true
		}
	}

	return false
}

func (v ValidationErrors) Error() string {
	if len(v) == 0 {
		return ""
	}

	err := strings.Builder{}
	for _, i := range v {
		str := i.Error()
		err.WriteString(str)
	}

	return err.String()
}

type ValidatableItem struct {
	oCmd      string
	vCmd      string
	vArg      string
	rVal      reflect.Value
	FieldName string
}

type ValidatableItems []ValidatableItem

func (item ValidatableItem) Regexp() error {
	if item.rVal.Kind() != reflect.String {
		return ErrVarlidationType
	}

	if matched, err := regexp.MatchString(item.vArg, item.rVal.String()); err != nil {
		return err
	} else if !matched {
		return ErrValidationRegExp
	}

	return nil
}

func (item ValidatableItem) In() error {
	switch item.rVal.Kind() { //nolint:exhaustive
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		sslice := strings.Split(item.vArg, ",")
		for _, s := range sslice {
			i, err := strconv.ParseInt(s, 10, 64)
			if err != nil {
				return err
			}
			if i == item.rVal.Int() {
				return nil
			}
		}
	case reflect.String:
		sslice := strings.Split(item.vArg, ",")
		for _, s := range sslice {
			if s == item.rVal.String() {
				return nil
			}
		}
	default:
		return ErrVarlidationType
	}

	return ErrValidationIn
}

func (item ValidatableItem) Len() error {
	switch item.rVal.Kind() { //nolint:exhaustive
	case reflect.Slice:
		slen, err := strconv.Atoi(item.vArg)
		if err != nil {
			return err
		}
		for i := 0; i < item.rVal.Len(); i++ {
			rF := item.rVal.Index(i)
			if reflect.TypeOf(rF.Interface()).Kind() != reflect.String {
				return ErrVarlidationType
			}

			if rF.Len() != slen {
				return ErrValidationLen
			}
		}
	case reflect.String:
		slen, err := strconv.Atoi(item.vArg)
		if err != nil {
			return err
		}
		if item.rVal.Len() != slen {
			return ErrValidationLen
		}
	default:
		return ErrVarlidationType
	}

	return nil
}

func (item ValidatableItem) Min() error {
	switch item.rVal.Kind() { //nolint:exhaustive
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		min, err := strconv.ParseInt(item.vArg, 10, 64)
		if err != nil {
			return err
		}

		if item.rVal.Int() < min {
			return ErrValidationMin
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		min, err := strconv.ParseUint(item.vArg, 10, 64)
		if err != nil {
			return err
		}

		if item.rVal.Uint() < min {
			return ErrValidationMin
		}
	case reflect.Float32, reflect.Float64:
		min, err := strconv.ParseFloat(item.vArg, 64)
		if err != nil {
			return err
		}

		// Maybe needs use something to check float-values. But I think this enough for our app
		if item.rVal.Float() < min {
			return ErrValidationMin
		}
	case reflect.Complex64, reflect.Complex128:
		min, err := strconv.ParseComplex(item.vArg, 128)
		if err != nil {
			return err
		}

		if real(item.rVal.Complex()) == real(min) && imag(item.rVal.Complex()) < imag(min) {
			return ErrValidationMin
		}
		if real(item.rVal.Complex()) < real(min) {
			return ErrValidationMin
		}
	default:
		return ErrVarlidationType
	}
	return nil
}

func (item ValidatableItem) Max() error {
	switch item.rVal.Kind() { //nolint:exhaustive
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		max, err := strconv.ParseInt(item.vArg, 10, 64)
		if err != nil {
			return err
		}

		if item.rVal.Int() > max {
			return ErrValidationMax
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		max, err := strconv.ParseUint(item.vArg, 10, 64)
		if err != nil {
			return err
		}

		if item.rVal.Uint() > max {
			return ErrValidationMax
		}
	case reflect.Complex64, reflect.Complex128:
		max, err := strconv.ParseComplex(item.vArg, 128)
		if err != nil {
			return err
		}

		if real(item.rVal.Complex()) == real(max) && imag(item.rVal.Complex()) > imag(max) {
			return ErrValidationMax
		}
		if real(item.rVal.Complex()) > real(max) {
			return ErrValidationMax
		}
	case reflect.Float32, reflect.Float64:
		max, err := strconv.ParseFloat(item.vArg, 64)
		if err != nil {
			return err
		}

		// Maybe needs use something to check float-values. But I think it should be enough for our app
		if item.rVal.Float() > max {
			return ErrValidationMax
		}
	default:
		return ErrVarlidationType
	}

	return nil
}

func doValidate(vItems ValidatableItems, vErr ValidationErrors) error {
	for _, vI := range vItems {
		if _, ok := reflect.TypeOf(vI).MethodByName(vI.vCmd); ok {
			if rVal := reflect.ValueOf(vI).MethodByName(vI.vCmd).Call(nil); len(rVal) == 1 {
				if err, ok := rVal[0].Interface().(error); ok && err != nil {
					vErr = append(vErr, ValidationError{
						Field: vI.FieldName,
						Err:   err,
					})
				}
			}
			continue
		}
		vErr = append(vErr, ValidationError{Field: vI.FieldName, Err: ErrValidationNImpl})
	}

	if len(vErr) == 0 {
		return nil
	}

	return vErr
}

func parseStruct(v interface{}) (ValidatableItems, ValidationErrors) {
	vItem := ValidatableItems{}
	vErr := ValidationErrors{}

	rValue := reflect.ValueOf(v)
	rType := reflect.TypeOf(v)

	for i := 0; i < rType.NumField(); i++ {
		fType := rType.Field(i)
		fValue := rValue.Field(i)

		if fType.Type.Kind() == reflect.Struct {
			vI, vE := parseStruct(fValue.Interface())
			vItem = append(vItem, vI...)
			vErr = append(vErr, vE...)
		}

		vtag := fType.Tag.Get("validate")
		if vtag != "" {
			sval := strings.Split(vtag, "|")
			for _, s := range sval {
				lines := strings.SplitN(s, ":", 2)
				if len(lines) != 2 {
					vErr = append(vErr, ValidationError{Field: fType.Name, Err: ErrValidationParse})
					continue
				}

				// Makes cmd-name with capitalize the first letter
				r := []rune(lines[0])
				cmdStr := string(append([]rune{unicode.ToUpper(r[0])}, r[1:]...))
				vItem = append(vItem, ValidatableItem{
					oCmd:      lines[0],
					vCmd:      cmdStr,
					vArg:      lines[1],
					rVal:      fValue,
					FieldName: fType.Name,
				})
			}
		}
	}

	return vItem, vErr
}

func Validate(v interface{}) error {
	if reflect.TypeOf(v).Kind() != reflect.Struct {
		return nil
	}

	return doValidate(parseStruct(v))
}
