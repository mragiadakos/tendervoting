package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"

	crypto "github.com/libp2p/go-libp2p-crypto"
)

const (
	CodeTypeOK            uint32 = 0
	CodeTypeEncodingError uint32 = 1
	CodeTypeBadNonce      uint32 = 2
	CodeTypeUnauthorized  uint32 = 3
	CodeTypeClientError   uint32 = 4
)

type KeyJson struct {
	PublicKey  string // hex
	PrivateKey []byte
}

func fileKey(filename string) (crypto.PrivKey, error) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, errors.New("Error: " + err.Error())
	}

	kj := KeyJson{}
	err = json.Unmarshal(b, &kj)
	if err != nil {
		return nil, errors.New("Error: json problem with the key " + err.Error())
	}

	edKey, err := crypto.UnmarshalPrivateKey(kj.PrivateKey)
	if err != nil {
		return nil, errors.New("Error: private key decoding problem with the key " + err.Error())
	}
	return edKey, nil
}
