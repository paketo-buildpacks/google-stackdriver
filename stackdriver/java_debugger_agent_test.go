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
	"github.com/paketo-buildpacks/google-stackdriver/stackdriver"
	"github.com/paketo-buildpacks/libpak"
	"github.com/sclevine/spec"
)

func testJavaDebuggerAgent(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect = NewWithT(t).Expect

		ctx libcnb.BuildContext
	)

	it.Before(func() {
		var err error

		ctx.Layers.Path, err = ioutil.TempDir("", "java-debugger-agent-layers")
		Expect(err).NotTo(HaveOccurred())
	})

	it.After(func() {
		Expect(os.RemoveAll(ctx.Layers.Path)).To(Succeed())
	})

	it("contributes Java agent", func() {
		dep := libpak.BuildpackDependency{
			URI:    "https://localhost/stub-stackdriver-debugger-agent.tar.gz",
			SHA256: "80ceb691b8b586e15dedae62564dea2cfe8e2f6ac44ec48fe4dc87599fa22cab",
		}
		dc := libpak.DependencyCache{CachePath: "testdata"}

		layer, err := ctx.Layers.Layer("test-layer")
		Expect(err).NotTo(HaveOccurred())

		layer, err = stackdriver.NewJavaDebuggerAgent(dep, dc, &libcnb.BuildpackPlan{}).Contribute(layer)
		Expect(err).NotTo(HaveOccurred())

		Expect(layer.Launch).To(BeTrue())

		Expect(filepath.Join(layer.Path, "cdbg_java_agent.so")).To(BeARegularFile())
		Expect(layer.Profile["debugger"]).To(Equal(fmt.Sprintf(`if [[ -z "${BPL_GOOGLE_STACKDRIVER_MODULE+x}" ]]; then
    MODULE="default-module"
else
	MODULE=${BPL_GOOGLE_STACKDRIVER_MODULE}
fi

if [[ -z "${BPL_GOOGLE_STACKDRIVER_VERSION+x}" ]]; then
	VERSION=""
else
	VERSION=${BPL_GOOGLE_STACKDRIVER_VERSION}
fi

printf "Google Stackdriver Debugger enabled for %%s" "${MODULE}"
export JAVA_OPTS="${JAVA_OPTS}
  -agentpath:%s=--logtostderr=1
  -Dcom.google.cdbg.auth.serviceaccount.enable=true
  -Dcom.google.cdbg.module=${MODULE}"

if [[ "${VERSION}" != "" ]]; then
	printf ":%%s" "${VERSION}"
	export JAVA_OPTS="${JAVA_OPTS} -Dcom.google.cdbg.version=${VERSION}"
fi

printf "\n"
`, filepath.Join(layer.Path, "cdbg_java_agent.so"))))
	})
}
