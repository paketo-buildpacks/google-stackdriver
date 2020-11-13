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

package profiler

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/buildpacks/libcnb"

	"github.com/paketo-buildpacks/google-cloud/internal/common"
	"github.com/paketo-buildpacks/libpak"
	"github.com/paketo-buildpacks/libpak/bard"
	"github.com/paketo-buildpacks/libpak/crush"
	"github.com/paketo-buildpacks/libpak/sherpa"
)

const AgentPath = "BPI_GOOGLE_CLOUD_PROFILER_AGENT_PATH"

type JavaBuild struct {
	LayerContributor libpak.DependencyLayerContributor
	Logger           bard.Logger
}

func NewJavaBuild(dependency libpak.BuildpackDependency, cache libpak.DependencyCache, plan *libcnb.BuildpackPlan) JavaBuild {
	return JavaBuild{LayerContributor: libpak.NewDependencyLayerContributor(dependency, cache, plan)}
}

func (j JavaBuild) Contribute(layer libcnb.Layer) (libcnb.Layer, error) {
	j.LayerContributor.Logger = j.Logger

	return j.LayerContributor.Contribute(layer, func(artifact *os.File) (libcnb.Layer, error) {
		j.Logger.Bodyf("Expanding to %s", layer.Path)

		if err := crush.ExtractTarGz(artifact, layer.Path, 0); err != nil {
			return libcnb.Layer{}, fmt.Errorf("unable to extract %s\n%w", artifact.Name(), err)
		}

		layer.LaunchEnvironment.Default(AgentPath, filepath.Join(layer.Path, "profiler_java_agent.so"))

		return layer, nil
	}, libpak.LaunchLayer)

}

func (j JavaBuild) Name() string {
	return j.LayerContributor.LayerName()
}

type JavaLaunch struct {
	CredentialSource common.CredentialSource
	Logger           bard.Logger
}

// https://cloud.google.com/profiler/docs/profiling-java#agent_configuration
// https://cloud.google.com/profiler/docs/profiling-external#linking_the_agent_to_a_project
func (j JavaLaunch) Execute() (map[string]string, error) {
	if j.CredentialSource == common.None {
		j.Logger.Info("Google Cloud Profiler disabled")
		return nil, nil
	}

	p, err := sherpa.GetEnvRequired(AgentPath)
	if err != nil {
		return nil, err
	}

	m, err := sherpa.GetEnvRequired(common.Module)
	if err != nil {
		return nil, err
	}

	v, err := sherpa.GetEnvRequired(common.Version)
	if err != nil {
		return nil, err
	}

	j.Logger.Infof("Google Cloud Profiler enabled (%s:%s)", m, v)

	var agent []string
	agent = append(agent,
		fmt.Sprintf("-agentpath:%s=-logtostderr=1", p),
		fmt.Sprintf("-cprof_service=%s", m),
		fmt.Sprintf("-cprof_service_version=%s", v),
	)

	if j.CredentialSource == common.Binding {
		projectId, err := sherpa.GetEnvRequired(common.ProjectID)
		if err != nil {
			return nil, err
		}

		agent = append(agent, fmt.Sprintf("-cprof_project_id=%s", projectId))
	}

	return map[string]string{
		"JAVA_TOOL_OPTIONS": sherpa.AppendToEnvVar("JAVA_TOOL_OPTIONS", " ", strings.Join(agent, ",")),
	}, nil
}
