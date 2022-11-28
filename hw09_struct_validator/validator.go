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
	ErrValidationParse  = errors.New("validation:? parse \"cmd:arg\" failure")
	ErrValidationMin    = errors.New("validation:min failure")
	ErrValidationMax    = errors.New("validation:max failure")
	ErrValidationLen    = errors.New("validation:len failure")
	ErrValidationIn     = errors.New("validation:in failure")
	ErrValidationRegExp = errors.New("validation:regexp failure")
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
				if errors.Is(t.Err, errors.Unwrap(i.Err)) {
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
		str := fmt.Sprintf("Field[%s] %s;", i.Field, i.Err)
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
	if item.rVal.Kind() == reflect.String && item.rVal.String() != "" {
		if matched, err := regexp.MatchString(item.vArg, item.rVal.String()); err != nil {
			return &ValidationError{Field: item.FieldName, Err: err}
		} else if !matched {
			return &ValidationError{
				Field: item.FieldName,
				Err:   fmt.Errorf("%w (%s  ∉ \"%s\")", ErrValidationRegExp, item.rVal.String(), item.vArg),
			}
		}
	} else {
		return &ValidationError{
			Field: item.FieldName,
			Err:   fmt.Errorf("%w (wrong field type:%s)", ErrValidationRegExp, item.rVal.Kind().String()),
		}
	}

	return nil
}

func (item ValidatableItem) In() error {
	in := false
	switch item.rVal.Kind() { //nolint:exhaustive
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if item.rVal.Int() == 0 {
			return nil
		}
		sslice := strings.Split(item.vArg, ",")
		for _, s := range sslice {
			i, err := strconv.ParseInt(s, 10, 64)
			if err != nil {
				return &ValidationError{Field: item.FieldName, Err: err}
			}
			if i == item.rVal.Int() {
				in = true
				break
			}
		}
		if !in {
			return &ValidationError{
				Field: item.FieldName,
				Err:   fmt.Errorf("%w (%d not in[%s])", ErrValidationIn, item.rVal.Int(), item.vArg),
			}
		}
	case reflect.String:
		if item.rVal.String() == "" {
			return nil
		}
		sslice := strings.Split(item.vArg, ",")
		for _, s := range sslice {
			if s == item.rVal.String() {
				in = true
				break
			}
		}
		if !in {
			return &ValidationError{
				Field: item.FieldName,
				Err:   fmt.Errorf("%w (%s not in[%s])", ErrValidationIn, item.rVal.String(), item.vArg),
			}
		}
	default:
		return &ValidationError{
			Field: item.FieldName,
			Err:   fmt.Errorf("%w (wrong field type:%s)", ErrValidationIn, item.rVal.Kind().String()),
		}
	}
	return nil
}

func (item ValidatableItem) Len() *ValidationError {
	switch item.rVal.Kind() { //nolint:exhaustive
	case reflect.Slice:
		slen, err := strconv.Atoi(item.vArg)
		if err != nil {
			return &ValidationError{Field: item.FieldName, Err: err}
		}
		for i := 0; i < item.rVal.Len(); i++ {
			rF := item.rVal.Index(i)
			if reflect.TypeOf(rF.Interface()).Kind() == reflect.String {
				if rF.Len() != slen {
					return &ValidationError{
						Field: item.FieldName,
						Err:   fmt.Errorf("%w (%d != %d)", ErrValidationLen, rF.Len(), slen),
					}
				}
			} else {
				return &ValidationError{
					Field: item.FieldName,
					Err:   fmt.Errorf("%w (wrong field type:%s)", ErrValidationLen, rF.Kind().String()),
				}
			}
		}
	case reflect.String:
		slen, err := strconv.Atoi(item.vArg)
		if err != nil {
			return &ValidationError{Field: item.FieldName, Err: err}
		}
		if item.rVal.Len() != slen {
			return &ValidationError{
				Field: item.FieldName,
				Err:   fmt.Errorf("%w (%d != %d)", ErrValidationLen, item.rVal.Len(), slen),
			}
		}
	default:
		return &ValidationError{
			Field: item.FieldName,
			Err:   fmt.Errorf("%w (wrong field type:%s)", ErrValidationLen, item.rVal.Kind().String()),
		}
	}

	return nil
}

