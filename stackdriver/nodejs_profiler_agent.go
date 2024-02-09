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

package stackdriver

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/buildpacks/libcnb"
	"github.com/paketo-buildpacks/libpak"
	"github.com/paketo-buildpacks/libpak/bard"
	"github.com/paketo-buildpacks/libpak/effect"
	"github.com/paketo-buildpacks/libpak/sherpa"
)

type NodeJSProfilerAgent struct {
	ApplicationPath  string
	Executor         effect.Executor
	LayerContributor libpak.DependencyLayerContributor
	Logger           bard.Logger
}

func NewNodeJSProfilerAgent(applicationPath string, dependency libpak.BuildpackDependency, cache libpak.DependencyCache) (NodeJSProfilerAgent, libcnb.BOMEntry) {
	contributor, entry := libpak.NewDependencyLayer(dependency, cache, libcnb.LayerTypes{Launch: true})
	return NodeJSProfilerAgent{
		ApplicationPath:  applicationPath,
		Executor:         effect.NewExecutor(),
		LayerContributor: contributor,
	}, entry
}

func (n NodeJSProfilerAgent) Contribute(layer libcnb.Layer) (libcnb.Layer, error) {
	n.LayerContributor.Logger = n.Logger

	layer, err := n.LayerContributor.Contribute(layer, func(artifact *os.File) (libcnb.Layer, error) {
		n.Logger.Bodyf("Installing to %s", layer.Path)

		if err := n.Executor.Execute(effect.Execution{
			Command: "npm",
			Args:    []string{"install", "--no-save", artifact.Name()},
			Dir:     layer.Path,
			Stdout:  n.Logger.InfoWriter(),
			Stderr:  n.Logger.InfoWriter(),
		}); err != nil {
			return libcnb.Layer{}, fmt.Errorf("unable to run npm install\n%w", err)
		}

		layer.LaunchEnvironment.Prepend("NODE_PATH", string(os.PathListSeparator), filepath.Join(layer.Path, "node_modules"))

		return layer, nil
	})
	if err != nil {
		return libcnb.Layer{}, fmt.Errorf("unable to install node module\n%w", err)
	}

	m, err := sherpa.NodeJSMainModule(n.ApplicationPath)
	if err != nil {
		return libcnb.Layer{}, fmt.Errorf("unable to find main module in %s\n%w", n.ApplicationPath, err)
	}

	file := filepath.Join(n.ApplicationPath, m)
	c, err := os.ReadFile(file)
	if err != nil {
		return libcnb.Layer{}, fmt.Errorf("unable to read contents of %s\n%w", file, err)
	}

	if !regexp.MustCompile(`require\(['"]@google-cloud/profiler['"]\)`).Match(c) {
		n.Logger.Header("Requiring '@google-cloud/profiler' module")

		if err := os.WriteFile(file, append([]byte("require('@google-cloud/profiler').start();\n"), c...), 0644); err != nil {
			return libcnb.Layer{}, fmt.Errorf("unable to write main module %s\n%w", file, err)
		}
	}

	return layer, nil
}

func (n NodeJSProfilerAgent) Name() string {
	return n.LayerContributor.LayerName()
}
