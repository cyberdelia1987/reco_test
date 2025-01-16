package services

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"go.uber.org/zap"

	"github.com/cyber/test-project/config"
	"github.com/cyber/test-project/logging"
	"github.com/cyber/test-project/models"
)

type Dumper interface {
	DumpList(ctx context.Context, resources []models.TypedResource)
}

type AsanaDataDumper struct {
	cfg config.DataDumperConfig
}

func NewAsanaDataDumper(cfg config.DataDumperConfig) *AsanaDataDumper {
	return &AsanaDataDumper{
		cfg: cfg,
	}
}

func (d AsanaDataDumper) DumpList(ctx context.Context, resources []models.TypedResource) {
	logger := logging.FromContext(ctx).With(zap.String("operation", "dump_resources"))

	for _, res := range resources {
		path := fmt.Sprintf("%s/%s", d.cfg.Path, res.GetResourceType())
		err := os.MkdirAll(path, 0755)
		if err != nil {
			logger.Warn("Failed to create directory", zap.String("path", path))
			continue
		}

		pathFn := fmt.Sprintf("%s/%s", path, res.GetResourceType())

		encoded, err := json.Marshal(res)
		if err != nil {
			logger.Warn("failed to marshal resource", zap.String("resource", res.GetGid()), zap.Error(err))
			continue
		}

		fh, err := os.Create(pathFn)
		if err != nil {
			logger.Error("failed to create file", zap.String("path", pathFn), zap.Error(err))
			continue
		}

		_, err = fh.Write(encoded)
		if err != nil {
			logger.Error("failed into write file", zap.String("path", pathFn), zap.Error(err))
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
