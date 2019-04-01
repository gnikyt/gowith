package gowith

// Simple return value from Enter function to accept any type.
type EnterReturn struct {
	Value interface{}
}

// With interface to match Python's enter/exit.
type EnterExiter interface {
	// On enter, you must return a pointer to EnterReturn and error value (if any).
	Enter() (*EnterReturn, error)

	// On exit, pointer of EnterReturn from Enter function is passed in, along with any error (if any).
	// You must return the error passed in as the second argument (if any).
	Exit(er *EnterReturn, err error) error
}

// Main function.
// Accept anything conforming to EnterExiter as first argument.
// Accept an anonymous function which fires betwen Enter and Exit.
// Return an error (if any happened).
func New(ee EnterExiter, act func(er *EnterReturn) error) error {
	var err error

	// Perform the enter function.
	ent, err := ee.Enter()

	if err == nil {
		// No errors yet, perform the action function and pass in the enter value.
		err = act(ent)
	}

	// Perform the exit function, pass in the enter value and any error.
	return ee.Exit(ent, err)
}
