package eerror

import (
	"fmt"

	"testing"
)

const E_TESTERROR = "E_TESTERROR"
const E_TESTERROR_WITH_ATTRIBUTES = "E_TESTERROR_WITH_ATTRIBUTES"

// TestSimpleError tests the NewError constructor
func TestSimpleError(t *testing.T) {
	const errorMessage = "This is a test error"

	err := NewError(E_TESTERROR, errorMessage)
	if err.Error() != E_TESTERROR+": "+errorMessage {
		t.Error("Invalid message")
	}
}

// TestEnhancedError ensures context capacities from a new error
func TestEnhancedError(t *testing.T) {
	const errorMessage = "This is a test enhanced error"

	err := NewError(E_TESTERROR, errorMessage)
	if expect := E_TESTERROR + ": " + errorMessage; err.Error() != expect {
		t.Error("Invalid basic message\n", err, " | Expected:\n", expect)
	}

	const subcontextMessage = "Some kind of test subcontext"
	err.InContext(subcontextMessage)
	if expect := E_TESTERROR + ": " + errorMessage + " (" + subcontextMessage + ")"; err.Error() != expect {
		t.Error("Invalid enhanced message with one context\n", err, " | Expected:\n", expect)
	}

	const contextMessage = "Text sup-context"
	err.InContext(contextMessage)
	if expect := E_TESTERROR + ": " + errorMessage + " (" + subcontextMessage + "; " + contextMessage + ")"; err.Error() != expect {
		t.Error("Invalid enhanced message with two contexts\n", err, " | Expected:\n", expect)
	}
}

// TestFromIs tests the ability to parse wichever given type into an enhanced error, with the ability to test further inheritance
func TestFromIs(t *testing.T) {
	const errorMessage = "New standard error initialization"
	var stdError = fmt.Errorf(errorMessage)

	eerr := From(stdError)
	if !eerr.Is(stdError) {
		t.Error("Enhanced error Is() method should return true after being formed by standard error")
	}

	const contextMessage = "Test context"
	eerr.InContext(contextMessage)
	if !eerr.Is(stdError) {
		t.Error("Enhanced error Is() method should still return true after context update")
	}

	nextEerr := From(eerr)
	if !nextEerr.Is(stdError) {
		t.Error("Enhanced error Is() method should return true as indirectly formed by the given standard error instance")
	}
	if !nextEerr.Is(eerr) {
		t.Error("Enhanced error Is() method should return true as indirectly formed by the same enhanced error")
	}
	if !eerr.Is(nextEerr) {
		t.Error("Enhanced error Is() method should return true as the given argument is formed from the instance")
	}

	lastEerr := From(nextEerr)
	if !lastEerr.Is(stdError) {
		t.Error("Enhanced error Is() method should return true as indirectly formed by the given standard error instance")
	}
	if !eerr.Is(lastEerr) {
		t.Error("Enhanced error Is() method should return true as the given argument is indirectly formed by this enhanced error instance")
	}

	var pureEerr = NewError(E_TESTERROR, errorMessage)
	if pureEerr.Is(stdError) || pureEerr.Is(lastEerr) {
		t.Error("Enhanced error Is() method should return false as the enhanced error was freshly instanciated")
	}

	nextEerr = From(pureEerr)
	lastEerr = From(nextEerr)

	if !pureEerr.Is(lastEerr) || !lastEerr.Is(pureEerr) {
		t.Error("Enhanced error Is() method should return true as the last enhanced error is indirectly formed by parsing the first enhanced error instanciation")
	}

	var nextStdError = fmt.Errorf(errorMessage)
	if eerr.Is(nextStdError) || nextEerr.Is(nextStdError) {
		t.Error("Enhanced error Is() method should return false since wasn't formed, directly or not, by this new standard error instance")
	}
}

// TestAttribute ensures attributes capacity on enhanced errors
func TestAttribute(t *testing.T) {
	const errorMessage = "This is a test error"

	err := NewError(E_TESTERROR, errorMessage)
	if len(err.GetAttributes()) > 0 {
		t.Error("Enhanced error shouldn't possess any attribute after creation\n", err.GetAttributes())
	}

	secondErr := err.Copy()
	secondErr.WithAttribute("attribute", "Hello world")
	if len(err.GetAttributes()) > 0 {
		t.Error("First enhanced error shouldn't possess any attribute\n", err.GetAttributes())
	}
	if len(secondErr.GetAttributes()) != 1 || secondErr.GetAttributes()["attribute"].(string) != "Hello world" {
		t.Error("Second enhanced error should possess a single attribute labeled \"attribute\" containing an \"Hello world\" string\n", secondErr.GetAttributes())
	}

	thirdErr := err.Copy()
	thirdErr.WithAttribute("attribute", "overwritten")
	if len(secondErr.GetAttributes()) != 1 || secondErr.GetAttributes()["attribute"].(string) != "Hello world" {
		t.Error("Third enhanced error should still possess a single attribute labeled \"attribute\" containing an \"Hello world\" string\n", thirdErr.GetAttributes())
	}
	if len(thirdErr.GetAttributes()) != 1 || thirdErr.GetAttributes()["attribute"].(string) != "overwritten" {
		t.Error("Third enhanced error should possess a single attribute labeled \"attribute\" containing an \"overwritten\" string\n", thirdErr.GetAttributes())
	}
}
