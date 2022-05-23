package vlc

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

type encodingTable map[rune]string

type BinaryChunks []BinaryChunk

type BinaryChunk string

type HeChunk string

type HexChunks []HeChunk

const chunksSize = 8

func Encode(str string) string {
	// prepare text: M-> !m
	// encode to binary: some text -> 10101010
	// split binary by chunks (8): bits to bytes -> '10000011 10001000 10011001'
	// bytes to hex -> '20 30 3c'
	// return hexChunksStr

	str = prepareText(str)

	chunks := splitByChunks(encodeBin(str), chunksSize)

	return chunks.ToHex().ToString()
}

func (hcs HexChunks) ToString() string {
	const sep = " "

	switch len(hcs) {
	case 0:
		return ""
	case 1:
		return string(hcs[0])
	}
	var buf strings.Builder

	buf.WriteString(string(hcs[0]))

	for _, hc := range hcs[1:] {
		buf.WriteString(sep)
		buf.WriteString(string(hc))
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

// encodeBin encode str into binary codes string without spaces.
func encodeBin(str string) string {
	var buf strings.Builder
	for _, ch := range str {
		buf.WriteString(bin(ch))
	}
	return buf.String()
}

func bin(ch rune) string {
	table := getEncodingTable()
	res, ok := table[ch]
	if !ok {
		panic("unknow character:" + string(ch))
	}
	return res
}

func getEncodingTable() encodingTable {
	return encodingTable{
		' ': "11",
		't': "1001",
		'n': "10000",
		's': "0101",
		'r': "01000",
		'd': "00101",
		'!': "001000",
		'c': "000101",
		'm': "000011",
		'g': "0000100",
		'b': "0000010",
		'v': "00000001",
		'k': "0000000001",
		'q': "000000000001",
		'e': "101",
		'o': "10001",
		'a': "011",
		'i': "01001",
		'h': "0011",
		'l': "001001",
		'u': "00011",
		'f': "000100",
		'p': "0000101",
		'w': "0000011",
		'y': "0000001",
		'j': "000000001",
		'x': "00000000001",
		'z': "000000000000"}
}

// prepareText prepares text to be fit for encode:
// changes upper case letters to: ! + lower case letter
// i.g.: My name is Ted -> !my name is !ted
func prepareText(str string) string {
	var buf strings.Builder

	for _, ch := range str {
		if unicode.IsUpper(ch) {
			buf.WriteRune('!')
			buf.WriteRune(unicode.ToLower(ch))
		} else {
			buf.WriteRune(ch)
		}
	}
	return buf.String()
}
