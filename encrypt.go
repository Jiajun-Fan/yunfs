package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"errors"
	"fmt"
	"io"
)

type Decryptor io.ReadCloser

type Encryptor io.WriteCloser

type AesDecryptor struct {
	Reader io.Reader
	Key    [aes.BlockSize]byte
	IV     [aes.BlockSize]byte
	Mode   cipher.BlockMode
	Buffer bytes.Buffer
}

type AesEncryptor struct {
	Writer io.Writer
	Key    [aes.BlockSize]byte
	IV     [aes.BlockSize]byte
	Mode   cipher.BlockMode
	Buffer bytes.Buffer
	IvDone bool
}

func NewAesDecryptor(key []byte, reader io.Reader) (Decryptor, error) {
	dec := &AesDecryptor{}
	dec.Reader = reader
	if len(key) > aes.BlockSize {
		return nil, errors.New(fmt.Sprintf("key size must not larger than %d", aes.BlockSize))
	}
	copy(dec.Key[:], key)

	block, err := aes.NewCipher(dec.Key[:])
	if err != nil {
		return nil, err
	}

	i := 0
	for {
		if i >= aes.BlockSize {
			break
		}
		bs, err := dec.Reader.Read(dec.IV[i:])
		if err != nil {
			if err == io.EOF {
				err = errors.New("not enought bytes for IV")
			}
			return nil, err
		}
		i += bs
	}
	if i != aes.BlockSize {
		panic("reader error")
	}
	dec.Mode = cipher.NewCBCDecrypter(block, dec.IV[:])

	return dec, nil
}

func NewAesEncryptor(key []byte, writer io.Writer) (Encryptor, error) {
	enc := &AesEncryptor{}
	enc.Writer = writer
	if len(key) > aes.BlockSize {
		return nil, errors.New(fmt.Sprintf("key size must not larger than %d", aes.BlockSize))
	}
	copy(enc.Key[:], key)

	block, err := aes.NewCipher(enc.Key[:])
	if err != nil {
		return nil, err
	}

	copy(enc.IV[:], RandStringBytes(aes.BlockSize))
	enc.Mode = cipher.NewCBCEncrypter(block, enc.IV[:])
	return enc, nil
}

func (dec *AesDecryptor) Read(output []byte) (int, error) {

	readSize := (len(output) / aes.BlockSize * aes.BlockSize) - dec.Buffer.Len()
	readBuff := make([]byte, readSize)

	var errRet error
	i := 0
	for {
		if i >= readSize {
			break
		}
		bs, err := dec.Reader.Read(readBuff[i:])
		if err != nil {
			if err == io.EOF {
				errRet = err
				break
			}
			// return any error other than io.EOF
			return 0, err
		}
		if bs == 0 {
			// if currently there is no available byte, return without blocking the process
			break
		}
		i += bs
	}

	if i > readSize {
		panic("reader error")
	}

	// ignore error as the error returned bytes.Buffer is always nil
	dec.Buffer.Write(readBuff[:i])

	decSize := dec.Buffer.Len() / aes.BlockSize * aes.BlockSize
	decBuff := make([]byte, decSize)
	dec.Buffer.Read(decBuff)

	dec.Mode.CryptBlocks(output[:decSize], decBuff)
	if decSize != 0 {
		return decSize, nil
	} else {
		return 0, errRet
	}
}

func (dec *AesDecryptor) Close() (err error) {
	if dec.Buffer.Len() != 0 {
		err = errors.New("There is unread bytes in buffer")
	}
	return
}

func (enc *AesEncryptor) write(input []byte) (int, error) {
	size := len(input)
	i := 0
	for {
		if i >= size {
			break
		}
		bs, err := enc.Writer.Write(input[i:])
		if err != nil {
			return 0, err
		}
		i += bs
	}

	if i > size {
		panic("writer error")
	}
	return i, nil
}

func (enc *AesEncryptor) Write(input []byte) (int, error) {

	enc.Buffer.Write(input)
	encSize := enc.Buffer.Len() / aes.BlockSize * aes.BlockSize
	encBuff := make([]byte, encSize)
	readBuff := make([]byte, encSize)

	enc.Buffer.Read(readBuff)

	enc.Mode.CryptBlocks(encBuff, readBuff)

	if !enc.IvDone {
		_, err := enc.write(enc.IV[:])
		if err != nil {
			return 0, err
		}
		enc.IvDone = true
	}

	return enc.write(encBuff)
}

func (enc *AesEncryptor) Close() (err error) {
	if enc.Buffer.Len() >= aes.BlockSize {
		panic("there is too much unwritten bytes")
	} else {
		buff := make([]byte, aes.BlockSize-enc.Buffer.Len())
		_, err = enc.Write(buff)
	}
	return
}

func MakeEncryptor(config EncryptConfig, writer io.Writer) Encryptor {
	if config.Type == "aes" {
		if aes, err := NewAesEncryptor([]byte(config.Key), writer); err != nil {
			Fatal(err.Error())
		} else {
			return aes
		}
	} else {
		Fatal(fmt.Sprintf("encryptor type '%s' is not implemented", config.Type))
	}
	return nil
}

func MakeDecryptor(config EncryptConfig, reader io.Reader) Decryptor {
	if config.Type == "aes" {
		if aes, err := NewAesDecryptor([]byte(config.Key), reader); err != nil {
			Fatal(err.Error())
		} else {
			return aes
		}
	} else {
		Fatal(fmt.Sprintf("decryptor type '%s' is not implemented", config.Type))
	}
	return nil
}
