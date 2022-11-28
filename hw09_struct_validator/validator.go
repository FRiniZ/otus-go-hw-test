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
	err := strings.Builder{}

	if len(v) == 0 {
		return err.String()
	}

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

func (item ValidatableItem) Regexp(vp *ValidationErrors) {
	switch item.rVal.Kind() {
	case reflect.String:
		if item.rVal.String() != "" {
			fmt.Printf("Checking %v %v\n", item.vArg, item.rVal.String())
			matched, err := regexp.MatchString(item.vArg, item.rVal.String())
			if err != nil {
				*vp = append(*vp, ValidationError{Field: item.FieldName, Err: err})
				break
			}

			if !matched {
				*vp = append(*vp, ValidationError{
					Field: item.FieldName,
					Err:   fmt.Errorf("%w (%s  ∉ \"%s\")", ErrValidationRegExp, item.rVal.String(), item.vArg),
				})
			}
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Array, reflect.Bool, reflect.Chan, reflect.Complex128, reflect.Complex64,
		reflect.Float32, reflect.Float64, reflect.Func, reflect.Interface, reflect.Invalid, reflect.Slice,
		reflect.Map, reflect.Ptr, reflect.Struct, reflect.Uint, reflect.Uint16, reflect.Uint32,
		reflect.Uint64, reflect.Uint8, reflect.Uintptr, reflect.UnsafePointer:
		*vp = append(*vp, ValidationError{
			Field: item.FieldName,
			Err:   fmt.Errorf("%w (wrong field type:%s)", ErrValidationRegExp, item.rVal.Kind().String()),
		})
	default:
		*vp = append(*vp, ValidationError{
			Field: item.FieldName,
			Err:   fmt.Errorf("%w (wrong field type:%s)", ErrValidationRegExp, item.rVal.Kind().String()),
		})
	}
}

func (item ValidatableItem) In(vp *ValidationErrors) {
	in := false
	switch item.rVal.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if item.rVal.Int() == 0 {
			break
		}
		sslice := strings.Split(item.vArg, ",")
		for _, s := range sslice {
			i, err := strconv.ParseInt(s, 10, 64)
			if err != nil {
				*vp = append(*vp, ValidationError{Field: item.FieldName, Err: err})
				break
			}
			if i == item.rVal.Int() {
				in = true
				break
			}
		}
		if !in {
			*vp = append(*vp, ValidationError{
				Field: item.FieldName,
				Err:   fmt.Errorf("%w (%d not in[%s])", ErrValidationIn, item.rVal.Int(), item.vArg),
			})
		}
	case reflect.String:
		if item.rVal.String() == "" {
			break
		}
		sslice := strings.Split(item.vArg, ",")
		for _, s := range sslice {
			if s == item.rVal.String() {
				in = true
				break
			}
		}
		if !in {
			*vp = append(*vp, ValidationError{
				Field: item.FieldName,
				Err:   fmt.Errorf("%w (%s not in[%s])", ErrValidationIn, item.rVal.String(), item.vArg),
			})
		}
	case reflect.Array, reflect.Bool, reflect.Chan, reflect.Complex128, reflect.Complex64,
		reflect.Float32, reflect.Float64, reflect.Func, reflect.Interface, reflect.Invalid, reflect.Slice,
		reflect.Map, reflect.Ptr, reflect.Struct, reflect.Uint, reflect.Uint16, reflect.Uint32,
		reflect.Uint64, reflect.Uint8, reflect.Uintptr, reflect.UnsafePointer:
		*vp = append(*vp, ValidationError{
			Field: item.FieldName,
			Err:   fmt.Errorf("%w (wrong field type:%s)", ErrValidationIn, item.rVal.Kind().String()),
		})
	default:
		*vp = append(*vp, ValidationError{
			Field: item.FieldName,
			Err:   fmt.Errorf("%w (wrong field type:%s)", ErrValidationIn, item.rVal.Kind().String()),
		})
	}
}

