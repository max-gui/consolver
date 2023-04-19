package confgen

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/gomodule/redigo/redis"
	"github.com/max-gui/consolver/internal/pkg/constset"
	"github.com/max-gui/consolver/internal/pkg/cypher"
	"github.com/max-gui/consulagent/pkg/consulhelp"
	"github.com/max-gui/fileconvagt/pkg/convertops"
	"github.com/max-gui/fileconvagt/pkg/fileops"
	"github.com/max-gui/logagent/pkg/logagent"
	"github.com/max-gui/redisagent/pkg/redisops"
	"gopkg.in/yaml.v2"
)

// var Getconfig = getOnlineConfig

var getFrom = GetOnlineConfig

func Getconfig(entityType interface{}, entityId interface{}, env string, c context.Context) map[string]interface{} {
	resstr := getFrom(convertops.StrValOfInterface(entityType), convertops.StrValOfInterface(entityId), env, c)
	infraInfo := ConvertMap4Json(convertops.ConvertYamlToMap(resstr, c), cypher.Decryptbyhex2str, c)

	return infraInfo //ConvertMap4Json(infraInfo, cypher.Decryptbyhex2str)
}

func GetOnlineConfig(entityType string, entityId string, env string, c context.Context) string {
	// resstr, _ := consulhelp.GetInfrastructureInfo(entityType, entityId, env)
	resstr := string(consulhelp.Getconfibytes(*constset.ConfResPrefix, entityType, entityId, env, c))

	infraInfo := make(map[string]interface{})
	err := yaml.Unmarshal([]byte(resstr), &infraInfo)

	logger := logagent.InstPlatform(c)
	if err == nil {
		if valid, ok := infraInfo["real-id"]; ok {
			if valenv, ok := infraInfo["real-env"]; ok {
				logger.Print(valid)
				logger.Print(valenv)
				if convertops.StrValOfInterface(valid) == entityId && convertops.StrValOfInterface(valenv) == env {
					logger.Panicf("refrence id and env cant be the same;entityid:%s,env:%s;real-id:%s,real-env:%s", entityId, env, valid, valenv)
				}
				resstr = GetOnlineConfig(entityType, convertops.StrValOfInterface(valid), convertops.StrValOfInterface(valenv), c)
				// resstr, err := GetFrom(convertops.StrValOfInterface(entityType), convertops.StrValOfInterface(entityId), env)
			} else {
				if convertops.StrValOfInterface(valid) == entityId {
					logger.Panicf("refrence id and env cant be the same;entityid:%s,env:%s;real-id:%s,real-env:%s", entityId, env, valid, valenv)
				}
				resstr = GetOnlineConfig(entityType, convertops.StrValOfInterface(valid), env, c)
			}
		}
	}
	return resstr
	// return ConvertMap4Json(convertops.ConvertYamlToMap(yamlString), cypher.Decryptbyhex2str)
}

// func GetFileConfig(entityType string, entityId string, env string) (string, error) {
// 	filepath := "yamls/" + entityType + "-" + entityId + "-" + env + ".yaml"
// 	content, err := ioutil.ReadFile(filepath)
// 	if err != nil {
// 		fmt.Println(filepath)
// 		panic(err)
// 	}
// 	return string(content), err
// 	// fmt.Println(string(content))

// 	// infraInfo := ConvertMap4Json(convertops.ConvertYamlToMap(string(content)), cypher.Decryptbyhex2str)
// 	// return infraInfo
// }

func GenerateConfigString(appTmplYaml string, env string, c context.Context) string {

	fmt.Println("--------------------generateConfigString-----------------------")
	resultMap := generateConfigContent(appTmplYaml, env, c)
	resultStr := convertops.ConvertStrMapToYaml(&resultMap, c)
	logger := logagent.InstPlatform(c)
	logger.Printf("%s config file content: \n%s", env, resultStr)
	return resultStr
}

func generateConfigContent(appTmplYaml string, env string, c context.Context) map[string]interface{} {
	contentMap, _ := readConfigContent(appTmplYaml, c)
	resultMap := configContentOutput(env, contentMap, c)
	// fmt.Println(resultMap)

	return resultMap
}

