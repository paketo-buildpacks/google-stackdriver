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

package credentials

import (
	"fmt"

	"github.com/buildpacks/libcnb"
	"github.com/paketo-buildpacks/libpak"
)

type Credentials struct {
	Bindings libcnb.Bindings
}

func (c Credentials) Execute() (string, error) {
	br := libpak.BindingResolver{Bindings: c.Bindings}

	if b, ok, err := br.Resolve("StackdriverDebugger"); err != nil {
		return "", fmt.Errorf("unable to resolve binding StackdriverDebugger\n%w", err)
	} else if ok {
		if p, ok := b.SecretFilePath("ApplicationCredentials"); ok {
			return fmt.Sprintf(`export GOOGLE_APPLICATION_CREDENTIALS="%s"`, p), nil
		}
	}

	if b, ok, err := br.Resolve("StackdriverProfiler"); err != nil {
		return "", fmt.Errorf("unable to resolve binding StackdriverProfiler\n%w", err)
	} else if ok {
		if p, ok := b.SecretFilePath("ApplicationCredentials"); ok {
			return fmt.Sprintf(`export GOOGLE_APPLICATION_CREDENTIALS="%s"`, p), nil
		}
	}

	return "", nil
}
