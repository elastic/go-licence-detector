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

package render

import (
	"testing"
)

func TestKeyValueFlags_Set(t *testing.T) {
	var kvs KeyValueFlags

	tests := []struct {
		input   string
		wantErr bool
		wantLen int
		wantKey string
		wantVal string
	}{
		{"foo=bar", false, 1, "foo", "bar"},
		{"baz=qux", false, 2, "baz", "qux"},
		{"", false, 2, "", ""},
		{"invalidpair", true, 2, "", ""},
		{"=novalue", true, 2, "", ""},
		{"nokey=", true, 2, "", ""},
	}

	for _, tt := range tests {
		err := kvs.Set(tt.input)
		if tt.wantErr && err == nil {
			t.Errorf("Set(%q) expected error, got nil", tt.input)
		}
		if !tt.wantErr && err != nil {
			t.Errorf("Set(%q) unexpected error: %v", tt.input, err)
		}
		if len(kvs) != tt.wantLen {
			t.Errorf("Set(%q) expected len %d, got %d", tt.input, tt.wantLen, len(kvs))
		}
		if tt.wantKey != "" && tt.wantVal != "" {
			found := false
			for _, kv := range kvs {
				if kv.Key == tt.wantKey && kv.Value == tt.wantVal {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Set(%q) expected key-value (%q,%q) not found", tt.input, tt.wantKey, tt.wantVal)
			}
		}
	}
}

func TestKeyValueFlags_Get(t *testing.T) {
	var kvs KeyValueFlags
	kvs = append(kvs, KeyValue{Key: "foo", Value: "bar"})
	kvs = append(kvs, KeyValue{Key: "baz", Value: "qux"})

	tests := []struct {
		name string
		kvs  *KeyValueFlags
		key  string
		want string
	}{
		{"existing key", &kvs, "foo", "bar"},
		{"another key", &kvs, "baz", "qux"},
		{"missing key", &kvs, "notfound", ""},
		{"nil receiver", nil, "foo", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.kvs.Get(tt.key)
			if got != tt.want {
				t.Errorf("Get(%q) = %q, want %q", tt.key, got, tt.want)
			}
		})
	}
}