func configContentOutput(env string, contentMap map[string]interface{}, c context.Context) map[string]interface{} {
	resultMap := make(map[string]interface{})

	afarch := make(map[string]interface{})
	resultMap["af-arch"] = afarch
	resource := make(map[string]interface{})
	resource["env"] = env
	afarch["resource"] = resource
	var entityids map[string]interface{}
	if val, ok := contentMap["af-arch"]; ok {
		if val != nil {
			if val, ok = val.(map[string]interface{})["resource"]; ok {
				if entityids, ok = val.(map[string]interface{}); ok {
					var d map[string]interface{}
					for k, v := range entityids {
						d = Getconfig(v, k, env, c)
						d["entityType"] = v

						k = strings.ReplaceAll(k, ".", "-")
						resource[k] = d
					}
				}
			}
		}
	}
	return resultMap
}

func getResEntity(appname string, c context.Context) map[string]interface{} { //}, string) {
	bytes := consulhelp.GetConfigFull(*constset.ConfArchPrefix+appname, c) //.HttpGetBytes("http://" + iacconfmap["url"].(string) + "/info/resource/project0/team0/" + appname)

	var resmap = map[string]interface{}{}
	json.Unmarshal(bytes, &resmap)

	logger := logagent.InstPlatform(c)
	logger.Print(resmap)

	if val, ok := resmap["Application"]; ok {
		if val != nil {
			if val, ok = val.(map[string]interface{})["Resource"]; ok {
				logger.Print(val)
				if entityids, ok := val.(map[string]interface{}); ok {
					// appyml := resmap["Deploy"].(map[string]interface{})["Build"].(map[string]interface{})["Appyml"].(string)
					return entityids //, appyml
				}
			}
		}
	}

	return nil //, ""
}

func configContentOutputwithResmap(env string, entityids map[string]interface{}, tags map[string]map[string]string, c context.Context) map[string]interface{} {
	resultMap := make(map[string]interface{})

	afarch := make(map[string]interface{})
	resultMap["af-arch"] = afarch
	resource := make(map[string]interface{})
	resource["env"] = env
	afarch["resource"] = resource
	logger := logagent.InstPlatform(c)
	logger.Print("ssssssssssssssssssssssssssssssssssssssssssssssssssssssssssss")
	logger.Print(entityids)

	var d map[string]interface{}
	for k, v := range entityids {
		d = Getconfig(v, k, env, c)
		d["entityType"] = v
		k = strings.ReplaceAll(k, ".", "-")
		resource[k] = d
		var confver string
		if ver, ok := d["version"]; ok {
			confver = ver.(string)
		} else {
			confver = "0"
		}

		tagkey := v.(string) + "_" + k
		if _, ok := tags[env]; ok {
			tags[env][tagkey] = "v" + confver
		} else {
			tags[env] = map[string]string{tagkey: "v" + confver}
		}
		// tags[tagkey] = v.(string) + "_" + k + "_" + confver
		// if val, ok := tags[tagkey]; ok {
		// 	tags[tagkey] = val + "_" + env + confver
		// } else {
		// 	tags[tagkey] = env + confver
		// }
	}

	return resultMap
}

func readConfigContent(appTmplYaml string, c context.Context) (map[string]interface{}, string) {
	logger := logagent.InstPlatform(c)
	logger.Printf("app config template content: \n%s", appTmplYaml)

	contentMap := ConvertMap4Json(convertops.ConvertYamlToMap(appTmplYaml, c), cypher.Decryptbyhex2str, c)

	appname := readAppname(contentMap, c)
	// appname := contentMap["arch"].(map[string]interface{})["appname"].(string)
	// team := contentMap["arch"].(map[string]interface{})["team"].(string)
	// proj := contentMap["arch"].(map[string]interface{})["proj"].(string)
	// if strings.Compare(appname, "") == 0 {
	// 	logger.Panic("spring.application.name is not found")
	// }
	return contentMap, appname
}

func readAppname(contentMap map[string]interface{}, c context.Context) string {

	logger := logagent.InstPlatform(c)
	defer func() {
		if e := recover(); e != nil {

			logger.Panic("missing spring.application.name\n" + fmt.Sprint(e))
		}
	}()

	appname := contentMap["spring"].(map[string]interface{})["application"].(map[string]interface{})["name"].(string)
	if strings.Compare(appname, "") == 0 {
		logger.Panic("spring.application.name is empty")
	}
	return appname
}

