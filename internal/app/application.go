package app

import (
	"github.com/cry1s/transport_layer/internal/http/delivery"
	"github.com/cry1s/transport_layer/internal/http/repository"
	"github.com/cry1s/transport_layer/internal/http/usecase"
)

type Application struct {
	repository *repository.Repository
	usecase    *usecase.UseCase
	handler    *delivery.Handler
}

func New() (*Application, error) {
	repo, err := repository.New()
	uc := usecase.New(repo)
	h := delivery.New(uc)

	if err != nil {
		return &Application{}, err
	}

	return &Application{
		repository: repo,
		usecase:    uc,
		handler:    h,
	}, nil
}
