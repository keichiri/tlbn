package bloom

import (
	"hash/fnv"
	"math"
	"math/big"
)

type BloomFilter struct {
	bitfield []byte
	rounds   int
	hashFunc func([]byte) []byte
}

func CalculateOptimalParameters(expectedItemCount int, falsePositivePercentage int) (int, int) {
	falsePositiveRate := float64(falsePositivePercentage) / 100
	bitCount := float64(-expectedItemCount) * math.Log(falsePositiveRate) / (math.Ln2 * math.Ln2)
	rounds := bitCount / float64(expectedItemCount) * math.Ln2
	return int(math.Ceil(bitCount)), int(math.Ceil(rounds))
}

func NewBloomFilter(bitfieldLength, rounds int, hashFunc func([]byte) []byte) *BloomFilter {
	byteCount := bitfieldLength / 8
	bitfield := make([]byte, byteCount)
	return &BloomFilter{
		bitfield: bitfield,
		rounds:   rounds,
		hashFunc: hashFunc,
	}
}

func (bf *BloomFilter) Add(input []byte) {
	for i := 0; i < bf.rounds; i++ {
		input = bf.hashFunc(input)
		bitPos := bf.getBitPos(input)
		bytePos, byteOffset := bitPos/8, bitPos%8
		bf.bitfield[bytePos] |= (1 << (7 - uint8(byteOffset)))
	}
}

func hashFunc1(input []byte) []byte {
	return fnv.New64().Sum(input)
}

func (bf *BloomFilter) IsMember(input []byte) bool {
	for i := 0; i < bf.rounds; i++ {
		input = bf.hashFunc(input)
		bitPos := bf.getBitPos(input)
		bytePos, byteOffset := bitPos/8, bitPos%8
		if bf.bitfield[bytePos]&(1<<(7-uint8(byteOffset))) == 0 {
			return false
		}
	}

	return true
}

func (bf *BloomFilter) getBitPos(hash []byte) uint64 {
	bigint := new(big.Int).SetBytes(hash)
	return bigint.Mod(bigint, big.NewInt(int64(len(bf.bitfield)*8))).Uint64()
}
