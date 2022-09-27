package main

import (
	"flag"

	"github.com/max-gui/consolver/internal/pkg/constset"
	"github.com/max-gui/consolver/router"
	"github.com/max-gui/consulagent/pkg/consulhelp"
	"github.com/max-gui/consulagent/pkg/consulsets"
	"github.com/max-gui/logagent/pkg/confload"
	"github.com/max-gui/logagent/pkg/logagent"
	"github.com/max-gui/logagent/pkg/logsets"
	ginprometheus "github.com/zsais/go-gin-prometheus"
)

func main() {

	// var Argsetmap = make(map[string]interface{})
	flag.Parse()
	c := logagent.GetRootContextWithTrace()
	bytes := confload.Load(c)
	constset.StartupInit(bytes, c)

	logger := logagent.InstArch(c)
	logger.Print(*consulsets.Acltoken)

	go consulhelp.StartWatch(*constset.ConfWatchPrefix, true, c)
	// if len(os.Args) > 2 {
	// 	port = os.Args[1]
	// }
	// if len(os.Args) > 2 {
	// 	consulhelp.Consulurl = os.Args[2]
	// }
	// if len(os.Args) > 3 {
	// 	consulhelp.AclToken = os.Args[3]
	// }

	// router.Envs
	//port := "4000"
	logger.Println(*constset.Servermode)
	if *constset.Servermode {
		r := router.SetupRouter()

		// defer func() {
		// 	if e := recover(); e != nil {
		// 		// err := e.(error)
		// 		log.Println("================main=================")
		// 		log.Println(e)
		// 	}
		// }()

		p := ginprometheus.NewPrometheus("gin")
		p.Use(r)
		r.Run(":" + *logsets.Port)
	} else {
		logger.Println(*constset.Filepaths)
		logger.Println(constset.EnvSet)
		logger.Println(constset.AppendSet)
		if *constset.Oniac {

			router.GenerateallconfigRemote(*constset.Filepaths, c)
		} else {

			router.Generateallconfig(*constset.Filepaths, c)
		}
	}
}
