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
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/buildpacks/libcnb"
	. "github.com/onsi/gomega"
	"github.com/paketo-buildpacks/libpak"
	"github.com/sclevine/spec"

	"github.com/paketo-buildpacks/google-stackdriver/stackdriver"
)

func testJavaProfilerAgent(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect = NewWithT(t).Expect

		ctx libcnb.BuildContext
	)

	it.Before(func() {
		var err error

		ctx.Layers.Path, err = ioutil.TempDir("", "java-profiler-agent-layers")
		Expect(err).NotTo(HaveOccurred())
	})

	it.After(func() {
		Expect(os.RemoveAll(ctx.Layers.Path)).To(Succeed())
	})

	it("contributes Java agent", func() {
		dep := libpak.BuildpackDependency{
			URI:    "https://localhost/stub-stackdriver-profiler-agent.tar.gz",
			SHA256: "a27bbf74fa913fe70273c38831bc1043a2da0e28766423b81d8e9a042b353797",
		}
		dc := libpak.DependencyCache{CachePath: "testdata"}

		layer, err := ctx.Layers.Layer("test-layer")
		Expect(err).NotTo(HaveOccurred())

		layer, err = stackdriver.NewJavaProfilerAgent(dep, dc, &libcnb.BuildpackPlan{}).Contribute(layer)
		Expect(err).NotTo(HaveOccurred())

		Expect(layer.Launch).To(BeTrue())

		Expect(filepath.Join(layer.Path, "profiler_java_agent.so")).To(BeARegularFile())
		Expect(layer.Profile["java-profiler.sh"]).To(Equal(fmt.Sprintf(`if [[ -z "${BPL_GOOGLE_STACKDRIVER_MODULE+x}" ]]; then
    MODULE="default-module"
else
	MODULE=${BPL_GOOGLE_STACKDRIVER_MODULE}
fi

if [[ -z "${BPL_GOOGLE_STACKDRIVER_VERSION+x}" ]]; then
	VERSION=""
else
	VERSION=${BPL_GOOGLE_STACKDRIVER_VERSION}
fi

printf "Google Stackdriver Profiler enabled for %%s" "${MODULE}"
AGENT="-agentpath:%s=-logtostderr=1,-cprof_service=${MODULE}"

if [[ "${VERSION}" != "" ]]; then
	printf ":%%s" "${VERSION}"
	AGENT="${AGENT},-cprof_service_version=${VERSION}"
fi

printf "\n"
export JAVA_OPTS="${JAVA_OPTS} ${AGENT}"
`, filepath.Join(layer.Path, "profiler_java_agent.so"))))
	})
}
