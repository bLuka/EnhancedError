package eerror

import (
	"fmt"
	"strings"
)

// Error formats the error to a human readable string, as described by the error interface
func (e Eerror) Error() string {
	const contextSeparator = "; "
	var contextString string

	if len(e.contexts) > 0 {
		contextString = " (" + strings.Join(e.contexts, contextSeparator) + ")"
	}

	return fmt.Sprintf("%s: %s%s", e.identifier, e.message, contextString)
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
