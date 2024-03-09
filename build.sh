#!/bin/bash
VERSION=""
IMAGE_NAME="usercore"
usage() {
    echo "Usage: $0 -t <version>"
    exit 1
}
while getopts ":t:" opt; do
  case ${opt} in
    t )
      VERSION=$OPTARG
      ;;
    \? )
      usage
      ;;
  esac
done
shift $((OPTIND -1))
if [ -z "$VERSION" ]; then
    echo -e "\n\n"
    echo -e "\033[1;31mError: Version tag is required.\033[0m"
    usage
fi
echo -e "\033[1m\033[34mBuilding Docker image with tag: $VERSION\033[0m\n"
docker build -t "$IMAGE_NAME":"$VERSION" .
echo -e "\n\033[1m\033[30m\033[42mBuild and Docker image creation completed.\033[0m\n"

FULL_IMAGE_HASH=$(docker inspect $IMAGE_NAME:"$VERSION" --format='{{.Id}}')
# Extract only the hash part after 'sha256:'
IMAGE_HASH=${FULL_IMAGE_HASH#"sha256:"}
echo -e "Docker image hash: \033[1m\033[30m\033[44m$IMAGE_HASH\033[0m\n"