package compression

type Encoder interface {
	Encoder(str string) []byte
}

type Decoder interface {
	Decoder(data []byte) string
}
