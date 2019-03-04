package eerror

import (
	"fmt"
	"reflect"
)

const E_EXTERNALERROR = "E_EXTERNALERROR"

var errorParsedAttributes = []interface{}{
	"_eerror_parsed", true,
}

/*
From takes any parameter to convert it as an enhanced error.
Returns the given parameter if it's already an enhanced error instance, or nil
*/
func From(e interface{}) Eerror {
	if eerr, ok := e.(Eerror); ok {
		return eerr
	}
	if eerr, ok := e.(*Eerror); ok {
		return *eerr
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
func (e *Eerror) Is(instance interface{}, log ...interface{}) bool {
	var initial = e.getInitialError()
	var instanceInitial interface{} = instance

	if instanceEerr, ok := instance.(*Eerror); ok {
		instanceInitial = instanceEerr.getInitialError()
	} else if instanceEerr, ok := instance.(Eerror); ok {
		instanceInitial = instanceEerr.getInitialError()
	}

	if len(log) > 0 {
		fmt.Printf("%s %p %p\n", e, initial, e.parent)
		fmt.Printf("%s %p %p %s\n", instance, instance, instanceInitial, instanceInitial)
		fmt.Println(reflect.TypeOf(initial))
		fmt.Println(reflect.TypeOf(instanceInitial))
	}

	testInstanceID := func(toEerr Eerror) bool {
		if withEerr, ok := initial.(Eerror); ok {
			if toEerr._instance == withEerr._instance {
				return true
			}
		}
		if toEerr._instance == e._instance {
			return true
		}
		return false
	}

	if toEerr, ok := instanceInitial.(Eerror); ok {
		if testInstanceID(toEerr) {
			return true
		}
	}
	if toEerr, ok := instanceInitial.(*Eerror); ok {
		if testInstanceID(*toEerr) {
			return true
		}
	}
	return initial == instanceInitial
}

// Dup ensures a copy of a given enhanced error, reinstanciating contexts and attributes
func (e Eerror) Dup() Eerror {
	err := Eerror{
		e.parent,

		e.identifier,
		e.message,
		make([]string, len(e.contexts)),
		make(map[string]interface{}, len(e.attributes)),
		e._instance,
	}

	copy(err.contexts, e.contexts)
	for key, value := range e.attributes {
		err.attributes[key] = value
	}

	return err
}

func (e *Eerror) getInitialError() interface{} {
	if parent, ok := e.parent.(*Eerror); ok {
		return parent.getInitialError()
	} else if parent, ok := e.parent.(Eerror); ok {
		return parent.getInitialError()
	}

	if e.parent != nil {
		return *(e.parent.(*interface{}))
	}
	return e
}

func fromError(err *interface{}) Eerror {
	if eerr, ok := parse(*err); ok {
		return eerr
	}

	eerr := Eerror{
		err,

		E_EXTERNALERROR,
		fmt.Sprint(*err),
		[]string{},
		make(map[string]interface{}, len(errorParsedAttributes)/2),
		generateUniqueID(),
	}

	eerr.WithAttributes(
		errorParsedAttributes...,
	)
	return eerr
}
