package util

// ReturnFirstErr returns the first not-nil error from a list of errors.
// Used to prevent many copies of if err != nil { return err } when they are indepedent
func ReturnFirstErr(errs ...error) error {
	for _, err := range errs {
		if err != nil {
			return err
		}
	}
	return nil
}