func (item ValidatableItem) Min() error {
	switch item.rVal.Kind() { //nolint:exhaustive
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		min, err := strconv.ParseInt(item.vArg, 10, 64)
		if err != nil {
			return &ValidationError{Field: item.FieldName, Err: err}
		}

		if item.rVal.Int() < min {
			return &ValidationError{
				Field: item.FieldName,
				Err:   fmt.Errorf("%w (%d < %d)", ErrValidationMin, item.rVal.Int(), min),
			}
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		min, err := strconv.ParseUint(item.vArg, 10, 64)
		if err != nil {
			return &ValidationError{Field: item.FieldName, Err: err}
		}

		if item.rVal.Uint() < min {
			return &ValidationError{
				Field: item.FieldName,
				Err:   fmt.Errorf("%w (%d < %d)", ErrValidationMin, item.rVal.Uint(), min),
			}
		}
	case reflect.Float32, reflect.Float64:
		min, err := strconv.ParseFloat(item.vArg, 64)
		if err != nil {
			return &ValidationError{Field: item.FieldName, Err: err}
		}

		// Maybe needs use something to check float-values. But I think this enough for our app
		if item.rVal.Float() < min {
			return &ValidationError{
				Field: item.FieldName,
				Err:   fmt.Errorf("%w (%g < %g)", ErrValidationMin, item.rVal.Float(), min),
			}
		}
	case reflect.Complex64, reflect.Complex128:
		min, err := strconv.ParseComplex(item.vArg, 128)
		if err != nil {
			return &ValidationError{Field: item.FieldName, Err: err}
		}

		if real(item.rVal.Complex()) == real(min) && imag(item.rVal.Complex()) < imag(min) {
			return &ValidationError{
				Field: item.FieldName,
				Err:   fmt.Errorf("%w failure (%v < %v)", ErrValidationMin, item.rVal.Complex(), min),
			}
		}
		if real(item.rVal.Complex()) < real(min) {
			return &ValidationError{
				Field: item.FieldName,
				Err:   fmt.Errorf("%w failure (%v < %v)", ErrValidationMin, item.rVal.Complex(), min),
			}
		}
	default:
		return &ValidationError{
			Field: item.FieldName,
			Err:   fmt.Errorf("%w (wrong field type:%s)", ErrValidationMin, item.rVal.Kind().String()),
		}
	}
	return nil
}

func (item ValidatableItem) Max() error {
	switch item.rVal.Kind() { //nolint:exhaustive
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		max, err := strconv.ParseInt(item.vArg, 10, 64)
		if err != nil {
			return &ValidationError{Field: item.FieldName, Err: err}
		}

		if item.rVal.Int() > max {
			return &ValidationError{
				Field: item.FieldName,
				Err:   fmt.Errorf("%w failure (%d > %d)", ErrValidationMax, item.rVal.Int(), max),
			}
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		max, err := strconv.ParseUint(item.vArg, 10, 64)
		if err != nil {
			return &ValidationError{Field: item.FieldName, Err: err}
		}

		if item.rVal.Uint() > max {
			return &ValidationError{
				Field: item.FieldName,
				Err:   fmt.Errorf("%w (%d < %d)", ErrValidationMax, item.rVal.Uint(), max),
			}
		}
	case reflect.Complex64, reflect.Complex128:
		max, err := strconv.ParseComplex(item.vArg, 128)
		if err != nil {
			return &ValidationError{Field: item.FieldName, Err: err}
		}

		if real(item.rVal.Complex()) == real(max) && imag(item.rVal.Complex()) > imag(max) {
			return &ValidationError{
				Field: item.FieldName,
				Err:   fmt.Errorf("%w failure (%v > %v)", ErrValidationMax, item.rVal.Complex(), max),
			}
		}
		if real(item.rVal.Complex()) > real(max) {
			return &ValidationError{
				Field: item.FieldName,
				Err:   fmt.Errorf("%w failure (%v > %v)", ErrValidationMax, item.rVal.Complex(), max),
			}
		}
	case reflect.Float32, reflect.Float64:
		max, err := strconv.ParseFloat(item.vArg, 64)
		if err != nil {
			return &ValidationError{Field: item.FieldName, Err: err}
		}

		// Maybe needs use something to check float-values. But I think it should be enough for our app
		if item.rVal.Float() > max {
			return &ValidationError{
				Field: item.FieldName,
				Err:   fmt.Errorf("%w (%g < %g)", ErrValidationMax, item.rVal.Float(), max),
			}
		}
	default:
		return &ValidationError{
			Field: item.FieldName,
			Err:   fmt.Errorf("%w (wrong field type:%s)", ErrValidationMax, item.rVal.Kind().String()),
		}
	}

	return nil
}

func doValidate(vItems ValidatableItems, vErr ValidationErrors) ValidationErrors {
	for _, vI := range vItems {
		if _, ok := reflect.TypeOf(vI).MethodByName(vI.vCmd); ok {
			if rVal := reflect.ValueOf(vI).MethodByName(vI.vCmd).Call(nil); len(rVal) == 1 {
				if err, ok := rVal[0].Interface().(*ValidationError); ok && err != nil {
					vErr = append(vErr, *err)
				}
			}
			continue
		}
		vErr = append(vErr, ValidationError{Field: vI.FieldName, Err: fmt.Errorf("%w \"%s\"", ErrValidationNImpl, vI.oCmd)})
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

	//	vItems, vErrs := parseStruct(v)
	//	doValidate(&vItems, &vErrs)

	vErrs := doValidate(parseStruct(v))

	if len(vErrs) == 0 {
		return nil
	}

	return vErrs
}
