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

package stackdriver

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/buildpacks/libcnb"
	"github.com/paketo-buildpacks/libpak"
	"github.com/paketo-buildpacks/libpak/bard"
	"github.com/paketo-buildpacks/libpak/crush"
)

type JavaDebuggerAgent struct {
	LayerContributor libpak.DependencyLayerContributor
	Logger           bard.Logger
}

func NewJavaDebuggerAgent(dependency libpak.BuildpackDependency, cache libpak.DependencyCache, plan *libcnb.BuildpackPlan) JavaDebuggerAgent {
	return JavaDebuggerAgent{LayerContributor: libpak.NewDependencyLayerContributor(dependency, cache, plan)}
}

func (j JavaDebuggerAgent) Contribute(layer libcnb.Layer) (libcnb.Layer, error) {
	j.Logger.Body(bard.FormatUserConfig("BPL_GOOGLE_STACKDRIVER_MODULE", "the name of the application", "default-module"))
	j.Logger.Body(bard.FormatUserConfig("BPL_GOOGLE_STACKDRIVER_VERSION", "the version of the application", "<EMPTY>"))

	j.LayerContributor.Logger = j.Logger

	return j.LayerContributor.Contribute(layer, func(artifact *os.File) (libcnb.Layer, error) {
		j.Logger.Body("Expanding to %s", layer.Path)

		if err := crush.ExtractTarGz(artifact, layer.Path, 0); err != nil {
			return libcnb.Layer{}, fmt.Errorf("unable to extract %s\n%w", artifact.Name(), err)
		}

		layer.Profile.Add("debugger", `if [[ -z "${BPL_GOOGLE_STACKDRIVER_MODULE+x}" ]]; then
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
`, filepath.Join(layer.Path, "cdbg_java_agent.so"))

		layer.Launch = true
		return layer, nil
	})
}

func (JavaDebuggerAgent) Name() string {
	return "java-debugger"
}