type Process4envFunc func(appname string, content map[string]interface{}, env string) (map[string]interface{}, error)

// func process4envSample(content map[string]interface{}, env string) (map[string]interface{}, error) {
// 	return nil, nil
// }

func GenerateConfigContentList(appTmplYaml string, envlist []string, Process4env Process4envFunc, c context.Context) (map[string]interface{}, error) {
	var resultlist = make(map[string]interface{})
	var d map[string]interface{}
	var err error
	// var appname string
	logger := logagent.InstPlatform(c)
	orgMap, appname := readConfigContent(appTmplYaml, c)
	for _, env := range envlist {
		d = configContentOutput(env, orgMap, c)
		// var envmap = make(map[string]interface{})
		// envmap["content"] = d
		resultlist[env] = d

		_, err = Process4env(appname, d, env)
		if err != nil {
			logger.Panic(err.Error())
			break
		}
	}

	return resultlist, err
}

func GetAppConfigContentList(appname string, envlist []string, Process4env Process4envFunc, c context.Context) (map[string]interface{}, error) {
	var resultlist = make(map[string]interface{})
	var d map[string]interface{}
	var err error
	rediscli := redisops.Pool().Get()

	defer rediscli.Close()

	logger := logagent.InstPlatform(c)
	// _, err := rediscli.Do("HSET", "confsolver-"+appname, filenamestr, writeContent)
	// contentMap := ConvertMap4Json(convertops.ConvertYamlToMap(appTmplYaml), cypher.Decryptbyhex2str)

	// sresult, err := redis.Values(rediscli.Do("HSCAN", "confsolver-"+appname, 0))

	var (
		cursor int64
		items  []string
	)

	// results := make([][]string, 0)
	confmap := map[string]string{}
	for {
		values, err := redis.Values(rediscli.Do("HSCAN", "confsolver-"+appname, cursor, "MATCH", "*", "COUNT", 1))

		if err != nil {
			logger.Panic(err)
		}

		_, err = redis.Scan(values, &cursor, &items)
		if err != nil {
			logger.Panic(err)
		}
		if len(items) > 0 && math.Mod(float64(len(items)), 2) == 0 {
			index := 0
			for {
				confmap[items[index]] = items[index+1]
				index = index + 2
				if index >= len(items) {
					break
				}
			}
		}
		// results = append(results, items)

		if cursor == 0 {
			break
		}
	}

	// for _, env := range envlist {
	// 	confstr, err := redis.String(rediscli.Do("HGET", "confsolver-"+appname, "application-"+env+".yml"))
	// 	rediscli.Do("EXPIRE", "confsolver-"+appname, 60*10)
	// 	// bytes, err := ioutil.ReadFile(constset.Confpath + appname + string(os.PathSeparator) + "application-" + env + ".yml")
	// 	if err != nil {
	// 		logger.Panic(err)
	// 	}
	// 	contentMap := ConvertMap4Json(convertops.ConvertYamlToMap(confstr, c), cypher.Decryptbyhex2str, c)
	// 	d = contentMap
	// 	// var envmap = make(map[string]interface{})
	// 	// envmap["content"] = d
	// 	resultlist[env] = d

	// 	_, err = Process4env(appname, d, env)
	// 	if err != nil {
	// 		logger.Panic(err.Error())
	// 		break
	// 	}
	// }

	rediscli.Do("EXPIRE", "confsolver-"+appname, 60*10)
	for env_key, conf_value := range confmap {

		contentMap := ConvertMap4Json(convertops.ConvertYamlToMap(conf_value, c), cypher.Decryptbyhex2str, c)
		d = contentMap
		// var envmap = make(map[string]interface{})
		// envmap["content"] = d
		resultlist[env_key] = d

		_, err = Process4env(appname, d, env_key)
		if err != nil {
			logger.Panic(err.Error())
			break
		}
	}

	return resultlist, err
}

type Process4envtagFunc func(appname string, content map[string]interface{}, tag map[string]map[string]string, env string, c context.Context) (map[string]interface{}, error)

