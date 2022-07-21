# `gcr.io/paketo-buildpacks/google-stackdriver`
The Paketo Google Stackdriver Buildpack is a Cloud Native Buildpack that contributes Stackdriver agents and configures them to connect to the service.

## Behavior
This buildpack will participate if any of the following conditions are met

* A binding exists with `type` of `StackdriverDebugger`
* A binding exists with `type` of `StackdriverProfiler`

The buildpack will do the following for Java applications:

* If `StackdriverDebugger` binding exists contributes a Java debugger agent to a layer and configures `$JAVA_TOOL_OPTIONS` to use it
* If `StackdriverProfiler` binding exists contributes a Java profiler agent to a layer and configures `$JAVA_TOOL_OPTIONS` to use it
* Sets `$GOOGLE_APPLICATION_CREDENTIALS` to the path of the `ApplicationCredentials` secret

The buildpack will do the following for NodeJS applications:

* If `StackdriverDebugger` binding exists
  * Contributes a NodeJS debugger agent to a layer and configures `$NODE_MODULES` to use it
  * If main module does not already require `@google-cloud/debug-agent` module, prepends the main module with `require('@google-cloud/debug-agent').start();`
* If `StackdriverProfiler` binding exists
  * Contributes a NodeJS profiler agent to a layer and configures `$NODE_MODULES` to use it
  * If main module does not already require `@google-cloud/profiler` module, prepends the main module with `require('@google-cloud/profiler').start();`
* Sets `$GOOGLE_APPLICATION_CREDENTIALS` to the path of the `ApplicationCredentials` secret

## Configuration
| Environment Variable | Description
| -------------------- | -----------
| `$BPL_GOOGLE_STACKDRIVER_MODULE` | Configure the name of the application.  Defaults to `default-module`.
| `$BPL_GOOGLE_STACKDRIVER_PROJECT_ID` | Configure the project id for the application.  Defaults to `<EMPTY>`.
| `$BPL_GOOGLE_STACKDRIVER_VERSION` | Configure the version of the application.  Defaults to `<EMPTY>`.

## Bindings
The buildpack optionally accepts the following bindings:

### Type: `dependency-mapping`
|Key                   | Value   | Description
|----------------------|---------|------------
|`<dependency-digest>` | `<uri>` | If needed, the buildpack will fetch the dependency with digest `<dependency-digest>` from `<uri>`

## License
This buildpack is released under version 2.0 of the [Apache License][a].

[a]: http://www.apache.org/licenses/LICENSE-2.0

