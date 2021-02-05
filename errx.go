package errx

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// Errx : all the libraries would be using this interface for sending out errors
type Errx interface {
	error
	HTTPStatusCode() int
	UserMessage() string
	Log()
}

// DigestErr : for any API route handler this helps to set the gin context and indicate if the error has occured
// this can be coupled with any function returning an error
// send back an integer indicating if the API needs to return or continue with further actions
func DigestErr(err error, c *gin.Context) int {
	if err == nil {
		return 0
	}
	x, ok := err.(Errx)
	if x != nil && ok {
		x.Log()
		c.AbortWithError(x.HTTPStatusCode(), fmt.Errorf(x.UserMessage()))
		return 1
	}
	// when the error object does not implement Errx interface
	log.Error(err)
	c.AbortWithError(http.StatusInternalServerError, err)
	return 1

}

type errx struct {
	UMsg     string // User Message
	Ctx      string
	InnerErr error // internal to 3rd party libraries, cannot be exposed to front end
}

// HTTPStatusCode : error type to a appropriate status code
func (e *errx) HTTPStatusCode() int {
	return http.StatusNotImplemented
}

// UserMessage :gets the user message approriate for front end user consumption
func (e *errx) UserMessage() string {
	return e.UMsg
}

func (e *errx) Error() string {
	return fmt.Sprintf("%s: %s-%s", e.Ctx, e.UMsg, e.InnerErr)
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

// ErrJSONBind : error binding request body from json to object
type ErrJSONBind struct {
	*errxBadRequest
}

// ErrNotFound : error binding request body from json to object
type ErrNotFound struct {
	*errxBadRequest
}

// ErrDuplicate : error binding request body from json to object
type ErrDuplicate struct {
	*errxBadRequest
}

// ErrInvalid : error binding request body from json to object
type ErrInvalid struct {
	*errxBadRequest
}

// ErrQuery : error binding request body from json to object
type ErrQuery struct {
	*errxGateway
}

// ErrCacheQuery : error binding request body from json to object
type ErrCacheQuery struct {
	*errxGateway
}

// ErrEncrypt : Error when encryption fails, typically happens on user account passwords
type ErrEncrypt struct {
	*errxIntServer
}

// ErrTokenExpired : jwt has expired, can be auth or refr ..
type ErrTokenExpired struct {
	*errxUnatuho
}

// ErrInsuffPrivlg : when the role of the user disallows certain actions on the api
type ErrInsuffPrivlg struct {
	*errxUnatuho
}

// ErrLogin : When the user credentials fail, and do not match that in the database
// typically after this tokens would not be generated
type ErrLogin struct {
	*errxUnatuho
}

// NewErr : generates a new custom error
// given the type of the error this can construct a new error of the desired type
// also enwraps the inner error within
// to be used from deep libraries underlying
func NewErr(t interface{}, e error, m, ct string) Errx {
	badrequest := &errxBadRequest{&errx{UMsg: m, Ctx: ct, InnerErr: e}}
	badgtway := &errxGateway{&errx{UMsg: m, Ctx: ct, InnerErr: e}}
	intsrv := &errxIntServer{&errx{UMsg: m, Ctx: ct, InnerErr: e}}
	unauth := &errxUnatuho{&errx{UMsg: m, Ctx: ct, InnerErr: e}}
	switch t.(type) {
	case *ErrJSONBind:
		return &ErrJSONBind{badrequest}
	case *ErrNotFound:
		return &ErrNotFound{badrequest}
	case *ErrDuplicate:
		return &ErrDuplicate{badrequest}
	case *ErrInvalid:
		return &ErrInvalid{badrequest}
	case *ErrQuery:
		return &ErrQuery{badgtway}
	case *ErrCacheQuery:
		return &ErrCacheQuery{badgtway}
	case *ErrEncrypt:
		return &ErrEncrypt{intsrv}
	case *ErrTokenExpired:
		return &ErrTokenExpired{unauth}
	case *ErrInsuffPrivlg:
		return &ErrInsuffPrivlg{unauth}
	case *ErrLogin:
		return &ErrLogin{unauth}
	}
	return nil
}
