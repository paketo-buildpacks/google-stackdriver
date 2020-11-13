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

package google

import (
	"fmt"

	"github.com/buildpacks/libcnb"

	"github.com/paketo-buildpacks/google-cloud/internal/common"
	"github.com/paketo-buildpacks/libpak"
)

const (
	DebuggerEnabled = "BP_GOOGLE_CLOUD_DEBUGGER_ENABLED"
	ProfilerEnabled = "BP_GOOGLE_CLOUD_PROFILER_ENABLED"
)

type Detect struct{}

func (d Detect) Detect(context libcnb.DetectContext) (libcnb.DetectResult, error) {
	cr, err := libpak.NewConfigurationResolver(context.Buildpack, nil)
	if err != nil {
		return libcnb.DetectResult{}, fmt.Errorf("unable to create configuration resolver\n%w", err)
	}

	debugger := cr.ResolveBool(DebuggerEnabled)
	profiler := cr.ResolveBool(ProfilerEnabled)

	if debugger && profiler {
		return libcnb.DetectResult{
			Pass: true,
			Plans: []libcnb.BuildPlan{
				{
					Provides: []libcnb.BuildPlanProvide{
						{Name: common.Credentials},
						{Name: common.DebuggerJava},
						{Name: common.ProfilerJava},
						{Name: common.DebuggerNodeJS},
						{Name: common.ProfilerNodeJS},
					},
					Requires: []libcnb.BuildPlanRequire{
						{Name: common.Credentials},
						{Name: common.DebuggerJava},
						{Name: common.ProfilerJava},
						{Name: "jvm-application"},
						{Name: common.DebuggerNodeJS},
						{Name: common.ProfilerNodeJS},
						{Name: "node", Metadata: map[string]interface{}{"build": true}},
					},
				},
				{
					Provides: []libcnb.BuildPlanProvide{
						{Name: common.Credentials},
						{Name: common.DebuggerJava},
						{Name: common.ProfilerNodeJS},
					},
					Requires: []libcnb.BuildPlanRequire{
						{Name: common.Credentials},
						{Name: common.DebuggerJava},
						{Name: "jvm-application"},
						{Name: common.ProfilerNodeJS},
						{Name: "node", Metadata: map[string]interface{}{"build": true}},
					},
				},
				{
					Provides: []libcnb.BuildPlanProvide{
						{Name: common.Credentials},
						{Name: common.ProfilerJava},
						{Name: common.DebuggerNodeJS},
					},
					Requires: []libcnb.BuildPlanRequire{
						{Name: common.Credentials},
						{Name: common.ProfilerJava},
						{Name: "jvm-application"},
						{Name: common.DebuggerNodeJS},
						{Name: "node", Metadata: map[string]interface{}{"build": true}},
					},
				},
				{
					Provides: []libcnb.BuildPlanProvide{
						{Name: common.Credentials},
						{Name: common.DebuggerJava},
						{Name: common.ProfilerJava},
					},
					Requires: []libcnb.BuildPlanRequire{
						{Name: common.Credentials},
						{Name: common.DebuggerJava},
						{Name: common.ProfilerJava},
						{Name: "jvm-application"},
					},
				},
				{
					Provides: []libcnb.BuildPlanProvide{
						{Name: common.Credentials},
						{Name: common.DebuggerNodeJS},
						{Name: common.ProfilerNodeJS},
					},
					Requires: []libcnb.BuildPlanRequire{
						{Name: common.Credentials},
						{Name: common.DebuggerNodeJS},
						{Name: common.ProfilerNodeJS},
						{Name: "node", Metadata: map[string]interface{}{"build": true}},
					},
				},
			},
		}, nil
	}

	if debugger {
		return libcnb.DetectResult{
			Pass: true,
			Plans: []libcnb.BuildPlan{
				{
					Provides: []libcnb.BuildPlanProvide{
						{Name: common.Credentials},
						{Name: common.DebuggerJava},
						{Name: common.DebuggerNodeJS},
					},
					Requires: []libcnb.BuildPlanRequire{
						{Name: common.Credentials},
						{Name: common.DebuggerJava},
						{Name: "jvm-application"},
						{Name: common.DebuggerNodeJS},
						{Name: "node", Metadata: map[string]interface{}{"build": true}},
					},
				},
				{
					Provides: []libcnb.BuildPlanProvide{
						{Name: common.Credentials},
						{Name: common.DebuggerJava},
					},
					Requires: []libcnb.BuildPlanRequire{
						{Name: common.Credentials},
						{Name: common.DebuggerJava},
						{Name: "jvm-application"},
					},
				},
				{
					Provides: []libcnb.BuildPlanProvide{
						{Name: common.Credentials},
						{Name: common.DebuggerNodeJS},
					},
					Requires: []libcnb.BuildPlanRequire{
						{Name: common.Credentials},
						{Name: common.DebuggerNodeJS},
						{Name: "node", Metadata: map[string]interface{}{"build": true}},
					},
				},
			},
		}, nil
	}

	if profiler {
		return libcnb.DetectResult{
			Pass: true,
			Plans: []libcnb.BuildPlan{
				{
					Provides: []libcnb.BuildPlanProvide{
						{Name: common.Credentials},
						{Name: common.ProfilerJava},
						{Name: common.ProfilerNodeJS},
					},
					Requires: []libcnb.BuildPlanRequire{
						{Name: common.Credentials},
						{Name: common.ProfilerJava},
						{Name: "jvm-application"},
						{Name: common.ProfilerNodeJS},
						{Name: "node", Metadata: map[string]interface{}{"build": true}},
					},
				},
				{
					Provides: []libcnb.BuildPlanProvide{
						{Name: common.Credentials},
						{Name: common.ProfilerJava},
					},
					Requires: []libcnb.BuildPlanRequire{
						{Name: common.Credentials},
						{Name: common.ProfilerJava},
						{Name: "jvm-application"},
					},
				},
				{
					Provides: []libcnb.BuildPlanProvide{
						{Name: common.Credentials},
						{Name: common.ProfilerNodeJS},
					},
					Requires: []libcnb.BuildPlanRequire{
						{Name: common.Credentials},
						{Name: common.ProfilerNodeJS},
						{Name: "node", Metadata: map[string]interface{}{"build": true}},
					},
				},
			},
		}, nil
	}

	return libcnb.DetectResult{Pass: false}, nil
}
