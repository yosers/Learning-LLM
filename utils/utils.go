// utils/db_error.go
package utils

import (
	"context"
	"strings"
)

func IsDBDown(err error) bool {
	if err == nil {
		return false
	}
	if strings.Contains(err.Error(), "connection refused") ||
		strings.Contains(err.Error(), "connection reset") ||
		strings.Contains(err.Error(), "no such host") ||
		strings.Contains(err.Error(), "i/o timeout") ||
		strings.Contains(err.Error(), "failed to connect") ||
		context.DeadlineExceeded == err {
		return true
	}
	return false
}

func StringOrEmpty(s *string) string {
	if s == nil || *s == "" {
		return ""
	}
	return *s
}