func (item ValidatableItem) Len(vp *ValidationErrors) {
	switch item.rVal.Kind() {
	case reflect.Slice:
		slen, err := strconv.Atoi(item.vArg)
		if err != nil {
			*vp = append(*vp, ValidationError{Field: item.FieldName, Err: err})
			break
		}
		for i := 0; i < item.rVal.Len(); i++ {
			rF := item.rVal.Index(i)
			if reflect.TypeOf(rF.Interface()).Kind() == reflect.String {
				if rF.Len() != slen {
					*vp = append(*vp, ValidationError{
						Field: item.FieldName,
						Err:   fmt.Errorf("%w (%d != %d)", ErrValidationLen, rF.Len(), slen),
					})
				}
			} else {
				*vp = append(*vp, ValidationError{
					Field: item.FieldName,
					Err:   fmt.Errorf("%w (wrong field type:%s)", ErrValidationLen, rF.Kind().String()),
				})
			}
		}
	case reflect.String:
		slen, err := strconv.Atoi(item.vArg)
		if err != nil {
			*vp = append(*vp, ValidationError{Field: item.FieldName, Err: err})
			break
		}
		if item.rVal.Len() != slen {
			*vp = append(*vp, ValidationError{
				Field: item.FieldName,
				Err:   fmt.Errorf("%w (%d != %d)", ErrValidationLen, item.rVal.Len(), slen),
			})
		}
	case reflect.Array, reflect.Bool, reflect.Chan, reflect.Complex128, reflect.Complex64,
		reflect.Float32, reflect.Float64, reflect.Func, reflect.Int, reflect.Int16,
		reflect.Int32, reflect.Int64, reflect.Int8, reflect.Interface, reflect.Invalid,
		reflect.Map, reflect.Ptr, reflect.Struct, reflect.Uint, reflect.Uint16, reflect.Uint32,
		reflect.Uint64, reflect.Uint8, reflect.Uintptr, reflect.UnsafePointer:
		*vp = append(*vp, ValidationError{
			Field: item.FieldName,
			Err:   fmt.Errorf("%w (wrong field type:%s)", ErrValidationLen, item.rVal.Kind().String()),
		})
	default:
		*vp = append(*vp, ValidationError{
			Field: item.FieldName,
			Err:   fmt.Errorf("%w (wrong field type:%s)", ErrValidationLen, item.rVal.Kind().String()),
		})
	}
}

func (item ValidatableItem) Min(vp *ValidationErrors) {
	switch item.rVal.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		min, err := strconv.ParseInt(item.vArg, 10, 64)
		if err != nil {
			*vp = append(*vp, ValidationError{Field: item.FieldName, Err: err})
			break
		}

		if item.rVal.Int() < min {
			*vp = append(*vp, ValidationError{
				Field: item.FieldName,
				Err:   fmt.Errorf("%w (%d < %d)", ErrValidationMin, item.rVal.Int(), min),
			})
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		min, err := strconv.ParseUint(item.vArg, 10, 64)
		if err != nil {
			*vp = append(*vp, ValidationError{Field: item.FieldName, Err: err})
			break
		}

		if item.rVal.Uint() < min {
			*vp = append(*vp, ValidationError{
				Field: item.FieldName,
				Err:   fmt.Errorf("%w (%d < %d)", ErrValidationMin, item.rVal.Uint(), min),
			})
		}
	case reflect.Float32, reflect.Float64:
		min, err := strconv.ParseFloat(item.vArg, 64)
		if err != nil {
			*vp = append(*vp, ValidationError{Field: item.FieldName, Err: err})
			break
		}

		// Maybe needs use something to check float-values. But I think this enough for our app
		if item.rVal.Float() < min {
			*vp = append(*vp, ValidationError{
				Field: item.FieldName,
				Err:   fmt.Errorf("%w (%g < %g)", ErrValidationMin, item.rVal.Float(), min),
			})
		}
	case reflect.Complex64, reflect.Complex128:
		min, err := strconv.ParseComplex(item.vArg, 128)
		if err != nil {
			*vp = append(*vp, ValidationError{Field: item.FieldName, Err: err})
			break
		}

		if real(item.rVal.Complex()) == real(min) && imag(item.rVal.Complex()) < imag(min) {
			*vp = append(*vp, ValidationError{
				Field: item.FieldName,
				Err:   fmt.Errorf("%w failure (%v < %v)", ErrValidationMin, item.rVal.Complex(), min),
			})
		}
		if real(item.rVal.Complex()) < real(min) {
			*vp = append(*vp, ValidationError{
				Field: item.FieldName,
				Err:   fmt.Errorf("%w failure (%v < %v)", ErrValidationMin, item.rVal.Complex(), min),
			})
		}
	case reflect.Array, reflect.Bool, reflect.Chan, reflect.Func, reflect.Interface, reflect.Invalid, reflect.Map,
		reflect.Ptr, reflect.Slice, reflect.String, reflect.Struct, reflect.Uintptr, reflect.UnsafePointer:
		*vp = append(*vp, ValidationError{
			Field: item.FieldName,
			Err:   fmt.Errorf("%w (wrong field type:%s)", ErrValidationMin, item.rVal.Kind().String()),
		})
	default:
		*vp = append(*vp, ValidationError{
			Field: item.FieldName,
			Err:   fmt.Errorf("%w (wrong field type:%s)", ErrValidationMin, item.rVal.Kind().String()),
		})
	}
}

