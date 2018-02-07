package nsq

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/nsqio/go-nsq"
	"github.com/tokopedia/gosample/redis"
	logging "gopkg.in/tokopedia/logging.v1"
)

type ServerConfig struct {
	Name string
}

type Config struct {
	Server ServerConfig
}

type NSQModule struct {
	cfg *Config
	q   *nsq.Consumer
}

func NewNSQModule() *NSQModule {

	var cfg Config

	ok := logging.ReadModuleConfig(&cfg, "config", "hello") || logging.ReadModuleConfig(&cfg, "files/etc/gosample", "hello")
	if !ok {
		// when the app is run with -e switch, this message will automatically be redirected to the log file specified
		log.Fatalln("failed to read config")
	}

	// this message only shows up if app is run with -debug option, so its great for debugging
	logging.Debug.Println("nsq init called", cfg.Server.Name)

	// contohnya: caranya ciptakan nsq consumer
	nsqCfg := nsq.NewConfig()
	q := createNewConsumer(nsqCfg, "test-nsq", "exchann", handler)
	q.SetLogger(log.New(os.Stderr, "nsq:", log.Ltime), nsq.LogLevelError)
	q.ConnectToNSQLookupd("http://devel-go.tkpd:4161")

	return &NSQModule{
		cfg: &cfg,
		q:   q,
	}

}

func handler(msg *nsq.Message) error {
	var jsonData string
	_ = json.Unmarshal(msg.Body, &jsonData)
	fmt.Println("consuming topic, message finish!")
	bjson := []byte(jsonData)
	err := redis.Set("order_list:juna", bjson)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("set message to redis order_list:juna!")
	err = redis.Expire("order_list:juna", 60)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("set key redis order_list:juna to expire 60 seconds!")
	msg.Finish()
	return nil
}

func createNewConsumer(nsqCfg *nsq.Config, topic string, channel string, handler nsq.HandlerFunc) *nsq.Consumer {
	q, err := nsq.NewConsumer(topic, channel, nsqCfg)
	if err != nil {
		log.Fatal("failed to create consumer for ", topic, channel, err)
	}
	q.AddHandler(handler)
	return q
}
