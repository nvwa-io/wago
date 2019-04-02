package wago

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

import (
	"github.com/gin-gonic/gin"
	"reflect"
)

type (
	IController interface {
		Init(*Context)
	}

	// Wrapper for controller functions
	// which is type of gin.HandlerFunc
	// WagoHandler gin.HandlerFunc
	WrapperFunc func()
)

// HandlerWrapper is a wrapper to transform controller'method to gin.HandlerFunc
// @TODO catch exceptions while method not exist
// encapsulate controller's method in gin.HandlerFunc
// means: while gin.HandlerFunc is invoked, the target controller's method will be invoked
func HandlerWrapper(controllerType reflect.Type, method string) gin.HandlerFunc {
	return func(c *gin.Context) {
		ct := reflect.New(controllerType)
		controller := ct.Interface().(IController)
		controller.Init(&Context{Context: c})
		m := ct.MethodByName(method)

		m.Call(nil)
	}
}
