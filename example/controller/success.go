package controller

import (
	"net/http"

	"github.com/payfazz/fz-sentry/logger"
	"go.uber.org/zap"
)

func Success() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log := logger.GetLogger(ctx)

		log.Debug("this is debug message")

		log.Info(
			"this is info from success controller",
			zap.String("status", "request processed successfully"),
		)

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)

		_, _ = w.Write([]byte("success"))
	}
}