package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/candbright/go-log/log"
	"github.com/gin-gonic/gin"
	"io"
	"strings"
	"time"
)

func LogHandler() gin.HandlerFunc {
	var notLogged []string
	var skip map[string]struct{}
	if length := len(notLogged); length > 0 {
		skip = make(map[string]struct{}, length)
		for _, path := range notLogged {
			skip[path] = struct{}{}
		}
	}
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery
		//get request body
		buff, _ := io.ReadAll(c.Request.Body)
		body := ""
		if buff != nil && len(buff) > 0 {
			t := io.NopCloser(bytes.NewBuffer(buff))
			c.Request.Body = t
			body = string(buff)
			var out bytes.Buffer
			err := json.Indent(&out, buff, "", "\t")
			if err == nil {
				body = out.String()
			}
		}
		c.Next()
		// Log only when path is not being skipped
		if _, ok := skip[path]; !ok {
			if raw != "" {
				path = path + "?" + raw
			}
			fields := map[string]interface{}{
				"ClientIP": c.ClientIP(),
				"Path":     path,
				"Latency":  strings.TrimSpace(fmt.Sprintf("%13v", time.Now().Sub(start))),
				"Method":   c.Request.Method,
				"Status":   c.Writer.Status(),
			}
			if c.Errors.ByType(gin.ErrorTypePrivate).String() != "" {
				fields["Error"] = c.Errors.ByType(gin.ErrorTypePrivate).String()
			}
			if body != "" {
				fields["RequestBody"] = body
			}
			log.WithFields(fields).Info("")
		}
	}
}
