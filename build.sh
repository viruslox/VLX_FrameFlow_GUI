#!/bin/bash
set -e

echo "Building Frontend..."
cd frontend
npm install
npm run build
cd ..

echo "Copying Frontend Build to Backend..."
rm -rf backend/ui/dist
mkdir -p backend/ui/dist
cp -r frontend/dist/* backend/ui/dist/

echo "Building Backend..."
cd backend
export CGO_ENABLED=0

# Build for linux/amd64
echo "Compiling for linux/amd64..."
GOOS=linux GOARCH=amd64 go build -o ../bin/vlx_frameflow_gui-amd64 cmd/controller/*.go

# Build for linux/arm64 (e.g. Raspberry Pi)
echo "Compiling for linux/arm64..."
GOOS=linux GOARCH=arm64 go build -o ../bin/vlx_frameflow_gui-arm64 cmd/controller/*.go

cd ..

echo "Build complete! Binaries are in the 'bin/' directory."
