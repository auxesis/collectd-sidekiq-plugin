package main

import (
	"encoding/json"
	"fmt"
	"github.com/jrallison/go-workers"
	"gopkg.in/alecthomas/kingpin.v2"
	"net/http"
	"net/http/httptest"
	"strconv"
	"time"
)

type stats struct {
	Processed int         `json:"processed"`
	Failed    int         `json:"failed"`
	Jobs      interface{} `json:"jobs"`
	Enqueued  interface{} `json:"enqueued"`
	Retries   int64       `json:"retries"`
}

func setup(server string, database string, pool string) {
	workers.Configure(map[string]string{
		// location of redis instance
		"server": server,
		// instance of the database
		"database": database,
		// number of connections to keep open with redis
		"pool": pool,
		// unique process id for this instance of workers (for proper recovery of inprogress jobs on crash)
		"process": "1",
	})

	workers.Process("scraper", func(message *workers.Msg) {}, 0)
}

func poll(t time.Time, hostname string, interval int) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Error: %+v\n", r)
		}
	}()
	req := http.Request{}
	res := httptest.NewRecorder()
	workers.Stats(res, &req)
	body := []byte(res.Body.String())
	var s stats
	err := json.Unmarshal(body, &s)
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Printf("PUTVAL %s/sidekiq/processed interval=%d %d:%d\n", hostname, interval, t.Unix(), s.Processed)
	fmt.Printf("PUTVAL %s/sidekiq/failed interval=%d %d:%d\n", hostname, interval, t.Unix(), s.Failed)
	fmt.Printf("PUTVAL %s/sidekiq/retries interval=%d %d:%d\n", hostname, interval, t.Unix(), s.Retries)
	for q, c := range s.Enqueued.(map[string]interface{}) {
		fmt.Printf("PUTVAL %s/sidekiq-%s/queue_depth interval=%d %d:%s\n", hostname, q, interval, t.Unix(), c)
	}
}

var (
	hostname = kingpin.Flag("hostname", "Hostname").OverrideDefaultFromEnvar("COLLECTD_HOSTNAME").Default("UNKNOWN").String()
	interval = kingpin.Flag("interval", "Interval").OverrideDefaultFromEnvar("COLLECTD_INTERVAL").Default("10").Int()
	server   = kingpin.Flag("redis-server", "Redis server in host:port format").Default("localhost:6379").String()
	database = kingpin.Flag("redis-database", "Redis database").Default("0").String()
	pool     = kingpin.Flag("redis-pool", "Redis pool size").Default("5").String()
)

func main() {
	kingpin.Parse()
	setup(*server, *database, *pool)

	duration, err := time.ParseDuration(strconv.Itoa(*interval) + "s")
	if err != nil {
		fmt.Printf("Error: %+v\n", err)
	}
	ticker := time.NewTicker(duration)
	for t := range ticker.C {
		poll(t, *hostname, *interval)
	}
}