func (item ValidatableItem) Max(vp *ValidationErrors) {
	switch item.rVal.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		max, err := strconv.ParseInt(item.vArg, 10, 64)
		if err != nil {
			*vp = append(*vp, ValidationError{Field: item.FieldName, Err: err})
			break
		}

		if item.rVal.Int() > max {
			*vp = append(*vp, ValidationError{
				Field: item.FieldName,
				Err:   fmt.Errorf("%w failure (%d > %d)", ErrValidationMax, item.rVal.Int(), max),
			})
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		max, err := strconv.ParseUint(item.vArg, 10, 64)
		if err != nil {
			*vp = append(*vp, ValidationError{Field: item.FieldName, Err: err})
			break
		}

		if item.rVal.Uint() > max {
			*vp = append(*vp, ValidationError{
				Field: item.FieldName,
				Err:   fmt.Errorf("%w (%d < %d)", ErrValidationMax, item.rVal.Uint(), max),
			})
		}
	case reflect.Complex64, reflect.Complex128:
		max, err := strconv.ParseComplex(item.vArg, 128)
		if err != nil {
			*vp = append(*vp, ValidationError{Field: item.FieldName, Err: err})
			break
		}

		if real(item.rVal.Complex()) == real(max) && imag(item.rVal.Complex()) > imag(max) {
			*vp = append(*vp, ValidationError{
				Field: item.FieldName,
				Err:   fmt.Errorf("%w failure (%v > %v)", ErrValidationMax, item.rVal.Complex(), max),
			})
		}
		if real(item.rVal.Complex()) > real(max) {
			*vp = append(*vp, ValidationError{
				Field: item.FieldName,
				Err:   fmt.Errorf("%w failure (%v > %v)", ErrValidationMax, item.rVal.Complex(), max),
			})
		}
	case reflect.Float32, reflect.Float64:
		max, err := strconv.ParseFloat(item.vArg, 64)
		if err != nil {
			*vp = append(*vp, ValidationError{Field: item.FieldName, Err: err})
			break
		}

		// Maybe needs use something to check float-values. But I think it should be enough for our app
		if item.rVal.Float() > max {
			*vp = append(*vp, ValidationError{
				Field: item.FieldName,
				Err:   fmt.Errorf("%w (%g < %g)", ErrValidationMax, item.rVal.Float(), max),
			})
		}
	case reflect.Array, reflect.Bool, reflect.Chan, reflect.Func, reflect.Interface, reflect.Invalid, reflect.Map,
		reflect.Ptr, reflect.Slice, reflect.String, reflect.Struct, reflect.Uintptr, reflect.UnsafePointer:
		*vp = append(*vp, ValidationError{
			Field: item.FieldName,
			Err:   fmt.Errorf("%w (wrong field type:%s)", ErrValidationMax, item.rVal.Kind().String()),
		})
	default:
		*vp = append(*vp, ValidationError{
			Field: item.FieldName,
			Err:   fmt.Errorf("%w (wrong field type:%s)", ErrValidationMax, item.rVal.Kind().String()),
		})
	}
}

func doValidate(vItems *ValidatableItems, vErr *ValidationErrors) {
	for _, vI := range *vItems {
		if _, ok := reflect.TypeOf(vI).MethodByName(vI.vCmd); ok {
			reflect.ValueOf(vI).MethodByName(vI.vCmd).Call([]reflect.Value{reflect.ValueOf(vErr)})
			continue
		}
		*vErr = append(*vErr, ValidationError{Field: vI.FieldName, Err: fmt.Errorf("%w \"%s\"", ErrValidationNImpl, vI.oCmd)})
	}
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

	vItems, vErrs := parseStruct(v)
	doValidate(&vItems, &vErrs)

	if len(vErrs) == 0 {
		return nil
	}

	return vErrs
}
