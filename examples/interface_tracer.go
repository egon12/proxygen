package main

import (
	"log"
	"time"
)

type SomeRepositoryTracer struct {
	real SomeRepository
}

func (t *SomeRepositoryTracer) Get(id int) (Some, error) {
	defer func(start time.Time) {
		end := time.Now()
		dif := end.Sub(start)
		log.Printf("Duration: *SomeRepositoryTracer.Get: %v", dif)
	}(time.Now())
	return t.real.Get(id)
}

func (t *SomeRepositoryTracer) Save(arg0 Some) error {
	defer func(start time.Time) {
		end := time.Now()
		dif := end.Sub(start)
		log.Printf("Duration: *SomeRepositoryTracer.Save: %v", dif)
	}(time.Now())
	return t.real.Save(arg0)
}

type RepTracer struct {
	real Rep
}

func (t *RepTracer) Get(id int) (Some, error) {
	defer func(start time.Time) {
		end := time.Now()
		dif := end.Sub(start)
		log.Printf("Duration: *RepTracer.Get: %v", dif)
	}(time.Now())
	return t.real.Get(id)
}

func (t *RepTracer) Save(arg0 Some) error {
	defer func(start time.Time) {
		end := time.Now()
		dif := end.Sub(start)
		log.Printf("Duration: *RepTracer.Save: %v", dif)
	}(time.Now())
	return t.real.Save(arg0)
}
