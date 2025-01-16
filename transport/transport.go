package transport

import (
	"context"
	"encoding/json"
	"net/http"

	"go.uber.org/zap"

	"github.com/cyber/test-project/logging"
)

func SendJson(ctx context.Context, w http.ResponseWriter, statusCode int, body any) {
	w.Header().Set("Content-Type", "application/json")

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		logging.Logger.Error("error marshalling body", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(statusCode)
	_, err = w.Write(bodyBytes)
	if err != nil {
		logging.Logger.Error("error writing body", zap.Error(err))
	}
}

func SendError(ctx context.Context, w http.ResponseWriter, err error) {
	SendJson(ctx, w, http.StatusInternalServerError, err)
}
