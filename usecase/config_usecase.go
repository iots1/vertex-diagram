package usecase

import (
	"context"
	"time"

	"github.com/iots1/vertex-diagram/domain"
)

type configUsecase struct {
	configRepo     domain.ConfigRepository
	contextTimeout time.Duration
}

func NewConfigUsecase(repo domain.ConfigRepository, timeout time.Duration) domain.ConfigUsecase {
	return &configUsecase{
		configRepo:     repo,
		contextTimeout: timeout,
	}
}

func (u *configUsecase) Get(ctx context.Context) (*domain.Config, error) {
	c, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()
	return u.configRepo.Get(c)
}

func (u *configUsecase) Save(ctx context.Context, c *domain.Config) error {
	c2, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()
	return u.configRepo.Upsert(c2, c)
}
