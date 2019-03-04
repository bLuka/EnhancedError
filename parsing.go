package eerror

import (
	"fmt"
	"runtime/debug"
	"strconv"
	"strings"
)

// parse unserializes an enhanced error from it's string representation to the Eerror format
func parse(err interface{}) (eerr Eerror, ok bool) {
	ok = false
	defer (func() {
		if e := recover(); e != nil {
			ok = false
		}
	})()

	var s = strings.TrimSpace(fmt.Sprint(err))

	errType, ok, endPos := parseEerrorType(s)
	if !ok {
		return
	}
	s = s[endPos:]
	message, ok, endPos := parseEerrorMessage(s)
	if !ok {
		return
	}
	s = s[endPos:]
	contexts, ok, endPos := parseEerrorContexts(s)
	if !ok {
		return
	}
	s = s[endPos:]
	attributes, ok, endPos := parseEerrorAttributes(s)
	if !ok {
		return
	}
	s = s[endPos:]

	eerr = Eerror{
		err,

		errType,
		message,
		contexts,
		attributes,

		generateUniqueID(),
	}
	if _, ok := attributes["stacktrace"]; !ok {
		eerr.WithAttribute("stacktrace", string(debug.Stack()))
	}
	return
}

func parseEerrorAttributes(s string) (attributes map[string]interface{}, ok bool, endPosition int) {
	ok = true
	if len(s) == 0 {
		return
	}

	ok = false
	if s[0] != '[' {
		return
	}

	parseAttribute := func(s string, index int) (endPosition int, name string, value interface{}, e error) {
		e = fmt.Errorf("Invalid")

		if s[index] == '"' {
			endPosition = strings.IndexByte(s[index+1:], '"')
			if endPosition == -1 {
				return
			}

			endPosition += index + 1
			name = s[index:endPosition]
		} else {
			endPosition = strings.IndexAny(s[index:], ":]")
			if endPosition == -1 {
				return
			}

			endPosition += index
			if s[endPosition] == ']' {
				return
			}

			name = s[index:endPosition]
		}

		if s[endPosition:endPosition+2] != ": " {
			return
		}
		index = endPosition + 2
		if s[index] == '"' {
			endPosition = strings.IndexByte(s[index+1:], '"')
			if endPosition == -1 {
				return
			}

			endPosition += index + 2
			value = s[index+1 : endPosition-1]
		} else {
			endPosition = strings.IndexAny(s[index:], ",]")
			if endPosition == -1 {
				return
			}

			endPosition += index
			value = s[index:endPosition]
		}

		if indexLastPar := strings.IndexByte(value.(string)[1:], ')'); value.(string)[0] == '(' && indexLastPar != -1 && indexLastPar < len(value.(string))-1 {
			t := value.(string)[1 : indexLastPar-1]

			switch t {
			case "string":
				value = value.(string)[indexLastPar+1:]
			case "bool":
				value, e = strconv.ParseBool(value.(string)[indexLastPar+1:])
				if e != nil {
					return
				}
			case "float", "float32", "float64":
				value, e = strconv.ParseFloat(value.(string)[indexLastPar+1:], 64)
				if e != nil {
					return
				}
			case "int", "uint":
				value, e = strconv.Atoi(value.(string)[indexLastPar+1:])
				if e != nil {
					return
				}
			}
		}

		e = nil
		if s[endPosition] != ']' {
			endPosition++
		}
		return
	}

	var index int = 1
	attributes = make(map[string]interface{})
	for s[index] != ']' {
		var name string
		var value interface{}
		var err error

		index, name, value, err = parseAttribute(s, index)
		attributes[name] = value
		if err != nil {
			return
		}
		for index < len(s) && s[index] == ' ' {
			index++
		}
		if index >= len(s) {
			return
		}
	}

	ok = true
	return
}

func parseEerrorContexts(s string) (contexts []string, ok bool, endPosition int) {
	ok = true
	if len(s) == 0 {
		return
	}

	ok = false
	if s[0:1] != "(" {
		return
	}

	parseContextName := func(s string, index int) (endPosition int, name string, e error) {
		e = fmt.Errorf("Invalid")

		if s[index] == ';' {
			return
		}
		if s[index] == '"' {
			endPosition = strings.IndexByte(s[index+1:], '"')
			if endPosition == -1 {
				return
			}

			endPosition += index + 1
			name = s[index+1 : endPosition]
		} else {
			endPosition = strings.IndexAny(s[index:], ";)")
			if endPosition == -1 {
				return
			}

			endPosition += index
			name = s[index:endPosition]
		}

		if s[endPosition] != ')' {
			endPosition++
		}
		if s[endPosition-1:endPosition+1] == ";)" {
			return
		}

		e = nil
		return
	}

	endPosition = 1
	for s[endPosition] != ')' {
		var name string
		var err error

		endPosition, name, err = parseContextName(s, endPosition)
		if err != nil {
			return
		}
		contexts = append(contexts, name)
		for endPosition < len(s) && s[endPosition] == ' ' {
			endPosition++
		}
		if endPosition >= len(s) {
			return
		}
	}

	endPosition++
	for endPosition < len(s) && s[endPosition] == ' ' {
		endPosition++
	}
	ok = true
	return
}

func parseEerrorMessage(s string) (message string, ok bool, endPosition int) {
	ok = false

	if len(s) < 1 {
		return
	}

	if s[0] == '"' {
		endPosition = strings.IndexByte(s[1:], '"')
		if endPosition == -1 {
			return
		}

		endPosition += 2
		if endPosition > 1 {
			message = s[1 : endPosition-1]
		}
		if endPosition < len(s) && s[endPosition] == ' ' {
			endPosition++
		}
	} else {
		endPosition = strings.IndexAny(s[1:], "([")
		if endPosition == -1 {
			endPosition = len(s) - 1
		}
		endPosition++
		message = s[:endPosition-1]
	}

	ok = true
	return
}

func parseEerrorType(s string) (errType string, ok bool, endPosition int) {
	ok = false

	if s[0] == '"' {
		endPosition = strings.IndexByte(s[1:], '"')
		if endPosition == -1 {
			return
		}
		if endPosition == len(s)-1 {
			return
		}

		endPosition += 2
		if s[endPosition] != ':' {
			return
		}
		errType = s[1 : endPosition-1]
	} else {
		endPosition = strings.IndexByte(s, ':')
		if endPosition < 1 {
			return
		}
		if endPosition == len(s)-1 {
			return
		}
		errType = s[:endPosition]
	}

	if endPosition+2 > len(s) || s[endPosition+1] != ' ' {
		return "", false, -1
	}

	ok = true
	endPosition += 2
	return
}

func parseString(s string) string {

	return s
}
