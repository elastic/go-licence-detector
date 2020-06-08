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

package detector

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoadRules(t *testing.T) {
	t.Run("embedded", func(t *testing.T) {
		rules, err := LoadRules("")

		require.NoError(t, err)
		require.NotNil(t, rules)
		require.True(t, len(rules.WhiteList) > 0)
	})

	t.Run("external", func(t *testing.T) {
		rules, err := LoadRules("testdata/rules.json")

		require.NoError(t, err)
		require.NotNil(t, rules)
		require.True(t, len(rules.WhiteList) > 0)
	})
}

func TestRulesWhiteList(t *testing.T) {
	rules, err := LoadRules("testdata/rules.json")

	require.NoError(t, err)
	require.True(t, rules.IsAllowed("Apache-2.0"))
	require.False(t, rules.IsAllowed("WTFPL"))
}

func TestRulesYellowList(t *testing.T) {
	rules, err := LoadRules("testdata/rules.json")

	require.NoError(t, err)
	require.True(t, rules.IsAllowed("GPL-3.0"))
	require.False(t, rules.IsAllowed("WTFPL"))
}
