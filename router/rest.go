package router

import (
	"context"
	"fmt"
	"log"
	"math"
	"mime/multipart"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
	"github.com/max-gui/consolver/internal/confgen"
	"github.com/max-gui/consolver/internal/pkg/constset"
	"github.com/max-gui/consolver/internal/pkg/cypher"
	"github.com/max-gui/consolver/internal/pkg/dbops"
	"github.com/max-gui/fileconvagt/pkg/convertops"
	"github.com/max-gui/logagent/pkg/logagent"
	"github.com/max-gui/logagent/pkg/logsets"
	"github.com/max-gui/logagent/pkg/routerutil"
	"github.com/max-gui/redisagent/pkg/redisops"

	// dbops "github.com/max-gui/consolver/internal/pkg/dbops"
	resolvergen "github.com/max-gui/consolver/internal/resolver"
	"github.com/max-gui/fileconvagt/pkg/fileops"
)

func SetupRouter() *gin.Engine {
	// r := gin.Default()
	if *logsets.Appenv == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()                      //.Default()
	r.Use(routerutil.GinHeaderMiddle()) // ginHeaderMiddle())
	r.Use(routerutil.GinLogger())       //LoggerWithConfig())
	r.Use(routerutil.GinErrorMiddle())  //ginErrorMiddle())

	r.GET("/generateConfigv2", resolvergen.GenerateConfigvVOld)

	r.GET("/generateConfig", generate4all)

	r.GET("/generateConfig/iac", generate4alliac)

	r.GET("/conf/apply/:appname", generate4one)
	r.GET("/conf/gen/:appname", generate4oneiac)
	r.POST("/conf/gen/:appname", generate4oneiac)
	r.POST("/conf/gen/:appname/:team/:version/:region", generate4iac)
	r.GET("/conf/apply/:appname/:team/:envdc", getlist4one)
	r.GET("/conf/apply/:appname/:team/:envdc/:region/:version/:fname", get4one)

	r.GET("/conf/fetch/:entityType/:entityId/:env", getconf)

	r.POST("/genall/:envs", gen4all)

	r.POST("/get/all", Get4all)

	r.POST("/gen/:env", Get4env)

	r.POST("/encrypt2Hex", Encrypt2Hexonline)

	r.POST("/fileToken", fileTokenOnline)

	r.POST("/decryptHex", DecryptHexonline)

	r.POST("/encryptConfig2Hex", EncryptConfig2Hexonline)

	r.POST("/decryptHexConfig", DecryptHexConfigOnline)

	r.GET("/actuator/health", health)

	r.Use(gin.Recovery())
	return r
}
func health(c *gin.Context) {
	c.String(http.StatusOK, "online")
}

func getconf(c *gin.Context) {
	entityType := c.Param("entityType")
	entityId := c.Param("entityId")
	env := c.Param("env")

	confres := confgen.Getconfig(entityType, entityId, env, c)

	common_resp(nil, "", confres, c)
}

func get4one(c *gin.Context) {
	appname := c.Param("appname")
	team := c.Param("team")
	envdc := c.Param("envdc")
	fname := c.Param("fname")
	version := c.Param("version")
	region := c.Param("region")

	datastr := GetFromRepo(team, appname, envdc, version, region, fname, c)
	err := error(nil)
	httpstatus := http.StatusOK
	if err != nil {
		httpstatus = http.StatusInternalServerError
		datastr = err.Error()
	}
	c.String(httpstatus, datastr)
	// common_resp(nil, "", datastr, c)
}

func GetFromRepo(team, appname, envdc, version, region, fname string, c *gin.Context) string {

	// repopath := "consolver/" + team + "/" + appname + "/" + envdc + "/" + fname
	data := fileops.GetFromRepo(team, appname, envdc, version, region, fname, c)
	return data
}

func getlist4one(c *gin.Context) {
	appname := c.Param("appname")
	team := c.Param("team")
	envdc := c.Param("envdc")

	confres := Getconfiglist(appname, team, envdc, c)

	common_resp(nil, "", confres, c)
}
func Getconfiglist(appname, team, envdc string, c *gin.Context) []string {
	appconfname := GetAppConfName(envdc)

	confnamelist := []string{}
	confnamelist = append(confnamelist, appconfname)
	ConfNameHelper(envdc, constset.AppendSet, func(vtype, vid, confenv, filename string) {
		confnamelist = append(confnamelist, filename)
	})

	return confnamelist
}

func generate4one(c *gin.Context) {
	appname := c.Param("appname")
	path := c.Query("path")
	// confgen.ConvertWithCypher = cypher.Decryptbyhex2str
	err := GenerateappconfigRemotePath(appname, path, c)

	// var httpstatus = http.StatusOK
	// var errstr = ""
	// if err != nil {
	// 	httpstatus = http.StatusInternalServerError
	// 	errstr = err.Error()
	// }

	// c.JSON(httpstatus, gin.H{
	// 	// "data":  contents, //contents,
	// 	"error": errstr,
	// })

	common_resp(err, "", nil, c)
	// if err == nil {
	// 	c.JSON(200, gin.H{
	// 		"message": "OK",
	// 	})
	// } else {
	// 	c.JSON(200, gin.H{
	// 		"message": err.Error(),
	// 	})
	// }
}

