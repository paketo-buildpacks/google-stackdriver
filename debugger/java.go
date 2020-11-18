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

package debugger

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/buildpacks/libcnb"

	"github.com/paketo-buildpacks/google-cloud/internal/common"
	"github.com/paketo-buildpacks/libpak"
	"github.com/paketo-buildpacks/libpak/bard"
	"github.com/paketo-buildpacks/libpak/crush"
	"github.com/paketo-buildpacks/libpak/sherpa"
)

const AgentPath = "BPI_GOOGLE_CLOUD_DEBUGGER_AGENT_PATH"

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

		layer.LaunchEnvironment.Default(AgentPath, filepath.Join(layer.Path, "cdbg_java_agent.so"))

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

// https://cloud.google.com/debugger/docs/setup/java#gke
// https://cloud.google.com/debugger/docs/setup/java#local
func (j JavaLaunch) Execute() (map[string]string, error) {
	if j.CredentialSource == common.None {
		j.Logger.Info("Google Cloud Debugger disabled")
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

	j.Logger.Infof("Google Cloud Debugger enabled (%s:%s)", m, v)

	var agent []string
	agent = append(agent,
		fmt.Sprintf("-agentpath:%s=--logtostderr=1", p),
		fmt.Sprintf("-Dcom.google.cdbg.module=%s", m),
		fmt.Sprintf("-Dcom.google.cdbg.version=%s", v),
	)

	if j.CredentialSource == common.Binding {
		agent = append(agent, "-Dcom.google.cdbg.auth.serviceaccount.enable=true")
	}

	return map[string]string{
		"JAVA_TOOL_OPTIONS": sherpa.AppendToEnvVar("JAVA_TOOL_OPTIONS", " ", agent...),
	}, nil
}
