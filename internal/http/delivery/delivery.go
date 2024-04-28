package delivery

import (
	"log"

	"github.com/cry1s/transport_layer/internal/http/usecase"
	"github.com/cry1s/transport_layer/internal/kafka"
)

type Handler struct {
	UseCase  *usecase.UseCase
	Producer *kafka.Producer
}

func New(uc *usecase.UseCase) *Handler {
	producer, err := kafka.NewProducer()
	if err != nil {
		log.Fatal("Error occured while creating producer: ", err)
	}

	return &Handler{
		UseCase:  uc,
		Producer: producer,
	}
}
