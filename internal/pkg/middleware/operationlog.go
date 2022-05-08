package middleware

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"regexp"
	"time"

	"github.com/gin-gonic/gin"

	"gobackend/pkg/log"
	metav1 "gobackend/pkg/meta/v1"

	"gobackend/internal/app/apiserver/store"
	"gobackend/internal/pkg/entity/apiserver/operationlog"
)

var regPattern = regexp.MustCompile(`^/operation-logs*`)

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)

	return w.ResponseWriter.Write(b)
}

func (w bodyLogWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)

	return w.ResponseWriter.WriteString(s)
}

// OperationLog is a middleware function that logs operation.
func OperationLog(storeIns store.Factory) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Read operations do not need to be logged.
		if c.Request.Method == http.MethodGet ||
			c.Request.Method == http.MethodOptions {
			return
		}

		// Operation log api does not need to be logged.
		if regPattern.MatchString(c.Request.URL.Path) {
			return
		}

		bodyLogWriter := &bodyLogWriter{
			body:           bytes.NewBufferString(""),
			ResponseWriter: c.Writer,
		}
		c.Writer = bodyLogWriter

		startTime := time.Now()

		requestBody := ""
		if c.Request.Body != nil {
			bodyBytes, _ := ioutil.ReadAll(c.Request.Body)
			c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
			requestBody = string(bodyBytes)
		}

		c.Next()

		responseBody := bodyLogWriter.body.String()

		latencyTime := time.Since(startTime).Seconds()

		requestID := getRequestID(c)
		requestURI := getRequestURI(c)
		username := c.GetString(UsernameKey)

		clientIP := c.ClientIP()
		httpStatusCode := c.Writer.Status()
		requestReferer := c.Request.Referer()
		requestUA := c.Request.UserAgent()

		operationLog := &operationlog.OperationLog{
			Username:   username,
			ClientIP:   clientIP,
			ReqMethod:  c.Request.Method,
			ReqPath:    requestURI,
			ReqBody:    requestBody,
			ReqReferer: requestReferer,
			UserAgent:  requestUA,
			ReqTime:    startTime,
			ReqLatency: latencyTime,
			HTTPStatus: httpStatusCode,
			ResData:    responseBody,
		}

		go func() {
			if err := storeIns.OperationLogs().Create(
				c,
				operationLog,
				metav1.CreateOptions{},
			); err != nil {
				log.Errorf(
					"request id %s: create an operation log error: %s",
					requestID,
					err,
				)
			}
		}()
	}
}

func getRequestID(c *gin.Context) string {
	requestID, ok := c.Get(XRequestIDKey)
	if !ok {
		requestID = ""
	}

	return requestID.(string)
}

func getRequestURI(c *gin.Context) string {
	requestURI := c.Request.URL.Path
	if c.Request.URL.RawQuery != "" {
		requestURI = c.Request.URL.Path + "?" + c.Request.URL.RawQuery
	}

	return requestURI
}
