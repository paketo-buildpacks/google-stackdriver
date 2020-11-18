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

package google

import (
	"fmt"

	"github.com/buildpacks/libcnb"

	"github.com/paketo-buildpacks/google-cloud/debugger"
	"github.com/paketo-buildpacks/google-cloud/internal/common"
	"github.com/paketo-buildpacks/google-cloud/profiler"
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

	var names []string

	if _, ok, err := pr.Resolve(common.Credentials); err != nil {
		return libcnb.BuildResult{}, fmt.Errorf("unable to resolve %s plan entry\n%w", common.Credentials, err)
	} else if ok {
		names = append(names, common.Credentials)
	}

	if _, ok, err := pr.Resolve(common.DebuggerJava); err != nil {
		return libcnb.BuildResult{}, fmt.Errorf("unable to resolve %s plan entry\n%w", common.DebuggerJava, err)
	} else if ok {
		dep, err := dr.Resolve(common.DebuggerJava, "")
		if err != nil {
			return libcnb.BuildResult{}, fmt.Errorf("unable to find dependency\n%w", err)
		}

		ja := debugger.NewJavaBuild(dep, dc, result.Plan)
		ja.Logger = b.Logger
		result.Layers = append(result.Layers, ja)

		names = append(names, common.DebuggerJava)
	}

	if _, ok, err := pr.Resolve(common.DebuggerNodeJS); err != nil {
		return libcnb.BuildResult{}, fmt.Errorf("unable to resolve %s plan entry\n%w", common.DebuggerNodeJS, err)
	} else if ok {
		dep, err := dr.Resolve(common.DebuggerNodeJS, "")
		if err != nil {
			return libcnb.BuildResult{}, fmt.Errorf("unable to find dependency\n%w", err)
		}

		ja := debugger.NewNodeJSBuild(dep, dc, result.Plan)
		ja.Logger = b.Logger
		result.Layers = append(result.Layers, ja)

		names = append(names, common.DebuggerNodeJS)
	}

	if _, ok, err := pr.Resolve(common.ProfilerJava); err != nil {
		return libcnb.BuildResult{}, fmt.Errorf("unable to resolve %s plan entry\n%w", common.ProfilerJava, err)
	} else if ok {
		dep, err := dr.Resolve(common.ProfilerJava, "")
		if err != nil {
			return libcnb.BuildResult{}, fmt.Errorf("unable to find dependency\n%w", err)
		}

		ja := profiler.NewJavaBuild(dep, dc, result.Plan)
		ja.Logger = b.Logger
		result.Layers = append(result.Layers, ja)

		names = append(names, common.ProfilerJava)
	}

	if _, ok, err := pr.Resolve(common.ProfilerNodeJS); err != nil {
		return libcnb.BuildResult{}, fmt.Errorf("unable to resolve %s plan entry\n%w", common.ProfilerNodeJS, err)
	} else if ok {
		dep, err := dr.Resolve(common.ProfilerNodeJS, "")
		if err != nil {
			return libcnb.BuildResult{}, fmt.Errorf("unable to find dependency\n%w", err)
		}

		ja := profiler.NewNodeJSBuild(dep, dc, result.Plan)
		ja.Logger = b.Logger
		result.Layers = append(result.Layers, ja)

		names = append(names, common.ProfilerNodeJS)
	}

	h := libpak.NewHelperLayerContributor(context.Buildpack, result.Plan, names...)
	h.Logger = b.Logger
	result.Layers = append(result.Layers, h)

	return result, nil
}
