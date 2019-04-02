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
	"github.com/gin-gonic/gin"
	"github.com/nvwa-io/wago/util"
	"os"
	"reflect"
	"strings"
)

var (
	CommentRouters = make(map[string][]CommentRouter)

	// HTTP_METHOD list the supported http methods.
	HTTP_METHOD = map[string]bool{
		"GET":       true,
		"POST":      true,
		"PUT":       true,
		"DELETE":    true,
		"PATCH":     true,
		"OPTIONS":   true,
		"HEAD":      true,
		"TRACE":     true,
		"CONNECT":   true,
		"MKCOL":     true,
		"COPY":      true,
		"MOVE":      true,
		"PROPFIND":  true,
		"PROPPATCH": true,
		"LOCK":      true,
		"UNLOCK":    true,
	}

	// while RouterMode = auto, don't register struct method to router
	EXCLUDE_ROUTER_METHOD = map[string]bool{
		"Init": true,
	}
)

type (
	MiddleWareHandler func(c *Context)
	CommentRouter     struct {
		Method     string
		Router     string
		HTTPMethod []string
	}
)

// add controller instance to set routers
type RouterGroup struct {
	prefix      string
	controllers []IController
	middleWares []MiddleWareHandler
}

func NewRouterGroup() *RouterGroup {
	return &RouterGroup{
		prefix:      "/",
		controllers: make([]IController, 0),
		middleWares: make([]MiddleWareHandler, 0),
	}
}

// set api prefixes
func (t *RouterGroup) Prefix(prefix string) *RouterGroup {
	t.prefix = prefix
	return t
}

// config controller routers
func (t *RouterGroup) Controller(controllers ...IController) *RouterGroup {
	t.controllers = controllers
	return t
}

// bind middleWares
func (t *RouterGroup) Use(middleWares ...MiddleWareHandler) *RouterGroup {
	t.middleWares = middleWares
	return t
}

// Config to WagoApp.Server.Group()
// run while wago app boot
func (t *RouterGroup) config() {
	group := WagoApp.Server.Group(t.prefix)
	if t.prefix == "" {
		t.prefix = "/"
	}

	// @TODO check
	// or use Pointer() to change
	//if valueOf.Type().Kind() == reflect.Ptr {
	//	fmt.Println(reflect.Indirect(valueOf).Type().Name())
	//} else {
	//	fmt.Println(valueOf.Type().Name())
	//}

	// register routers
	for _, c := range t.controllers {
		switch AppConfig.App.RouterMode {
		case ROUTER_MODE_COMMENT:
			t.registerRouterByComment(c, group)
		default:
			t.registerRouterByAuto(c, group)
		}
	}

	// use middleWares
	for _, f := range t.middleWares {
		group.Use(func(c *gin.Context) {
			f(&Context{c})
		})
	}
}

// comment mode: use comment to declare restful routers
// notice: golang doesn't have virtual machine (e.g: java -> jvm),
// so we aren't able to reflect to get comment info at runtime
// so here, we use golang's ast pkg to parse controllers' *.go file to auto-generate router configuration codes
// refer to comment.go
func (t *RouterGroup) registerRouterByComment(c IController, group *gin.RouterGroup) {
	v := reflect.ValueOf(c)
	vi := reflect.Indirect(v)

	// pkgPath:controllerName
	// e.g: github.com/nvwa-io/wago-example/controller/ExampleController
	fullPath := fmt.Sprintf("%s/%s", vi.Type().PkgPath(), vi.Type().Name())

	// because auto-generated router codes only keep last pkg path and file name
	// e.g: controller/ExampleController
	// or: controller/home/Example2Controller
	for key, rs := range CommentRouters {
		if !strings.HasSuffix(fullPath, key) {
			continue
		}

		for _, v := range rs {
			for _, hm := range v.HTTPMethod {
				group.Handle(hm, v.Router, HandlerWrapper(vi.Type(), v.Method))
			}
		}
	}
}

// auto mode: auto register routers by struct name and method name,
// [GET, POST] HTTP method are registered while method has not explicit the request method.
// e.g:
// func (t *ExampleController) HelloWorld() , [GET, POST] method will be registered.
// func (t *ExampleController) HelloWorld_POST() , [POST] method will be registered.

// while sep = 0 (means no config for router separator), use struct method name as router path
// while sep equal '-' or '_', use snake string as router path
func (t *RouterGroup) registerRouterByAuto(c IController, group *gin.RouterGroup) {
	v := reflect.ValueOf(c)
	vi := reflect.Indirect(v)
	typ := reflect.TypeOf(c)
	var sep byte
	if len(AppConfig.App.RouterSep) > 0 {
		if AppConfig.App.RouterSep[0] == '_' {
			sep = '_'
		} else {
			sep = '-'
		}
	}

	// pkg path
	routerPathPkg := ""
	cntlSep := fmt.Sprintf("%s%s%s", string(os.PathSeparator), AppConfig.App.ControllerPath, string(os.PathSeparator))
	arr := strings.Split(vi.Type().PkgPath(), cntlSep)
	if len(arr) <= 1 {
		routerPathPkg = "/"
	} else {
		routerPathPkg = "/" + arr[1]
	}

	// controller path
	controllerName := vi.Type().Name()
	routerPathCntl := ""
	if sep != 0 {
		routerPathCntl = util.Camel2Snake(strings.TrimSuffix(controllerName, "Controller"), sep)
	}

	// method path
	for i := 0; i < typ.NumMethod(); i++ {
		methodName := typ.Method(i).Name
		if _, ok := EXCLUDE_ROUTER_METHOD[methodName]; ok {
			continue
		}

		routerPathMethod := methodName
		mHttpMethods := make([]string, 0)

		// deal with explicit HTTP request METHOD
		// e.g: UpdateInfo_PUT()
		mArr := strings.Split(methodName, "_")
		if len(mArr) > 1 {
			last := mArr[len(mArr)-1]
			upLast := strings.ToUpper(last)
			if _, ok := HTTP_METHOD[upLast]; ok {
				mHttpMethods = append(mHttpMethods, upLast)
				routerPathMethod = strings.TrimSuffix(routerPathMethod, "_"+last)
			}
		}
		if len(mHttpMethods) == 0 {
			mHttpMethods = []string{"GET", "POST"}
		}

		// organize full request path: /{PKG NAME}/{CONTROLLER NAME}/{METHOD NAME}
		routerPathMethod = util.Camel2Snake(routerPathMethod, sep)
		requestPath := fmt.Sprintf("/%s/%s/%s",
			strings.Trim(routerPathPkg, "/"),
			strings.Trim(routerPathCntl, "/"),
			strings.Trim(routerPathMethod, "/"))
		requestPath = "/" + strings.TrimLeft(requestPath, "/") // maybe rootPathPkg = "/"
		for _, hm := range mHttpMethods {
			group.Handle(hm, requestPath, HandlerWrapper(vi.Type(), methodName))
		}
	}
}
