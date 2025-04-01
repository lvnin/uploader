package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"uploader/api"
	"uploader/core/logzap"
	"uploader/global"
	"uploader/utils"

	"github.com/fsnotify/fsnotify"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// 获取环境变量
func getEnvMode() string {
	var envMode string

	switch os.Getenv(gin.EnvGinMode) {
	case gin.DebugMode:
		envMode = "dev"
	case gin.TestMode:
		envMode = "test"
	case gin.ReleaseMode:
		envMode = "prod"
	default:
		envMode = "dev"
	}

	return envMode
}

func main() {
	// 载入配置文件
	loadConfig()

	// 日志初始化
	if global.Logger = initZap(); global.Logger == nil {
		os.Exit(1)
	}

	// 数据库初始化
	global.DB = initDB()

	// 路由初始化
	var router *gin.Engine
	if gin.IsDebugging() {
		router = gin.Default()
	} else {
		router = gin.New()
	}

	// 运行进程
	initRouter(router)
	router.Run(":" + strconv.Itoa(global.ServerConfig.System.Addr))
}

// @function: loadConfig
// @description: 载入配置文件
// @param: path string
// @return
func loadConfig() {

	path := fmt.Sprintf("config.%s.yaml", getEnvMode())

	v := viper.New()
	v.SetConfigFile(path)
	v.SetConfigType("yaml")
	err := v.ReadInConfig()
	if err != nil {
		log.Fatalf("Fatal error config file: %s \n", err)
	}
	v.WatchConfig()

	v.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("config file changed:", e.Name)
		if err = v.Unmarshal(&global.ServerConfig); err != nil {
			fmt.Println(err)
		}
	})
	if err = v.Unmarshal(&global.ServerConfig); err != nil {
		fmt.Println(err)
	}

	fmt.Println("Loading Config: " + path)
}

// @function: initZap
// @description: 初期化日志
// @param
// @return logger *zap.Logger
func initZap() (logger *zap.Logger) {
	fmt.Println("Initializing Logger")

	if ok, _ := utils.PathExists(global.ServerConfig.Zap.Director); !ok { // 判断是否有Director文件夹
		fmt.Printf("create %v directory\n", global.ServerConfig.Zap.Director)
		_ = os.Mkdir(global.ServerConfig.Zap.Director, os.ModePerm)
	}

	cores := logzap.Zap.GetZapCores()
	logger = zap.New(zapcore.NewTee(cores...))

	if global.ServerConfig.Zap.ShowLine {
		logger = logger.WithOptions(zap.AddCaller())
	}
	return logger
}

// @function: initDB
// @description: 初期化数据库
// @param
// @return
func initDB() *gorm.DB {
	m := global.ServerConfig.Mysql
	if m.Dbname == "" {
		fmt.Println("Database not enabled")
		return nil
	}

	fmt.Println("Initializing Database")

	mysqlConfig := mysql.Config{
		DSN:                       m.Dsn(), // DSN data source name
		DefaultStringSize:         256,     // string 类型字段的默认长度
		SkipInitializeWithVersion: false,   // 根据版本自动配置
	}

	if db, err := gorm.Open(mysql.New(mysqlConfig), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   m.Prefix,
			SingularTable: m.Singular,
		},
		DisableForeignKeyConstraintWhenMigrating: true,
	}); err != nil {
		return nil
	} else {
		db.InstanceSet("gorm:table_options", "ENGINE="+m.Engine)
		sqlDB, _ := db.DB()
		sqlDB.SetMaxIdleConns(m.MaxIdleConns)
		sqlDB.SetMaxOpenConns(m.MaxOpenConns)
		return db
	}
}

// @function: initRouter
// @description: 初期化路由
// @param: router *gin.Engine
// @return
func initRouter(router *gin.Engine) {
	uploadApi := api.ApiGroup.UploadApi
	router.Use(utils.JWTAuth())
	router.HEAD("/heartbeat", func(ctx *gin.Context) {})
	router.POST("/api/upload/single", uploadApi.UploadSingle)
	router.POST("/api/upload/multi", uploadApi.UploadMulti)
	router.POST("/api/upload/delete", uploadApi.UploadDelete)
	router.POST("/api/upload/covert", uploadApi.UploadCovert)
	router.GET("/api/upload/preview", uploadApi.UploadPreview)
	router.POST("/api/upload/generate/contract", uploadApi.UploadGenerateContract)
}
