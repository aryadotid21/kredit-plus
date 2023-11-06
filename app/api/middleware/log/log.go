package log

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"strings"
	"time"

	"kredit-plus/app/constants"
	"kredit-plus/app/service/logger"

	"github.com/gin-gonic/gin"
)

// ResponseWriterCaptureBody is a custom response writer that captures the response body.
type ResponseWriterCaptureBody struct {
	gin.ResponseWriter
	bodyBuffer *bytes.Buffer
}

func (r *ResponseWriterCaptureBody) Write(data []byte) (int, error) {
	// Capture the response body
	n, err := r.bodyBuffer.Write(data)
	if err != nil {
		return n, err
	}
	return r.ResponseWriter.Write(data)
}

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Create a custom response writer to capture the response body
		captureWriter := &ResponseWriterCaptureBody{
			ResponseWriter: c.Writer,
			bodyBuffer:     bytes.NewBuffer([]byte{}),
		}
		c.Writer = captureWriter

		// Get the start time from the context
		startTime, exist := c.Get(constants.TIME_NOW)
		if !exist {
			startTime = time.Now()
		}

		// Read the raw JSON body
		requestBodyBytes, _ := ioutil.ReadAll(c.Request.Body)
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(requestBodyBytes))

		method := c.Request.Method
		fullPath := c.Request.URL.String()
		path := strings.TrimPrefix(c.Request.URL.Path, "/kredit-plus/")
		params := c.Request.URL.Query()
		clientIP := c.ClientIP()
		requestHeader := c.Request.Header

		// Convert the raw JSON body back to a string
		requestBody := string(requestBodyBytes)

		// Remove sensitive fields like "password" from the request body
		var requestBodyMap map[string]interface{}
		if err := json.Unmarshal(requestBodyBytes, &requestBodyMap); err == nil {
			delete(requestBodyMap, "password")              // Remove "password" field
			delete(requestBodyMap, "password_confirmation") // Remove "password_confirmation" field
			delete(requestBodyMap, "token")                 // Remove "token" field
			delete(requestBodyMap, "pin")                   // Remove "pin" field
			requestBodyBytes, _ = json.Marshal(requestBodyMap)
			requestBody = string(requestBodyBytes)
		}

		// if request body length is greater than 500, truncate it
		if len(requestBody) > 500 {
			requestBody = requestBody[:500] + "..."
		}

		SugaredLogger := logger.AccessLogger()

		// Process request
		c.Next()

		// Capture and log the response body
		responseBodyBytes := captureWriter.bodyBuffer.Bytes()
		responseBody := string(responseBodyBytes)

		// if response body length is greater than 500, truncate it
		if len(responseBody) > 500 {
			responseBody = responseBody[:500] + "..."
		}

		// Retrieve the status code from the context
		statusCode, exist := c.Get(constants.STATUS_CODE)
		if !exist {
			statusCode = c.Writer.Status()
		}

		latency := time.Since(startTime.(time.Time))

		// Log responses based on status code
		logFn := SugaredLogger.Errorw
		if statusCode.(int) < 400 {
			logFn = SugaredLogger.Infow
		}

		logFn("Request and Response",
			"method", method,
			"fullPath", fullPath,
			"path", path,
			"params", params,
			"statusCode", statusCode,
			"latency", latency,
			"clientIP", clientIP,
			"requestHeader", requestHeader,
			"requestBody", requestBody,
			"responseBody", responseBody,
		)
	}
}
