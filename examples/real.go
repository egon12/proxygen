package main

import "time"

type RealSomeRepository struct{}

func (r *RealSomeRepository) Get(id int) (Some, error) {
	time.Sleep(time.Second)
	return Some{}, nil
}

func (r *RealSomeRepository) Save(s Some) error {
	time.Sleep(time.Second)
	return nil
}
