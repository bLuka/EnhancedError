package eerror

import "fmt"

const E_EXTERNALERROR = "E_EXTERNALERROR"

var errorParsedAttributes = []interface{}{
	"_eerror_parsed", true,
}

/*
From takes any parameter to convert it as an enhanced error.
Returns the given parameter if it's already an enhanced error instance, or nil
*/
func From(e interface{}) Eerror {
	if e == nil {
		return nil
	}

	if eerr, ok := e.(*eerror); ok {
		return eerr
	}
	if ptr, ok := e.(*interface{}); ok {
		return fromError(ptr)
	}

	return fromError(&e)
}

/*
Is tests relationship between an argument and an enhanced error instance, for error handling.
Useful to test if an enhanced error instance was formed from the given instance parameter

  const E_MY_ERROR_ID = "E_MY_ERROR_ID"

  var standardError = eerror.NewError(E_MY_ERROR_ID, "Some error")

  func errorFunction(myParameter interface{}) eerror.Eerror {
     err := standardError.Copy()
     err.WithAttribute("parameter", myParameter)

     return err
  }

  func main() {
     if eerr := errorFunction("hello world"); !eerr.Is(standardError) {
        panic(eerr)
     }
  }
*/
func (e *eerror) Is(instance interface{}) bool {
	var initial = e.getInitialError()
	var instanceInitial interface{} = instance

	if instanceEerr, ok := instance.(*eerror); ok {
		instanceInitial = instanceEerr.getInitialError()
	}

	return initial == instanceInitial
}

func (e eerror) Copy() Eerror {
	err := &eerror{
		e.parent,

		e.identifier,
		e.message,
		make([]string, len(e.contexts)),
		make(map[string]interface{}, len(e.attributes)),
	}

	copy(err.contexts, e.contexts)
	for key, value := range e.attributes {
		err.attributes[key] = value
	}

	return err
}

func (e *eerror) getInitialError() interface{} {
	if parent, ok := e.parent.(*eerror); ok {
		return parent.getInitialError()
	}

	if e.parent != nil {
		return *(e.parent.(*interface{}))
	}
	return e
}

func fromError(err *interface{}) Eerror {
	eerr := eerror{
		err,

		E_EXTERNALERROR,
		fmt.Sprint(*err),
		[]string{},
		make(map[string]interface{}, len(errorParsedAttributes)/2),
	}

	eerr.WithAttributes(
		errorParsedAttributes...,
	)
	return &eerr
}
