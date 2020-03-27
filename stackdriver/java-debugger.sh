if [[ -z "${BPL_GOOGLE_STACKDRIVER_MODULE+x}" ]]; then
    MODULE="default-module"
else
	MODULE=${BPL_GOOGLE_STACKDRIVER_MODULE}
fi

if [[ -z "${BPL_GOOGLE_STACKDRIVER_VERSION+x}" ]]; then
	VERSION=""
else
	VERSION=${BPL_GOOGLE_STACKDRIVER_VERSION}
fi

printf "Google Stackdriver Debugger enabled for %s" "${MODULE}"
export JAVA_OPTS="${JAVA_OPTS}
  -agentpath:{{.agentPath}}=--logtostderr=1
  -Dcom.google.cdbg.auth.serviceaccount.enable=true
  -Dcom.google.cdbg.module=${MODULE}"

if [[ "${VERSION}" != "" ]]; then
	printf ":%s" "${VERSION}"
	export JAVA_OPTS="${JAVA_OPTS} -Dcom.google.cdbg.version=${VERSION}"
fi

printf "\n"
