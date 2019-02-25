package eerror

import "testing"

const E_TEST_ERROR = "E_TEST_ERROR"
const E_TEST_ERROR_WITH_ATTRIBUTES = "E_TEST_ERROR_WITH_ATTRIBUTES"

func TestSimpleError(t *testing.T) {
	const errorMessage = "This is a test error"

	err := NewError(E_TEST_ERROR, errorMessage)
	if err.Error() != E_TEST_ERROR+": "+errorMessage {
		t.Error("Invalid message")
		return
	}
}

func TestEnhancedError(t *testing.T) {
	const errorMessage = "This is a test enhanced error"

	err := NewError(E_TEST_ERROR, errorMessage)
	if err.Error() != E_TEST_ERROR+": "+errorMessage {
		t.Error("Invalid basic message")
		return
	}

	const subcontextMessage = "Some kind of test subcontext"
	err.InContext(subcontextMessage)
	if err.Error() != E_TEST_ERROR+": "+errorMessage+"\n ->"+subcontextMessage {
		t.Error("Invalid enhanced message with one context", err.Error())
		return
	}

	const contextMessage = "Text sup-context"
	err.InContext(contextMessage)
	if err.Error() != E_TEST_ERROR+": "+errorMessage+"\n ->"+contextMessage+"\n ->"+subcontextMessage {
		t.Error("Invalid enhanced message with two contexts")
		return
	}
}
