package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/google/gops/agent"
	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	config "github.com/tokopedia/gosample/config"
	redis "github.com/tokopedia/gosample/redis"

	"github.com/tokopedia/gosample/handler"

	"github.com/tokopedia/gosample/hello"
	"github.com/tokopedia/logging/tracer"
	"gopkg.in/tokopedia/grace.v1"
	"gopkg.in/tokopedia/logging.v1"
)

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {

	flag.Parse()
	logging.LogInit()
	config.Init()
	redis.Init()
	defer config.Db.Close()
	debug := logging.Debug.Println

	debug("app started") // message will not appear unless run with -debug switch

	if err := agent.Listen(&agent.Options{}); err != nil {
		log.Fatal(err)
	}

	hwm := hello.NewHelloWorldModule()

	http.Handle("/metrics", promhttp.Handler())

	http.HandleFunc("/hello", hwm.SayHelloWorld)
	http.HandleFunc("/users", handler.GetUsers)
	http.HandleFunc("/search", handler.SearchUsers)
	go logging.StatsLog()

	tracer.Init(&tracer.Config{Port: 8700, Enabled: true})

	log.Fatal(grace.Serve(":9000", nil))
}
