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

package credentials_test

import (
	"testing"

	"github.com/buildpacks/libcnb"
	. "github.com/onsi/gomega"
	"github.com/sclevine/spec"

	"github.com/paketo-buildpacks/google-stackdriver/credentials"
)

func testCredentials(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect = NewWithT(t).Expect

		c credentials.Credentials
	)

	it("does not contribute properties if no binding exists", func() {
		Expect(c.Execute()).To(BeZero())
	})

	it("contributes credentials if debugger binding exists", func() {
		c.Bindings = libcnb.Bindings{
			{
				Name:   "test-binding",
				Path:   "/test/path/test-binding",
				Type:   "StackdriverDebugger",
				Secret: map[string]string{"ApplicationCredentials": "test-value"},
			},
		}

		Expect(c.Execute()).To(Equal(`export GOOGLE_APPLICATION_CREDENTIALS="/test/path/test-binding/ApplicationCredentials"`))
	})

	it("contributes credentials if profiler binding exists", func() {
		c.Bindings = libcnb.Bindings{
			{
				Name:   "test-binding",
				Path:   "/test/path/test-binding",
				Type:   "StackdriverProfiler",
				Secret: map[string]string{"ApplicationCredentials": "test-value"},
			},
		}

		Expect(c.Execute()).To(Equal(`export GOOGLE_APPLICATION_CREDENTIALS="/test/path/test-binding/ApplicationCredentials"`))
	})
}
