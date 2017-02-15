package main

import (
	"crypto/aes"
	"crypto/cipher"
	"errors"
	"fmt"
	"io"
)

const (
	kAesBlockSize = 16
)

type Cipher interface {
	Decrypt(key string, reader io.Reader) error
	Encrypt(key string, reader io.Reader) error
	io.Reader
}

type AesCipher struct {
	reader io.Reader
	data   *AesCipherData
}

type AesCipherData struct {
	Key   [kAesBlockSize]byte
	IV    [kAesBlockSize]byte
	Mode  cipher.BlockMode
	RetIV bool
}

func (ac *AesCipher) Decrypt(key []byte, reader io.Reader) error {
	return ac.init(true, key, reader)
}

func (ac *AesCipher) Encrypt(key []byte, reader io.Reader) error {
	return ac.init(false, key, reader)
}

func (ac *AesCipher) init(dec bool, key []byte, reader io.Reader) error {
	ac.data = &AesCipherData{}

	ac.reader = reader

	if len(key) > kAesBlockSize {
		return errors.New(fmt.Sprintf("key size must not larger than %d", kAesBlockSize))
	}
	copy(ac.data.Key[:], key)

	block, err := aes.NewCipher(ac.data.Key[:])
	if err != nil {
		return err
	}

	if dec {
		i := 0
		for {
			if i >= kAesBlockSize {
				break
			}
			bs, err := ac.reader.Read(ac.data.IV[i:])
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
		ac.data.Mode = cipher.NewCBCDecrypter(block, ac.data.IV[:])
	} else {
		copy(ac.data.IV[:], RandStringBytes(kAesBlockSize))
		ac.data.RetIV = true
		ac.data.Mode = cipher.NewCBCEncrypter(block, ac.data.IV[:])
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
		if ac.data.RetIV {
			copy(buff[i:i+kAesBlockSize], ac.data.IV[:])
			i += kAesBlockSize
			ac.data.RetIV = false
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
	if i > size {
		panic("reader error")
	}

	i = (i + kAesBlockSize - 1) / kAesBlockSize * kAesBlockSize
	ac.data.Mode.CryptBlocks(output[:i], buff[:i])
	if i != 0 {
		return i, nil
	} else {
		return 0, io.EOF
	}
}
