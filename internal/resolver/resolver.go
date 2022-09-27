package resolvergen

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/max-gui/consolver/internal/confgen"
	"github.com/max-gui/fileconvagt/pkg/convertops"
	"github.com/max-gui/fileconvagt/pkg/fileops"
	"github.com/max-gui/logagent/pkg/logagent"
)

// func(entityType interface{}, entityId interface{}, env string, resultMap map[interface{}]interface{}, tmplMap map[interface{}]interface{}) map[string]interface{}
// var resolveinfra = resolveInfrastructure
// var generateConfFromMap = generateConfigContentFromTmplMap

/*
*
根据模板yaml内容生产配置文件内容
*/
func generateConfigFileContent(tmplYaml string, c context.Context) (string, string, string) {
	m := convertops.ConvertYamlToMap(tmplYaml, c)

	prodResult := generateConfigContentFromTmplMap(m, "prod", c)
	uatResult := generateConfigContentFromTmplMap(m, "uat", c)
	sitResult := generateConfigContentFromTmplMap(m, "sit", c)

	prodSpring := prodResult["spring"].(map[interface{}]interface{})
	prodSpring["profiles"] = "prod"
	replaceRadm(prodResult, "prod")

	sitSpring := sitResult["spring"].(map[interface{}]interface{})
	sitSpring["profiles"] = "test"
	replaceRadm(sitResult, "sit")

	uatSpring := uatResult["spring"].(map[interface{}]interface{})
	uatSpring["profiles"] = "uat"
	replaceRadm(uatResult, "uat")

	return convertops.ConvertMapToYaml(&prodResult, c), convertops.ConvertMapToYaml(&uatResult, c), convertops.ConvertMapToYaml(&sitResult, c)
}

func replaceRadm(m map[interface{}]interface{}, env string) {
	if m["radm"] != nil {
		radmValue := m["radm"].(map[string]interface{})
		for Yek := range radmValue {
			if Yek == "env" {
				if env == "prod" {
					radmValue["env"] = "PRO"
				} else if env == "uat" {
					radmValue["env"] = "STG"
				} else if env == "sit" {
					radmValue["env"] = "STG"
				}
			}
			if Yek == "app" {
				appValue := radmValue["app"].(map[string]interface{})
				if appValue["name"] != nil {
					nameValue := appValue["name"].(string)
					if env == "prod" {
						appValue["name"] = nameValue + "-prod"
					} else if env == "uat" {
						appValue["name"] = nameValue + "-uat"
					} else if env == "sit" {
						appValue["name"] = nameValue + "-sit"
					}
				}
			}
		}
	}
}

func generateConfigContentFromTmplMap(tmplMap map[interface{}]interface{}, env string, c context.Context) map[interface{}]interface{} {
	resultMap := make(map[interface{}]interface{})

	if tmplMap["entityType"] != nil && tmplMap["entityId"] != nil {
		resolveInfrastructure(resultMap, tmplMap, env, c)
	} else {
		for Yek, value := range tmplMap {
			if !strings.HasPrefix(convertops.StrValOfType(value), "map") {
				resultMap[Yek] = value
			} else {
				resultMap[Yek] = generateConfigContentFromTmplMap(value.(map[interface{}]interface{}), env, c)
			}
		}
	}
	return resultMap
}

