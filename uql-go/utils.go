package uql

import (
	"encoding/base64"
	"fmt"
	"math"
	"math/rand"
	"reflect"
	"regexp"
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
	case FnSplit:
		if len(args) >= 1 {
			text := fmt.Sprintf("%v", args[0])
			separator := ""
			if len(args) >= 2 {
				separator = fmt.Sprintf("%v", args[1])
			}
			if separator == "" {
				// Split into characters
				result := make([]interface{}, 0, len(text))
				for _, ch := range text {
					result = append(result, string(ch))
				}
				return result
			}
			parts := strings.Split(text, separator)
			result := make([]interface{}, len(parts))
			for i, p := range parts {
				result[i] = p
			}
			return result
		}
		return []interface{}{}
	case FnExtract:
		if len(args) >= 3 {
			pattern := fmt.Sprintf("%v", args[0])
			index := 0
			if num, ok := toNumber(args[1]); ok {
				index = int(num)
			}
			text := fmt.Sprintf("%v", args[2])

			// Use regexp to extract
			re, err := regexp.Compile(pattern)
			if err != nil {
				return nil
			}
			matches := re.FindStringSubmatch(text)
			if matches == nil || index >= len(matches) {
				return nil
			}

			result := matches[index]

			// Check for type conversion
			if len(args) >= 4 {
				typeStr := fmt.Sprintf("%v", args[3])
				switch typeStr {
				case "number":
					if num, err := strconv.ParseFloat(result, 64); err == nil {
						return num
					}
				case "date":
					if t, err := time.Parse(time.RFC3339, result); err == nil {
						return t
					}
					// Try other common date formats
					formats := []string{
						"2006-01-02",
						"2006-01-02 15:04:05",
						"01/02/2006",
						"01-02-2006",
					}
					for _, format := range formats {
						if t, err := time.Parse(format, result); err == nil {
							return t
						}
					}
				}
			}
			return result
		}
		return nil

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
	case FnDiff:
		if len(args) >= 2 {
			a, ok1 := toNumber(args[0])
			b, ok2 := toNumber(args[1])
			if ok1 && ok2 {
				return a - b
			}
		}
	case FnMul:
		if len(args) == 0 {
			return 1.0
		}
		result := 1.0
		for _, arg := range args {
			if num, ok := toNumber(arg); ok {
				result *= num
			}
		}
		return result
	case FnDiv:
		if len(args) >= 2 {
			a, ok1 := toNumber(args[0])
			b, ok2 := toNumber(args[1])
			if ok1 && ok2 && b != 0 {
				return a / b
			}
		}
		return nil
	case FnPercentage:
		if len(args) >= 2 {
			part, ok1 := toNumber(args[0])
			whole, ok2 := toNumber(args[1])
			if ok1 && ok2 && whole != 0 {
				return (part / whole) * 100
			}
		}
		return nil

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
	case FnSinh:
		if len(args) > 0 {
			if num, ok := toNumber(args[0]); ok {
				return math.Sinh(num)
			}
		}
	case FnCosh:
		if len(args) > 0 {
			if num, ok := toNumber(args[0]); ok {
				return math.Cosh(num)
			}
		}
	case FnTanh:
		if len(args) > 0 {
			if num, ok := toNumber(args[0]); ok {
				return math.Tanh(num)
			}
		}
	case FnAtan2:
		if len(args) >= 2 {
			y, ok1 := toNumber(args[0])
			x, ok2 := toNumber(args[1])
			if ok1 && ok2 {
				return math.Atan2(y, x)
			}
		}
	case FnAsinh:
		if len(args) > 0 {
			if num, ok := toNumber(args[0]); ok {
				return math.Asinh(num)
			}
		}
	case FnAcosh:
		if len(args) > 0 {
			if num, ok := toNumber(args[0]); ok {
				return math.Acosh(num)
			}
		}
	case FnAtanh:
		if len(args) > 0 {
			if num, ok := toNumber(args[0]); ok {
				return math.Atanh(num)
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
	case FnParseURL:
		if len(args) > 0 {
			urlStr := fmt.Sprintf("%v", args[0])
			return parseURL(urlStr)
		}
	case FnParseURLQuery:
		if len(args) > 0 {
			queryStr := fmt.Sprintf("%v", args[0])
			return parseURLQuery(queryStr)
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
	case FnDistinct:
		if len(args) > 0 {
			// Get array argument
			arr, err := toSlice(args[0])
			if err != nil {
				return []interface{}{}
			}
			// Use map to track unique values
			seen := make(map[string]bool)
			result := make([]interface{}, 0)
			for _, item := range arr {
				key := fmt.Sprintf("%v", item)
				if !seen[key] {
					seen[key] = true
					result = append(result, item)
				}
			}
			return result
		}
		return []interface{}{}
	case FnPack:
		// Pack multiple values into an array
		return args
	case FnArrayFromEntries:
		if len(args) > 0 {
			arr, err := toSlice(args[0])
			if err != nil {
				return make(map[string]interface{})
			}
			result := make(map[string]interface{})
			for _, item := range arr {
				if m, ok := toMap(item); ok {
					if key, hasKey := m["key"]; hasKey {
						if value, hasValue := m["value"]; hasValue {
							keyStr := fmt.Sprintf("%v", key)
							result[keyStr] = value
						}
					}
				}
			}
			return result
		}
		return make(map[string]interface{})
	case FnArrayToMap:
		if len(args) >= 2 {
			arr, err := toSlice(args[0])
			if err != nil {
				return make(map[string]interface{})
			}
			keyField := fmt.Sprintf("%v", args[1])

			result := make(map[string]interface{})
			for _, item := range arr {
				if m, ok := toMap(item); ok {
					if keyVal, exists := m[keyField]; exists {
						keyStr := fmt.Sprintf("%v", keyVal)
						result[keyStr] = item
					}
				}
			}
			return result
		}
		return make(map[string]interface{})
	case FnBagPack:
		// Similar to pack but creates an object with field names
		if len(args)%2 != 0 {
			return nil
		}
		result := make(map[string]interface{})
		for i := 0; i < len(args); i += 2 {
			key := fmt.Sprintf("%v", args[i])
			value := args[i+1]
			result[key] = value
		}
		return result
	case FnKV:
		// Create key-value pair object
		if len(args) >= 2 {
			return map[string]interface{}{
				"key":   args[0],
				"value": args[1],
			}
		}
		return nil

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
	case FnUnixTimeMicrosecondsToDatetime:
		if len(args) > 0 {
			if num, ok := toNumber(args[0]); ok {
				us := int64(num)
				return time.Unix(0, us*int64(time.Microsecond))
			}
		}
	case FnUnixTimeNanosecondsToDatetime:
		if len(args) > 0 {
			if num, ok := toNumber(args[0]); ok {
				ns := int64(num)
				return time.Unix(0, ns)
			}
		}
	case FnStartOfMonth:
		if len(args) > 0 {
			if t, ok := args[0].(time.Time); ok {
				return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
			}
		}
	case FnStartOfWeek:
		if len(args) > 0 {
			if t, ok := args[0].(time.Time); ok {
				// Start of week (Sunday)
				weekday := int(t.Weekday())
				return time.Date(t.Year(), t.Month(), t.Day()-weekday, 0, 0, 0, 0, t.Location())
			}
		}
	case FnStartOfYear:
		if len(args) > 0 {
			if t, ok := args[0].(time.Time); ok {
				return time.Date(t.Year(), 1, 1, 0, 0, 0, 0, t.Location())
			}
		}
	case FnToDatetime:
		if len(args) > 0 {
			str := fmt.Sprintf("%v", args[0])
			// Try RFC3339 first
			if t, err := time.Parse(time.RFC3339, str); err == nil {
				return t
			}
			// Try other common formats
			formats := []string{
				"2006-01-02",
				"2006-01-02 15:04:05",
				"01/02/2006",
				"01-02-2006",
				time.RFC1123,
				time.RFC822,
			}
			for _, format := range formats {
				if t, err := time.Parse(format, str); err == nil {
					return t
				}
			}
		}
	case FnFormatDatetime:
		if len(args) >= 2 {
			if t, ok := args[0].(time.Time); ok {
				format := fmt.Sprintf("%v", args[1])
				// Convert from .NET/C# style format to Go's time format
				// This is a simplified conversion
				format = strings.ReplaceAll(format, "yyyy", "2006")
				format = strings.ReplaceAll(format, "yy", "06")
				format = strings.ReplaceAll(format, "MM", "01")
				format = strings.ReplaceAll(format, "dd", "02")
				format = strings.ReplaceAll(format, "HH", "15")
				format = strings.ReplaceAll(format, "mm", "04")
				format = strings.ReplaceAll(format, "ss", "05")
				return t.Format(format)
			}
		}
	case FnAddDatetime:
		if len(args) >= 3 {
			if t, ok := args[0].(time.Time); ok {
				value, ok1 := toNumber(args[1])
				unit := fmt.Sprintf("%v", args[2])
				if ok1 {
					duration := int(value)
					switch strings.ToLower(unit) {
					case "year", "years":
						return t.AddDate(duration, 0, 0)
					case "month", "months":
						return t.AddDate(0, duration, 0)
					case "day", "days":
						return t.AddDate(0, 0, duration)
					case "hour", "hours":
						return t.Add(time.Duration(duration) * time.Hour)
					case "minute", "minutes":
						return t.Add(time.Duration(duration) * time.Minute)
					case "second", "seconds":
						return t.Add(time.Duration(duration) * time.Second)
					}
				}
			}
		}
	}

	return nil
}

// parseURL parses a URL and returns components
func parseURL(urlStr string) map[string]interface{} {
	result := make(map[string]interface{})

	// Simple URL parsing without importing net/url to keep dependencies minimal
	// Protocol
	protocolEnd := strings.Index(urlStr, "://")
	if protocolEnd > 0 {
		result["scheme"] = urlStr[:protocolEnd]
		urlStr = urlStr[protocolEnd+3:]
	}

	// Host and path
	pathStart := strings.Index(urlStr, "/")
	if pathStart > 0 {
		hostPort := urlStr[:pathStart]
		result["host"] = hostPort

		// Port
		portStart := strings.LastIndex(hostPort, ":")
		if portStart > 0 {
			result["hostname"] = hostPort[:portStart]
			result["port"] = hostPort[portStart+1:]
		} else {
			result["hostname"] = hostPort
		}

		// Path and query
		remaining := urlStr[pathStart:]
		queryStart := strings.Index(remaining, "?")
		if queryStart > 0 {
			result["pathname"] = remaining[:queryStart]
			queryAndHash := remaining[queryStart+1:]

			// Hash
			hashStart := strings.Index(queryAndHash, "#")
			if hashStart > 0 {
				result["search"] = "?" + queryAndHash[:hashStart]
				result["query"] = queryAndHash[:hashStart]
				result["hash"] = "#" + queryAndHash[hashStart+1:]
			} else {
				result["search"] = "?" + queryAndHash
				result["query"] = queryAndHash
			}
		} else {
			hashStart := strings.Index(remaining, "#")
			if hashStart > 0 {
				result["pathname"] = remaining[:hashStart]
				result["hash"] = "#" + remaining[hashStart+1:]
			} else {
				result["pathname"] = remaining
			}
		}
	} else {
		result["host"] = urlStr
		result["hostname"] = urlStr
	}

	return result
}

// parseURLQuery parses URL query string
func parseURLQuery(queryStr string) map[string]interface{} {
	result := make(map[string]interface{})

	// Remove leading ? if present
	queryStr = strings.TrimPrefix(queryStr, "?")

	if queryStr == "" {
		return result
	}

	pairs := strings.Split(queryStr, "&")
	for _, pair := range pairs {
		kv := strings.SplitN(pair, "=", 2)
		if len(kv) == 2 {
			result[kv[0]] = kv[1]
		} else if len(kv) == 1 {
			result[kv[0]] = ""
		}
	}

	return result
}
