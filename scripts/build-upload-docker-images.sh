#!/bin/bash

set -euxf -o pipefail

make create-baseimg-debugimg

# build multi-arch binaries
make build-binaries-linux
make build-binaries-s390x
make build-binaries-ppc64le
make build-binaries-arm64

# build multi-arch docker images
platforms="linux/amd64,linux/s390x,linux/ppc64le,linux/arm64"

# build/upload images for release version of Jaeger backend components
for component in agent collector query ingester
do
	bash scripts/build-upload-a-docker-image.sh -b -c "jaeger-${component}" -d "cmd/${component}" -p "${platforms}" -t release
done

# build/upload images for jaeger-es-index-cleaner and jaeger-es-rollover
bash scripts/build-upload-a-docker-image.sh -b -c jaeger-es-index-cleaner -d cmd/es-index-cleaner -p "${platforms}" -t release
bash scripts/build-upload-a-docker-image.sh -b -c jaeger-es-rollover -d cmd/es-rollover  -p "${platforms}" -t release

# build/upload images for jaeger-tracegen and jaeger-anonymizer
for component in tracegen anonymizer
do
	bash scripts/build-upload-a-docker-image.sh -c "jaeger-${component}" -d "cmd/${component}" -p "${platforms}"
done 


# build amd64 docker images

# build/upload images for debug version of Jaeger backend components
for component in agent collector query ingester
do
	bash scripts/build-upload-a-docker-image.sh -b -c "jaeger-${component}-debug" -d "cmd/${component}" -t debug
done

# build/upload images for jaeger-cassandra-schema
bash scripts/build-upload-a-docker-image.sh -c jaeger-cassandra-schema -d plugin/storage/cassandra/
