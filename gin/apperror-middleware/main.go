package main

import ( 	
	"github.com/gin-gonic/gin"
)

//
// gin HTTP Server setup
//
func main() {
    router := gin.Default()
    router.Use(JSONAppErrorReporter())
    router.GET("/hostgroups/:id", fetchSingleHostGroup)
    router.Run(":3000")
}


//
// APP error definition
//
type appError struct {
    Code     int    `json:"code"`
    Message  string `json:"message"`
}


//
//  Report Error in app
//
func fetchSingleHostGroup(c *gin.Context) {
    hostgroupID := c.Param("id")                                                             
    // simulating an error
    err := appError{400, "Invalid argument. HostgroupID must be set!"}

    if err != nil {
        // put the Error to gin.Context.Errors
        c.Error(err)
        return
    }
    // return data of OK
    c.JSON(http.StatusOK, *hostGroupRes)
}

//
// Middleware Error Handler in server package
//
// Intercept the request, checks the context `Errors` field.
// Differentiates application errors `appError` from non-app type errors 
// then iterrupts the request and returns the error as JSON.  
func JSONAppErrorReporter() gin.HandlerFunc {
    return jsonAppErrorReporterT(gin.ErrorTypeAny)
}

func jsonAppErrorReporterT(errType gin.ErrorType) gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Next()
        detectedErrors := c.Errors.ByType(errType)
        // log.Println("Handle APP error")
        if len(detectedErrors) > 0 {
            err := detectedErrors[0].Err
            var parsedError *appError
            switch err.(type) {
            case *appError:
                parsedError = err.(*appError )
            default:
                parsedError = &appError{ 
                  code: http.StatusInternalServerError,
                  message: "Internal Server Error",
                }
            }
            // Put the error into response
            c.IndentedJSON(parsedError.Code, parsedError)
            c.Abort()
            // or c.AbortWithStatusJSON(parsedError.Code, parsedError)
            return
        }
    }
}