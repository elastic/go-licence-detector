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

package render // import "go.elastic.co/go-licence-detector/render"

import (
	"bytes"
	"fmt"
	"go/build"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"
	"time"

	"go.elastic.co/go-licence-detector/dependency"
)

var goModCache = filepath.Join(build.Default.GOPATH, "pkg", "mod")

func Template(dependencies *dependency.List, templatePath, outputPath string) error {
	funcMap := template.FuncMap{
		"currentYear": CurrentYear,
		"line":        Line,
		"licenceText": LicenceText,
	}
	tmpl, err := template.New(filepath.Base(templatePath)).Funcs(funcMap).ParseFiles(templatePath)
	if err != nil {
		return fmt.Errorf("failed to parse template at %s: %w", templatePath, err)
	}

	w, cleanup, err := mkWriter(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file %s: %w", outputPath, err)
	}
	defer cleanup()

	if err := tmpl.Execute(w, dependencies); err != nil {
		return fmt.Errorf("failed to render template: %w", err)
	}

	return nil
}

func mkWriter(path string) (io.Writer, func(), error) {
	if path == "-" {
		return os.Stdout, func() {}, nil
	}

	f, err := os.Create(path)
	return f, func() { f.Close() }, err
}

/* Template functions */

func CurrentYear() string {
	return strconv.Itoa(time.Now().Year())
}

func Line(ch string) string {
	return strings.Repeat(ch, 80)
}

func LicenceText(depInfo dependency.Info) string {
	if depInfo.LicenceFile == "" {
		return "No licence file provided."
	}

	var buf bytes.Buffer
	if depInfo.LicenceTextOverrideFile != "" {
		buf.WriteString("Contents of provided licence file")
	} else {
		buf.WriteString("Contents of probable licence file ")
		buf.WriteString(strings.Replace(depInfo.LicenceFile, goModCache, "$GOMODCACHE", -1))
	}
	buf.WriteString(":\n\n")

	f, err := os.Open(depInfo.LicenceFile)
	if err != nil {
		log.Fatalf("Failed to open licence file %s: %v", depInfo.LicenceFile, err)
	}
	defer f.Close()

	_, err = io.Copy(&buf, f)
	if err != nil {
		log.Fatalf("Failed to read licence file %s: %v", depInfo.LicenceFile, err)
	}

	return buf.String()
}
