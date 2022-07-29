package main

import (
	"context"
	"fmt"
	"go-jichu/controllers"
	"go-jichu/dao/mysql"
	"go-jichu/dao/redis"
	"go-jichu/logger"
	"go-jichu/pkg/snowflake"
	"go-jichu/routes"
	"go-jichu/settings"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func main() {

	//1.加载配置
	if err := settings.Init(); err != nil {
		fmt.Printf("init settings failed, err:%v\n", err)
		return
	}
	//2.初始化日志
	if err := logger.Init(settings.Conf.LogConfig, viper.GetString("app.mode")); err != nil {
		fmt.Printf("init logger failed, err:%v\n", err)
		return
	}
	defer zap.L().Sync()

	//3.初始化mysql
	if err := mysql.Init(settings.Conf.MySQLConfig); err != nil {
		fmt.Printf("init mysql failed, err:%v\n", err)
		return
	}
	defer mysql.Close()

	//4.初始化redis
	if err := redis.Init(settings.Conf.RedisConfig); err != nil {
		fmt.Printf("init redis failed, err:%v\n", err)
		return
	}
	defer redis.Close()
	if err := snowflake.Init(viper.GetString("app.start_time"), viper.GetInt64("app.machine_id")); err != nil {
		fmt.Printf("init snowflake failed, err:%v\n", err)
		return
	}

	if err := controllers.InitTrans("zh"); err != nil {
		fmt.Printf("init snowflake failed, err#{err}\n")
		return
	}

	//5.注册路由
	r := routes.SetUp(viper.GetString("app.mode"))

	//6.启动服务（优雅关机）
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", viper.GetInt("app.port")),
		Handler: r,
	}

	//开启一个goroutine启动服务
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen:%s\n", err)
		}
	}()

	//等待中断信号来优化的关闭服务器，为关闭服务器操作设置一个5秒的超时
	quit := make(chan os.Signal, 1) //创建一个接收信号的通道

	//kill 默认会发送 syscall.sigterm 信号
	//kill -2 发送 syscall.sigint 信号，通常用 ctrl+c 就是触发系统的sigint信号
	//kill -9 发送 syscall.sigkill 信号，但是不能被捕获，所以不需要添加它
	//signal notify 把收到的 syscall.sigint 或 syscall.sigterm 信号转发给quit
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) //此处不会阻塞
	<-quit                                               //阻塞在此，当接收到上述两种信号时才会往下继续执行
	log.Println("Shutdown Server ...")

	//创建一个5秒超市的context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	//5秒内优雅关闭服务，（将未处理完的请求处理完再关闭服务），超过5秒就超时退出
	if err := srv.Shutdown(ctx); err != nil {
		zap.L().Fatal("server shutdown:", zap.Error(err))
	}

	zap.L().Info("server exiting")

}