func GenerateAppConfigContentList(appname string, envlist []string, Process4env Process4envtagFunc, c context.Context) (map[string]interface{}, map[string]map[string]string, error) {
	var resultlist = make(map[string]interface{})
	var d map[string]interface{}
	var err error
	// var appname string
	// _, appname := readConfigContent(appTmplYaml)
	// iacconf, _ := consulhelp.GetInfrastructureInfo("iac", "server", *constset.Execenv)
	// iacconfmap := make(map[string]interface{})
	// yaml.Unmarshal([]byte(iacconf), &iacconfmap)

	// bytes := consulhelp.GetConfigFull(*constset.ConfArchPrefix + appname) //.HttpGetBytes("http://" + iacconfmap["url"].(string) + "/info/resource/project0/team0/" + appname)

	// var resmap = map[string]interface{}{}
	// json.Unmarshal(bytes, &resmap)

	var resmap = getResEntity(appname, c)
	var tags = map[string]map[string]string{}

	logger := logagent.InstPlatform(c)
	// contentMap := ConvertMap4Json(convertops.ConvertYamlToMap(appTmplYaml), cypher.Decryptbyhex2str)
	for _, env := range envlist {
		d = configContentOutputwithResmap(env, resmap, tags, c)
		// var envmap = make(map[string]interface{})
		// envmap["content"] = d
		resultlist[env] = d

		_, err = Process4env(appname, d, tags, env, c)
		if err != nil {
			logger.Panic(err.Error())
			break
		}
	}

	return resultlist, tags, err
}

func GenerateConfigContentListremote(appTmplYaml string, envlist []string, Process4env Process4envFunc, c context.Context) (map[string]interface{}, error) {
	var resultlist = make(map[string]interface{})
	var d map[string]interface{}
	var err error
	// var appname string
	_, appname := readConfigContent(appTmplYaml, c)
	// iacconf, _ := consulhelp.GetInfrastructureInfo("iac", "server", *constset.Execenv)
	// iacconfmap := make(map[string]interface{})
	// yaml.Unmarshal([]byte(iacconf), &iacconfmap)

	// bytes := consulhelp.GetConfigFull(*constset.ConfArchPrefix + appname) //.HttpGetBytes("http://" + iacconfmap["url"].(string) + "/info/resource/project0/team0/" + appname)

	// var resmap = map[string]interface{}{}
	// json.Unmarshal(bytes, &resmap)
	var resmap = getResEntity(appname, c)
	var tags = map[string]map[string]string{}
	logger := logagent.InstPlatform(c)
	// contentMap := ConvertMap4Json(convertops.ConvertYamlToMap(appTmplYaml), cypher.Decryptbyhex2str)
	for _, env := range envlist {
		d = configContentOutputwithResmap(env, resmap, tags, c)
		// var envmap = make(map[string]interface{})
		// envmap["content"] = d
		resultlist[env] = d

		_, err = Process4env(appname, d, env)
		if err != nil {
			logger.Panic(err.Error())
			break
		}
	}

	return resultlist, err
}

// type convert_inf struct {
// 	content        interface{}
// 	convert_cypher func(ciphertext string, Yek, Ecnon []byte) string
// }

func ConvertMap4Json(m interface{}, cypher_func convert_cypher, c context.Context) map[string]interface{} {
	return convertInterMap4Json(m, cypher_func, c).(map[string]interface{})
}

func getPostFileConfig(file io.Reader, cypher_func convert_cypher, c context.Context) (map[string]interface{}, error) {
	// err, configstr := fileops.ReadFrom(file)
	logger := logagent.InstPlatform(c)
	if configstr, err := fileops.ReadFrom(file, c); err != nil {

		logger.Panic(err.Error())
		//convertWithCypher= encryptStr2hex
		return nil, err
	} else {
		var contentMap = ConvertMap4Json(convertops.ConvertYamlToMap(configstr, c), cypher_func, c)
		return contentMap, err
	}
}

func GetPostFileConfigWithEncrypt(file io.Reader, c context.Context) (map[string]interface{}, error) {

	return getPostFileConfig(file, convertWihtEncypher, c)
}

