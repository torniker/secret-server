package app

import (
	"fmt"
	"net/http"
	"secret/log"
)

// Error checks error type and responses accordingly
func (c *Ctx) Error(err error) {
	switch err.(type) {
	case ErrorBadRequest:
		e := err.(ErrorBadRequest)
		c.Respond(http.StatusBadRequest, e)
	case ErrorStatusNotFound:
		e := err.(ErrorStatusNotFound)
		c.Respond(http.StatusNotFound, e)
	case ErrorInternalServerError:
		e := err.(ErrorInternalServerError)
		c.Respond(http.StatusInternalServerError, e)
	default:
		c.Respond(http.StatusInternalServerError, ErrorInternalServerError{Message: err.Error()})
	}
}

// NotFound response 404
func (ctx *Ctx) NotFound() error {
	e := ErrorStatusNotFound{
		Message:  "not found",
		Internal: fmt.Sprintf("url: %v not found", ctx.Req.URL.Path),
	}
	log.Caller(2).WithError(e).Warn("not found")
	return e
}

// InternalServerError response 500
func (ctx *Ctx) InternalServerError(err error) error {
	e := ErrorInternalServerError{
		Message:  "internal server error",
		Internal: err.Error(),
	}
	log.Caller(2).WithError(e).Error("internal server error")
	return e
}

// BadRequest response 400
func (ctx *Ctx) BadRequest(err error) error {
	log.Caller(2).WithError(err).Warn("bad request")
	return err
}

// ErrorBadRequest type for bad request
type ErrorBadRequest struct {
	Message  string `json:"message" xml:"message"`
	Internal string `json:"-" xml:"-"`
}

func (e ErrorBadRequest) Error() string {
	return e.Message
}

// ErrorStatusNotFound type for not found
type ErrorStatusNotFound struct {
	Message  string `json:"message" xml:"message"`
	Internal string `json:"-" xml:"-"`
}

func (e ErrorStatusNotFound) Error() string {
	return e.Message
}

// ErrorInternalServerError type for internal server error
type ErrorInternalServerError struct {
	Message  string `json:"message" xml:"message"`
	Internal string `json:"-" xml:"-"`
}

func (e ErrorInternalServerError) Error() string {
	return e.Message
}
