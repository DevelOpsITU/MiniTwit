package backgroundservices

import (
	"minitwit/config"
	"minitwit/database"
	"minitwit/log"
	"minitwit/metrics"
	"strconv"
	"time"
)

type program struct {
	cancellationToken chan struct{}
}

func start(p *program) {

	interval, err := strconv.Atoi(config.GetConfig().Services.ScrapeTimeInterval)

	if err != nil {
		log.Logger.Info().Msg("Could not parse SCRAPE_TIME_INTERVAL. Defaults to 60 secs")
		interval = 60
	}
	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	p.cancellationToken = make(chan struct{})
	go func() { // run async
		for {
			select {
			case <-ticker.C:
				observeCount((&database.User{}).TableName())
				observeCount((&database.Follower{}).TableName())
				observeCount((&database.Message{}).TableName())
			case <-p.cancellationToken:
				ticker.Stop()
				return
			}
		}
	}()
}

func stop(p *program) {
	close(p.cancellationToken)
}

func observeCount(tableName string) {
	var count int64
	database.GetGormDb().
		Table(tableName).
		Count(&count)
	metrics.PostgresData.WithLabelValues(tableName).Set(float64(count))
	log.Logger.Info().Int64("value", count).Str("table", tableName).Msg("Pulled count from database")
}
