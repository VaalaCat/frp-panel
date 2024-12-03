#!/usr/bin/env /bin/bash

# Function to print usage
usage() {
    echo "Usage: $0 [--platform <platform>] [--bintype <bintype>] [--arch <arch>]"
    echo "Platforms: windows, linux, darwin, all"
    echo "Binary Types: full, client, all"
    echo "Architectures: amd64, arm64, arm, all"
    echo "Example: $0 --platform linux --bintype full --arch amd64"
    exit 1
}

# Default values
PLATFORM="all"
BINTYPE="all"
ARCH="all"

# build variables
BUILD_DATE="$(date -u +'%Y-%m-%dT%H:%M:%SZ')"
GIT_COMMIT="$(git rev-parse HEAD)"
VERSION="$(git describe --tags --abbrev=0 | tr -d '\n')"

# Parse arguments
while [[ "$#" -gt 0 ]]; do
    case $1 in
        --platform) PLATFORM="$2"; shift ;;
        --bintype) BINTYPE="$2"; shift ;;
        --arch) ARCH="$2"; shift ;;
        --skip-frontend) SKIP_FRONTEND=true ;;
        --current) CURRENT=true ;;
        *) usage ;;
    esac
    shift
done

if [[ "$CURRENT" == "true" ]]; then
    PLATFORM=$(uname -s | tr '[:upper:]' '[:lower:]')
    ARCH=$(uname -m)
    if [[ "$ARCH" == "x86_64" ]]; then
        ARCH="amd64"
    fi
    if [[ "$ARCH" == "aarch64" ]]; then
        ARCH="arm64"
    fi
    BINTYPE="full"
fi

echo "Building for platform: $PLATFORM, binary type: $BINTYPE, architecture: $ARCH"
echo "Build Date: $BUILD_DATE"
echo "Git Commit: $GIT_COMMIT"
echo "Version: $VERSION"

BUILD_LD_FLAGS="-X 'github.com/VaalaCat/frp-panel/conf.buildDate=${BUILD_DATE}' -X 'github.com/VaalaCat/frp-panel/conf.gitCommit=${GIT_COMMIT}' -X 'github.com/VaalaCat/frp-panel/conf.gitVersion=${VERSION}'"

if [[ "$SKIP_FRONTEND" == "true" ]]; then
    echo "Skipping frontend build"
else
    echo "Building frontend"
    # Prepare build environment
    mkdir -p dist
    rm -rf dist/*

    # Build frontend
    cd www && pnpm install --no-frozen-lockfile && pnpm build && cd ..
fi

# Build function
build_binary() {
    local platform=$1
    local arch=$2
    local bintype=$3
    local output_name=""
    local source_path=""

    # Determine output name and source path
    if [[ "$bintype" == "full" ]]; then
        source_path="cmd/frpp/*.go"
        output_name="frp-panel"
    elif [[ "$bintype" == "client" ]]; then
        source_path="cmd/frppc/*.go"
        output_name="frp-panel-client"
    else
        echo "Invalid binary type"
        return 1
    fi

    # Set executable extension for Windows
    local exe_ext=""
    if [[ "$platform" == "windows" ]]; then
        exe_ext=".exe"
    fi

    # Special handling for ARM architectures
    local goarch="$arch"
    local goarm=""
    if [[ "$arch" == "arm" ]]; then
        goarch="arm"
        if [[ "$platform" == "linux" ]]; then
            # Build for ARMv7 and ARMv6
            for arm_version in 7 6; do
                local arm_output="${output_name}-${platform}-armv${arm_version}l${exe_ext}"
                CGO_ENABLED=0 GOOS="$platform" GOARCH="$goarch" GOARM="$arm_version" \
                go build -o "dist/${arm_output}" -ldflags "$BUILD_LD_FLAGS" $source_path
            done
            return 0
        fi
    fi

    # Standard build
    local output="${output_name}-${platform}-${arch}${exe_ext}"
    CGO_ENABLED=0 GOOS="$platform" GOARCH="$goarch" \
    go build -o "dist/${output}" -ldflags "$BUILD_LD_FLAGS" $source_path
}

# Platforms array
PLATFORMS=()
if [[ "$PLATFORM" == "all" ]]; then
    PLATFORMS=("windows" "linux" "darwin")
else
    PLATFORMS=("$PLATFORM")
fi

# Architectures array
ARCHS=()
if [[ "$ARCH" == "all" ]]; then
    ARCHS=("amd64" "arm64" "arm")
else
    ARCHS=("$ARCH")
fi

# Binary types array
BINTYPES=()
if [[ "$BINTYPE" == "all" ]]; then
    BINTYPES=("full" "client")
else
    BINTYPES=("$BINTYPE")
fi

# Build matrix
for platform in "${PLATFORMS[@]}"; do
    for arch in "${ARCHS[@]}"; do
        for bintype in "${BINTYPES[@]}"; do
            echo "Building $bintype binary for $platform-$arch"
            if [[ "$platform" == "darwin" && "$arch" == "arm" ]]; then continue; fi
            if [[ "$platform" == "windows" && "$arch" == "arm" ]]; then continue; fi
            build_binary "$platform" "$arch" "$bintype"
        done
    done
done

# Move to current directory if current enabled
if [[ "$CURRENT" == "true" ]]; then
    mv dist/frp* ./frp-panel
fi

echo "Build Done!"