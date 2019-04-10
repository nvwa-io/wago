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

import "github.com/nvwa-io/wago/logger"

type Controller struct {
	// context of current request
	Ctx *Context

	// controller logger
	// init request ID
	Logger *logger.Entry
}

// init request context
func (t *Controller) Init(c *Context) {
	t.Ctx = c

	t.Logger = logger.WithFields(logger.Fields{
		REQUEST_ID: c.GetString(REQUEST_ID),
	})
}

// get request id from context
func (t *Controller) RequestId() string {
	return t.Ctx.GetString(REQUEST_ID)
}
