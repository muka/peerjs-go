package util

import (
	"math"
	"math/rand"
	"time"
)

// Seed personal random source - faster and won't mess with global one
var tokenRand = rand.New(rand.NewSource(time.Now().UnixNano()))

const tokenChars = "1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandomToken() string {
	b := make([]byte, 11) // PeerJS random tokens are 11 chars long
	for i := range b {
		b[i] = tokenChars[tokenRand.Intn(len(tokenChars))]
	}
	return string(b)
}

//ChunckedData wraps a data slice with metadata to assemble back the whole data
type ChunckedData struct {
	PeerData int    `json:"__peerData"`
	N        int    `json:"n"`
	Total    int    `json:"total"`
	Data     []byte `json:"data"`
}

//Chunk slices a data payload in a list of ChunckedData
func Chunk(raw []byte) (chunks []ChunckedData) {
	s := slicer{
		chunks: chunks,
	}
	return s.chunk(raw)
}

type slicer struct {
	dataCount int
	chunks    []ChunckedData
}

func (s *slicer) chunk(raw []byte) []ChunckedData {
	size := len(raw)
	total := int(math.Ceil(float64(size) / ChunkedMTU))
	index := 0
	start := 0

	for start < size {
		end := math.Min(float64(size), float64(start)+ChunkedMTU)
		b := raw[start:int(end)]

		chunk := ChunckedData{
			PeerData: s.dataCount,
			N:        index,
			Data:     b,
			Total:    total,
		}

		s.chunks = append(s.chunks, chunk)

		start = int(end)
		index++
	}

	s.dataCount++

	return s.chunks
}
