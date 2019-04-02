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

import "strings"

// trans camel string to snake string, XxYy to xx_yy , XxYY to xx_yy
// @param c, custom snake separator, only '-' or '_' allowed
func Camel2Snake(s string, c ...byte) string {
	var sep byte = '_'
	if len(c) > 0 && c[0] == '-' {
		sep = c[0]
	}

	data := make([]byte, 0, len(s)*2)
	j := false
	num := len(s)

	for i := 0; i < num; i++ {
		d := s[i]
		if i > 0 && d >= 'A' && d <= 'Z' && j {
			data = append(data, sep)
		}

		if d != sep {
			j = true
		}

		data = append(data, d)
	}

	return strings.ToLower(string(data[:]))
}

// trans camel string to snake string, XxYy to xx_yy , XxYY to xx_yy
// @param c, custom snake separator, only '-' or '_' allowed
func Snake2Camel(s string, c ...byte) string {
	var sep byte = '_'
	if len(c) > 0 && c[0] == '-' {
		sep = c[0]
	}

	data := make([]byte, 0, len(s))
	j := false
	k := false
	num := len(s) - 1

	for i := 0; i <= num; i++ {
		d := s[i]
		if k == false && d >= 'A' && d <= 'Z' {
			k = true
		}

		if d >= 'a' && d <= 'z' && (j || k == false) {
			d = d - 32
			j = false
			k = true
		}

		if k && d == sep && num > i && s[i+1] >= 'a' && s[i+1] <= 'z' {
			j = true
			continue
		}

		data = append(data, d)
	}

	return string(data[:])
}
