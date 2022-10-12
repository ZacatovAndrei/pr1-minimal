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

const (
	ready int = iota
	waiting
)

var (
	idAccess    sync.Mutex
	QueueAccess sync.Mutex
)

type Generator struct {
	id, state int
}

func (g *Generator) Start(id int, inQ *list.List) {
	g.id = id
	for {
		switch g.state {

		case ready:
			g.Send(g.generate())
			time.Sleep(TimeUnit)

		case waiting:
			if res := g.findItem(inQ); res != nil {
				log.Printf("%v from consumer %v", res.Payload, res.ConsumerId)
			}
			time.Sleep(TimeUnit)
		}
	}
}

func (g *Generator) generate() Payload {
	g.state = waiting
	idAccess.Lock()
	SentItems++
	bufPl := Payload{
		Id:          SentItems,
		GeneratorId: g.id,
		ConsumerId:  -1, //just to not confuse with consumer #0 which will exist anyway
		Payload:     "PING",
	}
	idAccess.Unlock()
	return bufPl
}

func (g *Generator) findItem(inQ *list.List) *Payload {
	QueueAccess.Lock()
	defer QueueAccess.Unlock()
	for e := inQ.Front(); e != nil; e = e.Next() {
		if e.Value.(Payload).GeneratorId == g.id {
			cast := e.Value.(Payload)
			inQ.Remove(e)
			g.state = ready
			return &cast
		}
	}
	return nil
}

func (g *Generator) Send(p Payload) {
	var serialized []byte
	serialized, ok := json.Marshal(p)
	if ok != nil {
		log.Panicln(ok)
	}
	if _, ok := http.Post(S2URL, "text/json", bytes.NewBuffer(serialized)); ok != nil {
		log.Panicln(ok)
	}
	log.Printf("Sent %v to S2\n", string(serialized))
	g.state = waiting
}
