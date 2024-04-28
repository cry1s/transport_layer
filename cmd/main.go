package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"sort"
	"time"

	"github.com/IBM/sarama"
	"github.com/cry1s/transport_layer/internal/app"
	"github.com/cry1s/transport_layer/internal/http/constants"
	"github.com/cry1s/transport_layer/internal/model"
)

func main() {
	log.Println("Transport Layer Started")

	application, err := app.New()
	if err != nil {
		log.Println(err)
	}

	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	consumer, err := sarama.NewConsumer([]string{"localhost:29093"}, config)
	if err != nil {
		log.Fatalf("Failed to create consumer: %v", err)
	}
	defer consumer.Close()

	partConsumer, err := consumer.ConsumePartition("forum-topic", 0, sarama.OffsetNewest)
	if err != nil {
		log.Fatalf("Failed to consume partition: %v", err)
	}
	defer partConsumer.Close()

	go func() {
		type message struct {
			segments []model.Segment
			cycles   int
		}

		var messages = make(map[time.Time]*message)
		var ticker = time.NewTicker(constants.TIME_AWAIT)

		for {
			select {
			case <-ticker.C:
				for id, msg := range messages {
					msg.cycles++

					if msg.cycles < 1 {
						continue
					}

					var fullMsg string

					// sort msg.segments
					sort.Slice(msg.segments, func(i, j int) bool {
						return msg.segments[i].SegmentNumber < msg.segments[j].SegmentNumber
					})

					for _, segment := range msg.segments {
						fullMsg += segment.Payload
					}

					payload := model.CollectedMessage{
						Message:    fullMsg,
						SenderName: msg.segments[0].SenderName,
						Time:       msg.segments[0].ID,
						Error:      len(msg.segments) < int(msg.segments[0].TotalSegments),
					}

					payloadBytes, err := json.Marshal(payload)
					if err != nil {
						log.Printf("Error marshaling payload: %v", err)
						continue
					}

					resp, err := http.Post("http://localhost:8001/receive", "application/json", bytes.NewBuffer(payloadBytes))
					if err != nil {
						log.Printf("Error sending post request: %v", err)
						continue
					}

					if resp.StatusCode != http.StatusOK {
						log.Printf("Server responded with non-OK status: %v", resp.Status)
					}

					delete(messages, id) // remove processed message
				}
			case msg, ok := <-partConsumer.Messages():
				if !ok {
					log.Println("Channel closed, exiting goroutine")
					return
				}

				var segment model.Segment
				err := json.Unmarshal(msg.Value, &segment)
				if err != nil {
					log.Printf("Error unmarshaling segment: %v", err)
					continue
				}

				if _, ok := messages[segment.ID]; !ok {
					messages[segment.ID] = &message{}
				}

				messages[segment.ID].segments = append(messages[segment.ID].segments, segment)
			}
		}
	}()

	application.StartServer()

	log.Println("Transport Layer Shutting Down")
}
