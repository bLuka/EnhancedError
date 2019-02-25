package eerror

import (
	"fmt"
)

/*
Eerror interface type defines exported methods available to build an enhanced error type.
Enhanced errors allows strict error handling, ensure reproducable errors, and understandable error messages from context args.
*/
type Eerror interface {
	Error() string

	InContext(description string)
	WithAttribute(name string, value interface{})
	WithAttributes(attributeKeyValPairs ...interface{})
	GetAttributes() map[string]interface{}
}

type eerror struct {
	identifier string
	message    string
	context    []string
	attributes map[string]interface{}
}

func NewError(identifier, message string, attributeKeyValPairs ...interface{}) Eerror {
	e := eerror{
		identifier,
		message,
		[]string{},
		make(map[string]interface{}, len(attributeKeyValPairs)/2),
	}

	e.WithAttributes(attributeKeyValPairs...)
	return &e
}

func (e eerror) Error() string {
	const contextSeparator = "\n ->"
	var contextString string

	func() {
		for _, context := range e.context {
			defer func(context string) {
				contextString += contextSeparator + context
			}(context)
		}
	}()

	return fmt.Sprintf("%s: %s%s", e.identifier, e.message, contextString)
}

func (e *eerror) InContext(context string) {
	e.context = append(e.context, context)
}

func (e *eerror) WithAttribute(name string, value interface{}) {
	e.WithAttributes(name, value)
}

func (e *eerror) WithAttributes(attributeKeyValPairs ...interface{}) {
	for i, value := range attributeKeyValPairs {
		if i%2 != 0 {
			continue
		}

		key, ok := value.(string)
		if !ok {
			key = fmt.Sprint(value)
		}

		var value interface{} = nil
		if len(attributeKeyValPairs) > i+1 {
			value = attributeKeyValPairs[i+1]
		}
		e.attributes[key] = value
	}
}

func (e eerror) GetAttributes() map[string]interface{} {
	return e.attributes
}
