package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/google/gops/agent"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/tokopedia/gosample/hello"
	"github.com/tokopedia/gosample/redis"
	"github.com/tokopedia/gosample/tugas"
	"github.com/tokopedia/gosample/utils/db"
	"github.com/tokopedia/logging/tracer"
	"gopkg.in/tokopedia/grace.v1"
	"gopkg.in/tokopedia/logging.v1"
)

func main() {
	redis.Init()
	flag.Parse()
	logging.LogInit()
	db.InitDB()
	debug := logging.Debug.Println

	debug("app started") // message will not appear unless run with -debug switch

	if err := agent.Listen(&agent.Options{}); err != nil {
		log.Fatal(err)
	}

	hwm := hello.NewHelloWorldModule()
	http.Handle("/metrics", promhttp.Handler())

	http.HandleFunc("/hello", hwm.SayHelloWorld)
	http.HandleFunc("/get/orders", tugas.HandleGetOrderAndOrderDetail)
	http.HandleFunc("/get/redis", tugas.HandleGetRedisExample)

	go logging.StatsLog()

	tracer.Init(&tracer.Config{Port: 8700, Enabled: true})

	log.Fatal(grace.Serve(":9000", nil))
}
