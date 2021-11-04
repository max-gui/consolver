package router

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/max-gui/consolver/internal/confgen"
	"github.com/max-gui/consolver/internal/pkg/constset"
	"github.com/max-gui/fileconvagt/pkg/convertops"
	"github.com/max-gui/fileconvagt/pkg/fileops"
	"github.com/stretchr/testify/assert"
)

var plaintext, cryptedHexText, abstestpath, md5Hex string
var router *gin.Engine

func setup() {
	gin.SetMode(gin.TestMode)
	plaintext = "123"
	cryptedHexText = "1bda1896724a4521cfb7f38646824197929cd1"
	md5Hex = "202cb962ac59075b964b07152d234b70"
	c := context.Background()
	constset.StartupInit(nil, c)
	// constset.Consul_host = "http://127.0.0.1:8500"
	abstestpath = confgen.Makeconfiglist(context.Background())

	router = SetupRouter()
	// fmt.Println(config.AppSetting.JwtSecret)
	// fmt.Println("Before all tests")
}

func teardown() {

}

// func Test_Cases(t *testing.T) {
// 	// <setup code>
// 	setup()

// 	t.Run("Gen=all", Test_Gen4all)
// 	t.Run("Get=all", Test_Get4all)
// 	t.Run("Get=env", Test_Get4env)
// 	t.Run("fileToken=fileToken", Test_fileTokenOnline)
// 	t.Run("Encrypt=hex", Test_Encrypt2Hexonline)
// 	t.Run("Decrypt=hex", Test_DecryptHexonline)
// 	t.Run("Encrypt=config2hex", Test_EncryptConfig2Hexonline)
// 	t.Run("Decrypt=hexConfig", Test_DecryptHexConfigOnline)

// 	// <tear-down code>
// 	teardown()
// }

func body4PostFile(filefulname string, t *testing.T) (*bytes.Buffer, *multipart.Writer) {
	// filefulname := "yamls" + string(os.PathSeparator) + "orgconfig" + string(os.PathSeparator) + "pg-pgcypher-sit.yaml"

	file, err := os.Open(filefulname)
	if err != nil {
		t.Error(err)
	}
	defer file.Close()

	body := &bytes.Buffer{}

	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", filefulname)
	if err != nil {
		t.Error(err)
	}
	_, err = io.Copy(part, file)
	if err != nil {
		t.Error(err)
	}
	_ = writer.Close()
	return body, writer
}

func Test_Gen4all(t *testing.T) {

	//read test dir
	PthSep := string(os.PathSeparator)
	dirPth := abstestpath + PthSep + "pathtest"
	dir, err := ioutil.ReadDir(dirPth)
	if err != nil {
		log.Println(err.Error())
	}

	//read test files in test dir
	testpath := []string{}
	var filefullname, resstr, srcPaths string
	// router := SetupRouter()

	for _, fi := range dir {
		if fi.IsDir() {
			filefullname = dirPth + PthSep + fi.Name() + PthSep + "application-var.yml"
			srcPaths += filefullname + ","
			testpath = append(testpath, filefullname)
		}
	}
	srcPaths = strings.TrimSuffix(srcPaths, ",")

	//send request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/generateConfig?srcPaths="+srcPaths, nil)
	router.ServeHTTP(w, req)
	result := w.Result()
	defer result.Body.Close()
	resbody, _ := ioutil.ReadAll(result.Body)

	resstr = string(resbody)
	assert.Equal(t, http.StatusOK, result.StatusCode)

	//test result
	for _, fn := range testpath {
		fn0(filefullname, fn, t)
	}

	// target_env := "sit"
	// filefullname := "yamls" + string(os.PathSeparator) + "orgconfig" + string(os.PathSeparator) + "application-test-var.yml"

	// if !convertops.CompareTwoMapInterface(resstr, orgstr) {
	// 	t.Fatalf("Test_Get4env failed! config should be:\n%s \nget:\n%s", orgstr, resstr)
	// }
	t.Logf("Test_Gen4all result is:\n%s", resstr)
}

func fn0(filefullname string, fn string, t *testing.T) {
	// filefullname = fn

	content, _ := fileops.Read(filefullname)

	fntesthead := strings.TrimSuffix(filefullname, ".yml")
	c := context.Background()
	orgmap, _ := confgen.GenerateConfigContentList(content, constset.EnvSet, func(appnames string, content map[string]interface{}, env string) (map[string]interface{}, error) {
		return content, nil
	}, c)
	var fntest string
	for _, v := range constset.EnvSet {
		fntest = fntesthead + "-" + v + ".yml"
		log.Println("filename:" + fntest)
		content, _ := fileops.Read(fntest)
		mm := confgen.ConvertMap4Json(convertops.ConvertYamlToMap(content, c), func(ciphertext string, Yek, Ecnon []byte, c context.Context) string { return ciphertext }, c)
		assert.Equal(t, orgmap[v].(map[string]interface{})["af-arch"], mm["af-arch"])
	}
}

