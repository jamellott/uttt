#!/bin/sh

rm -rf dist
mkdir -p dist

set -e

echo "Building server..."
cd server
go build -o ../dist/server

echo "Building UI..."
cd ../ui
NODE_ENV=production npm run build

echo "Copying files..."
cp -R dist/ ../dist/ui

cat <<EOF > ../dist/config.yaml
port: 80
host:
acmetls: false
requestlogs: false
checkorigin: true
EOF

echo "Done"