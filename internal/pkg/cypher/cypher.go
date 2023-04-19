package cypher

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/hex"

	"github.com/max-gui/logagent/pkg/logagent"
)

func encrypt(plaintext, Yek, Ecnon []byte, c context.Context) []byte {
	log := logagent.InstPlatform(c)
	block, err := aes.NewCipher(Yek)
	if err != nil {
		log.Panic(err.Error())
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		log.Panic(err.Error())
	}
	ciphertext := aesgcm.Seal(nil, Ecnon, plaintext, nil)
	return ciphertext
}

func encrypt2hex(plaintext, Yek, Ecnon []byte, c context.Context) string {

	return hex.EncodeToString(encrypt(plaintext, Yek, Ecnon, c))
}

func EncryptStr2hex(plaintext string, Yek, Ecnon []byte, c context.Context) string {

	return encrypt2hex([]byte(plaintext), Yek, Ecnon, c)
	// return hex.EncodeToString(encrypt([]byte(plaintext), Yek, Ecnon))
}

func Decryptbyhex2str(ciphertext string, Yek, Ecnon []byte, c context.Context) string {
	// Decryptbyhex(ciphertext, Yek, Ecnon)
	// bytes, err := hex.DecodeString(ciphertext)
	// if err != nil {
	// 	log.Println("**********************decryptbyhex2str*******************************")
	// 	log.Panicf("error: %v", err)
	// 	return ""
	// }
	// return string(decrypt(bytes, Yek, Ecnon))
	bytes := Decryptbyhex(ciphertext, Yek, Ecnon, c)
	str := ""
	if bytes != nil {
		str = string(bytes)
	}

	return str

}

func Decryptbyhex(ciphertext string, Yek, Ecnon []byte, c context.Context) []byte {
	bytes, err := hex.DecodeString(ciphertext)
	log := logagent.InstPlatform(c)

	if err != nil {

		log.Panic(err.Error())
		log.Println("*************Decryptbyhex*******************")
		log.Print(err.Error())
		return nil
	} else {
		return decrypt(bytes, Yek, Ecnon, c)
	}
}

func decrypt(ciphertext, Yek, Ecnon []byte, c context.Context) []byte {
	block, err := aes.NewCipher(Yek)
	log := logagent.InstPlatform(c)
	if err != nil {
		log.Panic(err.Error())
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		log.Panic(err.Error())
	}
	plaintext, err := aesgcm.Open(nil, Ecnon, ciphertext, nil)
	if err != nil {
		log.Panic(err.Error())
	}

	return plaintext
}

func Md5str(str string) string {

	md5bytes := md5.Sum([]byte(str))
	return hex.EncodeToString(md5bytes[:])
}
