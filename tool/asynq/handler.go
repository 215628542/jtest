package asynqTool

import (
	"log"

	"github.com/hibiken/asynq"
)

const redisAddr2 = "127.0.0.1:6379"

func Handler() {
	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: redisAddr2},
		asynq.Config{
			// Specify how many concurrent workers to use
			Concurrency: 3,
			// Optionally specify multiple queues with different priority.
			Queues: map[string]int{
				"critical": 1,
				"default":  1,
				"low":      1,
			},
			// See the godoc for other configuration options
		},
	)

	// mux maps a type to a handler
	mux := asynq.NewServeMux()
	//mux.HandleFunc(TypeEmailDelivery, HandleEmailDeliveryTask)
	//mux.Handle(TypeImageResize, NewImageProcessor())
	mux.HandleFunc(TestDemo, HandleTestDemoTask)

	if err := srv.Run(mux); err != nil {
		log.Fatalf("could not run asynq server: %v", err)
	}
}
