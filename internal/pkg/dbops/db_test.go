package dbops

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/max-gui/consolver/internal/confgen"
	"github.com/max-gui/consolver/internal/pkg/constset"
	"github.com/stretchr/testify/assert"
)

// var Configlist []map[string]interface{}

// var args []string
var appnames = []string{"bpp", "cpp"}

func setup() {
	c := context.Background()
	constset.StartupInit(nil, c)
	confgen.Makeconfiglist(context.Background())

	orm := DB(c)
	var ar App_Resource
	// new_ar := App_Resource{Appname: appname, Resourceid: resourceid}
	res, err := orm.Table(&ar).Force().Delete()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("deleted: %d lines", res)

	for _, app := range appnames {
		Insert_appRes(app, "resid", "test", c)
	}
	// plaintext = "123"
	// cryptedHexText = "1bda1896724a4521cfb7f38646824197929cd1"
	// constset.StartupInit(os.Args)
	// testpath = Makeconfiglist()

	// orm.GetTable()

}

func teardown() {
	// orm := DB()

	// // orm.GetTable()
	// var ar App_Resource
	// // new_ar := App_Resource{Appname: appname, Resourceid: resourceid}
	// orm.Table(&ar).Force().Delete()
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

// func prepareTestConfigs() {
// 	Makeconfiglist(func(entitytype, entityid, env, configcontent string) {
// 		resp, err := consulhelp.Sendconfig2consul(entitytype, entityid, env, configcontent)
// 		if err != nil {
// 			fmt.Println(err.Error())
// 		}
// 		fmt.Println(resp)
// 	})
// }

func Test_Insert_appRes(t *testing.T) {
	c := context.Background()
	Insert_appRes("app", "resid", "test", c)
	// Insert_appRes("bpp", "resid")

	orm := DB(c)

	// orm.GetTable()
	var ar App_Resource
	// new_ar := App_Resource{Appname: appname, Resourceid: resourceid}
	err := orm.Table(&ar).Fields("appname,resourceid").Where("resourceid='resid'").Where("appname='app'").OrderBy("appname desc").Select()
	if err != nil {
		t.Log(err)
	}
	t.Log(orm.LastSql())
	assert.Equal(t, "app", ar.Appname)
	log.Printf("Test_Insert_appRes result is:\nport:%v", ar)
	orm.Table(&ar).Where("resourceid='resid'").Where("appname='app'").Delete()
}

func Test_Query_appRes(t *testing.T) {

	// orm := DB()

	// // orm.GetTable()
	// var ar App_Resource
	// // new_ar := App_Resource{Appname: appname, Resourceid: resourceid}
	// orm.Table(&ar).Fields("appname,resourceid").Where("appname", "=", "app", "resourceid", "=", "resid").OrderBy("appname desc").Select()
	// assert.Equal(t, "app", ar.Appname)
	var result []string
	c := context.Background()
	ar := Query_appRes(c)
	for _, e := range ar {
		result = append(result, e.Appname)
	}
	assert.Equal(t, appnames, result)
	log.Printf("Test_Query_appRes result is:\nport:%v", ar)
	// configstr := convertops.ConvertStrMapToYaml(&m)

	// assert.NoError(t, err, "read is ok")
	// assert.Equal(t, orgm, configm)
	// t.Logf("Test_Getconfig result is:\n%s", configm)

	// if !convertops.CompareTwoMapInterface(orgm, configm) {
	// 	t.Fatalf("Test_Getconfig failed! config should be:\n%s \nget:\n%s", orgm, configm)
	// }
	// t.Logf("Test_Getconfig is ok! result is:\n%s", configm)
}

func Test_Update_appRes(t *testing.T) {
	c := context.Background()
	Insert_appRes("app", "resid", "test", c)
	// Insert_appRes("bpp", "resid")
	resids := []string{"aa", "bb", "resid"}
	Update_appRes("bpp", "test", resids, c)
	orm := DB(c)

	// orm.GetTable()
	var ar []App_Resource
	// new_ar := App_Resource{Appname: appname, Resourceid: resourceid}
	err := orm.Table(&ar).Fields("appname,resourceid").Where("appname", "bpp").OrderBy("appname desc").Select()
	if err != nil {
		t.Log(err)
	}
	t.Log(orm.LastSql())
	assert.Equal(t, len(resids), len(ar))
	log.Printf("Test_Update_appRes result is:\nport:%v", ar)
	// orm.Table(&ar).Where("resourceid='resid'").Where("appname='app'").Delete()
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
