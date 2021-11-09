package web

import (
	"net/http"
	"strconv"
	"strings"
	"time"
)

func loggerMiddleware(next http.HandlerFunc, logger Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		now := time.Now()
		next(w, r)
		since := time.Since(now).Milliseconds()
		builder := strings.Builder{}
		builder.Write([]byte(r.RemoteAddr + " "))
		builder.Write([]byte(r.Method + " "))
		builder.Write([]byte(r.RequestURI + " "))
		builder.Write([]byte(r.Proto + " "))
		builder.Write([]byte(strconv.Itoa(http.StatusOK) + " "))
		builder.Write([]byte(strconv.FormatInt(since, 10) + "ms "))
		builder.Write([]byte(r.UserAgent()))
		logger.Info(builder.String())
	})
}
