// Copyright 2018-2026 the original author or authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"regexp"

	"github.com/Masterminds/semver/v3"
	"github.com/paketo-buildpacks/libdependency/retrieve"
	"github.com/paketo-buildpacks/libdependency/upstream"
	"github.com/paketo-buildpacks/libdependency/versionology"
	"github.com/paketo-buildpacks/packit/v2/cargo"
)

const (
	id   = "google-stackdriver-profiler-java"
	name = "Google Stackdriver Profiler Java Agent"
	purl = "pkg:generic/google-stackdriver-profiler-java"

	bucket = "cloud-profiler"
	prefix = "java/cloud-profiler-java-agent"
)

var namePattern = regexp.MustCompile(`^.*_(\d{8})_RC(\d+)\.tar\.gz$`)

type object struct {
	Name string `json:"name"`
}

type objectListing struct {
	NextPageToken string   `json:"nextPageToken"`
	Items         []object `json:"items"`
}

type version struct {
	version *semver.Version
	name    string
}

func (v version) Version() *semver.Version {
	return v.version
}

func main() {
	retrieve.NewMetadata(id, getAllVersions, generateMetadata)
}

func getAllVersions() (versionology.VersionFetcherArray, error) {
	objects, err := listObjects(bucket, prefix)
	if err != nil {
		return nil, fmt.Errorf("unable to list objects in %s\n%w", bucket, err)
	}

	var versions versionology.VersionFetcherArray
	for _, o := range objects {
		m := namePattern.FindStringSubmatch(o.Name)
		if m == nil {
			continue
		}

		date := m[1]
		versionString := fmt.Sprintf("%s.%s.%s%s", date[0:4], date[4:6], date[6:8], m[2])

		v, err := semver.NewVersion(versionString)
		if err != nil {
			fmt.Printf("Skipping %s: unable to parse version\n", versionString)
			continue
		}

		versions = append(versions, version{version: v, name: o.Name})
	}

	return versions, nil
}

func generateMetadata(versionFetcher versionology.VersionFetcher) ([]versionology.Dependency, error) {
	v, ok := versionFetcher.(version)
	if !ok {
		return nil, fmt.Errorf("unexpected version type %T", versionFetcher)
	}

	versionString := v.version.String()

	uri := fmt.Sprintf("https://storage.googleapis.com/%s/%s", bucket, v.name)
	checksum, err := upstream.GetSHA256OfRemoteFile(uri)
	if err != nil {
		return nil, fmt.Errorf("unable to checksum %s\n%w", uri, err)
	}

	dependency := cargo.ConfigMetadataDependency{
		Checksum: fmt.Sprintf("sha256:%s", checksum),
		CPE:      fmt.Sprintf("cpe:2.3:a:google:google-stackdriver-profiler-java:%s:*:*:*:*:*:*:*", versionString),
		ID:       id,
		Licenses: []interface{}{
			map[string]string{
				"type": "Apache-2.0",
				"uri":  "https://github.com/GoogleCloudPlatform/cloud-profiler-java/blob/master/LICENSE",
			},
		},
		Name:    name,
		PURL:    fmt.Sprintf("%s@%s", purl, versionString),
		SHA256:  checksum,
		Stacks:  []string{"*"},
		URI:     uri,
		Version: versionString,
	}

	return versionology.NewDependencyArray(dependency, "")
}

func listObjects(bucket, prefix string) ([]object, error) {
	var all []object
	pageToken := ""

	for {
		u := fmt.Sprintf("https://storage.googleapis.com/storage/v1/b/%s/o", bucket)
		q := url.Values{}
		q.Set("prefix", prefix)
		q.Set("maxResults", "1000")
		if pageToken != "" {
			q.Set("pageToken", pageToken)
		}

		resp, err := http.Get(fmt.Sprintf("%s?%s", u, q.Encode()))
		if err != nil {
			return nil, err
		}

		if resp.StatusCode != http.StatusOK {
			_ = resp.Body.Close()
			return nil, fmt.Errorf("storage API returned status %d", resp.StatusCode)
		}

		var listing objectListing
		if err := json.NewDecoder(resp.Body).Decode(&listing); err != nil {
			_ = resp.Body.Close()
			return nil, err
		}
		_ = resp.Body.Close()

		all = append(all, listing.Items...)

		if listing.NextPageToken == "" {
			break
		}
		pageToken = listing.NextPageToken
	}

	return all, nil
}
