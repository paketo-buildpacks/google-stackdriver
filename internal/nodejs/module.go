/*
 * Copyright 2018-2020 the original author or authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package nodejs

import (
	"bytes"
	"fmt"
	"regexp"
	"text/template"
)

func IsModuleRequired(module string, content []byte) (bool, error) {
	p := fmt.Sprintf(`require\(['"]%s['"]\)`, module)

	r, err := regexp.Compile(p)
	if err != nil {
		return false, fmt.Errorf("unable to compiler regex '%s'\n%w", p, err)
	}

	return r.Match(content), nil
}

func RequireModule(context ModuleContext) ([]byte, error) {
	t, err := template.New("require-module").
		Parse(`require('{{.Module}}').start({
  serviceContext: {
    service: '{{.Service}}',
    version: '{{.Version}}',
  },
});
`)
	if err != nil {
		return nil, fmt.Errorf("unable to parse template\n%w", err)
	}

	b := &bytes.Buffer{}
	if err := t.Execute(b, context); err != nil {
		return nil, fmt.Errorf("unable to execute template\n%w", err)
	}

	return b.Bytes(), nil
}

func RequireModuleExternal(context ModuleContext) ([]byte, error) {
	t, err := template.New("require-module-external").
		Parse(`require('{{.Module}}').start({
  projectId: '{{.ProjectId}}',
  serviceContext: {
    service: '{{.Service}}',
    version: '{{.Version}}',
  },
});
`)
	if err != nil {
		return nil, fmt.Errorf("unable to parse template\n%w", err)
	}

	b := &bytes.Buffer{}
	if err := t.Execute(b, context); err != nil {
		return nil, fmt.Errorf("unable to execute template\n%w", err)
	}

	return b.Bytes(), nil
}

type ModuleContext struct {
	Module    string
	ProjectId string
	Service   string
	Version   string
}
