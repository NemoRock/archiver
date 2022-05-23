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

type HeChunk string

type HexChunks []HeChunk

const chunksSize = 8

const hexChunksSeparator = " "

func NewHexChunks(str string) HexChunks {

	parts := strings.Split(str, hexChunksSeparator)

	res := make(HexChunks, 0, len(parts))

	for _, part := range parts {
		res = append(res, HeChunk(part))
	}
	return res
}

func (hcs HexChunks) ToString() string {

	switch len(hcs) {
	case 0:
		return ""
	case 1:
		return string(hcs[0])
	}
	var buf strings.Builder

	buf.WriteString(string(hcs[0]))

	for _, hc := range hcs[1:] {
		buf.WriteString(hexChunksSeparator)
		buf.WriteString(string(hc))
	}
	return buf.String()
}

func (hcs HexChunks) ToBinary() BinaryChunks {
	res := make(BinaryChunks, 0, len(hcs))

	for _, chunk := range hcs {
		bChunk := chunk.ToBinary()
		res = append(res, bChunk)
	}
	return res
}

func (hc HeChunk) ToBinary() BinaryChunk {
	num, err := strconv.ParseUint(string(hc), 16, chunksSize)
	if err != nil {
		panic("can't parse hex chunk: " + err.Error())
	}
	res := fmt.Sprintf("%08b", num)
	return BinaryChunk(res)
}

// Join joins chunks into one line and returns as string
func (bcs BinaryChunks) Join() string {
	var buf strings.Builder

	for _, bc := range bcs {
		buf.WriteString(string(bc))
	}
	return buf.String()
}

func (bcs BinaryChunks) ToHex() HexChunks {
	res := make(HexChunks, 0, len(bcs))

	for _, chunk := range bcs {
		// chunk -> hexChunk
		hexChunk := chunk.ToHex()

		res = append(res, hexChunk)
	}
	return res
}

func (bc BinaryChunk) ToHex() HeChunk {
	num, err := strconv.ParseUint(string(bc), 2, chunksSize)
	if err != nil {
		panic("can't parse binary chuck: " + err.Error())
	}
	res := strings.ToUpper(fmt.Sprintf("%x", num))

	if len(res) == 1 {
		res = "0" + res
	}
	return HeChunk(res)
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
