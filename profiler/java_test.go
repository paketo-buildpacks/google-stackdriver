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

	"github.com/paketo-buildpacks/google-cloud/internal/common"
	"github.com/paketo-buildpacks/google-cloud/profiler"
	"github.com/paketo-buildpacks/libpak"
)

func testJava(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect = NewWithT(t).Expect
	)

	context("Build", func() {
		var (
			ctx libcnb.BuildContext
		)

		it.Before(func() {
			var err error
			ctx.Layers.Path, err = ioutil.TempDir("", "profiler-java-build-layers")
			Expect(err).NotTo(HaveOccurred())
		})

		it.After(func() {
			Expect(os.RemoveAll(ctx.Layers.Path)).To(Succeed())
		})

		it("contributes Java agent", func() {
			dep := libpak.BuildpackDependency{
				URI:    "https://localhost/stub-cloud-profiler-agent.tar.gz",
				SHA256: "a27bbf74fa913fe70273c38831bc1043a2da0e28766423b81d8e9a042b353797",
			}
			dc := libpak.DependencyCache{CachePath: "testdata"}

			layer, err := ctx.Layers.Layer("test-layer")
			Expect(err).NotTo(HaveOccurred())

			layer, err = profiler.NewJavaBuild(dep, dc, &libcnb.BuildpackPlan{}).Contribute(layer)
			Expect(err).NotTo(HaveOccurred())

			Expect(layer.Launch).To(BeTrue())

			file := filepath.Join(layer.Path, "profiler_java_agent.so")
			Expect(file).To(BeARegularFile())
			Expect(layer.LaunchEnvironment[fmt.Sprintf("%s.default", profiler.AgentPath)]).To(Equal(file))
		})

	})


	context("Launch", func() {
		var (
			l = profiler.JavaLaunch{
				CredentialSource: common.MetadataServer,
			}
		)

		it.Before(func() {
			Expect(os.Setenv(profiler.AgentPath, "test-path")).To(Succeed())
			Expect(os.Setenv(common.Module, "test-module")).To(Succeed())
			Expect(os.Setenv(common.Version, "test-version")).To(Succeed())
		})

		it.After(func() {
			Expect(os.Unsetenv(profiler.AgentPath)).To(Succeed())
			Expect(os.Unsetenv(common.Module)).To(Succeed())
			Expect(os.Unsetenv(common.Version)).To(Succeed())
		})

		it("does not contribute if source is None", func() {
			l.CredentialSource = common.None

			Expect(l.Execute()).To(BeNil())
		})

		it("returns error if BPI_GOOGLE_CLOUD_PROFILER_AGENT_PATH is not set", func() {
			Expect(os.Unsetenv(profiler.AgentPath)).To(Succeed())

			_, err := l.Execute()
			Expect(err).To(MatchError("$BPI_GOOGLE_CLOUD_PROFILER_AGENT_PATH must be set"))
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

			it("contributes JAVA_TOOL_OPTIONS", func() {
				Expect(l.Execute()).To(Equal(map[string]string{
					"JAVA_TOOL_OPTIONS": strings.Join([]string{
						"-agentpath:test-path=-logtostderr=1",
						"-cprof_service=test-module",
						"-cprof_service_version=test-version",
						"-cprof_project_id=test-project-id",
					}, ","),
				}))

			})
		})

		context("metadata server", func() {

			it("contributes JAVA_TOOL_OPTIONS", func() {
				Expect(l.Execute()).To(Equal(map[string]string{
					"JAVA_TOOL_OPTIONS": strings.Join([]string{
						"-agentpath:test-path=-logtostderr=1",
						"-cprof_service=test-module",
						"-cprof_service_version=test-version",
					}, ","),
				}))
			})
		})

		context("existing $JAVA_TOOL_OPTIONS", func() {

			it.Before(func() {
				Expect(os.Setenv("JAVA_TOOL_OPTIONS", "test-java-tool-options")).To(Succeed())
			})

			it.After(func() {
				Expect(os.Unsetenv("JAVA_TOOL_OPTIONS")).To(Succeed())
			})

			it("contributes JAVA_TOOL_OPTIONS", func() {
				Expect(l.Execute()).To(Equal(map[string]string{
					"JAVA_TOOL_OPTIONS": strings.Join([]string{
						"test-java-tool-options",
						strings.Join([]string{
							"-agentpath:test-path=-logtostderr=1",
							"-cprof_service=test-module",
							"-cprof_service_version=test-version",
						}, ","),
					}, " "),
				}))
			})
		})
	})
}
