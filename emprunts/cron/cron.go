package cron

import (
	"fmt"

	"github.com/Bibliotheque-microservice/emprunts/handlers"
	"github.com/robfig/cron/v3"
)

func StartCron() {
	c := cron.New()

	// Every minute for test and demo purpose
	c.AddFunc("* * * * *", func() {
		result, err := handlers.CheckPenalities()
		if err != nil {
			fmt.Print("Erreur")
		} else {
			fmt.Print(result)
		}

	},
	)

	c.Start()

	select {}

}
