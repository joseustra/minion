package minion

import (
	"net/http"
	"time"
)

var (
	green   = string([]byte{27, 91, 57, 55, 59, 52, 50, 109})
	white   = string([]byte{27, 91, 57, 48, 59, 52, 55, 109})
	yellow  = string([]byte{27, 91, 57, 55, 59, 52, 51, 109})
	red     = string([]byte{27, 91, 57, 55, 59, 52, 49, 109})
	blue    = string([]byte{27, 91, 57, 55, 59, 52, 52, 109})
	magenta = string([]byte{27, 91, 57, 55, 59, 52, 53, 109})
	cyan    = string([]byte{27, 91, 57, 55, 59, 52, 54, 109})
	reset   = string([]byte{27, 91, 48, 109})
)

type statusWriter struct {
	http.ResponseWriter
	status int
	length int
}

func (w *statusWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func (w *statusWriter) Write(b []byte) (int, error) {
	if w.status == 0 {
		w.status = 200
	}
	w.length = len(b)
	return w.ResponseWriter.Write(b)
}

// Logger Logs the Http Status for all requests
func Logger(handle http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		start := time.Now()
		writer := statusWriter{rw, 0, 0}
		handle.ServeHTTP(&writer, req)
		end := time.Now()
		latency := end.Sub(start)
		statusCode := writer.status
		statusColor := colorForStatus(statusCode)
		methodColor := colorForMethod(req.Method)

		if req.URL.RawQuery != "" {
			l.Printf("%v |%s %-8s %s|%s %d %2s| %s | %12v | %s?%s",
				end.Format("2006/01/02 15:04:05"),
				methodColor, req.Method, reset,
				statusColor, statusCode, reset,
				req.RemoteAddr,
				latency,
				req.URL.Path,
				req.URL.RawQuery)
		} else {
			l.Printf("%v |%s %-8s %s|%s %d %2s| %s | %12v | %s",
				end.Format("2006/01/02 15:04:05"),
				methodColor, req.Method, reset,
				statusColor, statusCode, reset,
				req.RemoteAddr,
				latency,
				req.URL.Path)
		}
	})
}

func colorForStatus(code int) string {
	switch {
	case code >= 200 && code < 300:
		return green
	case code >= 300 && code < 400:
		return white
	case code >= 400 && code < 500:
		return yellow
	default:
		return red
	}
}

func colorForMethod(method string) string {
	switch method {
	case "GET":
		return blue
	case "POST":
		return cyan
	case "PUT":
		return yellow
	case "DELETE":
		return red
	case "PATCH":
		return green
	case "HEAD":
		return magenta
	case "OPTIONS":
		return white
	default:
		return reset
	}
}
