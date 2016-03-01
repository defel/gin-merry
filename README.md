[![forthebadge](http://forthebadge.com/images/badges/designed-in-ms-paint.svg)](http://forthebadge.com)

# gin-merry [![GoDoc](https://godoc.org/github.com/utrack/gin-merry?status.svg)](https://godoc.org/github.com/utrack/gin-merry)
Middleware that marries merry errors and Gin. 

It pretty-prints merry errors to the user with all the context embedded in the error.

This middleware is compatible with Golang's [Gin](https://github.com/gin-gonic/gin) HTTP router and [merry errors](https://github.com/ansel1/merry) with context™.

After enabling the middleware, if the handler returns an error to the gin.Context, it will be printed to the user with all the additional context that came with the error. The errors' queue is cleared, so the logs won't be cluttered with the useless errors.

However, if the error has code 500 - then the error is considered bad/not merry at all; some default text is printed to the user and the error is passed down the chain for logging.
