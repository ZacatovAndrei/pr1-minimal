package main

import (
	"container/list"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

const (
	TimeUnit    = 2 * time.Second
	ConsumerNum = 5
	S1URL       = "http://localhost:8080/receive"
	S2URL       = "localhost:8081"
)

var (
	consumers = make([]Consumer, ConsumerNum)
	inQueue   = list.New()
)

func main() {
	time.Sleep(2 * TimeUnit)
	for i := 0; i < ConsumerNum; i++ {
		log.Println("initialized consumer #", i)
		go consumers[i].Start(i, inQueue)
	}

	log.Println("starting server on port 8080")
	http.HandleFunc("/receive", receiver)
	if ok := http.ListenAndServe(S2URL, nil); ok != nil {
		log.Panic(ok)
	}

}

func receiver(w http.ResponseWriter, r *http.Request) {
	var (
		p Payload
		b []byte
	)
	b, ok := ioutil.ReadAll(r.Body)
	if ok != nil {
		log.Panic(ok)
	}

	if ok := json.Unmarshal(b, &p); ok != nil {
		log.Panic(ok)
	}

	QueueAccess.Lock()
	inQueue.PushBack(p)
	QueueAccess.Unlock()
	log.Printf("Got %v from S1\n", p)
	return
}
