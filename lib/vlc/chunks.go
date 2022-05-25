package vlc

import (
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"
)

type encodingTable map[rune]string

type BinaryChunks []BinaryChunk

type BinaryChunk string

const chunksSize = 8

func NewBinChunks(data []byte) BinaryChunks {

	res := make(BinaryChunks, 0, len(data))

	for _, code := range data {
		res = append(res, NewBinChunk(code))
	}
	return res
}

func NewBinChunk(code byte) BinaryChunk {
	return BinaryChunk(fmt.Sprintf("%08b", code))
}

func (bcs BinaryChunks) Bytes() []byte {
	res := make([]byte, 0, len(bcs))

	for _, bc := range bcs {
		res = append(res, bc.Byte())
	}
	return res
}

func (bc BinaryChunk) Byte() byte {
	num, err := strconv.ParseUint(string(bc), 2, chunksSize)
	if err != nil {
		panic("can't parse binary chuck: " + err.Error())
	}
	return byte(num)

}

// Join joins chunks into one line and returns as string
func (bcs BinaryChunks) Join() string {
	var buf strings.Builder

	for _, bc := range bcs {
		buf.WriteString(string(bc))
	}
	return buf.String()
}

// splitByChunks splits binary string by chunks with giver size,
// i.g.: '101010100101011001010100101010' -> '10101010 10101010 10101010'
func splitByChunks(bStr string, chunkSize int) BinaryChunks {
	strLen := utf8.RuneCountInString(bStr)

	chunkCount := strLen / chunkSize

	//TODO fix
	if chunkCount != 0 {
		chunkCount++
	}
	res := make(BinaryChunks, 0, chunkCount)

	var buf strings.Builder
	for i, ch := range bStr {
		buf.WriteString(string(ch))

		if (i+1)%chunkSize == 0 {
			res = append(res, BinaryChunk(buf.String()))
			buf.Reset()
		}
	}
	if buf.Len() != 0 {
		lastChunk := buf.String()
		lastChunk += strings.Repeat("0", chunkSize-len(lastChunk))

		res = append(res, BinaryChunk(lastChunk))
	}
	return res
}
