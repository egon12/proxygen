package main

type (
	Some struct {
		ID   int
		Name string
	}

	SomeRepository interface {
		Get(id int) (Some, error)
		Save(Some) error
		All() map[int]Some
	}

	Rep interface {
		Get(id int) (Some, error)
		Save(Some) error
	}
)
