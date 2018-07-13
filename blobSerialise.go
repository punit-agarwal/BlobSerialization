// Package utils defines independent utilities helpful for a sharding-enabled,
// Ethereum blockchain such as blob serialization as more.
package utils

import (
	"fmt"
	"math"
)

var (
	chunkSize      = int64(32)
	indicatorSize  = int64(1)
	chunkDataSize  = chunkSize - indicatorSize
	skipEvmBits    = byte(0x80)
	dataLengthBits = byte(0x1F)
)

// Flags to add to chunk delimiter.(These are stored in 3 MSB of 1st byte)
type Flags struct {
	skipEvmExecution bool
}

// RawBlob type which will contain flags and data for serialization.
type RawBlob struct {
	flags Flags
	data  []byte
}

// getNumChunks calculates the number of chunks that will be produced by a byte array of given length
func getNumChunks(dataSize int) int {
	numChunks := math.Ceil(float64(dataSize) / float64(chunkDataSize))
	return int(numChunks)
}

// getSerializedDatasize determines the number of bytes that will be produced by a byte array of given length
func getSerializedDatasize(dataSize int) int {
	return getNumChunks(dataSize) * int(chunkSize)
}

// getTerminalLength determines the length of the final chunk for a byte array of given length
func getTerminalLength(dataSize int) int {
	numChunks := getNumChunks(dataSize)
	return dataSize - ((numChunks - 1) * int(chunkDataSize))
}
