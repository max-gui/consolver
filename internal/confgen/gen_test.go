package confgen

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/garyburd/redigo/redis"
	"github.com/max-gui/consolver/internal/pkg/constset"
	"github.com/max-gui/consolver/internal/pkg/cypher"
	"github.com/max-gui/consulagent/pkg/consulsets"
	"github.com/max-gui/fileconvagt/pkg/convertops"
	"github.com/max-gui/fileconvagt/pkg/fileops"
	"github.com/max-gui/logagent/pkg/confload"
	"github.com/max-gui/logagent/pkg/logsets"
	"github.com/max-gui/redisagent/pkg/redisops"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

// var Configlist []map[string]interface{}

// var args []string
var testpath string

func setup() {
	*logsets.Apppath = "/Users/jimmy/Projects/hercules/consolver"
	*logsets.Port = "8080"
	*consulsets.Consul_host = "http://consul-stg.paic.com.cn"
	*logsets.DCENV = "test"
	flag.Parse()
	bytes := confload.Load(context.Background())
	constset.StartupInit(bytes, context.Background())
	// plaintext = "123"
	// cryptedHexText = "1bda1896724a4521cfb7f38646824197929cd1"
	// constset.StartupInit()
	// testpath = Makeconfiglist()

}

func teardown() {

}

// func Test_Cases(t *testing.T) {
// 	// <setup code>
// 	setup()

// 	t.Run("Getconfig=Getconfig", Test_Getconfig)
// 	t.Run("GenerateConfig=String", Test_GenerateConfigString)
// 	t.Run("GetPostFileConfig=Encrypt", Test_GetPostFileConfigWithEncrypt)
// 	t.Run("GetPostFileConfig=Decrypt", Test_GetPostFileConfigWithDecrypt)
// 	t.Run("GenerateConfig=ContentList", Test_GenerateConfigContentList)
// 	// <tear-down code>
// 	teardown()
// }

//	func prepareTestConfigs() {
//		Makeconfiglist(func(entitytype, entityid, env, configcontent string) {
//			resp, err := consulhelp.Sendconfig2consul(entitytype, entityid, env, configcontent)
//			if err != nil {
//				fmt.Println(err.Error())
//			}
//			fmt.Println(resp)
//		})
//	}
func Test_Getconfigdd(t *testing.T) {
	bytes, _ := os.ReadFile("/Users/max/Downloads/application.yml")
	readConfigContent(string(bytes), context.Background())

}

func Test_ddd(t *testing.T) {
	redisops.StartupInit("10.25.80.6:7617", "a01cfbde22a947e9b75b7aa027ac529e")
	rediscli := redisops.Pool().Get()
	// sresult, err := redis.Values(rediscli.Do("HSCAN", "confsolver-"+appname, 0))

	var (
		cursor int64
		items  []string
	)

	// results := make([][]string, 0)
	mm := map[string]string{}
	count := 1
	for {
		values, err := redis.Values(rediscli.Do("HSCAN", "aa", cursor, "MATCH", "*", "COUNT", 1))
		if err != nil {
			return
		}
		_, err = redis.Scan(values, &cursor, &items)
		if err != nil {
			return
		}

		mm[items[0]] = items[1]
		// results = append(results, items)

		if cursor == 0 {
			break
		}
		log.Print(count)
		count = count + 1
	}

	log.Print(mm)
}

func Test_GetOnlineConfig_sameid(t *testing.T) {
	var entityId, env = "b", "test"
	var valid, valenv = entityId, env
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		} else {
			// fmt.Print(r.(*logrus.Entry).Message)
			// fmt.Print(valid)
			// fmt.Print(valenv)
			msg := fmt.Sprintf("refrence id and env cant be the same;entityid:%s,env:%s;real-id:%s,real-env:%s", entityId, env, valid, valenv)
			assert.Equal(t, msg, r.(*logrus.Entry).Message)
		}

	}()

	c := context.Background()

	GetOnlineConfig("a", "b", "test", c)
}

func Test_GetOnlineConfig_sameid_sameenv(t *testing.T) {
	var entityId, env = "b", "test"
	var valid, valenv = entityId, env
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		} else {
			// fmt.Print(r.(*logrus.Entry).Message)
			// fmt.Print(valid)
			// fmt.Print(valenv)
			msg := fmt.Sprintf("refrence id and env cant be the same;entityid:%s,env:%s;real-id:%s,real-env:%s", entityId, env, valid, valenv)
			assert.Equal(t, msg, r.(*logrus.Entry).Message)
		}

	}()

	c := context.Background()

	GetOnlineConfig("a", entityId, env, c)
}

func Test_Getconfig(t *testing.T) {
	c := context.Background()
	orgstr, _ := fileops.Read(testpath + string(os.PathSeparator) + "consul-consul1-uat.yaml")
	orgm := ConvertMap4Json(convertops.ConvertYamlToMap(orgstr, c), cypher.Decryptbyhex2str, context.Background())
	// c := context.Background()
	configm := Getconfig("consul", "consul1", "uat", c)
	// configstr := convertops.ConvertStrMapToYaml(&m)

	// assert.NoError(t, err, "read is ok")
	assert.Equal(t, orgm, configm)
	t.Logf("Test_Getconfig result is:\n%s", configm)

	// if !convertops.CompareTwoMapInterface(orgm, configm) {
	// 	t.Fatalf("Test_Getconfig failed! config should be:\n%s \nget:\n%s", orgm, configm)
	// }
	// t.Logf("Test_Getconfig is ok! result is:\n%s", configm)
}

