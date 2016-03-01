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

// Handler returns the middleware func.
// Argument is a generic error string that is displayed
// if the error is internal (code 500 or not merry).
func Handler(generr string) gin.HandlerFunc {
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
			c.JSON(500, errOutput{Message: generr})
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

		c.JSON(merry.HTTPCode(err), out)
		// Clear errors, unclutter logs
		c.Errors = c.Errors[:0]
	}
}
