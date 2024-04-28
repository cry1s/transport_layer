package usecase

import "github.com/cry1s/transport_layer/internal/http/repository"

type UseCase struct {
	Repository *repository.Repository
}

func New(r *repository.Repository) *UseCase {
	return &UseCase{
		Repository: r,
	}
}
