package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"net/http"

	nsq "github.com/bitly/go-nsq"
)

func main() {
	request, _ := http.NewRequest("GET", "http://devel-go.tkpd:3002/v1/product/get_summary?product_id=84072", nil)
	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	response, err := client.Do(request)

	if err != nil {
		fmt.Println("Http Request error, detail = ", err.Error())
	}

	data, _ := ioutil.ReadAll(response.Body)
	jsonData, err := json.Marshal(string(data))

	if err != nil {
		fmt.Println(err.Error())
	}
	config := nsq.NewConfig()
	p, err := nsq.NewProducer("devel-go.tkpd:4150", config)
	if err != nil {
		log.Panic(err)
	}
	err = p.Publish("test-nsq", []byte(jsonData))
	if err != nil {
		log.Panic(err)
	}
	request, _ = http.NewRequest("GET", "http://devel-go.tkpd:3002/v1/shop/get_summary?shop_id=881", nil)
	request.Header.Set("Content-Type", "application/json")
	response, err = client.Do(request)

	if err != nil {
		fmt.Println("Http Request error, detail = ", err.Error())
	}

	data, _ = ioutil.ReadAll(response.Body)
	jsonData, err = json.Marshal(string(data))

	err = p.Publish("test-nsq", []byte(jsonData))
	if err != nil {
		log.Panic(err)
	}
	fmt.Println("success publish to test-nsq")
}
