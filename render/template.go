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

const (
	noticeTxt = `
{{- define "depInfo" -}}
{{- range $i, $dep := . }}
{{ "-" | line }}
Module  : {{ $dep.Name }}
Version : {{ $dep.Version }}
Time    : {{ $dep.VersionTime }}
Licence : {{ $dep.LicenceType }}

{{ $dep | licenceText }}
{{ end }}
{{- end -}}

Copyright 2019-{{ currentYear }} Elasticsearch BV

This product includes software developed by The Apache Software
Foundation (http://www.apache.org/).

{{ "=" | line }}
Third party libraries used by the go-licence-detector project
{{ "=" | line }}

{{ template "depInfo" .Direct }}

{{ if .Indirect }}
{{ "=" | line }}
Indirect dependencies

{{ template "depInfo" .Indirect }}
{{ end }}
`
)