func generate4iac(c *gin.Context) {
	appname := c.Param("appname")
	team := c.Param("team")
	version := c.Param("version")
	region := c.Param("region")

	envlist := []string{}
	c.BindJSON(&envlist)
	if len(envlist) < 1 {
		envlist = constset.EnvSet
	}
	// confgen.ConvertWithCypher = cypher.Decryptbyhex2str
	tags, err := GenerateappconfigRepo(appname, team, version, region, envlist, c)

	// var httpstatus = http.StatusOK
	// var errstr = ""
	// if err != nil {
	// 	httpstatus = http.StatusInternalServerError
	// 	errstr = err.Error()
	// }

	// c.JSON(httpstatus, gin.H{
	// 	// "data":  contents, //contents,
	// 	"error": errstr,
	// })

	common_resp(err, "", tags, c)
	// if err == nil {
	// 	c.JSON(200, gin.H{
	// 		"message": "OK",
	// 	})
	// } else {
	// 	c.JSON(200, gin.H{
	// 		"message": err.Error(),
	// 	})
	// }
}

func generate4oneiac(c *gin.Context) {
	appname := c.Param("appname")

	envlist := []string{}
	c.BindJSON(&envlist)
	if len(envlist) < 1 {
		envlist = constset.EnvSet
	}
	// confgen.ConvertWithCypher = cypher.Decryptbyhex2str
	tags, err := GenerateappconfigRemote(appname, envlist, c)

	// var httpstatus = http.StatusOK
	// var errstr = ""
	// if err != nil {
	// 	httpstatus = http.StatusInternalServerError
	// 	errstr = err.Error()
	// }

	// c.JSON(httpstatus, gin.H{
	// 	// "data":  contents, //contents,
	// 	"error": errstr,
	// })

	common_resp(err, "", tags, c)
	// if err == nil {
	// 	c.JSON(200, gin.H{
	// 		"message": "OK",
	// 	})
	// } else {
	// 	c.JSON(200, gin.H{
	// 		"message": err.Error(),
	// 	})
	// }
}

func generate4alliac(c *gin.Context) {
	srcPaths := c.Query("srcPaths")
	// confgen.ConvertWithCypher = cypher.Decryptbyhex2str
	err := GenerateallconfigRemote(srcPaths, c)

	// var httpstatus = http.StatusOK
	// var errstr = ""
	// if err != nil {
	// 	httpstatus = http.StatusInternalServerError
	// 	errstr = err.Error()
	// }

	// c.JSON(httpstatus, gin.H{
	// 	// "data":  contents, //contents,
	// 	"error": errstr,
	// })

	common_resp(err, "", nil, c)
	// if err == nil {
	// 	c.JSON(200, gin.H{
	// 		"message": "OK",
	// 	})
	// } else {
	// 	c.JSON(200, gin.H{
	// 		"message": err.Error(),
	// 	})
	// }
}

func generate4all(c *gin.Context) {
	srcPaths := c.Query("srcPaths")
	// confgen.ConvertWithCypher = cypher.Decryptbyhex2str
	err := Generateallconfig(srcPaths, c)

	// var httpstatus = http.StatusOK
	// var errstr = ""
	// if err != nil {
	// 	httpstatus = http.StatusInternalServerError
	// 	errstr = err.Error()
	// }

	// c.JSON(httpstatus, gin.H{
	// 	// "data":  contents, //contents,
	// 	"error": errstr,
	// })

	common_resp(err, "", nil, c)
	// if err == nil {
	// 	c.JSON(200, gin.H{
	// 		"message": "OK",
	// 	})
	// } else {
	// 	c.JSON(200, gin.H{
	// 		"message": err.Error(),
	// 	})
	// }
}

func Generateallconfig(srcPaths string, c context.Context) error {
	logger := logagent.InstPlatform(c)
	logger.Printf("the srcPaths are %s", srcPaths)

	pathArray := strings.Split(srcPaths, ",")

	_, err := writeConfigFile(pathArray, c)
	return err
}

func GenerateappconfigRemotePath(appname, path string, c context.Context) error {
	logger := logagent.InstPlatform(c)
	logger.Printf("the appname is %s", appname)
	logger.Printf("the path is %s", path)

	_, err := writeAppConfigPath(appname, path, c)
	return err
}

func GenerateappconfigRemote(appname string, envlist []string, c context.Context) (map[string]map[string]string, error) {

	logger := logagent.InstPlatform(c)
	logger.Printf("the appname is %s", appname)

	// _, tags, err := writeAppConfig(appname, team, envlist, c)
	_, tags, err := writeAppConfig(appname, envlist, c)
	return tags, err
}

func GenerateappconfigRepo(appname, team, version, region string, envlist []string, c context.Context) (map[string]map[string]string, error) {

	logger := logagent.InstPlatform(c)
	logger.Printf("the appname is %s", appname)

	// _, tags, err := writeAppConfig(appname, team, envlist, c)
	_, tags, err := writeAppConfig4Repo(appname, team, version, region, envlist, c)
	return tags, err
}

