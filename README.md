# Paketo Buildpack for Google StackDriver

## Buildpack ID: `paketo-buildpacks/google-stackdriver`
## Registry URLs: `docker.io/paketobuildpacks/google-stackdriver`

The Paketo Buildpack for Google Stackdriver is a Cloud Native Buildpack that contributes Stackdriver profiler agent and configure it to connect to the service.

## Behavior

This buildpack will participate if any of the following conditions are met

* A binding exists with `type` of `StackdriverProfiler`

The buildpack will do the following for Java applications:

* If `StackdriverProfiler` binding exists contributes a Java profiler agent to a layer and configures `$JAVA_TOOL_OPTIONS` to use it
* Sets `$GOOGLE_APPLICATION_CREDENTIALS` to the path of the `ApplicationCredentials` secret

The buildpack will do the following for NodeJS applications:

* If `StackdriverProfiler` binding exists
  * Contributes a NodeJS profiler agent to a layer and configures `$NODE_MODULES` to use it
  * If main module does not already require `@google-cloud/profiler` module, prepends the main module with `require('@google-cloud/profiler').start();`
* Sets `$GOOGLE_APPLICATION_CREDENTIALS` to the path of the `ApplicationCredentials` secret

## Notes

While this buildpack is packaged to support ARM64 there are no upstream bindings for Aternity on ARM64. If you try to use this buildpack on ARM64 it will fail attempting to download ARM64 binaries that do not exist. Please contact the vendor if you'd like to see ARM64 support. The buildpack is positioned to support it as soon as it is available upstream.

## Configuration

| Environment Variable                 | Description                                                           |
| ------------------------------------ | --------------------------------------------------------------------- |
| `$BPL_GOOGLE_STACKDRIVER_MODULE`     | Configure the name of the application.  Defaults to `default-module`. |
| `$BPL_GOOGLE_STACKDRIVER_PROJECT_ID` | Configure the project id for the application.  Defaults to `<EMPTY>`. |
| `$BPL_GOOGLE_STACKDRIVER_VERSION`    | Configure the version of the application.  Defaults to `<EMPTY>`.     |

## Bindings

The buildpack optionally accepts the following bindings:

### Type: `dependency-mapping`

| Key                   | Value   | Description                                                                                       |
| --------------------- | ------- | ------------------------------------------------------------------------------------------------- |
| `<dependency-digest>` | `<uri>` | If needed, the buildpack will fetch the dependency with digest `<dependency-digest>` from `<uri>` |

## License

This buildpack is released under version 2.0 of the [Apache License][a].

[a]: http://www.apache.org/licenses/LICENSE-2.0