func GetPostFileConfigWithDecrypt(file io.Reader, c context.Context) (map[string]interface{}, error) {
	return getPostFileConfig(file, cypher.Decryptbyhex2str, c)
}

// var ConvertWithCypher func(ciphertext string, Yek, Ecnon []byte) string

func convertWihtEncypher(ciphertext string, Yek, Ecnon []byte, c context.Context) string {
	return "cypher=" + cypher.EncryptStr2hex(ciphertext, Yek, Ecnon, c)
}

type convert_cypher func(ciphertext string, Yek, Ecnon []byte, c context.Context) string

func convertInterMap4Json(m interface{}, cypher_func convert_cypher, c context.Context) interface{} {
	var res = make(map[string]interface{})
	// var ok bool
	// var strvalue map[string]interface{}
	// var intermap map[interface{}]interface{}
	var strmap map[string]interface{}
	var mstr string

	logger := logagent.InstPlatform(c)
	if intermap, ok := m.(map[interface{}]interface{}); ok {
		for k, v := range intermap {
			res[k.(string)] = convertInterMap4Json(v, cypher_func, c)
		}
	} else {

		if strmap, ok = m.(map[string]interface{}); ok {
			for k, v := range strmap {
				res[k] = convertInterMap4Json(v, cypher_func, c)
			}
		} else {
			mstr, ok = m.(string)
			if ok && strings.HasPrefix(mstr, "cypher=") {
				cypherstr := strings.TrimPrefix(mstr, "cypher=")
				logger.Println("**********************convertInterMap4Json*******************************")

				logger.Printf("cypherstr: %v", cypherstr)
				m = cypher_func(cypherstr, constset.Yek, constset.Ecnon, c)
				logger.Printf("m: %v", m)
			}
			return m
		}
	}

	return res
}

// func converSubmap4Json(v interface{}, cypher_func convert_cypher) interface{} {
// 	var value interface{}
// 	var ok bool
// 	value, ok = v.(map[interface{}]interface{})
// 	if ok {
// 		v = convertInterMap4Json(value, cypher_func)
// 	} else {
// 		value, ok = v.(map[string]interface{})
// 		if ok {
// 			v = convertInterMap4Json(value, cypher_func)
// 		}
// 	}

// 	return v
// }

func Makeconfiglist(c context.Context) string { //f0 func(entitytype, entityid, env, configcontent string)) {

	// pathname := "yamls"
	pwd, _ := os.Getwd()
	pathname := strings.Split(pwd, "consolver")[0] + "consolver/testconf" + string(os.PathSeparator) + "yamls"
	abspath, _ := filepath.Abs(pathname)
	// consulhelp.Consulurl = "http://localhost:32771"
	// consulhelp.AclToken = ""
	files, err := os.ReadDir(abspath)
	logger := logagent.InstPlatform(c)
	if err != nil {
		logger.Panic(err)
	}

	PthSep := string(os.PathSeparator)
	var filename, entitytype, entityid, env, configfilepath string
	var entityinfos []string
	// var configlist []map[string]interface{}
	// var config = make(map[string]interface{})
	for _, file := range files {
		if path.Ext(file.Name()) == ".yaml" {
			filename = strings.Split(file.Name(), ".")[0]
			// fmt.Println(filename)
			entityinfos = strings.Split(filename, "-")
			entitytype = entityinfos[0]
			// fmt.Println(entitytype)
			entityid = entityinfos[1]
			// fmt.Println(entityid)
			env = entityinfos[2]
			// fmt.Println(env)

			configfilepath = pathname + PthSep + file.Name()
			// fmt.Println(configfilepath)

			content, _ := os.ReadFile(configfilepath)

			// config := make(map[string]interface{})
			// configlist = append(configlist, config)

			// fmt.Println(string(content))
			_, err := consulhelp.Sendconfig2consul(entitytype, entityid, env, string(content), c)
			if err != nil {
				logger.Panic(err.Error())
				fmt.Println(err.Error())
			}
			// f0(entitytype, entityid, env, string(content))
			// resp, err := consulhelp.Sendconfig2consul(entitytype, entityid, env, string(content))
			// if err != nil {
			// 	fmt.Println(err.Error())
			// }
			// fmt.Println(resp)
		}
	}

	return abspath
}
