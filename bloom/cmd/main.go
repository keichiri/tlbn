package main

import (
	"bloom"
	"bytes"
	"crypto/sha1"
	"fmt"
	"hash/fnv"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"time"
)

func main() {
	falsePositiveRate := 1
	input := getInput()
	bloomFilterDemo(input, falsePositiveRate)
}

func hashFuncFnv(input []byte) []byte {
	hasher := fnv.New64()
	hasher.Write(input)
	return hasher.Sum(nil)
}

func hashFuncSha1(input []byte) []byte {
	hasher := sha1.New()
	hasher.Write(input)
	return hasher.Sum(nil)
}

func getInput() [][]byte {
	f, err := os.Open("/usr/share/dict/words")
	if err != nil {
		log.Fatalf("Failed to open input file")
	}
	defer f.Close()

	input, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatalf("Failed to read input file")
	}
	if input[len(input)-1] == '\n' {
		input = input[:len(input)-1]
	}

	return bytes.Split(input, []byte{'\n'})
}

func bloomFilterDemo(input [][]byte, falsePositivePercentage int) {
	add := make([][]byte, 0)
	skip := make([][]byte, 0)
	for _, inputWord := range input {
		if rand.Intn(100) < falsePositivePercentage {
			skip = append(skip, inputWord)
		} else {
			add = append(add, inputWord)
		}
	}

	optimalLength, optimalRounds := bloom.CalculateOptimalParameters(len(add), falsePositivePercentage)
	fmt.Printf("Item count: %d. Target false positive percentage: %d. Optimal length: %d. Optimal rounds: %d\n",
		len(add), falsePositivePercentage, optimalLength, optimalRounds)

	start := time.Now()

	bf := bloom.NewBloomFilter(optimalLength, optimalRounds, hashFuncSha1)
	for _, word := range add {
		bf.Add(word)
	}

	hitCount := 0
	for _, word := range add {
		if !bf.IsMember(word) {
			log.Fatalf("Bloom filter returned negative result for existing word %s", string(word))
		} else {
			hitCount++
		}
	}

	falseHitCount := 0
	for _, word := range skip {
		if bf.IsMember(word) {
			falseHitCount++
		}
	}

	ratio := float64(falseHitCount) * 100 / float64(len(skip))
	fmt.Printf("Tested with %d items. False hit count: %d. False positive percentage: %f. Duration: %v\n",
		len(skip), falseHitCount, ratio, time.Since(start))

	// Testing with sizes and round counts that differ from the calculated optimal ones
	fmt.Printf("\n\n\nResults using FNV:\n")
	for sizeDiff := -2000000; sizeDiff <= 2000000; sizeDiff += 1000000 {
		for roundDiff := -6; roundDiff < 6; roundDiff += 2 {
			if sizeDiff == 0 && roundDiff == 0 {
				continue
			}
			demo(add, skip, optimalLength+sizeDiff, optimalRounds+roundDiff, hashFuncFnv)
		}
	}

	fmt.Printf("\n\n\nResults using SHA1:\n")
	for sizeDiff := -2000000; sizeDiff <= 2000000; sizeDiff += 1000000 {
		for roundDiff := -6; roundDiff < 6; roundDiff += 2 {
			if sizeDiff == 0 && roundDiff == 0 {
				continue
			}
			demo(add, skip, optimalLength+sizeDiff, optimalRounds+roundDiff, hashFuncSha1)
		}
	}

}

func demo(add [][]byte, skip [][]byte, bitfieldLength, rounds int, hashFunc func([]byte) []byte) {
	bf := bloom.NewBloomFilter(bitfieldLength, rounds, hashFunc)
	for _, word := range add {
		bf.Add(word)
	}

	hitCount := 0
	for _, word := range add {
		if !bf.IsMember(word) {
			log.Fatalf("Bloom filter returned negative result for existing word %s", string(word))
		} else {
			hitCount++
		}
	}

	falseHitCount := 0
	for _, word := range skip {
		if bf.IsMember(word) {
			falseHitCount++
		}
	}

	ratio := float64(falseHitCount) * 100 / float64(len(skip))
	fmt.Printf("Bit length: %d\nRounds: %d\nFalse positive percentage: %f\n\n",
		bitfieldLength, rounds, ratio)
}
