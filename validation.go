package main

import "fmt"

func validateInspection(value string) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("inspection value rejected")
		}
	}()
	for _, allowed := range []string{"clear", "needs_attention", "cleared"} {
		if value == allowed {
			return nil
		}
	}
	return fmt.Errorf("unsupported inspection %q", value)
}
