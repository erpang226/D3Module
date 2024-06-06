package tools

import (
	"errors"
	"main/common/global"
	"strconv"
)

// Convert string to other types
func Convert(valueType string, value string) (result interface{}, err error) {
	switch valueType {
	case global.TypeINT:
		return strconv.Atoi(value)
	case global.TypeINT64:
		return strconv.ParseInt(value, 10, 64)
	case global.TypeFLOAT:
		return strconv.ParseFloat(value, 32)
	case global.TypeDOUBLE:
		return strconv.ParseFloat(value, 64)
	case global.TypeBOOL:
		return strconv.ParseBool(value)
	case global.TypeSTRING:
		return value, nil
	default:
		return nil, errors.New("convert failed")
	}
}
