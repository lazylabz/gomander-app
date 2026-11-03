#!/bin/sh

ARCH=$(dpkg --print-architecture)

rm -rf cmd/gomander/frontend/node_modules && \
cd cmd/gomander && \
wails3 task linux:build && \
cd ../.. && \
cp "build/bin/gomander" "/app/output/gomander-linux-${ARCH}" && \
bash scripts/create-deb.sh "build/bin/gomander" "${ARCH}" && \
cp "build/bin/gomander_"*"_${ARCH}.deb" "/app/output/"