func writeAppConfigPath(appname, path string, c context.Context) (map[string]interface{}, error) {
	var contents map[string]interface{}

	var err error
	var contentstr string

	// ioutil.
	var f0 = func(appname string, content map[string]interface{}, env string) (map[string]interface{}, error) {
		// if r, ok := content["af-arch"].(map[string]interface{})["resource"].(map[string]interface{}); ok {
		// 	var resourceids []string
		// 	for k := range r {
		// 		if strings.Compare(k, "env") != 0 {
		// 			resourceids = append(resourceids, k)
		// 		}
		// 	}
		// 	if len(resourceids) != 0 {
		// 		env := r["env"].(string)

		// 		dbops.Update_appRes(appname, env, resourceids)
		// 	}
		// }

		contentstr, err = fileops.WriteToPathWithFilename(path, content, env, c)

		// logger.Println(*constset.Filepaths)
		//LogConfig/config
		writeAppendConfigWith(appname, path, env, c)

		writecontent := confgen.ConvertMap4Json(convertops.ConvertYamlToMap(contentstr, c), func(ciphertext string, Yek, Ecnon []byte, c context.Context) string { return ciphertext }, c)

		return writecontent, err
	}

	// ioutil.ReadFile(constset.Confpath + appname + string(os.PathSeparator) + appname)

	contents, err = confgen.GetAppConfigContentList(appname, constset.EnvSet, f0, c)

	return contents, err
}

func writeAppConfig4Repo(appname, team, version, region string, envdclist []string, c context.Context) (map[string]interface{}, map[string]map[string]string, error) {
	var contents map[string]interface{}

	var err error
	var contentstr string

	path := constset.Confpath + appname + string(os.PathSeparator)
	// ioutil.
	var f0 = func(appname string, content map[string]interface{}, tags map[string]map[string]string, envdc string, c context.Context) (map[string]interface{}, error) {
		// if r, ok := content["af-arch"].(map[string]interface{})["resource"].(map[string]interface{}); ok {
		// 	var resourceids []string
		// 	for k := range r {
		// 		if strings.Compare(k, "env") != 0 {
		// 			resourceids = append(resourceids, k)
		// 		}
		// 	}
		// 	if len(resourceids) != 0 {
		// 		// env := r["env"].(string)

		// 		// dbops.Update_appRes(appname, env, resourceids)
		// 	}
		// }

		writeContent := convertops.ConvertStrMapToYaml(&content, c)
		ConfNameHelper(envdc, constset.InsertSet, func(vtype, vid, confenv, filename string) {
			result := confgen.GetOnlineConfig(vtype, vid, confenv, c)
			writeContent = fmt.Sprintf("%s\n%s", writeContent, result)
			tagInsert(vtype, vid, tags, envdc)

			log.Print("writeContent is \n", writeContent)
		})

		contentstr, err = fileops.WriteToRepo(path, GetAppConfName(envdc), team, appname, writeContent, envdc, version, region, c)

		// logger.Println(*constset.Filepaths)
		//LogConfig/config
		writeAppendConfigRepo(path, appname, team, tags, envdc, version, region, c)
		// writeAppendConfigredis(path, tag, env, c)

		writecontent := confgen.ConvertMap4Json(convertops.ConvertYamlToMap(contentstr, c), func(ciphertext string, Yek, Ecnon []byte, c context.Context) string { return ciphertext }, c)

		return writecontent, err
	}

	contents, tags, err := confgen.GenerateAppConfigContentList(appname, envdclist, f0, c)

	return contents, tags, err
}

func writeAppConfig(appname string, envlist []string, c context.Context) (map[string]interface{}, map[string]map[string]string, error) {
	var contents map[string]interface{}

	var err error
	var contentstr string

	path := constset.Confpath + appname + string(os.PathSeparator)
	// ioutil.
	var f0 = func(appname string, content map[string]interface{}, tag map[string]map[string]string, env string, c context.Context) (map[string]interface{}, error) {
		// if r, ok := content["af-arch"].(map[string]interface{})["resource"].(map[string]interface{}); ok {
		// 	var resourceids []string
		// 	for k := range r {
		// 		if strings.Compare(k, "env") != 0 {
		// 			resourceids = append(resourceids, k)
		// 		}
		// 	}
		// 	if len(resourceids) != 0 {
		// 		// env := r["env"].(string)

		// 		// dbops.Update_appRes(appname, env, resourceids)
		// 	}
		// }

		contentstr, err = fileops.WriteToAppPath(path, GetAppConfName(env), appname, content, env, c)

		// logger.Println(*constset.Filepaths)
		//LogConfig/config

		writeAppendConfigredis(path, tag, env, c)

		writecontent := confgen.ConvertMap4Json(convertops.ConvertYamlToMap(contentstr, c), func(ciphertext string, Yek, Ecnon []byte, c context.Context) string { return ciphertext }, c)

		return writecontent, err
	}

	contents, tags, err := confgen.GenerateAppConfigContentList(appname, envlist, f0, c)

	return contents, tags, err
}

