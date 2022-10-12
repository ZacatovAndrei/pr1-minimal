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
	TimeUnit     = 2 * time.Second
	GeneratorNum = 5
	S2URL        = "http://localhost:8081/receive"
	S1URL        = "localhost:8080"
)

var (
	generators = make([]Generator, GeneratorNum)
	inQueue    = list.New()
)

var SentItems int

func main() {
	inQueue.Init()
	SentItems = 0

	for i := 0; i < GeneratorNum; i++ {
		log.Println("initialized generator #", i)
		go generators[i].Start(i, inQueue)
	}

	log.Println("starting server on port 8080")
	http.HandleFunc("/receive", receiver)
	if ok := http.ListenAndServe(S1URL, nil); ok != nil {
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
	return
}
