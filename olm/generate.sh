#!/bin/bash

DIR="${DIR:-$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )}"

helm -n placeholder template ${DIR}/../chart/k8gb \
              --name-template=k8gb \
              --set k8gb.securityContext.runAsUser=null \
              --set k8gb.log.format=simple \
              --set k8gb.log.level=info | olm-bundle --chart-file-path=${DIR}/../chart/k8gb/Chart.yaml --output-dir ${DIR}
