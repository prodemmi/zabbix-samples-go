package main

import (
	"fmt"
	"math/rand/v2"
	"os"
	"time"

	"github.com/go-faker/faker/v4"
	"github.com/sirupsen/logrus"
)

func randRange(min int64, max int64) int64 {
	return rand.Int64N(max) + min
}

func randate() time.Time {
	min := time.Date(1970, 1, 0, 0, 0, 0, 0, time.UTC).Unix()
	max := time.Date(2070, 1, 0, 0, 0, 0, 0, time.UTC).Unix()
	delta := max - min

	sec := rand.Int64N(delta) + min
	return time.Unix(sec, 0)
}

func ranbool() bool {
	return rand.IntN(2) == 1
}

func SetupUserSeeder(m *UserMetrics) {
	var log = logrus.New()
	log.Out = os.Stdout

	file, err := os.OpenFile("./logs/out.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		log.Out = file
	} else {
		log.Info("Failed to log to file, using default stderr")
	}

	for i := 0; i < 10_000; i++ {
		loginAt, _ := time.Parse("2006-01-02", faker.Date())
		coutries := []string{"Iran", "America", "China"}
		newUser := User{
			ID:          i,
			Username:    fmt.Sprintf("@%v", faker.Username()),
			OrderCount:  rand.Int64N(20),
			LastOrderAt: randate(),
			MaxPayment:  rand.Int64N(5_000_000),
			IsOnline:    ranbool(),
			IP:          faker.IPv4(),
			Country:     coutries[rand.IntN(len(coutries))],
			LoginAt:     loginAt,
			CreatedAt:   time.Now(),
		}

		log.WithFields(logrus.Fields{
			"ID":        newUser.ID,
			"Username":  newUser.Username,
			"IP":        newUser.IP,
			"Country":   newUser.Country,
			"LoginAt":   newUser.LoginAt,
			"CreatedAt": newUser.CreatedAt,
		}).Info("User created:")

		time.Sleep(time.Second * 1)
	}
}
