package gowith

// With interface to match Python's enter/exit.
type EnterExiter[T any] interface {
	// On enter, return enter value (T) to pass to callable function.
	Enter() (T, error)

	// On exit, accept the enter value (T), return error (if any).
	Exit(T, error) error
}

// Accept an implementation of EnterExiter (ee).
// Accept a function which fires between Enter and Exit.
// Return an error (if any).
func New[T any](ee EnterExiter[T], cf func(T) error) error {
	var err error

	// Perform the enter function.
	ent, err := ee.Enter()
	if err == nil {
		// No errors yet.
		// Perform the callable function and pass in the enter value.
		err = cf(ent)
	}

	// Perform the exit function, pass in the enter value and any error.
	return ee.Exit(ent, err)
}
