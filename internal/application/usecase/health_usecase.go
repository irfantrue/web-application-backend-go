package usecase

import (
	"context"
	"time"

	"github.com/irfantrue/web-application-backend-go/internal/application/dto"
)

var startTime = time.Now()

type HealthUsecase interface {
	Check(ctx context.Context) (*dto.HealthResponse, error)
}

type healthUsecaseImpl struct {
}

func NewHealthUsecase() HealthUsecase {
	return &healthUsecaseImpl{}
}

func (s *healthUsecaseImpl) Check(ctx context.Context) (*dto.HealthResponse, error) {
	res := &dto.HealthResponse{
		Timestamp: time.Now().UTC().Truncate(time.Second),
		UpTime:    int(time.Since(startTime).Seconds()),
	}

	return res, nil
}
