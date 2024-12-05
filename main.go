package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"
)

type Result struct {
	key    []byte
	result bool
}

func keyGenerator(knowPart string, keys chan<- []byte) {
	defer close(keys)

	for key := range GenerateKeys(knowPart) {
		keys <- []byte(key)
	}
}

func worker(input []byte, keys <-chan []byte, results chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()

	for key := range keys {
		valid := true

		for block := range len(input) / 16 {
			output := DecryptBlock(input, []byte(key), block)

			if !IsASCIIBytes(output) {
				valid = false
				break
			}
		}

		results <- Result{
			key:    key,
			result: valid,
		}
	}
}

func statusUpdater(result <-chan Result) {
	tried := int64(0)
	found := 0
	start := time.Now()

	go func() {
		for {
			var rate float64

			duration := time.Since(start).Seconds()

			if duration == 0 {
				rate = 0
			} else {
				rate = float64(tried) / duration
			}

			fmt.Printf("Tried %10d and found %d valid key (%.2f keys/sec)\r", tried, found, rate)
			time.Sleep(500 * time.Millisecond)
		}
	}()

	for r := range result {
		tried += 1
		if r.result {
			found += 1
			fmt.Printf("Found key %s\n", string(r.key))
		}
	}
}

func main() {
	var keyBase string

	flag.StringVar(&keyBase, "base", "", "known part of the key")
	flag.Parse()

	if len(keyBase) > 16 {
		panic("Key base too big")
	}

	bytes, err := io.ReadAll(os.Stdin)
	if err != nil {
		panic("Failed to read input")
	}

	data := string(bytes)
	data = strings.ReplaceAll(data, "\n", "")

	input, err := hex.DecodeString(data)
	if err != nil {
		panic("Failed to decode input")
	}

	var wg sync.WaitGroup

	keys := make(chan []byte, 15)
	results := make(chan Result, 15)

	go keyGenerator(keyBase, keys)

	for range 10 {
		wg.Add(1)
		go worker(input, keys, results, &wg)
	}

	go statusUpdater(results)

	wg.Wait()
	close(results)

	println()
}
