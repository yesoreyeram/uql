package uql

import (
	"encoding/base64"
	"fmt"
	"math"
	"math/rand"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func init() {
	// Seed the random number generator at package initialization
	rand.Seed(time.Now().UnixNano())
}

// getValue retrieves a value from an object using dot notation
func getValue(obj interface{}, path string) interface{} {
	if obj == nil {
		return nil
	}

	parts := strings.Split(path, ".")
	current := obj

	for _, part := range parts {
		m, ok := toMap(current)
		if !ok {
			return nil
		}
		val, exists := m[part]
		if !exists {
			return nil
		}
		current = val
	}

	return current
}

// toSlice converts an interface{} to []interface{}
func toSlice(v interface{}) ([]interface{}, error) {
	if slice, ok := v.([]interface{}); ok {
		return slice, nil
	}

	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Slice && val.Kind() != reflect.Array {
		return nil, fmt.Errorf("not a slice or array")
	}

	result := make([]interface{}, val.Len())
	for i := 0; i < val.Len(); i++ {
		result[i] = val.Index(i).Interface()
	}

	return result, nil
}

// compareValues compares two values, returns -1, 0, or 1
func compareValues(a, b interface{}) int {
	if a == nil && b == nil {
		return 0
	}
	if a == nil {
		return -1
	}
	if b == nil {
		return 1
	}

	// Try numeric comparison
	aNum, aOk := toNumber(a)
	bNum, bOk := toNumber(b)
	if aOk && bOk {
		if aNum < bNum {
			return -1
		}
		if aNum > bNum {
			return 1
		}
		return 0
	}

	// Try string comparison
	aStr := fmt.Sprintf("%v", a)
	bStr := fmt.Sprintf("%v", b)
	if aStr < bStr {
		return -1
	}
	if aStr > bStr {
		return 1
	}
	return 0
}

// toNumber converts a value to float64
func toNumber(v interface{}) (float64, bool) {
	switch val := v.(type) {
	case float64:
		return val, true
	case float32:
		return float64(val), true
	case int:
		return float64(val), true
	case int64:
		return float64(val), true
	case int32:
		return float64(val), true
	case string:
		if num, err := strconv.ParseFloat(val, 64); err == nil {
			return num, true
		}
	}
	return 0, false
}

// evaluateFunction evaluates a function with given arguments
func evaluateFunction(fn FunctionName, args []interface{}) interface{} {
	switch fn {
	// String functions
	case FnToupper:
		if len(args) > 0 {
			return strings.ToUpper(fmt.Sprintf("%v", args[0]))
		}
	case FnToLower:
		if len(args) > 0 {
			return strings.ToLower(fmt.Sprintf("%v", args[0]))
		}
	case FnTrim:
		if len(args) > 0 {
			return strings.TrimSpace(fmt.Sprintf("%v", args[0]))
		}
	case FnTrimStart:
		if len(args) > 0 {
			return strings.TrimLeft(fmt.Sprintf("%v", args[0]), " \t\n\r")
		}
	case FnTrimEnd:
		if len(args) > 0 {
			return strings.TrimRight(fmt.Sprintf("%v", args[0]), " \t\n\r")
		}
	case FnReverse:
		if len(args) > 0 {
			s := fmt.Sprintf("%v", args[0])
			runes := []rune(s)
			for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
				runes[i], runes[j] = runes[j], runes[i]
			}
			return string(runes)
		}
	case FnStrlen:
		if len(args) > 0 {
			return len(fmt.Sprintf("%v", args[0]))
		}
	case FnStrcat:
		var result strings.Builder
		for _, arg := range args {
			result.WriteString(fmt.Sprintf("%v", arg))
		}
		return result.String()
	case FnReplaceString:
		if len(args) >= 3 {
			text := fmt.Sprintf("%v", args[0])
			old := fmt.Sprintf("%v", args[1])
			new := fmt.Sprintf("%v", args[2])
			return strings.ReplaceAll(text, old, new)
		}
	case FnSubstring:
		if len(args) >= 2 {
			text := fmt.Sprintf("%v", args[0])
			start, _ := toNumber(args[1])
			startInt := int(start)
			if startInt < 0 || startInt >= len(text) {
				return ""
			}
			if len(args) >= 3 {
				length, _ := toNumber(args[2])
				lengthInt := int(length)
				end := startInt + lengthInt
				if end > len(text) {
					end = len(text)
				}
				return text[startInt:end]
			}
			return text[startInt:]
		}

	// Math functions
	case FnSum:
		var sum float64
		for _, arg := range args {
			if num, ok := toNumber(arg); ok {
				sum += num
			}
		}
		return sum
	case FnMin:
		if len(args) == 0 {
			return nil
		}
		min, ok := toNumber(args[0])
		if !ok {
			return nil
		}
		for _, arg := range args[1:] {
			if num, ok := toNumber(arg); ok && num < min {
				min = num
			}
		}
		return min
	case FnMax:
		if len(args) == 0 {
			return nil
		}
		max, ok := toNumber(args[0])
		if !ok {
			return nil
		}
		for _, arg := range args[1:] {
			if num, ok := toNumber(arg); ok && num > max {
				max = num
			}
		}
		return max
	case FnMean:
		if len(args) == 0 {
			return nil
		}
		var sum float64
		count := 0
		for _, arg := range args {
			if num, ok := toNumber(arg); ok {
				sum += num
				count++
			}
		}
		if count == 0 {
			return nil
		}
		return sum / float64(count)
	case FnAbs:
		if len(args) > 0 {
			if num, ok := toNumber(args[0]); ok {
				return math.Abs(num)
			}
		}
	case FnFloor:
		if len(args) > 0 {
			if num, ok := toNumber(args[0]); ok {
				return math.Floor(num)
			}
		}
	case FnCeil:
		if len(args) > 0 {
			if num, ok := toNumber(args[0]); ok {
				return math.Ceil(num)
			}
		}
	case FnRound:
		if len(args) > 0 {
			if num, ok := toNumber(args[0]); ok {
				return math.Round(num)
			}
		}
	case FnPow:
		if len(args) >= 2 {
			base, ok1 := toNumber(args[0])
			exp, ok2 := toNumber(args[1])
			if ok1 && ok2 {
				return math.Pow(base, exp)
			}
		}
	case FnLog:
		if len(args) > 0 {
			if num, ok := toNumber(args[0]); ok && num > 0 {
				return math.Log(num)
			}
		}
	case FnLog2:
		if len(args) > 0 {
			if num, ok := toNumber(args[0]); ok && num > 0 {
				return math.Log2(num)
			}
		}
	case FnLog10:
		if len(args) > 0 {
			if num, ok := toNumber(args[0]); ok && num > 0 {
				return math.Log10(num)
			}
		}

	// Trigonometric functions
	case FnSin:
		if len(args) > 0 {
			if num, ok := toNumber(args[0]); ok {
				return math.Sin(num)
			}
		}
	case FnCos:
		if len(args) > 0 {
			if num, ok := toNumber(args[0]); ok {
				return math.Cos(num)
			}
		}
	case FnTan:
		if len(args) > 0 {
			if num, ok := toNumber(args[0]); ok {
				return math.Tan(num)
			}
		}
	case FnAsin:
		if len(args) > 0 {
			if num, ok := toNumber(args[0]); ok {
				return math.Asin(num)
			}
		}
	case FnAcos:
		if len(args) > 0 {
			if num, ok := toNumber(args[0]); ok {
				return math.Acos(num)
			}
		}
	case FnAtan:
		if len(args) > 0 {
			if num, ok := toNumber(args[0]); ok {
				return math.Atan(num)
			}
		}

	// Conversion functions
	case FnToString:
		if len(args) > 0 {
			return fmt.Sprintf("%v", args[0])
		}
	case FnToInt, FnToLong:
		if len(args) > 0 {
			if num, ok := toNumber(args[0]); ok {
				return int64(num)
			}
		}
	case FnToDouble, FnToFloat, FnToNumber:
		if len(args) > 0 {
			if num, ok := toNumber(args[0]); ok {
				return num
			}
		}
	case FnToBool:
		if len(args) > 0 {
			switch v := args[0].(type) {
			case bool:
				return v
			case string:
				// Case-insensitive boolean parsing
				lower := strings.ToLower(strings.TrimSpace(v))
				return lower == "true" || lower == "1" || lower == "yes"
			case int, int64, float64:
				num, _ := toNumber(v)
				return num != 0
			}
		}

	// Encoding functions
	case FnAtob:
		if len(args) > 0 {
			str := fmt.Sprintf("%v", args[0])
			decoded, err := base64.StdEncoding.DecodeString(str)
			if err == nil {
				return string(decoded)
			}
		}
	case FnBtoa:
		if len(args) > 0 {
			str := fmt.Sprintf("%v", args[0])
			return base64.StdEncoding.EncodeToString([]byte(str))
		}

	// Other functions
	case FnRandom:
		return rand.Float64()
	case FnSign:
		if len(args) > 0 {
			if num, ok := toNumber(args[0]); ok {
				if num > 0 {
					return 1
				} else if num < 0 {
					return -1
				}
				return 0
			}
		}

	// Date/time functions
	case FnUnixTimeMillisecondsToDatetime:
		if len(args) > 0 {
			if num, ok := toNumber(args[0]); ok {
				ms := int64(num)
				return time.Unix(0, ms*int64(time.Millisecond))
			}
		}
	case FnUnixTimeSecondsToDatetime:
		if len(args) > 0 {
			if num, ok := toNumber(args[0]); ok {
				sec := int64(num)
				return time.Unix(sec, 0)
			}
		}
	case FnToUnixTime:
		if len(args) > 0 {
			if t, ok := args[0].(time.Time); ok {
				return t.Unix()
			}
		}
	case FnStartOfDay:
		if len(args) > 0 {
			if t, ok := args[0].(time.Time); ok {
				return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
			}
		}
	case FnStartOfHour:
		if len(args) > 0 {
			if t, ok := args[0].(time.Time); ok {
				return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), 0, 0, 0, t.Location())
			}
		}
	case FnStartOfMinute:
		if len(args) > 0 {
			if t, ok := args[0].(time.Time); ok {
				return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), 0, 0, t.Location())
			}
		}
	}

	return nil
}
