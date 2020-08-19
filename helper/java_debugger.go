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
	"os"
	"strings"

	"github.com/paketo-buildpacks/libpak/bard"
)

type JavaDebugger struct {
	Logger bard.Logger
}

func (j JavaDebugger) Execute() (map[string]string, error) {
	module, ok := os.LookupEnv("BPL_GOOGLE_STACKDRIVER_MODULE")
	if !ok {
		module = "default-module"
	}

	version := os.Getenv("BPL_GOOGLE_STACKDRIVER_VERSION")

	var values []string
	if s, ok := os.LookupEnv("JAVA_TOOL_OPTIONS"); ok {
		values = append(values, s)
	}

	values = append(values,
		"-Dcom.google.cdbg.auth.serviceaccount.enable=true",
		fmt.Sprintf("-Dcom.google.cdbg.module=%s", module),
	)

	if version != "" {
		values = append(values, fmt.Sprintf("-Dcom.google.cdbg.version=%s", version))
	}

	message := fmt.Sprintf("Google Stackdriver Debugger enabled for %s", module)
	if version != "" {
		message = fmt.Sprintf("%s:%s", message, version)
	}
	j.Logger.Info(message)

	return map[string]string{"JAVA_TOOL_OPTIONS": strings.Join(values, " ")}, nil
}
