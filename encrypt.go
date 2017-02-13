package main

import (
	"crypto/aes"
	"crypto/cipher"
	"errors"
	"fmt"
	"io"
)

const (
	kAesBlockSize = 32
)

type Cipher interface {
	Decrypt(key string, reader io.Reader)
	Encrypt(key string, reader io.Reader)
	io.Reader
}

type AesCipher struct {
	key    [kAesBlockSize]byte
	reader io.Reader
	iv     [kAesBlockSize]byte
	mode   cipher.BlockMode
	retIV  bool
}

func (ac AesCipher) Decrypt(key []byte, reader io.Reader) error {
	return ac.init(true, key, reader)
}

func (ac AesCipher) Encrypt(key []byte, reader io.Reader) error {
	return ac.init(false, key, reader)
}

func (ac AesCipher) init(dec bool, key []byte, reader io.Reader) error {
	ac.reader = reader

	if len(key) > kAesBlockSize {
		return errors.New(fmt.Sprintf("key size must not larger than %d", kAesBlockSize))
	}
	copy(ac.key[:], key)
	println(len(ac.key[:]))
	println(len(ac.iv[:]))

	block, err := aes.NewCipher(ac.key[:])
	if err != nil {
		return err
	}

	if dec {
		i := 0
		for {
			if i >= kAesBlockSize {
				break
			}
			bs, err := ac.reader.Read(ac.iv[i:])
			if err != nil {
				if err == io.EOF {
					err = errors.New("not enought bytes for IV")
				}
				return err
			}
			i += bs
		}
		if i != kAesBlockSize {
			panic("reader error")
		}
		ac.mode = cipher.NewCBCDecrypter(block, ac.iv[:])
	} else {
		copy(ac.iv[:], RandStringBytes(kAesBlockSize))
		ac.retIV = true
		ac.mode = cipher.NewCBCEncrypter(block, ac.iv[:])
	}
	return nil
}

func (ac AesCipher) Read(output []byte) (int, error) {
	size := len(output)
	buff := make([]byte, size)

	i := 0
	if i%kAesBlockSize != 0 {
		return 0, errors.New(fmt.Sprintf("output buff size must be multiple of %d", kAesBlockSize))
	}
	for {
		if i >= size {
			break
		}
		if ac.retIV {
			copy(buff[i:i+kAesBlockSize], ac.iv[:])
			i += kAesBlockSize
			ac.retIV = false
		} else {
			bs, err := ac.reader.Read(buff[i:])
			if err != nil {
				if err == io.EOF {
					break
				}
				return 0, err
			}
			i += bs
		}
	}
	if i != size {
		panic("reader error")
	}

	i = (i + kAesBlockSize - 1) / kAesBlockSize * kAesBlockSize
	ac.mode.CryptBlocks(output[:i], buff[:i])
	return i, nil
}
