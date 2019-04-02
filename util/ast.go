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

package util

import (
	"go/parser"
	"go/token"
	"os"
	"strings"
)

// use go ast pkg to get target path's pkg name
func ASTPkgName(targetPath string) (string, error) {
	astPkgs, err := parser.ParseDir(
		token.NewFileSet(),
		targetPath,
		func(info os.FileInfo) bool { // filter .go files
			name := info.Name()
			return !info.IsDir() && !strings.HasPrefix(name, ".") && strings.HasSuffix(name, ".go")
		},
		parser.ParseComments)
	if err != nil {
		return "", err
	}

	// get pkg name by ast
	pkgName := ""
	for _, pkg := range astPkgs {
		pkgName = pkg.Name
		break
	}

	if pkgName != "" {
		return pkgName, nil
	}

	// if there isn't *.go file in target path
	// then return target path's dir name
	arr := strings.Split(targetPath, string(os.PathSeparator))
	return arr[len(arr)-1], nil
}
