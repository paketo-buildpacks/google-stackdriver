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

package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/buildpacks/libcnb"

	"github.com/paketo-buildpacks/google-cloud/credentials"
	"github.com/paketo-buildpacks/google-cloud/debugger"
	"github.com/paketo-buildpacks/google-cloud/internal/common"
	"github.com/paketo-buildpacks/google-cloud/profiler"
	"github.com/paketo-buildpacks/libpak/bard"
	"github.com/paketo-buildpacks/libpak/bindings"
	"github.com/paketo-buildpacks/libpak/sherpa"
)

func main() {
	sherpa.Execute(func() error {
		l := bard.NewLogger(os.Stdout)

		a, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("unable to read working directory\n%w", err)
		}

		bs, err := libcnb.NewBindingsFromEnvironment()
		if err != nil {
			return fmt.Errorf("unable to read bindings from environment\n%w", err)
		}

		var (
			b libcnb.Binding
			c common.CredentialSource
		)

		if hasMetadataServer() {
			c = common.MetadataServer
		} else if g, ok, err := bindings.ResolveOne(bs, bindings.OfType("GoogleCloud")); err != nil {
			return fmt.Errorf("unable to resolve GoogleCloud binding\n%w", err)
		} else if ok {
			b = g
			c = common.Binding
		} else {
			c = common.None
		}

		return sherpa.Helpers(map[string]sherpa.ExecD{
			common.Credentials:    credentials.Launch{Binding: b, CredentialSource: c, Logger: l},
			common.DebuggerJava:   debugger.JavaLaunch{CredentialSource: c, Logger: l},
			common.DebuggerNodeJS: debugger.NodeJSLaunch{ApplicationPath: a, CredentialSource: c, Logger: l},
			common.ProfilerJava:   profiler.JavaLaunch{CredentialSource: c, Logger: l},
			common.ProfilerNodeJS: profiler.NodeJSLaunch{ApplicationPath: a, CredentialSource: c, Logger: l},
		})
	})
}

// hasMetadataServer detects whether the application has access to a Google Cloud metadata server.  Detection is based
// on the existence of the metadata server as defined in
// https://cloud.google.com/compute/docs/storing-retrieving-metadata#querying.
func hasMetadataServer() bool {
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: func(req *http.Request) (*url.URL, error) { return nil, nil },
		},
	}

	req, err := http.NewRequest("GET", "http://metadata.google.internal/computeMetadata/v1/", nil)
	if err != nil {
		return false
	}
	req.Header.Add("Metadata-Flavor", "Google")

	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == 200
}
