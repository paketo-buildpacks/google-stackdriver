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
	"os"
	"path/filepath"
	"testing"

	"github.com/buildpacks/libcnb"
	. "github.com/onsi/gomega"
	"github.com/paketo-buildpacks/libpak"
	"github.com/paketo-buildpacks/libpak/effect"
	"github.com/paketo-buildpacks/libpak/effect/mocks"
	"github.com/sclevine/spec"
	"github.com/stretchr/testify/mock"

	"github.com/paketo-buildpacks/google-stackdriver/v5/stackdriver"
)

func testNodeJSProfilerAgent(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect = NewWithT(t).Expect

		ctx      libcnb.BuildContext
		executor *mocks.Executor
	)

	it.Before(func() {
		var err error

		ctx.Application.Path = t.TempDir()
		Expect(err).NotTo(HaveOccurred())

		ctx.Layers.Path = t.TempDir()
		Expect(err).NotTo(HaveOccurred())

		executor = &mocks.Executor{}
		executor.On("Execute", mock.Anything).Return(nil)
	})

	it("contributes NodeJS agent", func() {
		Expect(os.WriteFile(filepath.Join(ctx.Application.Path, "package.json"), []byte(`{ "main": "main.js" }`),
			0644)).To(Succeed())
		Expect(os.WriteFile(filepath.Join(ctx.Application.Path, "main.js"), []byte{}, 0644)).To(Succeed())

		dep := libpak.BuildpackDependency{
			URI:    "https://localhost/stub-stackdriver-profiler-agent.tgz",
			SHA256: "eee0eca3815f2d4aeaa7e23c1150878ee42864d26ec7331d800b34e667714138",
		}
		dc := libpak.DependencyCache{CachePath: "testdata"}

		n, be := stackdriver.NewNodeJSProfilerAgent(ctx.Application.Path, dep, dc)
		Expect(be.Launch).To(BeTrue())
		Expect(be.Metadata["uri"]).To(Equal("https://localhost/stub-stackdriver-profiler-agent.tgz"))

		n.Executor = executor
		layer, err := ctx.Layers.Layer("test-layer")
		Expect(err).NotTo(HaveOccurred())

		layer, err = n.Contribute(layer)
		Expect(err).NotTo(HaveOccurred())

		Expect(layer.LayerTypes.Launch).To(BeTrue())

		execution := executor.Calls[0].Arguments[0].(effect.Execution)
		Expect(execution.Command).To(Equal("npm"))
		Expect(execution.Args).To(Equal([]string{"install", "--no-save",
			filepath.Join("testdata",
				"eee0eca3815f2d4aeaa7e23c1150878ee42864d26ec7331d800b34e667714138",
				"stub-stackdriver-profiler-agent.tgz"),
		}))

		Expect(layer.LaunchEnvironment["NODE_PATH.delim"]).To(Equal(string(os.PathListSeparator)))
		Expect(layer.LaunchEnvironment["NODE_PATH.prepend"]).To(Equal(filepath.Join(layer.Path, "node_modules")))
	})

	it("requires @google-cloud/profiler module", func() {
		Expect(os.WriteFile(filepath.Join(ctx.Application.Path, "package.json"), []byte(`{ "main": "main.js" }`),
			0644)).To(Succeed())
		Expect(os.WriteFile(filepath.Join(ctx.Application.Path, "main.js"), []byte("test"), 0644)).To(Succeed())

		dep := libpak.BuildpackDependency{
			URI:    "https://localhost/stub-stackdriver-profiler-agent.tgz",
			SHA256: "eee0eca3815f2d4aeaa7e23c1150878ee42864d26ec7331d800b34e667714138",
		}
		dc := libpak.DependencyCache{CachePath: "testdata"}

		n, _ := stackdriver.NewNodeJSProfilerAgent(ctx.Application.Path, dep, dc)

		n.Executor = executor
		layer, err := ctx.Layers.Layer("test-layer")
		Expect(err).NotTo(HaveOccurred())

		_, err = n.Contribute(layer)
		Expect(err).NotTo(HaveOccurred())

		Expect(os.ReadFile(filepath.Join(ctx.Application.Path, "main.js"))).To(Equal(
			[]byte("require('@google-cloud/profiler').start();\ntest")))
	})

	it("does not require @google-cloud/profiler module", func() {
		Expect(os.WriteFile(filepath.Join(ctx.Application.Path, "package.json"), []byte(`{ "main": "main.js" }`),
			0644)).To(Succeed())
		Expect(os.WriteFile(filepath.Join(ctx.Application.Path, "main.js"),
			[]byte("test\nrequire('@google-cloud/profiler')\ntest"), 0644)).To(Succeed())

		dep := libpak.BuildpackDependency{
			URI:    "https://localhost/stub-stackdriver-profiler-agent.tgz",
			SHA256: "eee0eca3815f2d4aeaa7e23c1150878ee42864d26ec7331d800b34e667714138",
		}
		dc := libpak.DependencyCache{CachePath: "testdata"}

		n, _ := stackdriver.NewNodeJSProfilerAgent(ctx.Application.Path, dep, dc)
		n.Executor = executor
		layer, err := ctx.Layers.Layer("test-layer")
		Expect(err).NotTo(HaveOccurred())

		_, err = n.Contribute(layer)
		Expect(err).NotTo(HaveOccurred())

		Expect(os.ReadFile(filepath.Join(ctx.Application.Path, "main.js"))).To(Equal(
			[]byte("test\nrequire('@google-cloud/profiler')\ntest")))
	})
}
