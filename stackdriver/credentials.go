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
	"github.com/paketo-buildpacks/libpak/sherpa"
)

type Credentials struct {
	LayerContributor libpak.HelperLayerContributor
	Logger           bard.Logger
}

func NewCredentials(buildpack libcnb.Buildpack, plan *libcnb.BuildpackPlan) Credentials {
	return Credentials{
		LayerContributor: libpak.NewHelperLayerContributor(filepath.Join(buildpack.Path, "bin", "google-application-credentials"),
			"Google Application Credentials", buildpack.Info, plan),
	}
}

func (c Credentials) Contribute(layer libcnb.Layer) (libcnb.Layer, error) {
	c.LayerContributor.Logger = c.Logger

	return c.LayerContributor.Contribute(layer, func(artifact *os.File) (libcnb.Layer, error) {
		c.Logger.Body("Copying to %s", layer.Path)

		file := filepath.Join(layer.Path, "bin", "google-application-credentials")
		if err := sherpa.CopyFile(artifact, file); err != nil {
			return libcnb.Layer{}, fmt.Errorf("unable to copy %s to %s\n%w", artifact.Name(), file, err)
		}

		layer.Profile.Add("properties", `printf "Configuring Google application credentials\n"

eval $(google-application-credentials)
`)

		layer.Launch = true
		return layer, nil
	})
}

func (Credentials) Name() string {
	return "credentials"
}
