package main

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

func main() {
	// 设置配置文件类型
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.SetConfigFile(".env")

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("配置文件读取失败: %v", err)
	}

	app_debug := viper.GetString("APP_DEBUG")

	fmt.Printf("dd %s\n ", app_debug)

	//r := gin.Default()
	//r.GET("/", func(c *gin.Context) {
	//	c.String(http.StatusOK, "hello world")
	//})
	//
	//r.Run(":8080")
}
