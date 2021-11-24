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

	dc, err := libpak.NewDependencyCache(context)
	if err != nil {
		return libcnb.BuildResult{}, fmt.Errorf("unable to create dependency cache\n%w", err)
	}
	dc.Logger = b.Logger

	names := []string{"credentials"}

	if _, ok, err := pr.Resolve("google-stackdriver-debugger-java"); err != nil {
		return libcnb.BuildResult{}, fmt.Errorf("unable to resolve google-stackdriver-debugger-java plan entry\n%w", err)
	} else if ok {
		dep, err := dr.Resolve("google-stackdriver-debugger-java", "")
		if err != nil {
			return libcnb.BuildResult{}, fmt.Errorf("unable to find dependency\n%w", err)
		}

		ja, be := NewJavaDebuggerAgent(dep, dc)
		ja.Logger = b.Logger
		result.Layers = append(result.Layers, ja)
		if be.Name != "" {
			result.BOM.Entries = append(result.BOM.Entries, be)
		}

		names = append(names, "java-debugger")
	}

	if _, ok, err := pr.Resolve("google-stackdriver-debugger-nodejs"); err != nil {
		return libcnb.BuildResult{}, fmt.Errorf("unable to resolve google-stackdriver-debugger-nodejs plan entry\n%w", err)
	} else if ok {
		dep, err := dr.Resolve("google-stackdriver-debugger-nodejs", "")
		if err != nil {
			return libcnb.BuildResult{}, fmt.Errorf("unable to find dependency\n%w", err)
		}

		ja, be := NewNodeJSDebuggerAgent(context.Buildpack.Path, dep, dc)
		ja.Logger = b.Logger
		result.Layers = append(result.Layers, ja)
		if be.Name != "" {
			result.BOM.Entries = append(result.BOM.Entries, be)
		}
	}

	if _, ok, err := pr.Resolve("google-stackdriver-profiler-java"); err != nil {
		return libcnb.BuildResult{}, fmt.Errorf("unable to resolve google-stackdriver-profiler-java plan entry\n%w", err)
	} else if ok {
		dep, err := dr.Resolve("google-stackdriver-profiler-java", "")
		if err != nil {
			return libcnb.BuildResult{}, fmt.Errorf("unable to find dependency\n%w", err)
		}

		ja, be := NewJavaProfilerAgent(dep, dc)
		ja.Logger = b.Logger
		result.Layers = append(result.Layers, ja)
		if be.Name != "" {
			result.BOM.Entries = append(result.BOM.Entries, be)
		}

		names = append(names, "java-profiler")
	}

	if _, ok, err := pr.Resolve("google-stackdriver-profiler-nodejs"); err != nil {
		return libcnb.BuildResult{}, fmt.Errorf("unable to resolve google-stackdriver-profiler-nodejs plan entry\n%w", err)
	} else if ok {
		dep, err := dr.Resolve("google-stackdriver-profiler-nodejs", "")
		if err != nil {
			return libcnb.BuildResult{}, fmt.Errorf("unable to find dependency\n%w", err)
		}

		ja, be := NewNodeJSProfilerAgent(context.Buildpack.Path, dep, dc)
		ja.Logger = b.Logger
		result.Layers = append(result.Layers, ja)
		if be.Name != "" {
			result.BOM.Entries = append(result.BOM.Entries, be)
		}
	}

	h, be := libpak.NewHelperLayer(context.Buildpack, names...)
	h.Logger = b.Logger
	result.Layers = append(result.Layers, h)
	if be.Name != "" {
		result.BOM.Entries = append(result.BOM.Entries, be)
	}

	return result, nil
}
