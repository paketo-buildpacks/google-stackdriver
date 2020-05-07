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

	"github.com/buildpacks/libcnb"
	"github.com/paketo-buildpacks/libpak"
	"github.com/paketo-buildpacks/libpak/bard"
)

type Build struct {
	Logger bard.Logger
}

func (b Build) Build(context libcnb.BuildContext) (libcnb.BuildResult, error) {
	b.Logger.Title(context.Buildpack)
	result := libcnb.NewBuildResult()

	_, err := libpak.NewConfigurationResolver(context.Buildpack, &b.Logger)
	if err != nil {
		return libcnb.BuildResult{}, fmt.Errorf("unable to create configuration resolver\n%w", err)
	}

	pr := libpak.PlanEntryResolver{Plan: context.Plan}

	dr, err := libpak.NewDependencyResolver(context)
	if err != nil {
		return libcnb.BuildResult{}, fmt.Errorf("unable to create dependency resolver\n%w", err)
	}

	dc := libpak.NewDependencyCache(context.Buildpack)
	dc.Logger = b.Logger

	if e, ok, err := pr.Resolve("google-stackdriver-debugger-java"); err != nil {
		return libcnb.BuildResult{}, fmt.Errorf("unable to resolve google-stackdriver-debugger-java plan entry\n%w", err)
	} else if ok {
		dep, err := dr.Resolve("google-stackdriver-debugger-java", e.Version)
		if err != nil {
			return libcnb.BuildResult{}, fmt.Errorf("unable to find dependency\n%w", err)
		}

		ja := NewJavaDebuggerAgent(dep, dc, result.Plan)
		ja.Logger = b.Logger
		result.Layers = append(result.Layers, ja)
	}

	if e, ok, err := pr.Resolve("google-stackdriver-debugger-nodejs"); err != nil {
		return libcnb.BuildResult{}, fmt.Errorf("unable to resolve google-stackdriver-debugger-nodejs plan entry\n%w", err)
	} else if ok {
		dep, err := dr.Resolve("google-stackdriver-debugger-nodejs", e.Version)
		if err != nil {
			return libcnb.BuildResult{}, fmt.Errorf("unable to find dependency\n%w", err)
		}

		ja := NewNodeJSDebuggerAgent(context.Buildpack.Path, dep, dc, result.Plan)
		ja.Logger = b.Logger
		result.Layers = append(result.Layers, ja)
	}

	if e, ok, err := pr.Resolve("google-stackdriver-profiler-java"); err != nil {
		return libcnb.BuildResult{}, fmt.Errorf("unable to resolve google-stackdriver-profiler-java plan entry\n%w", err)
	} else if ok {
		dep, err := dr.Resolve("google-stackdriver-profiler-java", e.Version)
		if err != nil {
			return libcnb.BuildResult{}, fmt.Errorf("unable to find dependency\n%w", err)
		}

		ja := NewJavaProfilerAgent(dep, dc, result.Plan)
		ja.Logger = b.Logger
		result.Layers = append(result.Layers, ja)
	}

	if e, ok, err := pr.Resolve("google-stackdriver-profiler-nodejs"); err != nil {
		return libcnb.BuildResult{}, fmt.Errorf("unable to resolve google-stackdriver-profiler-nodejs plan entry\n%w", err)
	} else if ok {
		dep, err := dr.Resolve("google-stackdriver-profiler-nodejs", e.Version)
		if err != nil {
			return libcnb.BuildResult{}, fmt.Errorf("unable to find dependency\n%w", err)
		}

		ja := NewNodeJSProfilerAgent(context.Buildpack.Path, dep, dc, result.Plan)
		ja.Logger = b.Logger
		result.Layers = append(result.Layers, ja)
	}

	c := NewCredentials(context.Buildpack, result.Plan)
	c.Logger = b.Logger
	result.Layers = append(result.Layers, c)

	return result, nil
}
