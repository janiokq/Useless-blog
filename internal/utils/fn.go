package utils

type vFn func() error

func V(fns ...vFn) error {
	for _, fn := range fns {
		err := fn()
		if err != nil {
			return err
		}
	}
	return nil
}
