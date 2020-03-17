uri() {
  if [[ "${DEPENDENCY}" == "google-stackdriver-debugger-java" ]]; then
    echo "https://github.com/GoogleCloudPlatform/cloud-debug-java/releases/download/v$(cat "${ROOT}"/dependency/version)/compute-java_debian-wheezy_cdbg_java_agent_gce.tar"
  else
    cat "${ROOT}"/dependency/uri
  fi
}

sha256() {
  if [[ "${DEPENDENCY}" == "google-stackdriver-debugger-java" ]]; then
    shasum -a 256 "${ROOT}"/dependency/compute-java_debian-wheezy_cdbg_java_agent_gce.tar | cut -f 1 -d ' '
  else
    cat "${ROOT}"/dependency/sha256
  fi
}
