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

package common

const (
	Module    = "BPL_GOOGLE_CLOUD_MODULE"
	ProjectID = "BPL_GOOGLE_CLOUD_PROJECT_ID"
	Version   = "BPL_GOOGLE_CLOUD_VERSION"
)

const (
	Credentials    = "google-cloud-credentials"
	DebuggerJava   = "google-cloud-debugger-java"
	DebuggerNodeJS = "google-cloud-debugger-nodejs"
	ProfilerJava   = "google-cloud-profiler-java"
	ProfilerNodeJS = "google-cloud-profiler-nodejs"
)

type CredentialSource uint8

const (
	Binding CredentialSource = iota
	MetadataServer
	None
)
