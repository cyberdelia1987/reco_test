package services

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"go.uber.org/zap"

	"github.com/cyber/test-project/logging"
	"github.com/cyber/test-project/models"
)

type AsanaDataDumper struct {
	pathToSave string
}

func NewAsanaDataDumper(pathToSave string) *AsanaDataDumper {
	return &AsanaDataDumper{
		pathToSave: pathToSave,
	}
}

func (d AsanaDataDumper) DumpList(ctx context.Context, resources []models.TypedResource) {
	logger := logging.FromContext(ctx).With(zap.String("operation", "dump_resources"))

	for _, res := range resources {
		path := fmt.Sprintf("%s/%s/%s.json", d.pathToSave, res.GetResourceType(), res.GetGid())

		encoded, err := json.Marshal(res)
		if err != nil {
			logger.Warn("failed to marshal resource", zap.String("resource", res.GetGid()), zap.Error(err))
			continue
		}

		fh, err := os.Create(path)
		if err != nil {
			logger.Error("failed to create file", zap.String("path", path), zap.Error(err))
			continue
		}

		_, err = fh.Write(encoded)
		if err != nil {
			logger.Error("failed into write file", zap.String("path", path), zap.Error(err))
		}

		closeFile(ctx, fh)
	}
}

func closeFile(ctx context.Context, fh *os.File) {
	err := fh.Close()
	if err == nil {
		return
	}

	logging.FromContext(ctx).Error("Failed to close file", zap.Error(err))
}
