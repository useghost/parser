package parser

import (
	"math"
	"strconv"
	"strings"
)

const (
	TYPE_UINT8 = iota
	TYPE_UINT16
	TYPE_UINT32
	TYPE_UINT64

	TYPE_INT8
	TYPE_INT32
	TYPE_INT64

	TYPE_FLOAT32
	TYPE_FLOAT64
)

func InferNumericType(s string) string {
	if strings.Contains(s, ".") {
		// Try parsing as float
		f, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return "not a number"
		}
		if f >= -math.MaxFloat32 && f <= math.MaxFloat32 {
			return "float32"
		}
		return "float64"
	}

	// Try parsing as integer
	i, err := strconv.ParseInt(s, 10, 64)
	if err == nil {
		switch {
		case i >= math.MinInt8 && i <= math.MaxInt8:
			return "int8"
		case i >= math.MinInt16 && i <= math.MaxInt16:
			return "int16"
		case i >= math.MinInt32 && i <= math.MaxInt32:
			return "int32"
		default:
			return "int64"
		}
	}

	// Try parsing as unsigned integer
	u, err := strconv.ParseUint(s, 10, 64)
	if err == nil {
		switch {
		case u <= math.MaxUint8:
			return "uint8"
		case u <= math.MaxUint16:
			return "uint16"
		case u <= math.MaxUint32:
			return "uint32"
		default:
			return "uint64"
		}
	}

	return "not a number"
}
