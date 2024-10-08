// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package validate

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"go.elastic.co/go-licence-detector/dependency"
)

func TestValidateURLs(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r.Body != nil {
				_, _ = io.Copy(io.Discard, r.Body)
				r.Body.Close()
			}
		}()

		if r.Method == http.MethodHead && r.URL.Query().Get("no_head") == "true" {
			http.Error(w, "method not supported", http.StatusMethodNotAllowed)
			return
		}

		if r.URL.Query().Get("valid") == "true" {
			fmt.Fprintln(w, "OK")
		} else {
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	mkDepInfo := func(name string, valid, noHead bool) dependency.Info {
		return dependency.Info{
			Name: name,
			URL:  fmt.Sprintf("%s/%s?valid=%t&no_head=%t", server.URL, name, valid, noHead),
		}
	}

	testCases := []struct {
		name    string
		deps    *dependency.List
		wantErr bool
	}{
		{
			name: "AllValid",
			deps: &dependency.List{
				Direct:   []dependency.Info{mkDepInfo("a", true, false), mkDepInfo("b", true, false)},
				Indirect: []dependency.Info{mkDepInfo("c", true, false), mkDepInfo("d", true, false)},
			},
		},
		{
			name: "AllValidWithUnsupportedMethod",
			deps: &dependency.List{
				Direct:   []dependency.Info{mkDepInfo("a", true, false), mkDepInfo("b", true, true)},
				Indirect: []dependency.Info{mkDepInfo("c", true, false), mkDepInfo("d", true, true)},
			},
		},
		{
			name: "InvalidDirectDep",
			deps: &dependency.List{
				Direct:   []dependency.Info{mkDepInfo("a", true, false), mkDepInfo("b", false, false)},
				Indirect: []dependency.Info{mkDepInfo("c", true, false), mkDepInfo("d", true, false)},
			},
			wantErr: true,
		},
		{
			name: "InvalidDirectDepWithUnsupportedMethod",
			deps: &dependency.List{
				Direct:   []dependency.Info{mkDepInfo("a", true, false), mkDepInfo("b", false, true)},
				Indirect: []dependency.Info{mkDepInfo("c", true, false), mkDepInfo("d", true, false)},
			},
			wantErr: true,
		},
		{
			name: "InvalidIndirectDep",
			deps: &dependency.List{
				Direct:   []dependency.Info{mkDepInfo("a", true, false), mkDepInfo("b", true, false)},
				Indirect: []dependency.Info{mkDepInfo("c", true, false), mkDepInfo("d", false, false)},
			},
			wantErr: true,
		},
		{
			name: "InvalidIndirectDepWithUnsupportedMethod",
			deps: &dependency.List{
				Direct:   []dependency.Info{mkDepInfo("a", true, false), mkDepInfo("b", true, false)},
				Indirect: []dependency.Info{mkDepInfo("c", true, false), mkDepInfo("d", false, true)},
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := validateURLs(tc.deps)
			if tc.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
