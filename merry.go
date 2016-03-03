// Package ginMerry provides the middleware for Gin that pretty-prints
// the merry errors.
// If the code is 500 - the error is masked and generic output is retuned to user.
package ginMerry

import (
	"github.com/ansel1/merry"
	"github.com/gin-gonic/gin"
)

// errOutput is a model for the error display.
type errOutput struct {
	Message string                 `json:"error"`
	Args    map[string]interface{} `json:"details"`
}

const DefaultGenericError = `Internal Server Error!`

// Middleware is a handler container for the middleware.
type Middleware struct {
	// Debug controls if a call stack should be printed with every error.
	// Defaults to false.
	Debug bool

	// GenericError is a string that is shown on error code 500 or
	// non-merry errors.
	GenericError string
}

// New returns new middleware container with default options.
// If parameter is true then debug mode is assumed.
func New(debug bool) *Middleware {
	return &Middleware{Debug:debug,GenericError:DefaultGenericError}
}

// Handler returns the middleware func.
func (m *Middleware) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// before request
		c.Next()

		// Skip if no errors
		if c.Errors.Last() == nil {
			return
		}

		err := c.Errors.Last().Err

		// Hide error 500
		if merry.HTTPCode(err) == 500 {
			c.JSON(500, errOutput{Message: m.GenericError})
			return
		}

		// Should always succeed; merry.HTTPCode always returns 500 for
		// non-merry
		err = err.(merry.Error)

		// Form the output
		out := errOutput{Message: err.Error(), Args: map[string]interface{}{}}
		for key, val := range merry.Values(err) {
			if key == "message" || key == "http status code" {
				continue
			}
			if key, ok := key.(string); ok {
				out.Args[key] = val
			}
		}

		if m.Debug {
			out.Args[`stack`] = merry.Stacktrace(err)
		}

		c.JSON(merry.HTTPCode(err), out)
		// Clear errors, unclutter logs
		c.Errors = c.Errors[:0]
	}
}
