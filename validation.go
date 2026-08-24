package main

import "fmt"

func validateInspection(value string) error {
	for _, allowed := range []string{"clear", "needs_attention", "cleared"} {
		if value == allowed {
			return nil
		}
	}
	return fmt.Errorf("unsupported inspection %q", value)
}
