package machinery

import (
	"os"
	"shortener/internal/tasks"

	"github.com/RichardKnop/machinery/v2"
	backendsAmqp "github.com/RichardKnop/machinery/v2/backends/amqp"
	brokersAmqp "github.com/RichardKnop/machinery/v2/brokers/amqp"
	"github.com/RichardKnop/machinery/v2/config"
	"github.com/RichardKnop/machinery/v2/locks/eager"

	"log"
)

func New(configEmail string, configPassword string) *machinery.Server {
	cnf := &config.Config{
		DefaultQueue:    "tasks",
		ResultsExpireIn: 3600,
		Broker:          "amqp://guest:guest@localhost:5672/",
		ResultBackend:   "amqp://guest:guest@localhost:5672/",
		AMQP: &config.AMQPConfig{
			Exchange:      "machinery_exchange",
			ExchangeType:  "direct",
			BindingKey:    "machinery_task",
			PrefetchCount: 3,
		},
	}

	broker := brokersAmqp.New(cnf)
	backend := backendsAmqp.New(cnf)
	lock := eager.New()

	server := machinery.NewServer(cnf, broker, backend, lock)

	err := server.RegisterTask("send_email", tasks.SendVerificationEmail(configEmail, configPassword))
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	return server
}
