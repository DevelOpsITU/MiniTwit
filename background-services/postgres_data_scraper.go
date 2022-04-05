package backgroundservices

import (
	"minitwit/config"
	"minitwit/database"
	"minitwit/metrics"
	"time"
)

type program struct {
	cancelationtoken chan struct{}
}

// background timer to scrape every 5 seconds
func start(p *program) {
	ticker := time.NewTicker(time.Duration(config.GetConfig().Services.ScrapeTimeInterval) * time.Second)
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
