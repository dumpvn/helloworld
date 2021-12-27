package main

import (
	"fmt"
	"log"
	"os"
	"sync"
)

type FooReader struct{}

func (fooReader *FooReader) Read(b []byte) (int, error) {
	fmt.Print("in > ")
	return os.Stdin.Read(b)
}

type FooWriter struct{}

func (fooWriter *FooWriter) Write(b []byte) (int, error) {
	fmt.Print("out> ")
	return os.Stdout.Write(b)
}

func main() {

	var wg sync.WaitGroup
	var ch = make(chan int, 10)
	go run(ch, &wg)
	for i := range ch {
		fmt.Println(i)
	}

	// Instantiate reader and writer.
	var (
		reader FooReader
		writer FooWriter
	)

	input := make([]byte, 4096)
	s, err := reader.Read(input)
	if err != nil {
		log.Fatalln("Unable to read data")
	}

	fmt.Printf("read %d bytes from stdin\n", s)
	s, err = writer.Write(input)
	if err != nil {
		log.Fatalln("Unable to write data")
	}

	fmt.Printf("Wrote %d bytes to stdout\n", s)
}

func run(ch chan int, wg *sync.WaitGroup) {
	defer close(ch)

	for i := 0; i < cap(ch); i++ {
		wg.Add(1)
		go doSomeThing(i, ch, wg)
	}
}

func doSomeThing(i int, ch chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	ch <- i * 2
}

func ProcessFiles(files []string) {
	// create a channel with a capacity of 10 to be used as a semaphore
	var sem = make(chan bool, 10)
	var wg sync.WaitGroup
	for _, file := range files {
		sem <- true
		wg.Add(1)
		go func(file string) {
			defer wg.Done()
			processFile(file)
			// read a value out of the channel to free the capacity
			<-sem
		}(file)
	}
	wg.Wait()
}

func processFile(file string) {

}
