package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"math"
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
		data := DecryptBlocks(input, key, 5)

		results <- Result{
			key:    key,
			result: IsASCIIBytes(data),
		}
	}
}

func statusUpdater(result <-chan Result) {
	tried := int64(0)
	found := make([]struct {
		Result
		elapsed time.Duration
	}, 0)
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

			fmt.Printf("Tried %12d and found %2d valid keys (%10.2f keys/sec)\r", tried, len(found), rate)
			time.Sleep(100 * time.Millisecond)
		}
	}()

	for r := range result {
		tried += 1
		if r.result {
			found = append(found, struct {
				Result
				elapsed time.Duration
			}{
				Result:  r,
				elapsed: time.Since(start),
			})
		}
	}

	fmt.Printf("\nTried %d keys (%.2f keys/sec)\n", tried, float64(tried)/time.Since(start).Seconds())

	for _, r := range found {
		fmt.Println("Found key:", string(r.key), "in", r.elapsed)
	}
}

func main() {
	var keyBase string

	flag.StringVar(&keyBase, "base", "", "known part of the key")
	flag.Parse()

	if len(keyBase) > 16 {
		panic("Key base too big")
	}

	fmt.Printf("Keys to try: %.0f\n", math.Pow(float64(len(characters)), float64(16-len(keyBase))))

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

	keys := make(chan []byte, 10000)
	results := make(chan Result, 10000)

	go keyGenerator(keyBase, keys)

	for range 18 {
		wg.Add(1)
		go worker(input, keys, results, &wg)
	}

	go statusUpdater(results)

	wg.Wait()
	close(results)

	println()
}
