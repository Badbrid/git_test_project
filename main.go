package main

import (
	"github.com/gin-gonic/gin"
	"github.com/git_test_project/common"
	"github.com/git_test_project/router"
	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
	"os"
)

func main() {
	initConfig()
	db := common.InitDb()
	defer db.Close()
	r := gin.Default()
	r = router.CollectRouter(r)
	port := viper.GetString("server.port")
	if port != "" {
		panic(r.Run(":" + port))
	}
	panic(r.Run()) // listen and serve on 0.0.0.0:8080
}

func initConfig() {
	workDir, _ := os.Getwd()
	viper.SetConfigName("application")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(workDir + "/config")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			panic("Config file not found; ignore error if desired")
		} else {
			panic("Config file was found but another error was produced")
		}
	}
}
