package backgroundservices

import (
	"minitwit/database"
	"minitwit/metrics"
	"time"
)

const seconds = 5 // interval

type program struct {
	cancelationtoken chan struct{}
}

// background timer to scrape every 5 seconds
func start(p *program) {
	ticker := time.NewTicker(time.Duration(seconds) * time.Second)
	p.cancelationtoken = make(chan struct{})
	go func() { // run async
		for {
			select {
			case <-ticker.C:
				observeCount((&database.User{}).TableName())
				observeCount((&database.Follower{}).TableName())
				observeCount((&database.Message{}).TableName())
			case <-p.cancelationtoken:
				ticker.Stop()
				return
			}
		}
	}()
}

func stop(p *program) {
	close(p.cancelationtoken)
}

func observeCount(tablename string) {
	var count int64
	database.GetGormDb().
		Table(tablename).
		Count(&count)
	metrics.PostgresData.WithLabelValues(tablename).Set(float64(count))
}
