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

// NewRawBlob builds a raw blob from any interface by using
// RLP encoding.
func NewRawBlob(i interface{}, skipEvm bool) (*RawBlob, error) {
	data, err := rlp.EncodeToBytes(i)
	if err != nil {
		return nil, fmt.Errorf("RLP encoding was a failure:%v", err)
	}
	return &RawBlob{data: data, flags: Flags{skipEvmExecution: skipEvm}}, nil
}

// ConvertFromRawBlob converts raw blob back from a byte array
// to its interface.
func ConvertFromRawBlob(blob *RawBlob, i interface{}) error {
	data := (*blob).data
	err := rlp.DecodeBytes(data, i)
	if err != nil {
		return fmt.Errorf("RLP decoding was a failure:%v", err)
	}

	return nil
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
