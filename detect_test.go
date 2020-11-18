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

package google_test

import (
	"os"
	"testing"

	"github.com/buildpacks/libcnb"
	. "github.com/onsi/gomega"
	"github.com/sclevine/spec"

	"github.com/paketo-buildpacks/google-cloud"
	"github.com/paketo-buildpacks/google-cloud/internal/common"
)

func testDetect(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect = NewWithT(t).Expect

		ctx    libcnb.DetectContext
		detect google.Detect
	)

	it("fails without any environment variables", func() {
		Expect(detect.Detect(ctx)).To(Equal(libcnb.DetectResult{}))
	})

	context("$BP_GOOGLE_CLOUD_DEBUGGER_ENABLED and $BP_GOOGLE_CLOUD_PROFILER_ENABLED", func() {
		it.Before(func() {
			Expect(os.Setenv(google.DebuggerEnabled, "true")).To(Succeed())
			Expect(os.Setenv(google.ProfilerEnabled, "true")).To(Succeed())
		})

		it.After(func() {
			Expect(os.Unsetenv(google.DebuggerEnabled)).To(Succeed())
			Expect(os.Unsetenv(google.ProfilerEnabled)).To(Succeed())
		})

		it("passes with debugger and profiler enabled", func() {
			Expect(detect.Detect(ctx)).To(Equal(libcnb.DetectResult{
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
			}))
		})
	})

	context("$BP_GOOGLE_CLOUD_DEBUGGER_ENABLED", func() {
		it.Before(func() {
			Expect(os.Setenv(google.DebuggerEnabled, "true")).To(Succeed())
		})

		it.After(func() {
			Expect(os.Unsetenv(google.DebuggerEnabled)).To(Succeed())
		})

		it("passes with debugger enabled", func() {
			Expect(detect.Detect(ctx)).To(Equal(libcnb.DetectResult{
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
			}))
		})

	})

	context("$BP_GOOGLE_CLOUD_PROFILER_ENABLED", func() {
		it.Before(func() {
			Expect(os.Setenv(google.ProfilerEnabled, "true")).To(Succeed())
		})

		it.After(func() {
			Expect(os.Unsetenv(google.ProfilerEnabled)).To(Succeed())
		})

		it("passes with profiler enabled", func() {
			Expect(detect.Detect(ctx)).To(Equal(libcnb.DetectResult{
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
			}))
		})
	})

}
