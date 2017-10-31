package main

import (
	"flag"
	"github.com/google/gops/agent"
	"log"
	"fmt"
	"net/http"

	"github.com/tokopedia/gosample/debugging"
	"gopkg.in/tokopedia/grace.v1"
	"gopkg.in/tokopedia/logging.v1"
	"github.com/tokopedia/gosample/calc"
)

func main() {

	flag.Parse()
	logging.LogInit()

	debug := logging.Debug.Println

	debug("app started") // message will not appear unless run with -debug switch

	if err = agent.Listen(agent.Options{
		ShutdownCleanup: true, // automatically closes on os.Interrupt
	}); err != nil {
		log.Fatal(err)
	}

	fmt.Println(calc.substract(4, 3))
	http.HandleFunc("/prob", debugging.ProbHandler)

	log.Fatal(grace.Serve(":9000", nil))
}
