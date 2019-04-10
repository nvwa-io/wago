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
	"github.com/sirupsen/logrus"
	"log"
)

const (
	RUN_MODE_DEBUG   = "debug"
	RUN_MODE_TEST    = "test"
	RUN_MODE_RELEASE = "release"

	ROUTER_MODE_AUTO    = "auto"
	ROUTER_MODE_COMMENT = "comment"

	// PanicLevel level, highest level of severity. Logs and then calls panic with the
	// message passed to Debug, Info, ...
	LevelPanic = logrus.PanicLevel
	// FatalLevel level. Logs and then calls `logger.Exit(1)`. It will exit even if the
	// logging level is set to Panic.
	LevelFatal = logrus.FatalLevel
	// ErrorLevel level. Logs. Used for errors that should definitely be noted.
	// Commonly used for hooks to send errors to an error tracking service.
	LevelError = logrus.ErrorLevel
	// WarnLevel level. Non-critical entries that deserve eyes.
	LevelWarn = logrus.WarnLevel
	// InfoLevel level. General operational entries about what's going on inside the
	// application.
	LevelInfo = logrus.InfoLevel
	// DebugLevel level. Usually only enabled when debugging. Very verbose logging.
	LevelDebug = logrus.DebugLevel
	// TraceLevel level. Designates finer-grained informational events than the Debug.
	LevelTrace = logrus.TraceLevel
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

	// app configuration
	App App

	// log configuration
	Log Log

	// HTTP server configuration
	Server Server
}

type (
	// application configuration
	App struct {
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

	// log configurations
	Log struct {
		// [json,text,logstash,fluentd] supported
		Formatter string

		// Log level
		Level uint32

		// logging method name
		LogMethodName bool

		// Filename is the file to write logs to.  Backup log files will be retained
		// in the same directory.  It uses <processname>-lumberjack.log in
		// os.TempDir() if empty.
		Filename string

		// MaxSize is the maximum size in megabytes of the log file before it gets
		// rotated. It defaults to 100 megabytes.
		MaxSize int

		// MaxBackups is the maximum number of old log files to retain.  The default
		// is to retain all old log files (though MaxAge may still cause them to get
		// deleted.)
		MaxBackups int

		// MaxAge is the maximum number of days to retain old log files based on the
		// timestamp encoded in their filename.  Note that a day is defined as 24
		// hours and may not exactly correspond to calendar days due to daylight
		// savings, leap seconds, etc. The default is not to remove old log files
		// based on age.
		MaxAge int

		// Compress determines if the rotated log files should be compressed
		// using gzip.
		Compress bool

		// print log to console
		Console bool
	}

	// TLSConfig *tls.Config
	// TLSNextProto map[string]func(*Server, *tls.Conn, Handler)
	// ConnState func(net.Conn, ConnState)
	// HTTP server configuration
	Server struct {
		Port              int
		Host              string
		ReadTimeout       int
		ReadHeaderTimeout int
		WhiteTimeout      int
		IdleTimeout       int
		MaxHeaderBytes    int

		// HTTP CORS configuration
		Cors Cors
	}

	// HTTP CORS configuration
	Cors struct {
		AllowOrigins     []string
		AllowMethods     []string
		AllowHeaders     []string
		ExposeHeaders    []string
		AllowCredentials bool
		AllowOriginFunc  func(string) bool
		MaxAge           int
	}
)
