package cypher

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/max-gui/consolver/internal/pkg/constset"
	"github.com/stretchr/testify/assert"
)

var plaintext, cryptedHexText, md5Hex string

func setup() {
	plaintext = "123"
	cryptedHexText = "1bda1896724a4521cfb7f38646824197929cd1"
	md5Hex = "202cb962ac59075b964b07152d234b70"
}

func teardown() {

}

// func Test_Cases(t *testing.T) {
// 	// <setup code>
// 	setup()

// 	t.Run("Encrypt=Str2hex", Test_EncryptStr2hex)
// 	t.Run("Decrypt=hex2str", Test_Decryptbyhex2str)
// 	t.Run("Decrypt=hex2byte", Test_Decryptbyhex)
// 	// t.Run("Write=ExistedFile", Test_Write)
// 	// t.Run("Write=WhiteToPath", Test_WhiteToPath)
// 	// <tear-down code>
// 	teardown()
// }

func Test_EncryptStr2hex(t *testing.T) {
	c := context.Background()
	str := EncryptStr2hex(plaintext, constset.Yek, constset.Ecnon, c)
	// assert.NoError(t, err, "read is ok")
	assert.Equal(t, cryptedHexText, str)
	log.Printf("Test_EncryptStr2hex result is:\n%s", str)

}

func Test_Decryptbyhex2str(t *testing.T) {
	c := context.Background()
	str := Decryptbyhex2str(cryptedHexText, constset.Yek, constset.Ecnon, c)

	assert.Equal(t, plaintext, str)
	log.Printf("Test_Decryptbyhex2str result is:\n%s", str)

}

func Test_Decryptbyhex(t *testing.T) {
	c := context.Background()
	bytes := Decryptbyhex(cryptedHexText, constset.Yek, constset.Ecnon, c)
	// if strings.Compare(plaintext, string(bytes)) != 0 {
	// 	t.Fatal("Test_Decryptbyhex failed!")
	// }
	assert.Equal(t, plaintext, string(bytes))
	log.Printf("Test_Decryptbyhex2str result is:\n%s", bytes)
}

func Test_md5str(t *testing.T) {

	str := Md5str(plaintext)
	// if strings.Compare(plaintext, string(bytes)) != 0 {
	// 	t.Fatal("Test_Decryptbyhex failed!")
	// }
	assert.Equal(t, md5Hex, str)
	log.Printf("Test_Decryptbyhex2str result is:\n%s", str)
}

func TestMain(m *testing.M) {
	setup()
	// constset.StartupInit()
	// sendconfig2consul()
	// configgen.Getconfig = getTestConfig

	exitCode := m.Run()
	teardown()
	// // 退出
	os.Exit(exitCode)
}