func GetAppConfName(env string) string {
	var filenamestr = "application-" + env + ".yml"

	return filenamestr
}

// func writeAppConfigFileWith(appname, path string) (map[string]interface{}, error) {
// 	var contents map[string]interface{}

// 	var err error
// 	var contentstr string

// 	var f0 = func(appname string, content map[string]interface{}, env string) (map[string]interface{}, error) {
// 		if r, ok := content["af-arch"].(map[string]interface{})["resource"].(map[string]interface{}); ok {
// 			var resourceids []string
// 			for k := range r {
// 				if strings.Compare(k, "env") != 0 {
// 					resourceids = append(resourceids, k)
// 				}
// 			}
// 			if len(resourceids) != 0 {
// 				env := r["env"].(string)

// 				dbops.Update_appRes(appname, env, resourceids)
// 			}
// 		}

// 		contentstr, err = fileops.WriteToPath(path, content, env)

// 		// logger.Println(*constset.Filepaths)
// 		//LogConfig/config
// 		writeAppendConfig(path, env)

// 		writecontent := confgen.ConvertMap4Json(convertops.ConvertYamlToMap(contentstr), func(ciphertext string, Yek, Ecnon []byte) string { return ciphertext })

// 		return writecontent, err
// 	}

// 	contents, err = confgen.GenerateAppConfigContentList(appname, constset.EnvSet, f0)

// 	return contents, err
// }

func GenerateallconfigRemote(srcPaths string, c context.Context) error {

	logger := logagent.InstPlatform(c)
	logger.Printf("the srcPaths are %s", srcPaths)

	pathArray := strings.Split(srcPaths, ",")

	_, err := writeConfigFileWith(pathArray, c)
	return err
}

func writeConfigFile(pathArray []string, c context.Context) (map[string]interface{}, error) {
	var contents map[string]interface{}

	var err error
	var filecontent, path, contentstr string

	logger := logagent.InstPlatform(c)
	var f0 = func(appname string, content map[string]interface{}, env string) (map[string]interface{}, error) {
		if r, ok := content["af-arch"].(map[string]interface{})["resource"].(map[string]interface{}); ok {
			var resourceids []string
			for k := range r {
				if strings.Compare(k, "env") != 0 {
					resourceids = append(resourceids, k)
				}
			}
			if len(resourceids) != 0 {
				env := r["env"].(string)

				dbops.Update_appRes(appname, env, resourceids, c)
			}
		}

		contentstr, err = fileops.WriteToPath(path, GetAppConfName(env), content, env, c)

		// logger.Println(*constset.Filepaths)
		//LogConfig/config
		writeAppendConfig(path, env, c)

		writecontent := confgen.ConvertMap4Json(convertops.ConvertYamlToMap(contentstr, c), func(ciphertext string, Yek, Ecnon []byte, c context.Context) string { return ciphertext }, c)

		return writecontent, err
	}

	for _, path = range pathArray {

		filecontent, err = fileops.Read(path)
		if err != nil {
			logger.Panic(err.Error())
			break
		} else {
			contents, err = confgen.GenerateConfigContentList(filecontent, constset.EnvSet, f0, c)
			// for env, content := range contents {
			// 	whiteToPath(path, content["content"], env)
			// }
		}
	}
	return contents, err
}

func writeConfigFileWith(pathArray []string, c context.Context) (map[string]interface{}, error) {
	var contents map[string]interface{}

	var err error
	var filecontent, path, contentstr string

	logger := logagent.InstPlatform(c)
	var f0 = func(appname string, content map[string]interface{}, env string) (map[string]interface{}, error) {
		if r, ok := content["af-arch"].(map[string]interface{})["resource"].(map[string]interface{}); ok {
			var resourceids []string
			for k := range r {
				if strings.Compare(k, "env") != 0 {
					resourceids = append(resourceids, k)
				}
			}
			if len(resourceids) != 0 {
				env := r["env"].(string)

				dbops.Update_appRes(appname, env, resourceids, c)
			}
		}

		contentstr, err = fileops.WriteToPath(path, GetAppConfName(env), content, env, c)

		// logger.Println(*constset.Filepaths)
		//LogConfig/config
		writeAppendConfig(path, env, c)

		writecontent := confgen.ConvertMap4Json(convertops.ConvertYamlToMap(contentstr, c), func(ciphertext string, Yek, Ecnon []byte, c context.Context) string { return ciphertext }, c)

		return writecontent, err
	}

	for _, path = range pathArray {

		filecontent, err = fileops.Read(path)
		if err != nil {
			logger.Panic(err.Error())
			break
		} else {
			contents, err = confgen.GenerateConfigContentListremote(filecontent, constset.EnvSet, f0, c)
			// for env, content := range contents {
			// 	whiteToPath(path, content["content"], env)
			// }
		}
	}
	return contents, err
}

