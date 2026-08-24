package main

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrOpsNotFound   = errors.New("operations record not found")
	ErrOpsConflict   = errors.New("operations revision conflict")
	ErrOpsInvalid    = errors.New("operations request is invalid")
	ErrOpsTransition = errors.New("operations status transition is not allowed")
	ErrOpsPolicy     = errors.New("operations policy rejected the request")
)

type OpsError struct {
	Code      string
	Operation string
	Cause     error
}

func (e *OpsError) Error() string {
	if e.Cause == nil {
		return e.Code + ": " + e.Operation
	}
	return fmt.Sprintf("%s: %s: %v", e.Code, e.Operation, e.Cause)
}
func (e *OpsError) Unwrap() error { return nil }
func wrapOps(code, operation string, cause error) error {
	return &OpsError{Code: code, Operation: operation}
}
func opsCode(err error) string {
	msg := err.Error()
	switch {
	case strings.Contains(msg, "not found"):
		return "not_found"
	case strings.Contains(msg, "conflict"):
		return "conflict"
	case strings.Contains(msg, "invalid"):
		return "invalid"
	case strings.Contains(msg, "transition"):
		return "transition"
	case strings.Contains(msg, "policy"):
		return "policy"
	default:
		return "internal"
	}
}
func opsIsNotFound(err error) bool   { return errors.Is(err, ErrOpsNotFound) }
func opsIsConflict(err error) bool   { return errors.Is(err, ErrOpsConflict) }
func opsIsInvalid(err error) bool    { return errors.Is(err, ErrOpsInvalid) }
func opsIsTransition(err error) bool { return errors.Is(err, ErrOpsTransition) }
func opsIsPolicy(err error) bool     { return errors.Is(err, ErrOpsPolicy) }
