#!/usr/bin/env bash
set -eou pipefail

GOPKG="github.com/szuecs/routegroup-client"
SCRIPT_ROOT="$(dirname "${BASH_SOURCE[0]}")/.."

rm -rf "${SCRIPT_ROOT}/generated"

go run k8s.io/code-generator/cmd/deepcopy-gen --input-dirs ${GOPKG}/apis/zalando.org/v1 \
  -O zz_generated.deepcopy \
  --bounding-dirs ${GOPKG}/apis \
  --go-header-file "${SCRIPT_ROOT}/hack/boilerplate.go.txt" \
  -o "${SCRIPT_ROOT}/generated"

go run k8s.io/code-generator/cmd/client-gen --clientset-name versioned \
  --input-base '' \
  --input ${GOPKG}/apis/zalando.org/v1 \
  --output-package ${GOPKG}/client/clientset \
  --go-header-file "${SCRIPT_ROOT}/hack/boilerplate.go.txt" \
  -o "${SCRIPT_ROOT}/generated"

cp -rv "${SCRIPT_ROOT}/generated/${GOPKG}"/* .
rm -rf "${SCRIPT_ROOT}/generated"
