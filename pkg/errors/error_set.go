package errors

type ErrorSet []error

// Err returns the error set as an error if error set length is greater then 0.
// Otherwise, it returns nil.
func (e ErrorSet) Err() error {
	if len(e) > 0 {
		return e
	}

	return nil
}

func (e ErrorSet) Error() string {
	var s string
	for _, err := range e {
		s += err.Error() + "\n"
	}
	return s
}
