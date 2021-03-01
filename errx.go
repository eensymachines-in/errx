package errx

import (
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type errx struct {
	UMsg     string // User Message
	Ctx      string
	InnerErr error // internal to 3rd party libraries, cannot be exposed to front end
	uid      string
}

// HTTPStatusCode : error type to a appropriate status code
func (e *errx) HTTPStatusCode() int {
	return http.StatusNotImplemented
}

// UserMessage :gets the user message approriate for front end user consumption
func (e *errx) UserMessage() string {
	// will be shown on the ui
	// the user can quote this number
	return fmt.Sprintf("%s\n%s", e.UMsg, e.uid)
}

func (e *errx) Error() string {
	return fmt.Sprintf("%s:%s: %s-%s", e.uid, e.Ctx, e.UMsg, e.InnerErr)
}

// Log : logs the error to whatever tty is assigned
func (e *errx) Log() {
	log.Error(e.Error())
}

type errxBadRequest struct {
	*errx
}

// HTTPStatusCode : json binding errors are often the result of body of the request being disfugured
func (errbdr *errxBadRequest) HTTPStatusCode() int {
	return http.StatusBadRequest
}

type errxGateway struct {
	*errx
}

// HTTPStatusCode : json binding errors are often the result of body of the request being disfugured
func (errgtwy *errxGateway) HTTPStatusCode() int {
	return http.StatusBadGateway
}

type errxServUnavail struct {
	*errx
}

// HTTPStatusCode : json binding errors are often the result of body of the request being disfugured
func (errsa *errxServUnavail) HTTPStatusCode() int {
	return http.StatusServiceUnavailable
}

type errxIntServer struct {
	*errx
}

// HTTPStatusCode : json binding errors are often the result of body of the request being disfugured
func (errints *errxIntServer) HTTPStatusCode() int {
	return http.StatusInternalServerError
}

type errxUnatuho struct {
	*errx
}

// HTTPStatusCode : json binding errors are often the result of body of the request being disfugured
func (errunau *errxUnatuho) HTTPStatusCode() int {
	return http.StatusUnauthorized
}

type errxNotFound struct {
	*errx
}

func (enf *errxNotFound) HTTPStatusCode() int {
	return http.StatusNotFound
}

type errxForbid struct {
	*errx
}

// HTTPStatusCode : json binding errors are often the result of body of the request being disfugured
func (errfbd *errxForbid) HTTPStatusCode() int {
	return http.StatusForbidden
}
