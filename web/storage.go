package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"time"
)

type trigger struct {
	LastTriggerTime time.Time
	c               context.Context
	cancel          context.CancelFunc
}

type database struct {
	path string
	D    map[string]trigger
}

func (s *database) load() {
	d, err := ioutil.ReadFile(s.path)
	if err != nil {
		log.Fatal("database file open failed.")
	}
	err = json.Unmarshal(d, s)
	if err != nil {
		log.Fatal("database load failed.")
	}
	if s.D == nil {
		s.D = make(map[string]trigger)
	}
	for k, v := range s.D {
		v.c, v.cancel = context.WithCancel(context.Background())
		s.D[k] = v
	}
}

func (s *database) save() {
	data, err := json.Marshal(s)
	if err != nil {
		log.Fatal("database marshal failed.")
	}
	err = ioutil.WriteFile(s.path, data, 500)
	if err != nil {
		log.Fatal("database write failed.")
	}
}
