package controllers

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"io"
)

type AuthJsonEncoded interface {
	RegisterRequestData | RegisterResponseData | LoginRequestData | LoginResponseData | Login2RequestData | Login2ResponseData
}

func AuthDecodeString(input string) []byte {
	ret, err := base64.StdEncoding.DecodeString(input)
	if err != nil {
		panic(err)
	}
	return ret
}

func AuthEncodeBytes(input []byte) string {
	return base64.StdEncoding.EncodeToString(input)
}

func AuthDecodeHexString(input string) []byte {
	ret, err := hex.DecodeString(input)
	if err != nil {
		panic(err)
	}
	return ret
}

func AuthEncodeHexBytes(input []byte) string {
	return hex.EncodeToString(input)
}

func AuthEncodeJson[D AuthJsonEncoded](input D) *bytes.Buffer {
	body := bytes.NewBufferString("")
	encoder := json.NewEncoder(body)
	if err := encoder.Encode(input); err != nil {
		panic(err)
	}
	return body
}

func AuthEncodeAndWriteJson[D AuthJsonEncoded](w io.Writer, input D) {
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(input); err != nil {
		panic(err)
	}
}

func AuthDecodeJson[D AuthJsonEncoded](input io.Reader, errCallback func(error)) *D {
	var result D
	decoder := json.NewDecoder(input)
	if err := decoder.Decode(&result); err != nil {
		errCallback(err)
		return nil
	}
	return &result
}
