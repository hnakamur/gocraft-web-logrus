package gocraft_web_logrus

import (
	"strconv"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/gocraft/web"
)

func LoggerMiddlewareFactory(accessLogMsg string) func(web.ResponseWriter, *web.Request, web.NextMiddlewareFunc) {
	return func(rw web.ResponseWriter, req *web.Request, next web.NextMiddlewareFunc) {
		startTime := time.Now()

		next(rw, req)

		duration := time.Since(startTime).Nanoseconds()
		var durationUnits string
		switch {
		case duration > 2000000:
			durationUnits = "ms"
			duration /= 1000000
		case duration > 1000:
			durationUnits = "μs"
			duration /= 1000
		default:
			durationUnits = "ns"
		}

		log.WithFields(log.Fields{
			"duration": strconv.FormatInt(duration, 10) + durationUnits,
			"status":   strconv.Itoa(rw.StatusCode()),
			"path":     req.URL.Path,
		}).Info(accessLogMsg)
	}
}