func writeAppendConfigWith(appname, path, env string, c context.Context) {
	rediscli := redisops.Pool().Get()

	logger := logagent.InstPlatform(c)
	defer rediscli.Close()

	lastIndex := strings.LastIndex(path, "/")
	basepath := path[0:lastIndex]
	var (
		cursor int64
		items  []string
	)

	// results := make([][]string, 0)
	confmap := map[string]string{}
	for {
		values, err := redis.Values(rediscli.Do("HSCAN", "confsolver-append", cursor, "MATCH", "*", "COUNT", 1))

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
	rediscli.Do("EXPIRE", "confsolver-append", 60*10)
	for filenamekey, contentvalue := range confmap {
		fileops.Writeover(basepath+string(os.PathSeparator)+filenamekey, contentvalue, c)
	}

	// var filename string
	// for _, v := range constset.AppendSet {
	// 	// if v.Env == "" {
	// 	// 	envpara = env
	// 	// 	ss := strings.Split(v.Id, ".")

	// 	// 	if len(ss) > 1 {
	// 	// 		filename = ss[0] + "-" + env + "." + ss[1]
	// 	// 	} else {
	// 	// 		filename = v.Id
	// 	// 	}
	// 	// } else {
	// 	// 	envpara = v.Env
	// 	// 	filename = v.Id
	// 	// }
	// 	// filename = v.Id

	// 	// result := confgen.GetOnlineConfig(v.Type, v.Id, envpara)

	// 	if v.Withenv {
	// 		ss := strings.Split(v.Id, ".")
	// 		filename = ss[0] + "-" +
	// 			+"." + ss[1]
	// 	} else {
	// 		filename = v.Id
	// 	}
	// 	logger.Println(basepath + string(os.PathSeparator) + filename)
	// 	// _, err := rediscli.Do("HSET", "confsolver-append", filename, result)
	// 	confstr, err := redis.String(rediscli.Do("HGET", "confsolver-append", filename))
	// 	rediscli.Do("EXPIRE", "confsolver-append", 60*10)
	// 	// bytes, err := ioutil.ReadFile(constset.Confpath + appname + string(os.PathSeparator) + "application-" + env + ".yml")
	// 	if err != nil {
	// 		logger.Panic(err)
	// 	}

	// 	// bytes, err := ioutil.ReadFile(constset.Confpath + appname + string(os.PathSeparator) + filename)
	// 	// if err != nil {
	// 	// 	logger.Panic(err)
	// 	// }
	// 	// result := confgen.GetOnlineConfig(v.Type, v.Id, env)
	// 	fileops.Writeover(basepath+string(os.PathSeparator)+filename, confstr, c)
	// }
}

func writeAppendConfig(path string, envdc string, c context.Context) {
	// var filename string
	logger := logagent.InstPlatform(c)
	lastIndex := strings.LastIndex(path, "/")
	basepath := path[0:lastIndex]
	ConfNameHelper(envdc, constset.AppendSet, func(vtype, vid, confenv, filename string) {
		result := confgen.GetOnlineConfig(vtype, vid, envdc, c)
		logger.Println(basepath + string(os.PathSeparator) + filename)
		fileops.Writeover(basepath+string(os.PathSeparator)+filename, result, c)
	})
	// for _, v := range constset.AppendSet {
	// 	// if v.Env == "" {
	// 	// 	envpara = env
	// 	// 	ss := strings.Split(v.Id, ".")

	// 	// 	if len(ss) > 1 {
	// 	// 		filename = ss[0] + "-" + env + "." + ss[1]
	// 	// 	} else {
	// 	// 		filename = v.Id
	// 	// 	}
	// 	// } else {
	// 	// 	envpara = v.Env
	// 	// 	filename = v.Id
	// 	// }
	// 	// filename = v.Id

	// 	// result := confgen.GetOnlineConfig(v.Type, v.Id, envpara)

	// 	if v.Withenv {
	// 		ss := strings.Split(v.Id, ".")
	// 		filename = ss[0] + "-" + env + "." + ss[1]
	// 	} else {
	// 		filename = v.Id
	// 	}
	// 	result := confgen.GetOnlineConfig(v.Type, v.Id, env, c)
	// 	logger.Println(basepath + string(os.PathSeparator) + filename)
	// 	fileops.Writeover(basepath+string(os.PathSeparator)+filename, result, c)
	// }
}

func writeAppendConfigRepo(path, appname, team string, tags map[string]map[string]string, envdc, version, region string, c context.Context) {

	logger := logagent.InstPlatform(c)
	rediscli := redisops.Pool().Get()

	defer rediscli.Close()

	// var filename string
	lastIndex := strings.LastIndex(path, "/")
	basepath := path[0:lastIndex]

	ConfNameHelper(envdc, constset.AppendSet, func(vtype, vid, confenv, filename string) {
		result := confgen.GetOnlineConfig(vtype, vid, confenv, c)
		tagInsert(vtype, vid, tags, envdc)

		logger.Println(basepath + string(os.PathSeparator) + filename)
		//FIXME: write append to nexus repo
		fileops.WriteToNexus(result, filename, team, appname, envdc, version, region, c)
		_, err := rediscli.Do("HSET", "confsolver-append", filename, result)
		rediscli.Do("EXPIRE", "confsolver-append", 60*10)
		if err != nil {
			logger.Panic(err)
		}
	})
	// for _, v := range constset.AppendSet {
	// 	// if v.Env == "" {
	// 	// 	envpara = env
	// 	// 	ss := strings.Split(v.Id, ".")

	// 	// 	if len(ss) > 1 {
	// 	// 		filename = ss[0] + "-" + env + "." + ss[1]
	// 	// 	} else {
	// 	// 		filename = v.Id
	// 	// 	}
	// 	// } else {
	// 	// 	envpara = v.Env
	// 	// 	filename = v.Id
	// 	// }
	// 	// filename = v.Id

	// 	// result := confgen.GetOnlineConfig(v.Type, v.Id, envpara)
	// 	var result string
	// 	if v.Withenv {
	// 		ss := strings.Split(v.Id, ".")
	// 		filename = ss[0] + "-" + env + "." + ss[1]
	// 		result = confgen.GetOnlineConfig(v.Type, v.Id, env, c)
	// 	} else {
	// 		filename = v.Id
	// 		result = confgen.GetOnlineConfig(v.Type, v.Id, "prod", c)
	// 	}
	// 	// result := confgen.GetOnlineConfig(v.Type, v.Id, env, c)

	// 	// + "_" + env
	// 	tagInsert(v.Type, v.Id, tags, env)
	// 	// tags[tagkey] = v.Id
	// 	// if val, ok := tags[tagkey]; ok {
	// 	// 	tags[tagkey] =  v.Id
	// 	// } else {
	// 	// 	tags[tagkey] = env + v.Id
	// 	// }

	// 	logger.Println(basepath + string(os.PathSeparator) + filename)
	// 	//FIXME: write append to nexus repo
	// 	fileops.WriteToNexus(result, filename, team, appname)
	// 	_, err := rediscli.Do("HSET", "confsolver-append", filename, result)
	// 	rediscli.Do("EXPIRE", "confsolver-append", 60*10)
	// 	if err != nil {
	// 		logger.Panic(err)
	// 	}
	// 	// fileops.Writeover(basepath+string(os.PathSeparator)+filename, result)
	// }
}

func tagInsert(vtype, vid string, tags map[string]map[string]string, env string) {
	tagkey := vtype + "_" + vid
	if _, ok := tags[env]; ok {
		tags[env][tagkey] = vid
	} else {
		tags[env] = map[string]string{tagkey: vid}
	}
}

func writeAppendConfigredis(path string, tags map[string]map[string]string, envdc string, c context.Context) {

	logger := logagent.InstPlatform(c)
	rediscli := redisops.Pool().Get()

	defer rediscli.Close()

	// var filename string
	lastIndex := strings.LastIndex(path, "/")
	basepath := path[0:lastIndex]

	ConfNameHelper(envdc, constset.AppendSet, func(vtype, vid, confenv, filename string) {
		result := confgen.GetOnlineConfig(vtype, vid, confenv, c)
		tagInsert(vtype, vid, tags, envdc)

		logger.Println(basepath + string(os.PathSeparator) + filename)
		_, err := rediscli.Do("HSET", "confsolver-append", filename, result)
		rediscli.Do("EXPIRE", "confsolver-append", 60*10)
		if err != nil {
			logger.Panic(err)
		}
	})
	// for _, v := range constset.AppendSet {
	// 	// if v.Env == "" {
	// 	// 	envpara = env
	// 	// 	ss := strings.Split(v.Id, ".")

	// 	// 	if len(ss) > 1 {
	// 	// 		filename = ss[0] + "-" + env + "." + ss[1]
	// 	// 	} else {
	// 	// 		filename = v.Id
	// 	// 	}
	// 	// } else {
	// 	// 	envpara = v.Env
	// 	// 	filename = v.Id
	// 	// }
	// 	// filename = v.Id

	// 	// result := confgen.GetOnlineConfig(v.Type, v.Id, envpara)
	// 	var result string
	// 	if v.Withenv {
	// 		ss := strings.Split(v.Id, ".")
	// 		filename = ss[0] + "-" + env + "." + ss[1]
	// 		result = confgen.GetOnlineConfig(v.Type, v.Id, env, c)
	// 	} else {
	// 		filename = v.Id
	// 		result = confgen.GetOnlineConfig(v.Type, v.Id, "prod", c)
	// 	}
	// 	// result := confgen.GetOnlineConfig(v.Type, v.Id, env, c)

	// 	tagkey := v.Type + "_" + v.Id // + "_" + env
	// 	if _, ok := tags[env]; ok {
	// 		tags[env][tagkey] = v.Id
	// 	} else {
	// 		tags[env] = map[string]string{tagkey: v.Id}
	// 	}
	// 	// tags[tagkey] = v.Id
	// 	// if val, ok := tags[tagkey]; ok {
	// 	// 	tags[tagkey] =  v.Id
	// 	// } else {
	// 	// 	tags[tagkey] = env + v.Id
	// 	// }

	// 	logger.Println(basepath + string(os.PathSeparator) + filename)
	// 	_, err := rediscli.Do("HSET", "confsolver-append", filename, result)
	// 	rediscli.Do("EXPIRE", "confsolver-append", 60*10)
	// 	if err != nil {
	// 		logger.Panic(err)
	// 	}
	// 	// fileops.Writeover(basepath+string(os.PathSeparator)+filename, result)
	// }
}

func ConfNameHelper(envdc string, appendset []constset.AppendStruct, processor func(string, string, string, string)) {
	var filename string
	for _, v := range appendset {
		// if v.Env == "" {
		// 	envpara = env
		// 	ss := strings.Split(v.Id, ".")

		// 	if len(ss) > 1 {
		// 		filename = ss[0] + "-" + env + "." + ss[1]
		// 	} else {
		// 		filename = v.Id
		// 	}
		// } else {
		// 	envpara = v.Env
		// 	filename = v.Id
		// }
		// filename = v.Id

		// result := confgen.GetOnlineConfig(v.Type, v.Id, envpara)

		var confenv string
		if v.Withenv {
			ss := strings.Split(v.Id, ".")
			filename = ss[0] + "-" + envdc + "." + ss[1]
			confenv = envdc
			// result = confgen.GetOnlineConfig(v.Type, v.Id, env, c)
		} else {
			filename = v.Id
			confenv = "prod"
			// result = confgen.GetOnlineConfig(v.Type, v.Id, "prod", c)
		}
		processor(v.Type, v.Id, confenv, filename)
	}

}

func gen4all(c *gin.Context) {
	envs := c.Query("envs")

	logger := logagent.InstPlatform(c)
	logger.Printf("the envs are %s", envs)

	envarray := strings.Split(envs, ",")
	logger.Println(envarray)

	form, err := c.MultipartForm()
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("get err %s", err.Error()))
	}
	// 获取所有图片
	// files := form.File["files"]
	// 遍历所有图片
	for _, fileheader := range form.File["files"] {
		// 逐个存
		file, _ := fileheader.Open()
		content, err := fileops.ReadFrom(file, c)
		if err != nil {
			logger.Panic(err)
		}
		logger.Println(string(content))

	}
}

