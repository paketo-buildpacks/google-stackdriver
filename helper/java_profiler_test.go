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

package helper_test

import (
	"testing"

	. "github.com/onsi/gomega"
	"github.com/sclevine/spec"

	"github.com/paketo-buildpacks/google-stackdriver/v5/helper"
)

func testJavaProfiler(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect = NewWithT(t).Expect

		j helper.JavaProfiler
	)

	it("fails without $BPI_GOOGLE_STACKDRIVER_PROFILER_JAVA_AGENT_PATH", func() {
		_, err := j.Execute()
		Expect(err).To(MatchError("$BPI_GOOGLE_STACKDRIVER_PROFILER_JAVA_AGENT_PATH must be set"))
	})

	context("$BPI_GOOGLE_STACKDRIVER_PROFILER_JAVA_AGENT_PATH", func() {

		it.Before(func() {
			t.Setenv("BPI_GOOGLE_STACKDRIVER_PROFILER_JAVA_AGENT_PATH", "test-path")
		})

		it("uses defaults", func() {
			Expect(j.Execute()).To(Equal(map[string]string{
				"JAVA_TOOL_OPTIONS": "-agentpath:test-path=-logtostderr=1,-cprof_project_id=,-cprof_service=default-module",
			}))
		})

		context("$BPL_GOOGLE_STACKDRIVER_MODULE", func() {
			it.Before(func() {
				t.Setenv("BPL_GOOGLE_STACKDRIVER_MODULE", "test-module")
			})

			it("uses configured module", func() {
				Expect(j.Execute()).To(Equal(map[string]string{
					"JAVA_TOOL_OPTIONS": "-agentpath:test-path=-logtostderr=1,-cprof_project_id=,-cprof_service=test-module",
				}))
			})
		})

		context("$BPL_GOOGLE_STACKDRIVER_VERSION", func() {
			it.Before(func() {
				t.Setenv("BPL_GOOGLE_STACKDRIVER_VERSION", "test-version")
			})

			it("uses configured version", func() {
				Expect(j.Execute()).To(Equal(map[string]string{
					"JAVA_TOOL_OPTIONS": "-agentpath:test-path=-logtostderr=1,-cprof_project_id=,-cprof_service=default-module,-cprof_service_version=test-version",
				}))
			})
		})

		context("$JAVA_TOOL_OPTIONS", func() {
			it.Before(func() {
				t.Setenv("JAVA_TOOL_OPTIONS", "test-java-tool-options")
			})

			it("uses configured JAVA_TOOL_OPTIONS", func() {
				Expect(j.Execute()).To(Equal(map[string]string{
					"JAVA_TOOL_OPTIONS": "test-java-tool-options -agentpath:test-path=-logtostderr=1,-cprof_project_id=,-cprof_service=default-module",
				}))
			})
		})
	})
}
