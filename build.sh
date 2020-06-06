#!/bin/sh

mkdir -p dist

set -e

echo "Building server..."
cd server
go build -o ../dist/server

echo "Building UI..."
cd ../ui
npm run build

echo "Copying files..."
cp -R dist/ ../dist/ui

echo "Done"