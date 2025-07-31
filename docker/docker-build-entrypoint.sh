#!/bin/sh

rm -rf cmd/gomander/frontend/node_modules && \
cd cmd/gomander && \
wails build -tags webkit2_41  && \
cd ../.. && \
cp "build/bin/gomander" "/app/output/gomander-linux-$(dpkg --print-architecture)"
