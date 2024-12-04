package middleware

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/segmentfault/pacman/log"
	"io/ioutil"
	"time"
)

//copy官方的logger

const (
	green   = "\033[97;42m"
	white   = "\033[90;47m"
	yellow  = "\033[90;43m"
	red     = "\033[97;41m"
	blue    = "\033[97;44m"
	magenta = "\033[97;45m"
	cyan    = "\033[97;46m"
	reset   = "\033[0m"
)

//var DefaultWriter io.Writer = os.Stdout

// DefaultErrorWriter is the default io.Writer used by Gin to debug errors
//var DefaultErrorWriter io.Writer = os.Stderr

var defaultLogFormatter = func(param gin.LogFormatterParams, requestBody string) string {
	var statusColor, methodColor, resetColor string
	if param.IsOutputColor() {
		statusColor = param.StatusCodeColor()
		methodColor = param.MethodColor()
		resetColor = param.ResetColor()
	}

	if param.Latency > time.Minute {
		param.Latency = param.Latency.Truncate(time.Second)
	}
	return fmt.Sprintf("[HTTP API] %v |%s %3d %s| %13v | %15s |%s %-7s %s %#v\n%s   %s",
		param.TimeStamp.Format("2006/01/02 - 15:04:05"),
		statusColor, param.StatusCode, resetColor,
		param.Latency,
		param.ClientIP,
		methodColor, param.Method, resetColor,
		param.Path,
		param.ErrorMessage,

		requestBody, //@my add
	)
}

func LoggerWithConfig() gin.HandlerFunc {

	//if formatter == nil {
	formatter := defaultLogFormatter
	//}
	//
	//out := conf.Output
	//out:=nil
	//if out == nil {
	//	out = DefaultWriter
	//}
	//out := DefaultWriter
	//
	//notlogged := conf.SkipPaths

	//isTerm := true

	//if w, ok := out.(*os.File); !ok || os.Getenv("TERM") == "dumb" ||
	//	(!isatty.IsTerminal(w.Fd()) && !isatty.IsCygwinTerminal(w.Fd())) {
	//	isTerm = false
	//}

	var skip map[string]struct{}

	//if length := len(notlogged); length > 0 {
	//	skip = make(map[string]struct{}, length)
	//
	//	for _, path := range notlogged {
	//		skip[path] = struct{}{}
	//	}
	//}

	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		//https://blog.csdn.net/impressionw/article/details/84194783
		rawData, err := c.GetRawData()
		if err != nil {
			fmt.Println(err.Error())
		}
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(rawData)) // 关键点

		// Process request
		c.Next()

		// Log only when path is not being skipped
		if _, ok := skip[path]; !ok {
			param := gin.LogFormatterParams{
				Request: c.Request,
				//isTerm:  isTerm,
				Keys: c.Keys,
			}

			// Stop timer
			param.TimeStamp = time.Now()
			param.Latency = param.TimeStamp.Sub(start)

			param.ClientIP = c.ClientIP()
			param.Method = c.Request.Method
			param.StatusCode = c.Writer.Status()
			param.ErrorMessage = c.Errors.ByType(gin.ErrorTypePrivate).String()

			param.BodySize = c.Writer.Size()

			if raw != "" {
				path = path + "?" + raw
			}

			param.Path = path

			//	fmt.Fprint(out, formatter(param, string(rawData)))
			lastStr := formatter(param, string(rawData))
			if param.StatusCodeColor() == green {
				log.Info(lastStr)
			} else {
				log.Error(lastStr)
			}

		}
	}
}
