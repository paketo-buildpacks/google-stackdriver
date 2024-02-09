/*
 * Copyright 2018-2024 the original author or authors.
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

package main

import (
	"fmt"
	"os"

	"github.com/buildpacks/libcnb"
	"github.com/paketo-buildpacks/libpak/bard"
	"github.com/paketo-buildpacks/libpak/sherpa"

	"github.com/paketo-buildpacks/google-stackdriver/v5/helper"
)

func main() {
	sherpa.Execute(func() error {
		var (
			err error
			l   = bard.NewLogger(os.Stdout)
			c   = helper.Credentials{Logger: l}
			p   = helper.JavaProfiler{Logger: l}
		)

		c.Bindings, err = libcnb.NewBindingsFromEnvironment()
		if err != nil {
			return fmt.Errorf("unable to read bindings from environment\n%w", err)
		}

		return sherpa.Helpers(map[string]sherpa.ExecD{
			"credentials":   c,
			"java-profiler": p,
		})
	})
}
