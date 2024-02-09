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

package helper

import (
	"fmt"

	"github.com/buildpacks/libcnb"
	"github.com/paketo-buildpacks/libpak/bard"
	"github.com/paketo-buildpacks/libpak/bindings"
)

type Credentials struct {
	Bindings libcnb.Bindings
	Logger   bard.Logger
}

func (c Credentials) Execute() (map[string]string, error) {
	if b, ok, err := bindings.ResolveOne(c.Bindings, bindings.OfType("StackdriverProfiler")); err != nil {
		return nil, fmt.Errorf("unable to resolve binding StackdriverProfiler\n%w", err)
	} else if ok {
		if p, ok := b.SecretFilePath("ApplicationCredentials"); ok {
			c.Logger.Info("Configuring Google application credentials")
			return map[string]string{"GOOGLE_APPLICATION_CREDENTIALS": p}, nil
		}
	}

	return nil, nil
}
