package cmd

import (
	"context"
	"fmt"
	"mvtable/internal/storage/db"
	"mvtable/internal/storage/redis"
	"mvtable/pkg/jwt"
	"mvtable/pkg/mail"
	"mvtable/pkg/oss"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"mvtable/internal/di"
	"mvtable/internal/pkg/config"
	"mvtable/pkg/log"
)

// serverCmd 服务器命令
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "启动HTTP服务器",
	Long:  `启动HTTP服务器，提供API服务`,
	Run:   runServer,
}

func init() {
	serverCmd.Flags().StringVarP(&configFile, "config", "c", "", "配置文件路径")
	serverCmd.Flags().StringVarP(&env, "env", "e", "dev", "环境 (dev|prod)")
}

func runServer(cmd *cobra.Command, args []string) {
	// 确定配置文件路径
	if configFile == "" {
		configFile = fmt.Sprintf("configs/config.%s.yaml", env)
	}

	// 加载配置
	cfg, err := config.Load(configFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "加载配置失败: %v\n", err)
		os.Exit(1)
	}

	// 初始化日志
	if err = log.Init(&cfg.Log); err != nil {
		fmt.Fprintf(os.Stderr, "初始化日志失败: %v\n", err)
		os.Exit(1)
	}

	// 初始化数据库
	if err = db.Init(&cfg.Database); err != nil {
		fmt.Fprintf(os.Stderr, "初始化数据库失败: %v\n", err)
		os.Exit(1)
	}

	// 初始化 redis
	if err = redis.Init(&cfg.Redis); err != nil {
		fmt.Fprintf(os.Stderr, "初始化 Redis 失败: %v\n", err)
		os.Exit(1)
	}

	// 初始化邮件
	mail.Init(&cfg.Mail)

	// 初始化oss
	if err = oss.Init(&cfg.OSS); err != nil {
		fmt.Fprintf(os.Stderr, "初始化 OSS 失败: %v\n", err)
		os.Exit(1)
	}

	log.Info("Starting server",
		zap.String("config", configFile),
		zap.String("env", env),
	)

	// 初始化 JWT
	jwt.Init(cfg.JWT)

	// 初始化依赖注入
	app, cleanup, err := di.InitializeApp(cfg)
	if err != nil {
		log.Fatal("初始化应用失败", zap.Error(err))
	}
	defer cleanup()

	// 启动服务器
	go func() {
		if err = app.HTTPServer.Start(); err != nil {
			log.Fatal("启动HTTP服务器失败", zap.Error(err))
		}
	}()

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down server...")

	// 优雅关闭
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err = app.HTTPServer.Stop(ctx); err != nil {
		log.Error("停止HTTP服务器失败", zap.Error(err))
	}

	log.Info("Server stopped")
}
