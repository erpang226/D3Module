package adminModule

import (
	"context"
	"errors"
	"log"
	core "main/pkg/actor/core"
	actorContext "main/pkg/actor/core/context"
	"main/pkg/admin/handler"
	"main/pkg/admin/router/module"
	"main/pkg/admin/router/plugin"
	"main/pkg/admin/router/system"
	"main/pkg/config"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
	"k8s.io/klog/v2"
)

type AdminModule struct {
	enable bool
	name   string
	server *http.Server
}

func newAdminModule(config *config.AdminServerConfig) *AdminModule {
	return &AdminModule{
		enable: config.Enable,
		name:   config.Name,
	}
}
func Register(config *config.AdminServerConfig) {
	adminModule := newAdminModule(config)

	core.Register(adminModule)
}

func (a *AdminModule) Name() string {
	return "admin"
}

func (a *AdminModule) Group() string {
	return "admin"
}

func (a *AdminModule) Enable() bool {
	return a.enable
}

func (a *AdminModule) SetEnable(enable bool) {
	a.enable = enable
}

func (a *AdminModule) Start() {
	router := gin.Default()
	router.Use(handler.Auth())
	adminGroup := router.Group("/admin")
	system.InitRouter(adminGroup)
	module.InitRouter(adminGroup)
	plugin.InitRouter(adminGroup)

	a.server = &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		// 服务连接
		if err := a.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	klog.Info("shutdown gin server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := a.server.Shutdown(ctx); err != nil {
		klog.Fatal("gin server shutdown:", err)
	}
	klog.Info("gin server exit")
}

func (a *AdminModule) Stop() {
	actorContext.Cleanup(a.Name())
	// todo
}
