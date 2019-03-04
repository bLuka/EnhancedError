package eerror

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
)

/*
Error formats the error to a human readable string, as described by the error interface.

Eg: `E_SOMEERROR: My error message (context 1; "context 2 with; (special) chars") [some attribute: some value, some other attribute: (int)1]
*/
func (e Eerror) Error() string {
	const contextSeparator = "; "
	var contextString string
	var attributesString string

	if len(e.contexts) > 0 {
		contextString = " ("
		for i, context := range e.contexts {
			if i > 0 {
				contextString += contextSeparator
			}
			contextString += escapeString(context, "();")
		}
		contextString += ")"
	}
	if len(e.attributes) > 0 {
		attributesString = " ["

		keys := make([]string, 0)
		for key := range e.attributes {
			keys = append(keys, key)
		}
		sort.Strings(keys)

		prependSeparator := false
		for _, key := range keys {
			value := e.attributes[key]

			if prependSeparator {
				attributesString += ", "
			}
			prependSeparator = true

			var serializedValue = serialize(value)
			if reflect.TypeOf(value).Kind() == reflect.String {
				serializedValue = escapeString(serializedValue, "[]:,")
			}
			attributesString += escapeString(key, "[]:,") + ": " + serializedValue
		}
		attributesString += "]"
	}

	return fmt.Sprintf("%s: %s%s%s", escapeString(e.identifier, ":"), escapeString(e.message, ":()[]"), contextString, attributesString)
}

// Map formats the error to a protocol-aware object, marshable without data loss
func (e Eerror) Map() map[string]interface{} {
	return map[string]interface{}{
		"error":      e.Error(),
		"code":       e.identifier,
		"message":    e.message,
		"contexts":   e.contexts,
		"attributes": e.attributes,
	}
}

func escapeString(s string, chars string) string {
	if len(s) == 0 || strings.IndexAny(s, chars+"\"") != -1 {
		return fmt.Sprintf("\"%s\"", strings.Replace(s, "\"", "\\\"", -1))
	}
	return s
}

func serialize(value interface{}) string {
	var valueType = reflect.TypeOf(value)

	if valueType.Kind() == reflect.String {
		return value.(string)
	}
	return fmt.Sprintf("(%s)%v", valueType, value)
}
