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

package stackdriver_test

import (
	"testing"

	"github.com/buildpacks/libcnb"
	. "github.com/onsi/gomega"
	"github.com/paketo-buildpacks/libpak"
	"github.com/sclevine/spec"

	"github.com/paketo-buildpacks/google-stackdriver/v5/stackdriver"
)

func testBuild(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect = NewWithT(t).Expect

		ctx libcnb.BuildContext
	)

	it("contributes Java profiler agent", func() {
		ctx.Plan.Entries = append(ctx.Plan.Entries, libcnb.BuildpackPlanEntry{Name: "google-stackdriver-profiler-java"})
		ctx.Buildpack.Metadata = map[string]interface{}{
			"dependencies": []map[string]interface{}{
				{
					"id":      "google-stackdriver-profiler-java",
					"version": "1.1.1",
					"stacks":  []interface{}{"test-stack-id"},
					"cpes":    []string{"cpe:2.3:a:google:google-stackdriver-profiler-java:1.1.0:*:*:*:*:*:*:*"},
					"purl":    "pkg:generic/google-stackdriver-profiler-java@2021.11.1500",
				},
			},
		}
		ctx.StackID = "test-stack-id"
		ctx.Buildpack.API = "0.7"

		result, err := stackdriver.Build{}.Build(ctx)
		Expect(err).NotTo(HaveOccurred())

		Expect(result.Layers).To(HaveLen(2))
		Expect(result.Layers[0].Name()).To(Equal("google-stackdriver-profiler-java"))
		Expect(result.Layers[1].Name()).To(Equal("helper"))
		Expect(result.Layers[1].(libpak.HelperLayerContributor).Names).To(Equal([]string{"credentials", "java-profiler"}))

		Expect(len(result.BOM.Entries)).To(Equal(2))
		Expect(result.BOM.Entries[0].Name).To(Equal("google-stackdriver-profiler-java"))
		Expect(result.BOM.Entries[1].Name).To(Equal("helper"))
	})

	it("contributes NodeJS profiler agent", func() {
		ctx.Plan.Entries = append(ctx.Plan.Entries, libcnb.BuildpackPlanEntry{Name: "google-stackdriver-profiler-nodejs"})
		ctx.Buildpack.Metadata = map[string]interface{}{
			"dependencies": []map[string]interface{}{
				{
					"id":      "google-stackdriver-profiler-nodejs",
					"version": "1.1.1",
					"stacks":  []interface{}{"test-stack-id"},
					"cpes":    []string{"cpe:2.3:a:google:google-stackdriver-profiler-nodejs:1.1.0:*:*:*:*:*:*:*"},
					"purl":    "pkg:generic/google-stackdriver-profiler-nodejs@2021.11.1500",
				},
			},
		}
		ctx.StackID = "test-stack-id"
		ctx.Buildpack.API = "0.7"

		result, err := stackdriver.Build{}.Build(ctx)
		Expect(err).NotTo(HaveOccurred())

		Expect(result.Layers).To(HaveLen(2))
		Expect(result.Layers[0].Name()).To(Equal("google-stackdriver-profiler-nodejs"))
		Expect(result.Layers[1].Name()).To(Equal("helper"))
		Expect(result.Layers[1].(libpak.HelperLayerContributor).Names).To(Equal([]string{"credentials"}))

		Expect(len(result.BOM.Entries)).To(Equal(2))
		Expect(result.BOM.Entries[0].Name).To(Equal("google-stackdriver-profiler-nodejs"))
		Expect(result.BOM.Entries[1].Name).To(Equal("helper"))
	})
}
