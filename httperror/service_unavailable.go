package httperror

import "net/http"

// ServiceUnavailable is a constructor to create ServiceUnavailableError instance
func ServiceUnavailable(err error) error {
	return New(http.StatusServiceUnavailable, err)
}

// IsServiceUnavailableError check whether given error is a ServiceUnavailableError
func IsServiceUnavailableError(err error) bool {
	return GetInstance(err).Code == http.StatusServiceUnavailable
}