package eerror

import (
	"fmt"
	"testing"
)

func TestSerialization(t *testing.T) {
	const errorMessage = "This is a test error"
	const contextMessage = "This is some context"

	err := NewError(E_TESTERROR, errorMessage)

	if expected := E_TESTERROR + ": " + errorMessage; err.Error() != expected {
		t.Error("Invalid error formatting (result, expected)\n", err.Error()+"\n", expected)
	}

	err.WithAttribute("some attribute", "some value")
	err.WithAttribute("some typed attribute", 42)
	if expected := E_TESTERROR + ": " + errorMessage + " [some attribute: some value, some typed attribute: (int)42]"; err.Error() != expected {
		t.Error("Invalid error formatting (result, expected)\n", err.Error()+"\n", expected)
	}

	err.InContext(contextMessage)
	if expected := E_TESTERROR + ": " + errorMessage + " (" + contextMessage + ") [some attribute: some value, some typed attribute: (int)42]"; err.Error() != expected {
		t.Error("Invalid error formatting (result, expected)\n", err.Error()+"\n", expected)
	}

	eerr := From(err.Error())
	if fmt.Sprint(eerr.Map()) != fmt.Sprint(err.Map()) {
		t.Error("Bad parsing (result, expected)\n", err.Error()+"\n", eerr.Error())
	}
}