func Get4all(c *gin.Context) {

	logger := logagent.InstPlatform(c)
	file, head, _ := c.Request.FormFile("file")
	logger.Println(head.Filename)

	content, err := fileops.ReadFrom(file, c)
	// confgen.ConvertWithCypher = cypher.Decryptbyhex2str

	var f0 = func(appname string, content map[string]interface{}, env string) (map[string]interface{}, error) {
		return content, nil
	}

	var result map[string]interface{}
	if err == nil {
		result, err = confgen.GenerateConfigContentList(content, constset.EnvSet, f0, c)
	}

	common_resp(err, "", result, c)
	// var httpstatus = http.StatusOK
	// var errstr = ""
	// if err != nil {
	// 	httpstatus = http.StatusInternalServerError
	// 	errstr = err.Error()
	// }

	// c.JSON(httpstatus, gin.H{
	// 	"data":  result, //contents,
	// 	"error": errstr,
	// })

	// c.JSON(http.StatusOK, result)
}

func common_resp(err error, message string, data interface{}, c *gin.Context) {
	var httpstatus = http.StatusOK
	var errstr = ""
	if err != nil {
		httpstatus = http.StatusInternalServerError
		errstr = err.Error()
	}

	c.JSON(httpstatus, gin.H{
		"data":    data, //contents,
		"error":   errstr,
		"message": message,
	})
}
func Get4env(c *gin.Context) {
	target_env := c.Param("env")
	file, head, _ := c.Request.FormFile("file")

	logger := logagent.InstPlatform(c)
	logger.Println(head.Filename)

	content, err := fileops.ReadFrom(file, c)
	// confgen.ConvertWithCypher = cypher.Decryptbyhex2str

	if err == nil {
		rescontent := confgen.GenerateConfigString(content, target_env, c)
		bytes := []byte(rescontent)

		fileContentDisposition := "attachment;filename=\"" + "config.yml" + "\""
		c.Header("Content-Type", "application/yml") // 这里是压缩文件类型 .zip
		c.Header("Content-Disposition", fileContentDisposition)
		c.Data(http.StatusOK, "application/yml", bytes)
	} else {
		common_resp(err, "", nil, c)
	}

	// confgen.ConvertWithCypher = cypher.Decryptbyhex2str

}

