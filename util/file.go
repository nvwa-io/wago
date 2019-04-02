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
	"io/ioutil"
	"os"
	"strings"
)

func ScanLevel1Dirs(path string) ([]string, error) {
	dir, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}

	dirs := make([]string, 0)
	for _, fi := range dir {
		if fi.IsDir() && !strings.HasPrefix(fi.Name(), ".") {
			dirs = append(dirs, path+string(os.PathSeparator)+fi.Name())
		}
	}

	return dirs, nil
}

func ScanFilesAndDirs(dirPth string) (files []string, dirs []string, err error) {
	dir, err := ioutil.ReadDir(dirPth)
	if err != nil {
		return nil, nil, err
	}

	sep := string(os.PathSeparator)
	for _, fi := range dir {
		if fi.IsDir() {
			dirs = append(dirs, dirPth+sep+fi.Name())
			ScanFilesAndDirs(dirPth + sep + fi.Name())
		} else {
			if ok := strings.HasSuffix(fi.Name(), ".go"); ok {

				files = append(files, dirPth+sep+fi.Name())
			}
		}
	}

	return files, dirs, nil
}
