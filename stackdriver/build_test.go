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

package stackdriver_test

import (
	"testing"

	"github.com/buildpacks/libcnb"
	. "github.com/onsi/gomega"
	"github.com/sclevine/spec"

	"github.com/paketo-buildpacks/google-stackdriver/stackdriver"
)

func testBuild(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect = NewWithT(t).Expect

		ctx libcnb.BuildContext
	)

	it("contributes Java debugger agent", func() {
		ctx.Plan.Entries = append(ctx.Plan.Entries, libcnb.BuildpackPlanEntry{Name: "google-stackdriver-debugger-java"})
		ctx.Buildpack.Metadata = map[string]interface{}{
			"dependencies": []map[string]interface{}{
				{
					"id":      "google-stackdriver-debugger-java",
					"version": "1.1.1",
					"stacks":  []interface{}{"test-stack-id"},
				},
			},
		}
		ctx.StackID = "test-stack-id"

		result, err := stackdriver.Build{}.Build(ctx)
		Expect(err).NotTo(HaveOccurred())

		Expect(result.Layers).To(HaveLen(2))
		Expect(result.Layers[0].Name()).To(Equal("google-stackdriver-debugger-java"))
		Expect(result.Layers[1].Name()).To(Equal("google-application-credentials"))
	})

	it("contributes NodeJS debugger agent", func() {
		ctx.Plan.Entries = append(ctx.Plan.Entries, libcnb.BuildpackPlanEntry{Name: "google-stackdriver-debugger-nodejs"})
		ctx.Buildpack.Metadata = map[string]interface{}{
			"dependencies": []map[string]interface{}{
				{
					"id":      "google-stackdriver-debugger-nodejs",
					"version": "1.1.1",
					"stacks":  []interface{}{"test-stack-id"},
				},
			},
		}
		ctx.StackID = "test-stack-id"

		result, err := stackdriver.Build{}.Build(ctx)
		Expect(err).NotTo(HaveOccurred())

		Expect(result.Layers).To(HaveLen(2))
		Expect(result.Layers[0].Name()).To(Equal("google-stackdriver-debugger-nodejs"))
		Expect(result.Layers[1].Name()).To(Equal("google-application-credentials"))
	})

	it("contributes Java profiler agent", func() {
		ctx.Plan.Entries = append(ctx.Plan.Entries, libcnb.BuildpackPlanEntry{Name: "google-stackdriver-profiler-java"})
		ctx.Buildpack.Metadata = map[string]interface{}{
			"dependencies": []map[string]interface{}{
				{
					"id":      "google-stackdriver-profiler-java",
					"version": "1.1.1",
					"stacks":  []interface{}{"test-stack-id"},
				},
			},
		}
		ctx.StackID = "test-stack-id"

		result, err := stackdriver.Build{}.Build(ctx)
		Expect(err).NotTo(HaveOccurred())

		Expect(result.Layers).To(HaveLen(2))
		Expect(result.Layers[0].Name()).To(Equal("google-stackdriver-profiler-java"))
		Expect(result.Layers[1].Name()).To(Equal("google-application-credentials"))
	})

	it("contributes NodeJS profiler agent", func() {
		ctx.Plan.Entries = append(ctx.Plan.Entries, libcnb.BuildpackPlanEntry{Name: "google-stackdriver-profiler-nodejs"})
		ctx.Buildpack.Metadata = map[string]interface{}{
			"dependencies": []map[string]interface{}{
				{
					"id":      "google-stackdriver-profiler-nodejs",
					"version": "1.1.1",
					"stacks":  []interface{}{"test-stack-id"},
				},
			},
		}
		ctx.StackID = "test-stack-id"

		result, err := stackdriver.Build{}.Build(ctx)
		Expect(err).NotTo(HaveOccurred())

		Expect(result.Layers).To(HaveLen(2))
		Expect(result.Layers[0].Name()).To(Equal("google-stackdriver-profiler-nodejs"))
		Expect(result.Layers[1].Name()).To(Equal("google-application-credentials"))
	})
}
