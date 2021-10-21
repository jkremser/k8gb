#!/bin/bash
[ "${DEBUG}" == 1 ] && set -x

TOOL_VERSION=${TOOL_VERSION:-"0.5.3"}

main() {
    # checks
    [[ $# != 1 ]] && echo "Usage: $0 <version> # provide version in x.y.z format" && exit 1
    _VERSION=$1
    _VERSION=${_VERSION#"v"}

    git checkout v${_VERSION}

    _OS=${OS:-`uname | tr '[:upper:]' '[:lower:]'`}
    _ARCH=""
    case $(uname -m) in
        x86_64)     _ARCH="amd64" ;;
        i386 | i686) _ARCH="386" ;;
        arm)         dpkg --print-architecture | grep -q "arm64" && _ARCH="arm64" || _ARCH="arm" ;;
        *)           echo "Unknown architecture: $(uname -m)" ; exit 1 ;;
    esac

    DIR="${DIR:-$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )}"

    if ! which olm-bundle > /dev/null; then
        [ -f ${DIR}/olm-bundle ] || dlForked
        OLM_BINARY="${DIR}/olm-bundle"
    else
        OLM_BINARY="olm-bundle"
        # OLM_BINARY="/Users/ab017z6/workspace/olm-bundle/_output/bin/darwin_amd64/olm-bundle"
    fi

    generate
}

generate() {
    cd ${DIR}/../chart/k8gb && helm dependency update && cd -
    helm -n placeholder template ${DIR}/../chart/k8gb \
        --name-template=k8gb \
        --set k8gb.securityContext.runAsUser=null \
        --set k8gb.log.format=simple \
        --set k8gb.log.level=info | ${OLM_BINARY} \
            --chart-file-path=${DIR}/../chart/k8gb/Chart.yaml \
            --version=${_VERSION} \
            --output-dir ${DIR}
}

dlUpstream() {
    # upstream
    _VERSION="0.5.2"
    curl -Lo ${DIR}/olm-bundle https://github.com/upbound/olm-bundle/releases/download/v${_VERSION}/olm-bundle_${_OS}-${_ARCH}
    chmod +x ${DIR}/olm-bundle
}

dlForked() {
    # our fork
    curl -Ls https://github.com/AbsaOSS/olm-bundle/releases/download/v${_VERSION}/olm-bundle_${_VERSION}_${_OS}_${_ARCH}.tar.gz | tar -xz
    mv ${DIR}/bin/olm-bundle ${DIR}/olm-bundle
}

main $@
