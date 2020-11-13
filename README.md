# `paketo-buildpacks/google-cloud`
The Paketo Google Cloud Buildpack is a Cloud Native Buildpack that contributes Google Cloud agents and configures them to connect to their services.

## Behavior

* If `$BP_GOOGLE_CLOUD_DEBUGGER_ENABLED` is set to `true` and the application is Java
  * At build time, contributes an agent to a layer
  * At launch time, if credentials are available, configures the application to use the agent
* If `$BP_GOOGLE_CLOUD_DEBUGGER_ENABLED` is set to `true` and the application is NodeJS
  * At build time, contributes an agent to a layer
  * At launch time, if credentials are available, configures `$NODE_MODULES` with the agent path.  If the main module does not already require `@google-cloud/debug-agent`, prepends the main module with `require('@google-cloud/debug-agent').start({...});`.

* If `$BP_GOOGLE_CLOUD_PROFILER_ENABLED` is set to `true` and the application is Java
  * At build time, contributes an agent to a layer
  * At launch time, if credentials are available, configures the application to use the agent
* If `$BP_GOOGLE_CLOUD_PROFILER_ENABLED` is set to `true` and the application is NodeJS
  * At build time, contributes an agent to a layer
  * At launch time, if credentials are available, configures `$NODE_MODULES` with the agent path.  If the main module does not already require `@google-cloud/profiler`, prepends the main module with `require('@google-cloud/profiler').start({...});`.

### Credential Availability
If the applications runs within Google Cloud and the [Google Metadata Service][m] is accessible, those credentials will be used.  If the application runs within any other environment, credentials must be provided with a service binding as described below.

[m]: https://cloud.google.com/compute/docs/storing-retrieving-metadata

## Configuration
| Environment Variable | Description
| -------------------- | -----------
| `$BP_GOOGLE_CLOUD_DEBUGGER_ENABLED` | Whether to add Google Cloud Debugger during build
| `$BP_GOOGLE_CLOUD_PROFILER_ENABLED` | Whether to add Google Cloud Profiler during build 
| `$BPL_GOOGLE_CLOUD_MODULE` | Configure the name of the application (required)
| `$BPL_GOOGLE_CLOUD_PROJECT_ID` | Configure the project id for the application (required if running outside of Google Cloud)
| `$BPL_GOOGLE_CLOUD_VERSION` | Configure the version of the application (required)

## Bindings
The buildpack optionally accepts the following bindings:

### Type: `GoogleCloud`
|Key                      | Value            | Description
|-------------------------|------------------|------------
|`ApplicationCredentials` | `<JSON Payload>` | Google Cloud Application Credentials in JSON form

### Type: `dependency-mapping`
|Key                   | Value   | Description
|----------------------|---------|------------
|`<dependency-digest>` | `<uri>` | If needed, the buildpack will fetch the dependency with digest `<dependency-digest>` from `<uri>`

## License
This buildpack is released under version 2.0 of the [Apache License][a].

[a]: http://www.apache.org/licenses/LICENSE-2.0
