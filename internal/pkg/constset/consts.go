package constset

import (
	"context"
	"flag"
	"os"
	"strconv"
	"strings"

	"github.com/max-gui/consulagent/pkg/consulsets"
	"github.com/max-gui/logagent/pkg/logagent"
	"github.com/max-gui/logagent/pkg/logsets"
	"github.com/max-gui/redisagent/pkg/redisops"
	"gopkg.in/yaml.v2"
)

// type Resp struct {
// 	message string
// 	err     error
// 	data    map[string]interface{}
// }

// func GetResp(message string, errstr string, data map[string]interface{}) map[string]interface{} {
// 	var resp = make(map[string]interface{})
// 	resp[MESSAGE] = message
// 	resp[ERRORSTRING] = errstr
// 	resp[DATA] = data
// 	return resp
// }

// const MESSAGE string = "message"
// const ERRORSTRING string = "error"
// const DATA string = "data"
// const HTTPOK bool = true

const SYSTEM_ERROR string = "500"

const CONSUL_ERROR string = "400"

const IO_ERROR string = "300"
const PthSep = string(os.PathSeparator)

var Yek = []byte{74, 103, 115, 173, 168, 227, 72, 68, 25, 245, 63, 49, 136, 236, 197, 236}
var Ecnon = []byte{9, 65, 48, 149, 170, 165, 84, 222, 74, 84, 4, 106}

// var argsetmap = make(map[string]string)

// const port = "port"
// const consul_host = "consulhost"
// const acltoken = "acltoken"
// const envset = "envset"
// const apppath = "apppath"
// const iocurl = "iocurl"
// const templurl = "templurl"
// const archurl = "archurl"

// const repo = "repo"

// var Apppath, IocUrl, Repopathname string
// var EnvSet []string

// var ConfigFolder string

var (
	envset, Filepaths, ConfArchPrefix, ConfResPrefix, ConfWatchPrefix *string
	appendSet                                                         *string
	Confpath                                                          string
	Oniac, Servermode                                                 *bool
	EnvSet                                                            []string
	AppendSet                                                         []struct {
		Type    string
		Id      string
		Withenv bool
	}
)

func StartupInit(bytes []byte, c context.Context) {
	logger := logagent.InstArch(c)
	confmap := map[string]interface{}{}
	yaml.Unmarshal(bytes, confmap)
	*consulsets.Acltoken = confmap["af-arch"].(map[interface{}]interface{})["resource"].(map[interface{}]interface{})["private"].(map[interface{}]interface{})["acl-token"].(string)
	consulsets.StartupInit(*consulsets.Acltoken)
	Confpath = *logsets.Apppath + "confs" + string(os.PathSeparator)

	// config := consulhelp.Getconfaml(*ConfResPrefix, "redis", "redis-cluster-predixy", *logsets.Appenv, c)
	redisopsUrl := confmap["af-arch"].(map[interface{}]interface{})["resource"].(map[interface{}]interface{})["redis-cluster-predixy"].(map[interface{}]interface{})["url"].(string)
	// redisops.Url = confmap["url"].(string)
	redisopsPwd := confmap["af-arch"].(map[interface{}]interface{})["resource"].(map[interface{}]interface{})["redis-cluster-predixy"].(map[interface{}]interface{})["password"].(string)
	// redisops.Pwd = confmap["password"].(string)
	redisops.StartupInit(redisopsUrl, redisopsPwd)

	EnvSet = strings.Split(*envset, ",")
	var appenditem = struct {
		Type    string
		Id      string
		Withenv bool
	}{}
	var appendv []string
	// var envpara = ""
	for _, v := range strings.Split(*appendSet, ",") {
		// envpara = ""
		appendv = strings.Split(v, ":")
		if len(appendv) <= 1 {
			continue
		}
		withenv, err := strconv.ParseBool(appendv[2])
		if err != nil {
			logger.Panic(err)
		}
		// withenv := false
		// if appendv[2] == "withenv" {
		// 	withenv = true
		// }else if appendv[2] == "noenv"{
		// 	withenv = false
		// }else

		appenditem = struct {
			Type    string
			Id      string
			Withenv bool
		}{Type: appendv[0], Id: appendv[1], Withenv: withenv}

		AppendSet = append(AppendSet, appenditem)
	}

	// var argset []string

	// for _, arg := range args {
	// 	argset = strings.Split(arg, "=")
	// 	if len(argset) > 1 {
	// 		argsetmap[argset[0]] = argset[1]

	// 	}
	// }

	// EnvSet = strings.Split(argsetmap[envset], ",")
	// // Port = argsetmap[port]
	// Consul_host = argsetmap[consul_host]
	// Acltoken = argsetmap[acltoken]
	// Apppath = argsetmap[apppath]
	// IocUrl = argsetmap[iocurl]
	// Repopathname = argsetmap[repo]

	// // log.Println(port)
	// // log.Println(Port)
	// log.Println(consul_host)
	// log.Println(Consul_host)
	// log.Println(acltoken)
	// log.Println(Acltoken)
	// log.Println(Apppath)
	// log.Println(IocUrl)
	// log.Println(Repopathname)
	// log.Println(envset)
	// log.Println(argsetmap[envset])
	// log.Println("EnvSet")
	// log.Println(EnvSet)
}

func init() {
	envset = flag.String("envset", "test,uat,prod", "envset spilt by ','")
	Oniac = flag.Bool("oniac", false, "iac or not")
	// Appenv = flag.String("appenv", "prod", "the exec env")
	ConfArchPrefix = flag.String("confArchPrefix", "ops/iac/arch/", "arch prefix for consul")
	ConfResPrefix = flag.String("ConfResPrefix", "ops/resource/", "resource prefix for consul")
	ConfWatchPrefix = flag.String("ConfWatchPrefix", "ops/", "watch prefix for consul")
	// Apppath = flag.String("apppath", "/Users/jimmy/Downloads/consolver/", "app root path")
	// Acltoken = flag.String("acltoken", "", "consul acltoken")
	Servermode = flag.Bool("servermode", true, "true: run as httpserver;false: run as commond line")
	Filepaths = flag.String("filepaths", "", "(in nonserver mode)filepath split by ','")
	appendSet = flag.String("appendset", "bootload:bootstrap.yml:true,LogConfig:logback-spring.xml:false", "append confg, format is type:id:withenv, split by ',' e.g: 'type:id:true,type:id:false'")

	// Consul_host = flag.String("consulhost", "http://consul-prod.kube.com", "consul url")
	// Port = flag.String("port", "8181", "this app's port")
}
