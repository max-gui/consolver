package dbops

import (
	"context"
	"strings"
	"sync"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gohouse/gorose/v2"
	"github.com/max-gui/consolver/internal/confgen"
	"github.com/max-gui/logagent/pkg/logagent"
	"github.com/max-gui/logagent/pkg/logsets"
)

type App_Resource struct {
	Id          int64  `gorose:"id"`
	Create_time string `gorose:"create_time"`
	Update_time string `gorose:"update_time"`
	Removed     bool   `gorose:"removed"`
	Appname     string `gorose:"appname"`
	Resourceid  string `gorose:"resourceid"`
	Env         string `gorose:"env"`
}

// 设置表名, 如果没有设置, 默认使用struct的名字
func (u *App_Resource) TableName() string {
	return "app_resource"
}

var engin *gorose.Engin
var once sync.Once

func dbInit(c context.Context) *gorose.Engin {
	// 全局初始化数据库,并复用
	// 这里的engin需要全局保存,可以用全局变量,也可以用单例
	// 配置&gorose.Config{}是单一数据库配置
	// 如果配置读写分离集群,则使用&gorose.ConfigCluster{}
	var err error
	log := logagent.InstArch(c)
	once.Do(func() {

		c := logagent.GetRootContextWithTrace()
		config := confgen.Getconfig("mysql", "AFDISCZ", *logsets.Appenv, c) // "test")
		un := config["username"].(string)
		pw := config["password"].(string)

		urlstr := config["url"].(string)
		urlvar := strings.Split(strings.TrimPrefix(urlstr, "jdbc:mysql://"), "/")
		dsnstr := un + ":" + pw + "@tcp(" + urlvar[0] + ")/" + strings.Split(urlvar[1], "?")[0]
		log.Println(dsnstr)

		engin, err = gorose.Open(&gorose.Config{Driver: "mysql", Dsn: dsnstr})
		if err != nil {
			log.Panic(err)
		}
	})
	return engin
}

func DB(c context.Context) gorose.IOrm {
	if engin == nil {
		dbInit(c)
	}
	return engin.NewOrm()
}

func Insert_appRes(appname, resourceid, env string, c context.Context) {
	orm := DB(c)
	log := logagent.InstArch(c)

	// orm.GetTable()
	var ar App_Resource
	new_ar := App_Resource{Appname: appname, Resourceid: resourceid, Env: env}
	_, err := orm.Table(&ar).Data(new_ar).Insert()
	if err != nil {
		log.Panic(err)
	}
}

func Update_appRes(appname, env string, resourceid []string, c context.Context) {
	orm := DB(c)
	log := logagent.InstArch(c)

	err := orm.Transaction(
		// 第一个业务
		func(db gorose.IOrm) error {
			var ar App_Resource
			_, err := orm.Table(&ar).Where("appname", appname).Where("env", env).Delete()
			// err = errors.New("test delete")
			if err != nil {
				return err
			}
			return nil
		},
		// 第二个业务
		func(db gorose.IOrm) error {
			var new_ar []App_Resource
			for _, e := range resourceid {
				new_ar = append(new_ar, App_Resource{Appname: appname, Resourceid: e, Env: env})

				// err = errors.New("test")
			}
			// var multi_data = []map[string]interface{}{ {"age":17, "job":"it3"},{"age":17, "job":"it4"} }
			// // insert into user (age, job) values (17, 'it3') (17, 'it4')
			// db.Table("user").Data(multi_data).Insert()
			_, err := orm.Table(&new_ar).Data(&new_ar).Insert()

			if err != nil {
				return err
			}
			return nil
		},
	)
	// defer func() {
	// 	if err := recover(); err != nil {
	// 		brr := orm.Rollback()
	// 		if brr != nil {
	// 			log.Print(brr)
	// 		}
	// 		fmt.Println(err)
	// 	}
	// }()

	// err = orm.Begin()
	// if err != nil {
	// 	log.Panic(err)
	// }
	// // orm.GetTable()
	// var ar App_Resource
	// var new_ar []App_Resource
	// _, err = orm.Table(&ar).Where("appname", appname).Delete()
	// // err = errors.New("test delete")
	// if err != nil {
	// 	orm.Rollback()
	// 	log.Println(err)
	// 	return
	// }
	// // log.Panic("test")
	// for _, e := range resourceid {
	// 	new_ar = append(new_ar, App_Resource{Appname: appname, Resourceid: e, Env: env})

	// 	// err = errors.New("test")
	// }
	// // var multi_data = []map[string]interface{}{ {"age":17, "job":"it3"},{"age":17, "job":"it4"} }
	// // // insert into user (age, job) values (17, 'it3') (17, 'it4')
	// // db.Table("user").Data(multi_data).Insert()
	// _, err = orm.Table(&new_ar).Data(&new_ar).Insert()
	// // err = errors.New("test insert")
	// if err != nil {
	// 	orm.Rollback()
	// 	log.Println(err)
	// 	return
	// }
	// err = orm.Commit()
	// // err = errors.New("test commit")
	// if err != nil {
	// 	orm.Rollback()
	// 	log.Println(err)
	// 	return
	// }
	if err != nil {
		log.Panic(err)
	}
}

func Query_appRes(c context.Context) []App_Resource {
	orm := DB(c)
	log := logagent.InstArch(c)

	// orm.GetTable()
	var ar []App_Resource
	// new_ar := App_Resource{Appname: appname, Resourceid: resourceid}
	err := orm.Table(&ar).Select()
	if err != nil {
		log.Panic(err)
	}
	return ar
}
