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
	"flag"
	"github.com/pelletier/go-toml"
	"log"
)

const (
	RUN_MODE_DEBUG   = "debug"
	RUN_MODE_TEST    = "test"
	RUN_MODE_RELEASE = "release"

	ROUTER_MODE_AUTO    = "auto"
	ROUTER_MODE_COMMENT = "comment"
)

var (
	AppConfig = &Config{}
)

func init() {
	configFile := flag.String("c", "config/app.toml", "configuration file path")
	flag.Parse()

	tree, err := toml.LoadFile(*configFile)
	if err != nil {
		log.Printf("fail to read file: %s \n", err.Error())
		return
	}
	err = tree.Unmarshal(AppConfig)
	if err != nil {
		log.Printf("fail to unmarshal wago config: %s \n", err.Error())
		return
	}
	AppConfig.Tree = *tree

	fullFillConfig()
}

// @TODO Full fill more fields
func fullFillConfig() {
	if AppConfig.App.ControllerPath == "" {
		AppConfig.App.ControllerPath = "controller"
	}

	// RunMode / RouterMode
}

type Config struct {
	toml.Tree

	// app config
	App App

	// HTTP server configuration
	Server Server
}

type App struct {
	App     string
	RunMode string

	// router mode, [auto, comment] supported.
	RouterMode string

	// only when RouterMode=comment, ControllerPath is valid.
	ControllerPath string

	// only when RouterMode=auto, RouterSep is valid.
	// only ''(empty), '-' or '_' allowed.
	// '' means use struct/method name as router path, eg: /v1/HomeTest/HelloWorld
	// '-' means use snake string as router path, eg: /v1/home-test/hello-world
	// '_' means use snake string as router path, eg: /v1/home_test/hello_world
	RouterSep string
}

// TLSConfig *tls.Config
// TLSNextProto map[string]func(*Server, *tls.Conn, Handler)
// ConnState func(net.Conn, ConnState)
// HTTP server configuration
type Server struct {
	Port              int
	Host              string
	ReadTimeout       int
	ReadHeaderTimeout int
	WhiteTimeout      int
	IdleTimeout       int
	MaxHeaderBytes    int
}
