#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../../.." && pwd)"
SPEC_REL_PATH="packages/contracts/openapi/api-server.yaml"
JS_OUT_REL_PATH="packages/contracts/sdk-js/api-server"
DART_OUT_REL_PATH="packages/contracts/sdk-dart/api_server"
GENERATOR_IMAGE="${OPENAPI_GENERATOR_IMAGE:-openapitools/openapi-generator-cli:v7.16.0}"

if ! command -v docker >/dev/null 2>&1; then
  echo "未检测到 docker，无法执行 SDK 生成。"
  echo "请先安装并启动 Docker，或在环境变量 OPENAPI_GENERATOR_IMAGE 中指定可用镜像。"
  exit 1
fi

echo ">>> 清理旧产物"
rm -rf "${ROOT_DIR}/${JS_OUT_REL_PATH}" "${ROOT_DIR}/${DART_OUT_REL_PATH}"

echo ">>> 生成 JS SDK (typescript-axios)"
docker run --rm \
  -u "$(id -u):$(id -g)" \
  -v "${ROOT_DIR}:/local" \
  "${GENERATOR_IMAGE}" generate \
  -i "/local/${SPEC_REL_PATH}" \
  -g typescript-axios \
  -o "/local/${JS_OUT_REL_PATH}" \
  --additional-properties=npmName=@weouc/api-server-sdk,npmVersion=0.1.0,supportsES6=true,withSeparateModelsAndApi=true,apiPackage=apis,modelPackage=models,sortParamsByRequiredFlag=true,sortModelPropertiesByRequiredFlag=true

echo ">>> 生成 Dart SDK (dart-dio)"
docker run --rm \
  -u "$(id -u):$(id -g)" \
  -v "${ROOT_DIR}:/local" \
  "${GENERATOR_IMAGE}" generate \
  -i "/local/${SPEC_REL_PATH}" \
  -g dart-dio \
  -o "/local/${DART_OUT_REL_PATH}" \
  --additional-properties=pubName=weouc_api_server_sdk,pubVersion=0.1.0,sortParamsByRequiredFlag=true,sortModelPropertiesByRequiredFlag=true

echo ">>> SDK 生成完成"
echo "JS:   ${JS_OUT_REL_PATH}"
echo "Dart: ${DART_OUT_REL_PATH}"
