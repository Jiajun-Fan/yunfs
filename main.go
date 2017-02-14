package main

import (
	"strings"
	"fmt"
)

func main() {
	reader := strings.NewReader("This is good text!")
	encryptor := AesCipher{}
	encryptor.Encrypt([]byte(`fjj30891`), reader)

	decryptor := AesCipher{}
	decryptor.Decrypt([]byte(`fjj30891`), encryptor)

	buff := make([]byte, 4096)
	n, err := decryptor.Read(buff)
	if err != nil {
		Fatal(err.Error())
	}
	fmt.Printf("%s", string(buff))
	println(n)
}