/*
*
解析资源，将tmplMap中的内容处理后赋值到resultMap中
*/
func resolveInfrastructure(resultMap map[interface{}]interface{}, tmplMap map[interface{}]interface{}, env string, c context.Context) {
	entityType := tmplMap["entityType"]
	entityId := tmplMap["entityId"]
	fmt.Println(entityType, entityId)

	infraInfo := confgen.Getconfig(entityType, entityId, env, c)

	deepCopyMap(resultMap, tmplMap)

	switch entityType {
	case "rabbitmq":
		resultMap["host"] = infraInfo["host"]
		resultMap["port"] = infraInfo["port"]
		resultMap["username"] = infraInfo["username"]
		resultMap["password"] = infraInfo["password"]
		resultMap["virtual-host"] = infraInfo["virtual-host"]
	case "consul":
		resultMap["host"] = infraInfo["host"]
		resultMap["port"] = infraInfo["port"]
		resultMap["acl-token"] = infraInfo["acl-token"]
	case "eureka":
		// serviceUrl := make(map[string]string)
		// serviceUrl["defaultZone"] = infraInfo["eurekaUrl"].(string)
		// client := make(map[string]interface{})
		// client["serviceUrl"] = serviceUrl
		// resultMap["client"] = client
		resultMap["client"] = infraInfo["client"]
		resultMap["instance"] = infraInfo["instance"]

	case "pg", "oracle", "mysql":
		// if val, ok := infraInfo["url"]; ok {
		// 	resultMap["url"] = val
		// }else{
		// 	resultMap["jdbc-url"] = infraInfo["jdbc-url"]
		// }
		resultMap["url"] = infraInfo["url"]
		resultMap["jdbc-url"] = infraInfo["url"]
		resultMap["password"] = infraInfo["password"]
		resultMap["username"] = infraInfo["username"]
	case "neo4j":
		resultMap["url"] = infraInfo["url"]
		resultMap["uri"] = infraInfo["uri"]
		resultMap["password"] = infraInfo["password"]
		resultMap["username"] = infraInfo["username"]
		resultMap["urls"] = infraInfo["urls"]
	case "xxl":
		// job := resultMap["job"].(map[interface{}]interface{})
		// admin := job["admin"].(map[interface{}]interface{})
		// admin["addresses"] = infraInfo["addresses"]
		// admin["appId"] = infraInfo["appId"]
		// admin["Yek"] = infraInfo["Yek"]
		resultMap["addresses"] = infraInfo["addresses"]
		resultMap["appId"] = infraInfo["appId"]
		resultMap["Yek"] = infraInfo["Yek"]
	case "redis":
		// resultMap["host"] = infraInfo["host"]
		// resultMap["port"], _ = strconv.Atoi(infraInfo["port"].(string))
		// resultMap["password"] = infraInfo["password"]
		// if resultMap["sentinel"] != nil {
		// 	resultMap["sentinel"] = make(map[interface{}]interface{})
		// 	sentinelMap := resultMap["sentinel"].(map[interface{}]interface{})
		// 	sentinelMap["master"] = infraInfo["sentinel.master"]
		// 	sentinelMap["nodes"] = infraInfo["sentinel.nodes"]
		// }
		resultMap["host"] = infraInfo["host"]
		resultMap["port"] = infraInfo["port"]
		resultMap["password"] = infraInfo["password"]
		if resultMap["sentinel"] != nil {
			resultMap["sentinel"] = infraInfo["sentinel"]
		}
	}
}

/*
*
copy原有的map，顺便排查掉entityType和entityId这两个Yek
*/
func deepCopyMap(dest map[interface{}]interface{}, src map[interface{}]interface{}) {
	for Yek, value := range src {
		if convertops.StrValOfInterface(Yek) != "entityType" && convertops.StrValOfInterface(Yek) != "entityId" {
			dest[Yek] = value
		}
	}
}

func generateConfigFile(path string, c context.Context) error {
	appTmplYaml, _ := fileops.Read(path)
	log := logagent.InstArch(c)
	log.Printf("app config template content: \n%s", appTmplYaml)

	prodConfigFileContent, uatConfigFileContent, sitConfigFileContent := generateConfigFileContent(appTmplYaml, c)
	log.Printf("prod config file content: \n%s", prodConfigFileContent)
	log.Printf("uat config file content: \n%s", uatConfigFileContent)
	log.Printf("sit config file content: \n%s", sitConfigFileContent)

	lastIndex := strings.LastIndex(path, "/")
	configFilePath := path[0:lastIndex]
	fileName := path[lastIndex+1:]

	dotIndex := strings.LastIndex(fileName, ".")

	fileops.Write(configFilePath, fileName[0:dotIndex]+"-prod.yml", prodConfigFileContent, c)
	fileops.Write(configFilePath, fileName[0:dotIndex]+"-uat.yml", uatConfigFileContent, c)
	fileops.Write(configFilePath, fileName[0:dotIndex]+"-test.yml", sitConfigFileContent, c)
	return nil
}

func GenerateConfigvVOld(c *gin.Context) {
	srcPaths := c.Query("srcPaths")
	log := logagent.InstArch(c)
	log.Printf("the srcPaths are %s", srcPaths)

	pathArray := strings.Split(srcPaths, ",")

	var err error = nil
	for _, path := range pathArray {
		err = generateConfigFile(path, c)
		if err != nil {
			log.Panic(err.Error())
			break
		}
	}

	if err == nil {
		c.JSON(http.StatusOK, gin.H{
			"message": "ok",
		})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
	}
}
