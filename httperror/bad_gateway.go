package httperror

import "net/http"

// BadGateway is a constructor to create BadGatewayError instance
func BadGateway(err error) error {
	return New(http.StatusBadGateway, err)
}

// IsBadGatewayError check whether given error is a BadGatewayError
func IsBadGatewayError(err error) bool {
	return GetInstance(err).Code == http.StatusBadGateway
}