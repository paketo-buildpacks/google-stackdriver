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

package stackdriver

import (
	"fmt"

	"github.com/buildpacks/libcnb"
	"github.com/paketo-buildpacks/libpak/bindings"
	"github.com/paketo-buildpacks/libpak/bard"
)

type Detect struct{
	Logger bard.Logger
}

func (d Detect) Detect(context libcnb.DetectContext) (libcnb.DetectResult, error) {
	result := libcnb.DetectResult{Pass: false}

	if _, ok, err := bindings.ResolveOne(context.Platform.Bindings, bindings.OfType("StackdriverDebugger")); err != nil {
		return libcnb.DetectResult{}, fmt.Errorf("unable to resolve binding StackdriverDebugger\n%w", err)
	} else if ok {
		result.Pass = true
		result.Plans = append(result.Plans,
			libcnb.BuildPlan{
				Provides: []libcnb.BuildPlanProvide{
					{Name: "google-stackdriver-debugger-java"},
				},
				Requires: []libcnb.BuildPlanRequire{
					{Name: "google-stackdriver-debugger-java"},
					{Name: "jvm-application"},
				},
			},
			libcnb.BuildPlan{
				Provides: []libcnb.BuildPlanProvide{
					{Name: "google-stackdriver-debugger-nodejs"},
				},
				Requires: []libcnb.BuildPlanRequire{
					{Name: "google-stackdriver-debugger-nodejs"},
					{Name: "node", Metadata: map[string]interface{}{"build": true}},
					{Name: "node_modules"},
				},
			},
		)
		d.Logger.Info("PASSED: binding of type 'StackdriverDebugger' found")
	}

	if _, ok, err := bindings.ResolveOne(context.Platform.Bindings, bindings.OfType("StackdriverProfiler")); err != nil {
		return libcnb.DetectResult{}, fmt.Errorf("unable to resolve binding StackdriverProfiler\n%w", err)
	} else if ok {
		result.Pass = true
		result.Plans = append(result.Plans,
			libcnb.BuildPlan{
				Provides: []libcnb.BuildPlanProvide{
					{Name: "google-stackdriver-profiler-java"},
				},
				Requires: []libcnb.BuildPlanRequire{
					{Name: "google-stackdriver-profiler-java"},
					{Name: "jvm-application"},
				},
			},
			libcnb.BuildPlan{
				Provides: []libcnb.BuildPlanProvide{
					{Name: "google-stackdriver-profiler-nodejs"},
				},
				Requires: []libcnb.BuildPlanRequire{
					{Name: "google-stackdriver-profiler-nodejs"},
					{Name: "node", Metadata: map[string]interface{}{"build": true}},
					{Name: "node_modules"},
				},
			},
		)
		d.Logger.Info("PASSED: binding of type 'StackdriverProfiler' found")
	}

	if result.Pass != true {
		d.Logger.Info("SKIPPED: no bindings of type 'StackdriverDebugger' or 'StackdriverProfiler' found")
	}
	return result, nil
}
