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
		t.Error("Bad parsing, both should be equals (result, expected)\n", err.Error()+"\n", eerr.Error())
	}

	failedTest, failedError, invalidTestsFailed := (func() (test string, eerr Eerror, ok bool) {
		for _, str := range []string{
			"E_SOMEERROR:",
			"E_SOMEERROR: ",
			"E_SOMEERROR:\"\"",
			": \"\"",
			"\"\"",
			"E: \"",
			"E some type: \"",
			"\"E some type: \":",
			"E_SOMEERROR: message (",
			"E_SOMEERROR: message ()",
			"E_SOMEERROR: message (context,)",
			"E_SOMEERROR: message (,context)",
			"E_SOMEERROR: message (context) [",
			"E_SOMEERROR: message (context) ]",
			"E_SOMEERROR: message (context) string",
			"E_SOMEERROR: message (context",
			"E_SOMEERROR: message[ (context)",
			"E_SOMEERROR: message[] (context)",
			"E_SOMEERROR: message (context) []",
			"E_SOMEERROR: message [attribute: value] (context)",
			"E_SOMEERROR: message (context) [attribute: value]",
			"E_SOMEERROR: message (context) [attribute: value",
			"E_SOMEERROR: message (context) [attribute: \"value]",
			"E_SOMEERROR: message (context) [\"attribute: \"value]",
			"E_SOMEERROR: message (context) \"attribute: \": value]",
			"E_SOMEERROR: message [\"attribute: \" value]",
		} {
			if eerr, ok = parse(str); ok {
				return
			}
		}
		ok = false
		return
	})()
	if invalidTestsFailed {
		t.Error("Parsing test should've failed, but succeeded (test, built error)", failedTest+"\n", failedError)
	}

	failedTest, failedError, validTestsSuccess := (func() (test string, eerr Eerror, ok bool) {
		for _, str := range []string{
			"E_SOMEERROR: msg",
			"E_SOMEERROR: \"\"",
			"\":\": \"\"",
			"E                       : \"\"",
			"\t\n    E: \"\" \t\t\n",
			"E: m",
			"E some type: error message",
			"\"E some type: \": error message",
			"E_SOMEERROR: message ",
			"E_SOMEERROR: message (context)",
			"E_SOMEERROR: message (a,context)",
			"E_SOMEERROR: message (context, \"another, very long, context\")",
			"E_SOMEERROR: message (context, \"a (special) one\")",
			"E_SOMEERROR: message (context, another, more, contexts)",
			"E_SOMEERROR: message (the same unique context but with spaces)",
			"E_SOMEERROR: message (context) [attribute]",
			"E_SOMEERROR: message (context) [attribtue: and its value]",
			"E_SOMEERROR: message (context) [attribute: \"and, its value\"]",
			"E_SOMEERROR: message (context) [attribute: value, another attribute: value]",
			"E_SOMEERROR: message (context) [attribute: (int)-1]",
			"E_SOMEERROR: message (context) [attribute: \"(int)string value\"]",
			"E_SOMEERROR: message [attribute]",
		} {
			if eerr, ok = parse(str); !ok {
				return
			}
		}
		ok = true
		return
	})()
	if validTestsSuccess {
		t.Error("Parsing test should've succeed, but faild (test, built error)", failedTest+"\n", failedError)
	}
}
