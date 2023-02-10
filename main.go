package main

import (
	"github.com/ashoreDove/common"
	"github.com/ashoreDove/parasite-song/handler"
	song "github.com/ashoreDove/parasite-song/proto/song"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/micro/go-micro/v2"
	log "github.com/micro/go-micro/v2/logger"
	ratelimit "github.com/micro/go-plugins/wrapper/ratelimiter/uber/v2"
	opentracing2 "github.com/micro/go-plugins/wrapper/trace/opentracing/v2"
	"github.com/opentracing/opentracing-go"
)

var QPS = 100

func main() {
	cfg, err := common.Init(true)
	if err != nil {
		log.Error(err)
		return
	}
	//ftp服务
	defer cfg.FtpConn.Quit()

	//链路追踪
	t, io, err := common.NewTracer("go.micro.service.song", "localhost:6831")
	if err != nil {
		log.Fatal(err)
	}
	defer io.Close()
	opentracing.SetGlobalTracer(t)
	// New Service
	service := micro.NewService(
		micro.Name("go.micro.service.song"),
		micro.Version("latest"),
		//设置地址和需要暴露的端口
		micro.Address("127.0.0.1:8083"),
		//添加consul 作为注册中心
		micro.Registry(*cfg.ConsulRegister),
		//绑定链路追踪
		micro.WrapHandler(opentracing2.NewHandlerWrapper(opentracing.GlobalTracer())),
		//添加限流
		//QPS：每秒处理请求数量
		micro.WrapHandler(ratelimit.NewHandlerWrapper(QPS)),
	)
	defer cfg.DB.Close()
	//禁止副表
	cfg.DB.SingularTable(true)

	// Initialise service
	service.Init()

	// Register Handler
	err = song.RegisterSongHandler(service.Server(), handler.NewSongService(cfg.DB, cfg.FtpConn))
	if err != nil {
		log.Fatal(err)
	}

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
