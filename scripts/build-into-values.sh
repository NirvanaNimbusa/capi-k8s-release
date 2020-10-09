#!/usr/bin/env bash

set -ex

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
REPO_BASE_DIR="${SCRIPT_DIR}/.."

# set defaults - any of these should be configrable from outside this script
IMAGE_DESTINATION_CCNG="${IMAGE_DESTINATION_CCNG:-gcr.io/cf-capi-arya/dev-ccng}"
IMAGE_DESTINATION_CF_API_CONTROLLERS="${IMAGE_DESTINATION_CF_API_CONTROLLERS:-gcr.io/cf-capi-arya/dev-controllers}"
IMAGE_DESTINATION_PACKAGE_IMAGE_UPLOADER="${IMAGE_DESTINATION_PACKAGE_IMAGE_UPLOADER:-gcr.io/cf-capi-arya/dev-package-image-uploader}"
CF_FOR_K8s_DIR="${CF_FOR_K8s_DIR:-${HOME}/workspace/cf-for-k8s/}"
CAPI_RELEASE_DIR="${CAPI_RELEASE_DIR:-${HOME}/workspace/capi-release/}"
CCNG_DIR="${CAPI_RELEASE_DIR}src/cloud_controller_ng"
CF_API_CONTROLLERS_DIR="${CAPI_CF_API_CONTROLLERS_DIR:-${REPO_BASE_DIR}/src/cf-api-controllers}"
PACKAGE_IMAGE_UPLOADER_DIR="${PACKAGE_IMAGE_UPLOADER_DIR:-${REPO_BASE_DIR}/src/package-image-uploader}"

if [ -z "$1" ]
  then
    echo "No values file was supplied!"
    exit 1
fi

# template the image destination into the kbld yml
KBLD_TMP="$(mktemp)"
ytt -f "${REPO_BASE_DIR}/dev-templates/" \
    -v "kbld.destination.ccng=${IMAGE_DESTINATION_CCNG}" \
    -v "kbld.destination.cf_api_controllers=${IMAGE_DESTINATION_CF_API_CONTROLLERS}" \
    -v "kbld.destination.package_image_uploader=${IMAGE_DESTINATION_PACKAGE_IMAGE_UPLOADER}" \
    -v "src_dirs.ccng=${CCNG_DIR}" \
    -v "src_dirs.cf_api_controllers=${CF_API_CONTROLLERS_DIR}" \
    -v "src_dirs.package_image_uploader=${PACKAGE_IMAGE_UPLOADER_DIR}" \
     > "${KBLD_TMP}"

# build a new values file with kbld
pushd "${REPO_BASE_DIR}"
  VALUES_TMP="$(mktemp)"
  echo "#@data/values" > "${VALUES_TMP}"
  kbld --images-annotation=false -f "${KBLD_TMP}" -f "$1" >> "${VALUES_TMP}"
popd

cat "${VALUES_TMP}" > "$1"
