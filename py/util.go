package py

import (
	"errors"
	"strconv"
)

func GetLen(obj Object) (Int, error) {
	getlen, ok := obj.(I__len__)
	if !ok {
		return 0, nil
	}

	lenObj, err := getlen.M__len__()
	if err != nil {
		return 0, err
	}

	return GetInt(lenObj)
}

func GetInt(obj Object) (Int, error) {
	toIdx, ok := obj.(I__index__)
	if !ok {
		_, err := cantConvert(obj, "int")
		return 0, err
	}

	return toIdx.M__index__()
}

func LoadTuple(args Tuple, vars []interface{}) error {

	if len(args) > len(vars) {
		return ExceptionNewf(RuntimeError, "%d args given, expected %d", len(args), len(vars))
	}

	if len(vars) > len(args) {
		vars = vars[:len(args)]
	}

	for i, rval := range vars {
		err := loadValue(args[i], rval)
		if err == ErrUnsupportedObjType {
			return ExceptionNewf(TypeError, "arg %d has unsupported object type: %s", i, args[i].Type().Name)
		}
	}

	return nil
}

var (
	ErrUnsupportedObjType = errors.New("unsupported obj type")
)

func LoadAttr(obj Object, attrName string, data interface{}) error {
	attr, err := GetAttrString(obj, attrName)
	if err != nil {
		return err
	}
	err = loadValue(attr, data)
	if err == ErrUnsupportedObjType {
		return ExceptionNewf(TypeError, "attribute \"%s\" has unsupported object type: %s", attrName, attr.Type().Name)
	}
	return nil
}

func loadValue(src Object, data interface{}) error {

	var (
		v_str   string
		v_float float64
		v_int   int64
	)

	if b, ok := src.(Bool); ok {
		if b {
			v_int = 1
			v_float = 1
			v_str = "True"
		} else {
			v_str = "False"
		}
	} else if val, ok := src.(Int); ok {
		v_int = int64(val)
		v_float = float64(val)
	} else if val, ok := src.(Float); ok {
		v_int = int64(val)
		v_float = float64(val)
	} else if val, ok := src.(String); ok {
		v_str = string(val)
		intval, _ := strconv.Atoi(v_str)
		v_int = int64(intval)
	} else if _, ok := src.(NoneType); ok {
		// No-op
	} else {
		return ErrUnsupportedObjType
	}

	switch dst := data.(type) {
	case *bool:
		*dst = v_int != 0
	case *int8:
		*dst = int8(v_int)
	case *uint8:
		*dst = uint8(v_int)
	case *int16:
		*dst = int16(v_int)
	case *uint16:
		*dst = uint16(v_int)
	case *int32:
		*dst = int32(v_int)
	case *uint32:
		*dst = uint32(v_int)
	case *int:
		*dst = int(v_int)
	case *uint:
		*dst = uint(v_int)
	case *int64:
		*dst = v_int
	case *uint64:
		*dst = uint64(v_int)
	case *float32:
		*dst = float32(v_float)
	case *float64:
		*dst = v_float
	case *string:
		*dst = v_str

	// case []uint64:
	// 	for i := range data {
	// 		dst[i] = order.Uint64(bs[8*i:])
	// 	}
	// case []float32:
	// 	for i := range data {
	// 		dst[i] = math.Float32frombits(order.Uint32(bs[4*i:]))
	// 	}
	// case []float64:
	// 	for i := range data {
	// 		dst[i] = math.Float64frombits(order.Uint64(bs[8*i:]))
	// 	}

	default:
		return ExceptionNewf(NotImplementedError, "%s", "unsupported Go data type")
	}
	return nil
}