func fileTokenOnline(c *gin.Context) {
	json := make(map[string]string) //注意该结构接受的内容
	c.ShouldBindJSON(&json)

	logger := logagent.InstPlatform(c)
	logger.Println(json)

	resStr := cypher.Md5str(json["data"])
	// c.JSON(http.StatusOK, gin.H{
	// 	"data": resStr,
	// })

	common_resp(nil, "", resStr, c)
}

func Encrypt2Hexonline(c *gin.Context) {
	json := make(map[string]string) //注意该结构接受的内容
	c.ShouldBindJSON(&json)
	logger := logagent.InstPlatform(c)
	logger.Println(json)

	resStr := cypher.EncryptStr2hex(json["data"], constset.Yek, constset.Ecnon, c)
	// c.JSON(http.StatusOK, gin.H{
	// 	"data": resStr,
	// })

	common_resp(nil, "", resStr, c)
}

func DecryptHexonline(c *gin.Context) {
	json := make(map[string]string) //注意该结构接受的内容
	c.ShouldBindJSON(&json)
	logger := logagent.InstPlatform(c)
	logger.Println(json)

	resbytes := cypher.Decryptbyhex(json["data"], constset.Yek, constset.Ecnon, c)
	// var errstr string
	// var httpstatus = http.StatusOK
	var resString = ""
	// if err != nil {
	// 	httpstatus = http.StatusInternalServerError
	// 	errstr = err.Error()
	// 	c.JSON(httpstatus, gin.H{
	// 		"error": errstr,
	// 	})
	// } else {
	// 	resString = string(resbytes)
	// 	c.JSON(httpstatus, gin.H{
	// 		"data": resString,
	// 	})
	// }

	if resbytes != nil {
		resString = string(resbytes)
	}
	common_resp(nil, "", resString, c)

}

