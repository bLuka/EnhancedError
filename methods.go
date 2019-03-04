package eerror

import (
	"fmt"
	"runtime/debug"
)

/*
NewEerror instanciates a new enhanced error given its unique identifier, message, and potential attributes.
Error types should be declared as string constants, preferably in a separated package of yours.
For readability reasons, you should divide error types in multiple small files, grouping them by categories, and prefixing their symbols accordingly.

  const E_MY_ERROR_ID = "E_MY_ERROR_ID"

  func errorFunction(myParameter interface{}) Eerror {
     return NewError(E_MY_ERROR_ID, "This function panics",
                     "parameter", myParameter,
     )
  }

  func main() {
     if eerr := errorFunction("hello world"); eerr != nil {
        panic(eerr)
     }
  }
*/
func NewError(identifier, message string, attributeKeyValPairs ...interface{}) Eerror {
	e := Eerror{
		nil,
		identifier,
		message,
		[]string{},
		make(map[string]interface{}, len(attributeKeyValPairs)/2),
		generateUniqueID(),
	}

	e.WithAttribute("stacktrace", string(debug.Stack()))
	e.WithAttributes(attributeKeyValPairs...)
	return e
}

// InContext appends a new context to the error stack. Useful to describe context during error forwarding.
func (e *Eerror) InContext(context string) {
	e.contexts = append(e.contexts, context)
}

// WithAttribute allows attribute set to an error. If any attribute with the same name exists, it will be reset
func (e *Eerror) WithAttribute(name string, value interface{}) {
	e.WithAttributes(name, value)
}

// WithAttributes allow setting multiple attributes at once. If any attribute with the same name exists, they will be reset
func (e *Eerror) WithAttributes(attributeKeyValPairs ...interface{}) {
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
		if e.attributes == nil {
			e.attributes = make(map[string]interface{})
		}
		e.attributes[key] = value
	}
}

// GetAttributes retrieves the attributes map copy
func (e Eerror) GetAttributes() map[string]interface{} {
	return e.attributes
}

// Id returns the identifier of the error
func (e Eerror) Id() string {
	return e.identifier
}
