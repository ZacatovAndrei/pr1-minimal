package main

import (
	"bytes"
	"container/list"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"
)

var (
	QueueAccess sync.Mutex
)

type Consumer struct {
	id int
}

func (c *Consumer) Init(id int) {
	c.id = id
}

func (c *Consumer) Start(id int, inQ *list.List) {
	c.Init(id)
	for {
		if item := c.findItem(inQ); item != nil {
			*item = c.process(*item)
			c.Send(*item)
		} else {
			time.Sleep(TimeUnit)
		}
	}
}

func (c *Consumer) process(pl Payload) Payload {
	log.Println("processing item #", pl.Id)
	pl.ConsumerId = c.id
	pl.Payload = "PONG!"
	return pl
}

func (c *Consumer) findItem(inQ *list.List) *Payload {
	if inQ.Len() == 0 {
		return nil
	}
	QueueAccess.Lock()
	cast := inQ.Remove(inQ.Front()).(Payload)
	QueueAccess.Unlock()
	return &cast
}

func (c *Consumer) Send(p Payload) {
	var serialized []byte
	serialized, ok := json.Marshal(p)
	if ok != nil {
		log.Panicln(ok)
	}
	if _, ok := http.Post(S1URL, "text/json", bytes.NewBuffer(serialized)); ok != nil {
		log.Println(ok)
	}
	log.Printf("Sent %v back to S1\n", string(serialized))
}
