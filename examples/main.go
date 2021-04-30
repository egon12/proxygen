package main

func main() {

	var r SomeRepository
	r = &SomeRepositoryTracer{&RealSomeRepository{}}

	r.Get(20)
	r.Save(Some{})

}
