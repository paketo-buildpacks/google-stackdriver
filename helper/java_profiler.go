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

type JavaProfiler struct {
	Logger bard.Logger
}

func (j JavaProfiler) Execute() (map[string]string, error) {
	module, ok := os.LookupEnv("BPL_GOOGLE_STACKDRIVER_MODULE")
	if !ok {
		module = "default-module"
	}

	projectId := os.Getenv("BPL_GOOGLE_STACKDRIVER_PROJECT_ID")

	version := os.Getenv("BPL_GOOGLE_STACKDRIVER_VERSION")

	agentPath, ok := os.LookupEnv("BPI_GOOGLE_STACKDRIVER_PROFILER_JAVA_AGENT_PATH")
	if !ok {
		return nil, fmt.Errorf("$BPI_GOOGLE_STACKDRIVER_PROFILER_JAVA_AGENT_PATH must be set")
	}

	var values []string
	if s, ok := os.LookupEnv("JAVA_TOOL_OPTIONS"); ok {
		values = append(values, s)
	}

	agent := fmt.Sprintf("-agentpath:%s=-logtostderr=1,-cprof_project_id=%s,-cprof_service=%s",
		agentPath, projectId, module)
	if version != "" {
		agent = fmt.Sprintf("%s,-cprof_service_version=%s", agent, version)
	}
	values = append(values, agent)

	message := fmt.Sprintf("Google Stackdriver Profiler enabled for %s", module)
	if version != "" {
		message = fmt.Sprintf("%s:%s", message, version)
	}
	j.Logger.Info(message)

	return map[string]string{"JAVA_TOOL_OPTIONS": strings.Join(values, " ")}, nil
}
