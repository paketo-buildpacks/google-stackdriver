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
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/buildpacks/libcnb"
	. "github.com/onsi/gomega"
	"github.com/paketo-buildpacks/google-stackdriver/stackdriver"
	"github.com/paketo-buildpacks/libpak"
	"github.com/paketo-buildpacks/libpak/effect"
	"github.com/paketo-buildpacks/libpak/effect/mocks"
	"github.com/sclevine/spec"
	"github.com/stretchr/testify/mock"
)

func testNodeJSDebuggerAgent(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect = NewWithT(t).Expect

		ctx      libcnb.BuildContext
		executor *mocks.Executor
	)

	it.Before(func() {
		var err error

		ctx.Application.Path, err = ioutil.TempDir("", "nodejs-debugger-agent-application")
		Expect(err).NotTo(HaveOccurred())

		ctx.Layers.Path, err = ioutil.TempDir("", "nodejs-debugger-agent-layers")
		Expect(err).NotTo(HaveOccurred())

		executor = &mocks.Executor{}
		executor.On("Execute", mock.Anything).Return(nil)
	})

	it.After(func() {
		Expect(os.RemoveAll(ctx.Application.Path)).To(Succeed())
		Expect(os.RemoveAll(ctx.Layers.Path)).To(Succeed())
	})

	it("contributes NodeJS agent", func() {
		Expect(ioutil.WriteFile(filepath.Join(ctx.Application.Path, "package.json"), []byte(`{ "main": "main.js" }`),
			0644)).To(Succeed())
		Expect(ioutil.WriteFile(filepath.Join(ctx.Application.Path, "main.js"), []byte{}, 0644)).To(Succeed())

		dep := libpak.BuildpackDependency{
			URI:    "https://localhost/stub-stackdriver-debugger-agent.tgz",
			SHA256: "c3ecfa1e2daa29db419b063dec9ea20108923e406d9ab7a35318f6f14f615dc6",
		}
		dc := libpak.DependencyCache{CachePath: "testdata"}

		n := stackdriver.NewNodeJSDebuggerAgent(ctx.Application.Path, dep, dc, &libcnb.BuildpackPlan{})
		n.Executor = executor
		layer, err := ctx.Layers.Layer("test-layer")
		Expect(err).NotTo(HaveOccurred())

		layer, err = n.Contribute(layer)
		Expect(err).NotTo(HaveOccurred())

		Expect(layer.Launch).To(BeTrue())

		execution := executor.Calls[0].Arguments[0].(effect.Execution)
		Expect(execution.Command).To(Equal("npm"))
		Expect(execution.Args).To(Equal([]string{"install", "--no-save",
			filepath.Join("testdata",
				"c3ecfa1e2daa29db419b063dec9ea20108923e406d9ab7a35318f6f14f615dc6",
				"stub-stackdriver-debugger-agent.tgz"),
		}))

		Expect(layer.LaunchEnvironment["NODE_PATH"]).To(Equal(filepath.Join(layer.Path, "node_modules")))
	})

	it("requires @google-cloud/debug-agent module", func() {
		Expect(ioutil.WriteFile(filepath.Join(ctx.Application.Path, "package.json"), []byte(`{ "main": "main.js" }`),
			0644)).To(Succeed())
		Expect(ioutil.WriteFile(filepath.Join(ctx.Application.Path, "main.js"), []byte("test"), 0644)).To(Succeed())

		dep := libpak.BuildpackDependency{
			URI:    "https://localhost/stub-stackdriver-debugger-agent.tgz",
			SHA256: "c3ecfa1e2daa29db419b063dec9ea20108923e406d9ab7a35318f6f14f615dc6",
		}
		dc := libpak.DependencyCache{CachePath: "testdata"}

		n := stackdriver.NewNodeJSDebuggerAgent(ctx.Application.Path, dep, dc, &libcnb.BuildpackPlan{})
		n.Executor = executor
		layer, err := ctx.Layers.Layer("test-layer")
		Expect(err).NotTo(HaveOccurred())

		layer, err = n.Contribute(layer)
		Expect(err).NotTo(HaveOccurred())

		Expect(ioutil.ReadFile(filepath.Join(ctx.Application.Path, "main.js"))).To(Equal(
			[]byte("require('@google-cloud/debug-agent').start();\ntest")))
	})

	it("does not require @google-cloud/debug-agent module", func() {
		Expect(ioutil.WriteFile(filepath.Join(ctx.Application.Path, "package.json"), []byte(`{ "main": "main.js" }`),
			0644)).To(Succeed())
		Expect(ioutil.WriteFile(filepath.Join(ctx.Application.Path, "main.js"),
			[]byte("test\nrequire('@google-cloud/debug-agent')\ntest"), 0644)).To(Succeed())

		dep := libpak.BuildpackDependency{
			URI:    "https://localhost/stub-stackdriver-debugger-agent.tgz",
			SHA256: "c3ecfa1e2daa29db419b063dec9ea20108923e406d9ab7a35318f6f14f615dc6",
		}
		dc := libpak.DependencyCache{CachePath: "testdata"}

		n := stackdriver.NewNodeJSDebuggerAgent(ctx.Application.Path, dep, dc, &libcnb.BuildpackPlan{})
		n.Executor = executor
		layer, err := ctx.Layers.Layer("test-layer")
		Expect(err).NotTo(HaveOccurred())

		layer, err = n.Contribute(layer)
		Expect(err).NotTo(HaveOccurred())

		Expect(ioutil.ReadFile(filepath.Join(ctx.Application.Path, "main.js"))).To(Equal(
			[]byte("test\nrequire('@google-cloud/debug-agent')\ntest")))
	})
}
