package main

import (
	"encoding/gob"
	"log"
	"os"
	"os/user"
	"path"
	"time"
)

type Person struct {
	Name string
	Age  int
}

func main() {
	p := Person{
		"Alice",
		20,
	}

	data := cache("test", func() Person {
		log.Println("not cached")
		time.Sleep(3 * time.Second)
		return p
	})

	log.Println(data)
}

type Fun[T any] func() T

func cache[T any, F func() T](key string, fun F) T {
	usr, _ := user.Current()

	dir := path.Join(usr.HomeDir, ".cache")

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.Mkdir(dir, os.ModeDir|0755)
	}

	file := path.Join(dir, key+".gob")

	if _, err := os.Stat(file); err == nil {
		f, err := os.Open(file)
		if err != nil {
			log.Fatal(err)
		}

		var q T
		dec := gob.NewDecoder(f)
		if err := dec.Decode(&q); err != nil {
			log.Fatal("decode error:", err)
		}

		return q
	}

	f, err := os.Create(path.Join(dir, key+".gob"))

	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	enc := gob.NewEncoder(f)

	p := fun()

	if err := enc.Encode(p); err != nil {
		log.Fatal(err)
	}

	return p
}

func ucache[T any, F func() T](key string, fun F) T {
	usr, _ := user.Current()

	dir := path.Join(usr.HomeDir, ".cache")

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.Mkdir(dir, os.ModeDir|0755)
	}

	file := path.Join(dir, key+".gob")

	p := fun()

	if _, err := os.Stat(file); err == nil {
		os.Remove(file)
		return p
	}

	return p
}
