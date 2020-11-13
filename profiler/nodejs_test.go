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

package profiler_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/buildpacks/libcnb"
	. "github.com/onsi/gomega"
	"github.com/sclevine/spec"
	"github.com/stretchr/testify/mock"

	"github.com/paketo-buildpacks/google-cloud/internal/common"
	"github.com/paketo-buildpacks/google-cloud/profiler"
	"github.com/paketo-buildpacks/libpak"
	"github.com/paketo-buildpacks/libpak/effect"
	"github.com/paketo-buildpacks/libpak/effect/mocks"
)

func testNodeJS(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect = NewWithT(t).Expect
	)

	context("Build", func() {
		var (
			ctx      libcnb.BuildContext
			executor *mocks.Executor
		)

		it.Before(func() {
			var err error
			ctx.Layers.Path, err = ioutil.TempDir("", "profiler-nodejs-build-layers")
			Expect(err).NotTo(HaveOccurred())

			executor = &mocks.Executor{}
			executor.On("Execute", mock.Anything).Return(nil)
		})

		it.After(func() {
			Expect(os.RemoveAll(ctx.Layers.Path)).To(Succeed())
		})

		it("contributes NodeJS agent", func() {
			dep := libpak.BuildpackDependency{
				URI:    "https://localhost/stub-cloud-profiler-agent.tgz",
				SHA256: "eee0eca3815f2d4aeaa7e23c1150878ee42864d26ec7331d800b34e667714138",
			}
			dc := libpak.DependencyCache{CachePath: "testdata"}

			n := profiler.NewNodeJSBuild(dep, dc, &libcnb.BuildpackPlan{})
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
					"eee0eca3815f2d4aeaa7e23c1150878ee42864d26ec7331d800b34e667714138",
					"stub-cloud-profiler-agent.tgz"),
			}))

			Expect(layer.LaunchEnvironment[fmt.Sprintf("%s.default", profiler.NodePath)]).
				To(Equal(filepath.Join(layer.Path, "node_modules")))
		})
	})

	context("Launch", func() {
		var (
			l = profiler.NodeJSLaunch{
				CredentialSource: common.MetadataServer,
			}
		)

		it.Before(func() {
			var err error
			l.ApplicationPath, err = ioutil.TempDir("", "profiler-nodejs-launch-application")
			Expect(err).NotTo(HaveOccurred())

			Expect(ioutil.WriteFile(filepath.Join(l.ApplicationPath, "package.json"), []byte(`{ "main": "main.js" }`), 0644)).
				To(Succeed())
			Expect(ioutil.WriteFile(filepath.Join(l.ApplicationPath, "main.js"), []byte("test"), 0644)).
				To(Succeed())

			Expect(os.Setenv(profiler.NodePath, "test-path")).To(Succeed())
			Expect(os.Setenv(common.Module, "test-module")).To(Succeed())
			Expect(os.Setenv(common.Version, "test-version")).To(Succeed())
		})

		it.After(func() {
			Expect(os.RemoveAll(l.ApplicationPath)).To(Succeed())

			Expect(os.Unsetenv(profiler.NodePath)).To(Succeed())
			Expect(os.Unsetenv(common.Module)).To(Succeed())
			Expect(os.Unsetenv(common.Version)).To(Succeed())
		})

		it("does not contribute if source is None", func() {
			l.CredentialSource = common.None

			Expect(l.Execute()).To(BeNil())
		})

		it("returns error if BPI_GOOGLE_CLOUD_PROFILER_NODE_PATH is not set", func() {
			Expect(os.Unsetenv(profiler.NodePath)).To(Succeed())

			_, err := l.Execute()
			Expect(err).To(MatchError("$BPI_GOOGLE_CLOUD_PROFILER_NODE_PATH must be set"))
		})

		it("returns error if BPL_GOOGLE_CLOUD_MODULE is not set", func() {
			Expect(os.Unsetenv(common.Module)).To(Succeed())

			_, err := l.Execute()
			Expect(err).To(MatchError("$BPL_GOOGLE_CLOUD_MODULE must be set"))
		})

		it("returns error if BPL_GOOGLE_CLOUD_VERSION is not set", func() {
			Expect(os.Unsetenv(common.Version)).To(Succeed())

			_, err := l.Execute()
			Expect(err).To(MatchError("$BPL_GOOGLE_CLOUD_VERSION must be set"))
		})

		context("binding", func() {

			it.Before(func() {
				l.CredentialSource = common.Binding
				Expect(os.Setenv(common.ProjectID, "test-project-id")).To(Succeed())
			})

			it.After(func() {
				Expect(os.Unsetenv(common.ProjectID)).To(Succeed())
			})

			it("returns error if BPL_GOOGLE_CLOUD_PROJECT_ID is not set", func() {
				Expect(os.Unsetenv(common.ProjectID)).To(Succeed())

				_, err := l.Execute()
				Expect(err).To(MatchError("$BPL_GOOGLE_CLOUD_PROJECT_ID must be set"))
			})

			context("@google-cloud/profiler already required", func() {

				it.Before(func() {
					Expect(ioutil.WriteFile(filepath.Join(l.ApplicationPath, "main.js"), []byte(`test
require('@google-cloud/profiler').start()
test`), 0644)).
						To(Succeed())
				})

				it("does not contribute require('@google-cloud/profiler')", func() {
					_, err := l.Execute()
					Expect(err).NotTo(HaveOccurred())

					b, err := ioutil.ReadFile(filepath.Join(l.ApplicationPath, "main.js"))
					Expect(err).NotTo(HaveOccurred())

					Expect(string(b)).To(Equal(`test
require('@google-cloud/profiler').start()
test`))
				})
			})

			context("@google-cloud/profiler not already required", func() {

				it("contributes require('@google-cloud/profiler')", func() {
					_, err := l.Execute()
					Expect(err).NotTo(HaveOccurred())

					b, err := ioutil.ReadFile(filepath.Join(l.ApplicationPath, "main.js"))
					Expect(err).NotTo(HaveOccurred())

					Expect(string(b)).To(Equal(`require('@google-cloud/profiler').start({
  projectId: 'test-project-id',
  serviceContext: {
    service: 'test-module',
    version: 'test-version',
  },
});
test`))
				})
			})

			it("contributes NODE_PATH", func() {
				Expect(l.Execute()).To(Equal(map[string]string{
					"NODE_PATH": "test-path",
				}))
			})
		})

		context("metadata server", func() {

			context("@google-cloud/profiler already required", func() {

				it.Before(func() {
					Expect(ioutil.WriteFile(filepath.Join(l.ApplicationPath, "main.js"), []byte(`test
require('@google-cloud/profiler').start()
test`), 0644)).
						To(Succeed())
				})

				it("does not contribute require('@google-cloud/profiler')", func() {
					_, err := l.Execute()
					Expect(err).NotTo(HaveOccurred())

					b, err := ioutil.ReadFile(filepath.Join(l.ApplicationPath, "main.js"))
					Expect(err).NotTo(HaveOccurred())

					Expect(string(b)).To(Equal(`test
require('@google-cloud/profiler').start()
test`))
				})
			})

			context("@google-cloud/profiler not already required", func() {

				it("contributes require('@google-cloud/profiler')", func() {
					_, err := l.Execute()
					Expect(err).NotTo(HaveOccurred())

					b, err := ioutil.ReadFile(filepath.Join(l.ApplicationPath, "main.js"))
					Expect(err).NotTo(HaveOccurred())

					Expect(string(b)).To(Equal(`require('@google-cloud/profiler').start({
  serviceContext: {
    service: 'test-module',
    version: 'test-version',
  },
});
test`))
				})
			})

			it("contributes NODE_PATH", func() {
				Expect(l.Execute()).To(Equal(map[string]string{
					"NODE_PATH": "test-path",
				}))
			})
		})

		context("existing $NODE_PATH", func() {

			it.Before(func() {
				Expect(os.Setenv("NODE_PATH", "test-node-path")).To(Succeed())
			})

			it.After(func() {
				Expect(os.Unsetenv("NODE_PATH")).To(Succeed())
			})

			it("contributes NODE_PATH", func() {
				Expect(l.Execute()).To(Equal(map[string]string{
					"NODE_PATH": strings.Join([]string{
						"test-node-path",
						"test-path",
					}, string(os.PathListSeparator)),
				}))
			})
		})
	})
}
