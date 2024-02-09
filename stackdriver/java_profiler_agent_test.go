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
	"path/filepath"
	"testing"

	"github.com/buildpacks/libcnb"
	. "github.com/onsi/gomega"
	"github.com/paketo-buildpacks/libpak"
	"github.com/sclevine/spec"

	"github.com/paketo-buildpacks/google-stackdriver/v5/stackdriver"
)

func testJavaProfilerAgent(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect = NewWithT(t).Expect

		ctx libcnb.BuildContext
	)

	it.Before(func() {
		var err error

		ctx.Layers.Path = t.TempDir()
		Expect(err).NotTo(HaveOccurred())
	})

	it("contributes Java agent", func() {
		dep := libpak.BuildpackDependency{
			URI:    "https://localhost/stub-stackdriver-profiler-agent.tar.gz",
			SHA256: "a27bbf74fa913fe70273c38831bc1043a2da0e28766423b81d8e9a042b353797",
		}
		dc := libpak.DependencyCache{CachePath: "testdata"}

		layer, err := ctx.Layers.Layer("test-layer")
		Expect(err).NotTo(HaveOccurred())

		contrib, be := stackdriver.NewJavaProfilerAgent(dep, dc)
		Expect(be.Launch).To(BeTrue())
		Expect(be.Metadata["uri"]).To(Equal("https://localhost/stub-stackdriver-profiler-agent.tar.gz"))

		layer, err = contrib.Contribute(layer)
		Expect(err).NotTo(HaveOccurred())

		Expect(layer.LayerTypes.Launch).To(BeTrue())

		file := filepath.Join(layer.Path, "profiler_java_agent.so")
		Expect(file).To(BeARegularFile())
		Expect(layer.LaunchEnvironment["BPI_GOOGLE_STACKDRIVER_PROFILER_JAVA_AGENT_PATH.default"]).To(Equal(file))
	})
}