func Test_GenerateConfigString(t *testing.T) {
	envstr := "sit"
	filestr, _ := fileops.Read(testpath + string(os.PathSeparator) + "orgconfig" + string(os.PathSeparator) + "application-test-singalvar.yml")

	c := context.Background()
	configstr := GenerateConfigString(filestr, envstr, c)
	configm := ConvertMap4Json(convertops.ConvertYamlToMap(configstr, c), cypher.Decryptbyhex2str, context.Background())

	var d map[string]interface{}
	for k, v := range configm["af-arch"].(map[string]interface{})["resource"].(map[string]interface{}) {
		if m, ok := v.(map[string]interface{}); ok {
			d = Getconfig(m["entityType"].(string), k, envstr, c)
			// configstr := convertops.ConvertStrMapToYaml(&m)
			d["entityType"] = m["entityType"]
			assert.Equal(t, v, d)
			// if !convertops.CompareTwoMapInterface(d, v.(map[string]interface{})) {
			// 	t.Fatalf("Test_Getconfig failed! config should be:\n%s \nget:\n%s", d, v)
			// }
		}
	}

	t.Logf("Test_GenerateConfigString result is:\n%s", configm)
	// t.Logf("Test_GenerateConfigString is ok! result is:\n%s", configstr)

}

func Test_GetPostFileConfigWithEncrypt(t *testing.T) {
	c := context.Background()
	f, _ := os.OpenFile(testpath+string(os.PathSeparator)+"orgconfig"+string(os.PathSeparator)+"pg-plain-sit.yaml", os.O_RDONLY, 0644)
	orgf, _ := fileops.Read(testpath + string(os.PathSeparator) + "orgconfig" + string(os.PathSeparator) + "pg-pgcypher-sit.yaml")
	ortfm := ConvertMap4Json(convertops.ConvertYamlToMap(orgf, c), func(ciphertext string, Yek, Ecnon []byte, c context.Context) string { return "cypher=" + ciphertext }, context.Background())
	// ReadFrom(f)
	m, _ := GetPostFileConfigWithEncrypt(f, context.Background())

	defer f.Close()
	assert.Equal(t, ortfm, m)
	t.Logf("Test_GetPostFileConfigWithEncrypt result is:\n%s", m)
	// if !convertops.CompareTwoMapInterface(m, ortfm) {
	// 	t.Fatalf("GetPostFileConfigWithEncrypt failed! config should be:\n%s \nget:\n%s", m, ortfm)
	// }
	// t.Logf("GetPostFileConfigWithEncrypt is ok! result is:\n%s", m)

}

func Test_GetPostFileConfigWithDecrypt(t *testing.T) {
	c := context.Background()
	f, _ := os.OpenFile(testpath+string(os.PathSeparator)+"orgconfig"+string(os.PathSeparator)+"pg-pgcypher-sit.yaml", os.O_RDONLY, 0644)
	orgf, _ := fileops.Read(testpath + string(os.PathSeparator) + "orgconfig" + string(os.PathSeparator) + "pg-pg2-sit.yaml")
	ortfm := ConvertMap4Json(convertops.ConvertYamlToMap(orgf, c), func(ciphertext string, Yek, Ecnon []byte, c context.Context) string { return ciphertext }, context.Background())
	// ReadFrom(f)
	m, _ := GetPostFileConfigWithDecrypt(f, context.Background())

	defer f.Close()
	assert.Equal(t, ortfm, m)
	t.Logf("Test_GetPostFileConfigWithDecrypt result is:\n%s", m)
	// if !convertops.CompareTwoMapInterface(m, ortfm) {
	// 	t.Fatalf("GetPostFileConfigWithDecrypt failed! config should be:\n%s \nget:\n%s", m, ortfm)
	// }
	// t.Logf("GetPostFileConfigWithDecrypt is ok! result is:\n%s", m)

}

func Test_GenerateConfigContentList(t *testing.T) {
	// envstr := "sit"
	filestr, _ := fileops.Read(testpath + string(os.PathSeparator) + "orgconfig" + string(os.PathSeparator) + "application-test-singalvar.yml")
	var f0 = func(appname string, content map[string]interface{}, env string) (map[string]interface{}, error) {
		return content, nil
	}

	c := context.Background()
	m, _ := GenerateConfigContentList(filestr, []string{"sit", "prod"}, f0, c)

	// configstr := GenerateConfigString(filestr, envstr)
	// configm := convertMap4Json(convertops.ConvertYamlToMap(configstr), cypher.Decryptbyhex2str)

	var d map[string]interface{}
	for k, v := range m {
		if m, ok := v.(map[string]interface{}); ok {
			d = ConvertMap4Json(convertops.ConvertYamlToMap(GenerateConfigString(filestr, k, c), c), cypher.Decryptbyhex2str, context.Background())

			// d = Getconfig(k, m["entityType"].(string), envstr)
			// configstr := convertops.ConvertStrMapToYaml(&m)

			assert.Equal(t, d, m)
			// if !convertops.CompareTwoMapInterface(d, m) {
			// 	t.Fatalf("Test_GenerateConfigContentList failed! config should be:\n%s \nget:\n%s", d, m)
			// }
		}
	}

	t.Logf("Test_GenerateConfigContentList result is:\n%s", m)
	// t.Logf("Test_GenerateConfigContentList is ok! result is:\n%s", m)

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
