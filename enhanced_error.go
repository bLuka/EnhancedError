/*
Package eerror (enhanced error) describes a new type allowing to handle, process and describe enhanced errors in a way every error should be handled.

Errors should always be composed by four major components:
 - an identifier ("E_PERMISSIONDENIED"), to filter an error by its kind
 - a message (too often considered as sufficient alone for error handling), humanly understandable
 - contexts, from wich error was triggered (stacking a new context each time we forward the error)
 - attributes, essential for error reproducing purposes

This package ensures the ability to manage errors following this pattern painlessly.
*/
package eerror

/*
Eerror type defines attributes available to build an enhanced error type.
Enhanced errors allows strict error handling, ensure reproducable errors, and understandable error messages from context args.
*/
type Eerror struct {
	parent interface{}

	identifier string
	message    string
	contexts   []string
	attributes map[string]interface{}

	_instance uint
}
