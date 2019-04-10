// Copyright 2019 - now The https://github.com/nvwa-io/wago Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package wago

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/nvwa-io/wago/logger"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	WagoApp *Wago
)

func init() {
	WagoApp = NewWago()
}

func NewWago() *Wago {
	w := &Wago{
		Server: gin.New(),
	}

	return w
}

type Wago struct {
	// Use gin as http server
	Server *gin.Engine

	// router groups for configuring HTTP request handler
	routerGroups []*RouterGroup
}

func (t *Wago) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	WagoApp.Server.ServeHTTP(w, req)
}

func config() {
	// config gin engine running mode.
	gin.SetMode(AppConfig.App.RunMode)

	// generate comment routers while RunMode is debug and RouterMode is "comment"
	if AppConfig.App.RunMode == RUN_MODE_DEBUG &&
		AppConfig.App.RouterMode == ROUTER_MODE_COMMENT {
		ParseRouter(AppConfig.App.ControllerPath)
	}

	// register routers
	for _, rg := range WagoApp.routerGroups {
		rg.config()
	}

	// config logger
	logger.SetFormatterByString(AppConfig.Log.Formatter)
	logger.SetLevelByUint32(AppConfig.Log.Level)
	logger.SetReportCaller(AppConfig.Log.LogMethodName)
	if AppConfig.Log.Console {
		w := io.MultiWriter(
			os.Stdout,
			&lumberjack.Logger{
				Filename:   AppConfig.Log.Filename,
				MaxSize:    AppConfig.Log.MaxSize,
				MaxBackups: AppConfig.Log.MaxBackups,
				MaxAge:     AppConfig.Log.MaxAge,
				Compress:   AppConfig.Log.Compress,
			},
		)
		logger.SetOutput(w)
	} else {
		logger.SetOutput(&lumberjack.Logger{
			Filename:   AppConfig.Log.Filename,
			MaxSize:    AppConfig.Log.MaxSize,
			MaxBackups: AppConfig.Log.MaxBackups,
			MaxAge:     AppConfig.Log.MaxAge,
			Compress:   AppConfig.Log.Compress,
		})
	}

	// config HTTP Server CORS
	WagoApp.Server.Use(cors.New(cors.Config{
		AllowOrigins:     AppConfig.Server.Cors.AllowOrigins,
		AllowMethods:     AppConfig.Server.Cors.AllowMethods,
		AllowHeaders:     AppConfig.Server.Cors.AllowHeaders,
		ExposeHeaders:    AppConfig.Server.Cors.ExposeHeaders,
		AllowCredentials: AppConfig.Server.Cors.AllowCredentials,
		AllowOriginFunc:  AppConfig.Server.Cors.AllowOriginFunc,
		MaxAge:           time.Duration(AppConfig.Server.Cors.MaxAge) * time.Hour,
	}))
}

// Boot Wago app
func Serve() {

	// app configuration
	config()

	// start gin server
	server := &http.Server{
		Addr: fmt.Sprintf("%s:%d", AppConfig.Server.Host, AppConfig.Server.Port),
		//Handler: WagoApp.Server,
		Handler: WagoApp,
		// @TODO Config http server
	}

	//engine.Use(Logger(), Recovery())
	err := server.ListenAndServe()
	if err != nil {
		log.Fatalln(AppConfig.App, " finished, err=", err.Error())
	}
}

// Use attaches a global middleware to the router. ie. the middleware attached though Use() will be
// included in the handlers chain for every single request. Even 404, 405, static files...
// For example, this is the right place for a logger or error management middleware.
func Use(middleware ...MiddleWareHandler) gin.IRoutes {
	return WagoApp.Server.Use(middleware...)
}

// Add router groups
func AddRouterGroups(rgs ...*RouterGroup) {
	WagoApp.routerGroups = append(WagoApp.routerGroups, rgs...)
}

// set CORS AllowOriginFunc
func AllowOriginFunc(f func(string) bool) {
	AppConfig.Server.Cors.AllowOriginFunc = f
}
