package errx

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

// Errx : all the libraries would be using this interface for sending out errors
type Errx interface {
	error
	HTTPStatusCode() int
	UserMessage() string
	Log()
}

// ErrJSONBind : error binding request body from json to object
type ErrJSONBind struct {
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

// ErrNotFound : error binding request body from json to object
type ErrNotFound struct {
	*errxNotFound
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

// ErrLogin : When the user credentials fail, and do not match that in the database
// typically after this tokens would not be generated
type ErrLogin struct {
	*errxUnatuho
}

// ErrInsuffPrivlg : when the role of the user disallows certain actions on the api
type ErrInsuffPrivlg struct {
	*errxForbid
}

// ErrConnFailed : failed database or cache connection
type ErrConnFailed struct {
	*errxServUnavail
}

// NewErr : generates a new custom error
// given the type of the error this can construct a new error of the desired type
// also enwraps the inner error within
// to be used from deep libraries underlying
func NewErr(t interface{}, e error, m, ct string) Errx {
	u := uuid.NewString()[24:] //32 bit id with 4 -, we want the last 12 unique number to identify the error
	badrequest := &errxBadRequest{&errx{UMsg: m, Ctx: ct, InnerErr: e, uid: u}}
	badgtway := &errxGateway{&errx{UMsg: m, Ctx: ct, InnerErr: e, uid: u}}
	intsrv := &errxIntServer{&errx{UMsg: m, Ctx: ct, InnerErr: e, uid: u}}
	unauth := &errxUnatuho{&errx{UMsg: m, Ctx: ct, InnerErr: e, uid: u}}
	forbid := &errxForbid{&errx{UMsg: m, Ctx: ct, InnerErr: e, uid: u}}
	notfnd := &errxNotFound{&errx{UMsg: m, Ctx: ct, InnerErr: e, uid: u}}
	srvunav := &errxServUnavail{&errx{UMsg: m, Ctx: ct, InnerErr: e, uid: u}}
	switch t.(type) {
	case *ErrJSONBind:
		return &ErrJSONBind{badrequest}
	case *ErrNotFound:
		return &ErrNotFound{notfnd}
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
		return &ErrInsuffPrivlg{forbid}
	case *ErrLogin:
		return &ErrLogin{unauth}
	case *ErrConnFailed:
		return &ErrConnFailed{srvunav}
	}
	return nil
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
		c.AbortWithStatusJSON(x.HTTPStatusCode(), gin.H{"message": x.UserMessage()})
		// c.AbortWithError(x.HTTPStatusCode(), fmt.Errorf(x.UserMessage()))
		return 1
	}
	// when the error object does not implement Errx interface
	log.Error(err)
	c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
	// c.AbortWithError(http.StatusInternalServerError, err)
	return 1
}