func Test_Get4all(t *testing.T) {
	// router := SetupRouter()
	// target_env := "sit"
	filefullname := abstestpath + string(os.PathSeparator) + "orgconfig" + string(os.PathSeparator) + "application-test-var.yml"
	body, writer := body4PostFile(filefullname, t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/genall", body)
	req.Header.Add("Content-type", writer.FormDataContentType())

	router.ServeHTTP(w, req)
	result := w.Result()
	defer result.Body.Close()
	resbody, _ := ioutil.ReadAll(result.Body)

	assert.Equal(t, http.StatusOK, result.StatusCode)
	// pg-pgcypher-sit
	// log.Println(string(resbody))

	// orgf, _ := fileops.Read("yamls" + string(os.PathSeparator) + "orgconfig" + string(os.PathSeparator) + "pg-pgcypher-sit.yaml")
	// ortfm := confgen.ConvertMap4Json(convertops.ConvertYamlToMap(orgf), func(ciphertext string, Yek, Ecnon []byte) string { return "cypher=" + ciphertext })

	resstr := string(resbody)
	resmap := make(map[string]interface{})
	json.Unmarshal(resbody, &resmap)

	content, _ := fileops.Read(filefullname)

	c := context.Background()
	orgstr, _ := confgen.GenerateConfigContentList(content, constset.EnvSet, func(appname string, content map[string]interface{}, env string) (map[string]interface{}, error) {
		return content, nil
	}, c)
	// resdata := resmap["data"].(map[string]interface{})
	// resmapdata := confgen.ConvertMap4Json(convertops.ConvertStrMapToYaml(&resdata), func(ciphertext string, Yek, Ecnon []byte) string { return ciphertext })
	bytes, _ := json.Marshal(orgstr)
	orgmap := make(map[string]interface{})
	json.Unmarshal(bytes, &orgmap)
	assert.Equal(t, orgmap, resmap["data"])
	// if !convertops.CompareTwoMapInterface(resstr, orgstr) {
	// 	t.Fatalf("Test_Get4env failed! config should be:\n%s \nget:\n%s", orgstr, resstr)
	// }
	t.Logf("Test_Get4all result is:\n%s", resstr)
}

func Test_Get4env(t *testing.T) {

	target_env := "sit"
	filefullname := abstestpath + string(os.PathSeparator) + "orgconfig" + string(os.PathSeparator) + "application-test-var.yml"
	body, writer := body4PostFile(filefullname, t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/gen/"+target_env, body)
	req.Header.Add("Content-type", writer.FormDataContentType())

	router.ServeHTTP(w, req)
	result := w.Result()
	defer result.Body.Close()
	resbody, _ := ioutil.ReadAll(result.Body)

	assert.Equal(t, http.StatusOK, result.StatusCode)
	// pg-pgcypher-sit
	// log.Println(string(resbody))

	// orgf, _ := fileops.Read("yamls" + string(os.PathSeparator) + "orgconfig" + string(os.PathSeparator) + "pg-pgcypher-sit.yaml")
	// ortfm := confgen.ConvertMap4Json(convertops.ConvertYamlToMap(orgf), func(ciphertext string, Yek, Ecnon []byte) string { return "cypher=" + ciphertext })

	resstr := string(resbody)
	// m := confgen.ConvertMap4Json(convertops.ConvertYamlToMap(resstr), func(ciphertext string, Yek, Ecnon []byte) string { return "cypher=" + ciphertext })
	// ReadFrom(f)
	content, _ := fileops.Read(filefullname)

	c := context.Background()
	orgstr := confgen.GenerateConfigString(content, target_env, c)

	assert.Equal(t, orgstr, resstr)
	// if !convertops.CompareTwoMapInterface(resstr, orgstr) {
	// 	t.Fatalf("Test_Get4env failed! config should be:\n%s \nget:\n%s", orgstr, resstr)
	// }
	t.Logf("Test_Get4env result is:\n%s", resstr)
}

func Test_Encrypt2Hexonline(t *testing.T) {
	// router := SetupRouter()
	jsonmap := make(map[string]interface{})
	jsonmap["data"] = plaintext
	jsonByte, _ := json.Marshal(jsonmap)
	// body, writer := body4PostFile("yamls"+string(os.PathSeparator)+"orgconfig"+string(os.PathSeparator)+"pg-pgcypher-sit.yaml", t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/encrypt2Hex", bytes.NewReader(jsonByte))
	// req.Header.Add("Content-type", writer.FormDataContentType())

	router.ServeHTTP(w, req)
	result := w.Result()
	defer result.Body.Close()
	resbody, _ := ioutil.ReadAll(result.Body)

	resstr := string(resbody)
	resjsonmap := make(map[string]interface{})
	json.Unmarshal(resbody, &resjsonmap)
	assert.Equal(t, http.StatusOK, result.StatusCode)
	assert.Equal(t, cryptedHexText, resjsonmap["data"])

	t.Logf("Test_DecryptHexonline result is:\n%s", resstr)
}

func Test_fileTokenOnline(t *testing.T) {
	// router := SetupRouter()
	jsonmap := make(map[string]interface{})
	jsonmap["data"] = plaintext
	jsonByte, _ := json.Marshal(jsonmap)
	// body, writer := body4PostFile("yamls"+string(os.PathSeparator)+"orgconfig"+string(os.PathSeparator)+"pg-pgcypher-sit.yaml", t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/fileToken", bytes.NewReader(jsonByte))
	// req.Header.Add("Content-type", writer.FormDataContentType())

	router.ServeHTTP(w, req)
	result := w.Result()
	defer result.Body.Close()
	resbody, _ := ioutil.ReadAll(result.Body)

	resstr := string(resbody)
	resjsonmap := make(map[string]interface{})
	json.Unmarshal(resbody, &resjsonmap)
	assert.Equal(t, http.StatusOK, result.StatusCode)
	assert.Equal(t, md5Hex, resjsonmap["data"])

	t.Logf("Test_DecryptHexonline result is:\n%s", resstr)
}

func Test_DecryptHexonline(t *testing.T) {
	// router := SetupRouter()
	jsonmap := make(map[string]interface{})
	jsonmap["data"] = cryptedHexText
	jsonByte, _ := json.Marshal(jsonmap)
	// body, writer := body4PostFile("yamls"+string(os.PathSeparator)+"orgconfig"+string(os.PathSeparator)+"pg-pgcypher-sit.yaml", t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/decryptHex", bytes.NewReader(jsonByte))
	// req.Header.Add("Content-type", writer.FormDataContentType())

	router.ServeHTTP(w, req)
	result := w.Result()
	defer result.Body.Close()
	resbody, _ := ioutil.ReadAll(result.Body)

	resstr := string(resbody)
	resjsonmap := make(map[string]interface{})
	json.Unmarshal(resbody, &resjsonmap)
	assert.Equal(t, http.StatusOK, result.StatusCode)
	assert.Equal(t, plaintext, resjsonmap["data"])

	t.Logf("Test_DecryptHexonline result is:\n%s", resstr)
}

func Test_EncryptConfig2Hexonline(t *testing.T) {
	// router := SetupRouter()
	c := context.Background()
	body, writer := body4PostFile(abstestpath+string(os.PathSeparator)+"orgconfig"+string(os.PathSeparator)+"pg-plain-sit.yaml", t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/encryptConfig2Hex", body)
	req.Header.Add("Content-type", writer.FormDataContentType())

	router.ServeHTTP(w, req)
	result := w.Result()
	defer result.Body.Close()
	resbody, _ := ioutil.ReadAll(result.Body)

	assert.Equal(t, http.StatusOK, result.StatusCode)
	// pg-pgcypher-sit
	// log.Println(string(resbody))

	orgf, _ := fileops.Read(abstestpath + string(os.PathSeparator) + "orgconfig" + string(os.PathSeparator) + "pg-pgcypher-sit.yaml")
	ortfm := confgen.ConvertMap4Json(convertops.ConvertYamlToMap(orgf, c), func(ciphertext string, Yek, Ecnon []byte, c context.Context) string { return "cypher=" + ciphertext }, c)

	resstr := string(resbody)
	m := confgen.ConvertMap4Json(convertops.ConvertYamlToMap(resstr, c), func(ciphertext string, Yek, Ecnon []byte, c context.Context) string { return "cypher=" + ciphertext }, c)
	// ReadFrom(f)
	if !convertops.CompareTwoMapInterface(m, ortfm) {
		t.Fatalf("Test_EncryptConfig2Hexonline failed! config should be:\n%s \nget:\n%s", ortfm, m)
	}
	t.Logf("Test_EncryptConfig2Hexonline result is:\n%s", resstr)
}

func Test_DecryptHexConfigOnline(t *testing.T) {
	// router := SetupRouter()
	c := context.Background()
	body, writer := body4PostFile(abstestpath+string(os.PathSeparator)+"orgconfig"+string(os.PathSeparator)+"pg-pgcypher-sit.yaml", t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/decryptHexConfig", body)
	req.Header.Add("Content-type", writer.FormDataContentType())

	router.ServeHTTP(w, req)
	result := w.Result()
	defer result.Body.Close()
	resbody, _ := ioutil.ReadAll(result.Body)

	assert.Equal(t, http.StatusOK, result.StatusCode)
	// pg-pgcypher-sit
	// log.Println(string(resbody))

	orgf, _ := fileops.Read(abstestpath + string(os.PathSeparator) + "orgconfig" + string(os.PathSeparator) + "pg-plain-sit.yaml")
	ortfm := confgen.ConvertMap4Json(convertops.ConvertYamlToMap(orgf, c), func(ciphertext string, Yek, Ecnon []byte, c context.Context) string { return ciphertext }, c)

	resstr := string(resbody)
	m := confgen.ConvertMap4Json(convertops.ConvertYamlToMap(resstr, c), func(ciphertext string, Yek, Ecnon []byte, c context.Context) string { return ciphertext }, c)
	// ReadFrom(f)
	if !convertops.CompareTwoMapInterface(m, ortfm) {
		t.Fatalf("Test_decryptHexConfig failed! config should be:\n%s \nget:\n%s", ortfm, m)
	}
	t.Logf("Test_decryptHexConfig result is:\n%s", resstr)

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
