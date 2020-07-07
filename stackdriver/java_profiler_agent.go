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
	"github.com/paketo-buildpacks/libpak/sherpa"

	_ "github.com/paketo-buildpacks/google-stackdriver/stackdriver/statik"
)

type JavaProfilerAgent struct {
	LayerContributor libpak.DependencyLayerContributor
	Logger           bard.Logger
}

func NewJavaProfilerAgent(dependency libpak.BuildpackDependency, cache libpak.DependencyCache, plan *libcnb.BuildpackPlan) JavaProfilerAgent {
	return JavaProfilerAgent{LayerContributor: libpak.NewDependencyLayerContributor(dependency, cache, plan)}
}

//go:generate statik -src . -include *.sh

func (j JavaProfilerAgent) Contribute(layer libcnb.Layer) (libcnb.Layer, error) {
	j.LayerContributor.Logger = j.Logger

	return j.LayerContributor.Contribute(layer, func(artifact *os.File) (libcnb.Layer, error) {
		j.Logger.Bodyf("Expanding to %s", layer.Path)

		if err := crush.ExtractTarGz(artifact, layer.Path, 0); err != nil {
			return libcnb.Layer{}, fmt.Errorf("unable to extract %s\n%w", artifact.Name(), err)
		}

		s, err := sherpa.TemplateFile("/java-profiler.sh", map[string]interface{}{
			"agentpath": filepath.Join(layer.Path, "profiler_java_agent.so"),
		})
		if err != nil {
			return libcnb.Layer{}, fmt.Errorf("unable to load java-profiler.sh\n%w", err)
		}

		layer.Profile.Add("java-profiler.sh", s)

		layer.Launch = true
		return layer, nil
	})

}

func (JavaProfilerAgent) Name() string {
	return "java-profiler"
}
