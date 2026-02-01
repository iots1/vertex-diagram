package usecase

import (
	"context"
	"time"

	"github.com/iots1/vertex-diagram/domain"
)

type diagramUsecase struct {
	diagramRepo domain.DiagramRepository
	contextTimeout time.Duration
}

func NewDiagramUsecase(d domain.DiagramRepository, timeout time.Duration) domain.DiagramUsecase {
	return &diagramUsecase{
		diagramRepo:    d,
		contextTimeout: timeout,
	}
}

func (u *diagramUsecase) GetAll(c context.Context) ([]domain.Diagram, error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()
	return u.diagramRepo.Fetch(ctx)
}

func (u *diagramUsecase) GetOne(c context.Context, id string) (*domain.Diagram, error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()
	return u.diagramRepo.GetByID(ctx, id)
}

func (u *diagramUsecase) Save(c context.Context, d *domain.Diagram) (*domain.Diagram, error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	// Logic: ถ้า ID ว่าง หรือ เป็น Zero Value -> Create
	// ถ้ามี ID -> Update
	if d.ID.IsZero() {
        if d.Name == "" { d.Name = "Untitled Diagram" }
		err := u.diagramRepo.Store(ctx, d)
		return d, err
	} else {
		err := u.diagramRepo.Update(ctx, d)
		return d, err
	}
}