package signer

import (
	"io"

	"github.com/tendermint/go-amino"
	cryptoAmino "github.com/tendermint/tendermint/crypto/encoding/amino"
	"github.com/tendermint/tendermint/privval"
)

var codec = amino.NewCodec()

// InitSerialization initalizes the private codec encoder/decoder
func InitSerialization() {
	cryptoAmino.RegisterAmino(codec)
	privval.RegisterRemoteSignerMsg(codec)
}

// ReadMsg reads a message from an io.Reader
func ReadMsg(reader io.Reader) (msg privval.SignerMessage, err error) {
	const maxRemoteSignerMsgSize = 1024 * 10
	_, err = codec.UnmarshalBinaryLengthPrefixedReader(reader, &msg, maxRemoteSignerMsgSize)
	return
}

// WriteMsg writes a message to an io.Writer
func WriteMsg(writer io.Writer, msg interface{}) (err error) {
	_, err = codec.MarshalBinaryLengthPrefixedWriter(writer, msg)
	return
}
