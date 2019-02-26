/*
Package eerror (enhanced error) describes a new type allowing to handle, process and describe enhanced errors in a way every error should be handled.

Errors should always be composed by four major components:
 - an identifier ("E_PERMISSIONDENIED"), to filter an error by its kind
 - a message (too often considered as sufficient for error description), humanly understandable
 - contexts, from wich error was triggered (stacking a new context each time we forward the error)
 - attributes, essential for error reproducing purposes

This package ensure the ability to manage errors following this pattern painlessly.
*/
package eerror

/*
Eerror interface defines exported methods available to build an enhanced error type.
Enhanced errors allows strict error handling, ensure reproducable errors, and understandable error messages from context args.
*/
type Eerror interface {
	error

	Id() string
	GetAttributes() map[string]interface{}
	Map() map[string]interface{}
	Is(interface{}) bool

	InContext(description string)
	WithAttribute(name string, value interface{})
	WithAttributes(attributeKeyValPairs ...interface{})

	Copy() Eerror
}

// eerror struct defines attributes required to implement Eerror interface.
type eerror struct {
	parent interface{}

	identifier string
	message    string
	contexts   []string
	attributes map[string]interface{}
}