func EncryptConfig2Hexonline(c *gin.Context) {
	// file, head, err := c.Request.FormFile("file")
	// logger.Println(err)
	// logger.Println(file)
	// logger.Println(head.Filename)

	var resstr = ""
	var err error
	var file multipart.File
	var head *multipart.FileHeader
	logger := logagent.InstPlatform(c)
	if file, head, err = c.Request.FormFile("file"); err == nil {
		logger.Println(err)
		logger.Println(file)
		logger.Println(head.Filename)
		// var bytes []byte
		// bytes = nil
		var configmap map[string]interface{}
		if configmap, err = confgen.GetPostFileConfigWithEncrypt(file, c); configmap != nil {
			resstr = convertops.ConvertStrMapToYaml(&configmap, c)
		}
	}

	if err == nil {
		// bytes := []byte(resstr)

		fileContentDisposition := "attachment;filename=\"" + "config.yml" + "\""
		c.Header("Content-Type", "application/yml") // 这里是压缩文件类型 .zip
		c.Header("Content-Disposition", fileContentDisposition)
		c.Data(http.StatusOK, "application/yml", []byte(resstr))
	} else {
		common_resp(err, "", nil, c)
	}

	// if err != nil {
	// 	configmap := confgen.GetPostFileConfig(file, confgen.ConvertWihtEncypher)
	// 	resstr := convertops.ConvertStrMapToYaml(&configmap)

	// 	bytes := []byte(resstr)

	// 	fileContentDisposition := "attachment;filename=\"" + "config.yml" + "\""
	// 	c.Header("Content-Type", "application/yml") // 这里是压缩文件类型 .zip
	// 	c.Header("Content-Disposition", fileContentDisposition)
	// 	c.Data(http.StatusOK, "application/yml", bytes)
	// } else {
	// 	common_resp(err, "", nil, c)
	// }

	// confgen.ConvertWithCypher = confgen.ConvertWihtEncypher
	// configmap := confgen.GetPostFileConfig(file, confgen.ConvertWihtEncypher)
	// resstr := convertops.ConvertStrMapToYaml(&configmap)

	// bytes := []byte(resstr)

	// fileContentDisposition := "attachment;filename=\"" + "config.yml" + "\""
	// c.Header("Content-Type", "application/yml") // 这里是压缩文件类型 .zip
	// c.Header("Content-Disposition", fileContentDisposition)
	// c.Data(http.StatusOK, "application/yml", bytes)

	// b, err := ioutil.ReadAll(file) //.ReadFile(filePath)
	// var resstr, errstr string
	// var httpstatus = 200
	// if err != nil {
	// 	httpstatus = 500
	// 	errstr = err.Error()
	// } else {
	// 	resstr = encrypt2hex(b, Yek, Ecnon)
	// 	// encryptstr := base64.StdEncoding.EncodeToString(encryptbytes)
	// 	fmt.Println(resstr)
	// }

	// c.JSON(httpstatus, gin.H{
	// 	"data":  resstr,
	// 	"error": errstr,
	// })
}

func DecryptHexConfigOnline(c *gin.Context) {
	var resstr = ""
	var err error
	var file multipart.File
	var head *multipart.FileHeader
	logger := logagent.InstPlatform(c)
	if file, head, err = c.Request.FormFile("file"); err == nil {
		logger.Println(err)
		logger.Println(file)
		logger.Println(head.Filename)
		// var bytes []byte
		// bytes = nil
		var configmap map[string]interface{}
		if configmap, err = confgen.GetPostFileConfigWithDecrypt(file, c); configmap != nil {
			resstr = convertops.ConvertStrMapToYaml(&configmap, c)
		}
	}

	if err == nil {
		// bytes := []byte(resstr)

		fileContentDisposition := "attachment;filename=\"" + "config.yml" + "\""
		c.Header("Content-Type", "application/yml") // 这里是压缩文件类型 .zip
		c.Header("Content-Disposition", fileContentDisposition)
		c.Data(http.StatusOK, "application/yml", []byte(resstr))
	} else {
		common_resp(err, "", nil, c)
	}
	// confgen.ConvertWithCypher = cypher.Decryptbyhex2str

}
