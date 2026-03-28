package bytes2hex

import (
	"bufio"
	"encoding/hex"
	"errors"
	"io"
	"os"
)

const PageSize int = 4096
const DoublePage int = PageSize << 1

type PageToHex func(page []byte, hx []byte)

func EncodePage(page []byte, hx []byte) {
	_ = hex.Encode(hx, page)
}

//nolint:gochecknoglobals
var PageToHexDefault PageToHex = EncodePage

type Encoder func(encoded []byte, original []byte) int

//nolint:gochecknoglobals
var EncoderDefault Encoder = hex.Encode

type BulkEncoder struct {
	Encoder
	PageToHex
}

func (b BulkEncoder) ReaderToWriter(rdr io.Reader, wtr io.Writer) error {
	var original [PageSize]byte
	var encoded [DoublePage]byte

	for {
		read, err := io.ReadFull(rdr, original[:])

		// Full page available
		if nil == err {
			b.PageToHex(original[:], encoded[:]) // read == len(original[:])
			_, err = wtr.Write(encoded[:])
			if nil != err {
				return err
			}
			continue
		}

		// Partial page may be available
		if errors.Is(err, io.ErrUnexpectedEOF) {
			if 0 == read {
				return err
			}

			var partial []byte = original[:read]
			var size int = b.Encoder(encoded[:], partial)
			_, err = wtr.Write(encoded[:size])
			if nil != err {
				return err
			}
			continue
		}

		if errors.Is(err, io.EOF) {
			return nil
		}
	}
}

func (b BulkEncoder) StdinToHexToStdout() error {
	var bwtr *bufio.Writer = bufio.NewWriter(os.Stdout)
	return errors.Join(
		b.ReaderToWriter(
			bufio.NewReader(os.Stdin),
			bwtr,
		),
		bwtr.Flush(),
	)
}

//nolint:gochecknoglobals
var BulkEncoderDefault BulkEncoder = BulkEncoder{
	Encoder:   EncoderDefault,
	PageToHex: PageToHexDefault,
}
