#!/bin/sh
set -e

# Bedrock installer
# Usage: curl -sSfL https://raw.githubusercontent.com/XRPL-Commons/Bedrock/main/install.sh | sh

REPO="XRPL-Commons/Bedrock"
INSTALL_DIR="/usr/local/bin"
BINARY_NAME="bedrock"

main() {
  need_cmd curl
  need_cmd uname
  need_cmd tar

  get_platform
  get_latest_version

  echo "Installing Bedrock ${VERSION} for ${OS}/${ARCH}..."

  TMPDIR=$(mktemp -d)
  trap "rm -rf ${TMPDIR}" EXIT

  ARCHIVE_NAME="bedrock_${VERSION}_${OS}_${ARCH}.tar.gz"
  DOWNLOAD_URL="https://github.com/${REPO}/releases/download/${VERSION}/${ARCHIVE_NAME}"

  echo "Downloading ${DOWNLOAD_URL}..."
  curl -sSfL "${DOWNLOAD_URL}" -o "${TMPDIR}/${ARCHIVE_NAME}" || {
    echo ""
    echo "Error: Failed to download release binary."
    echo ""
    echo "This could mean:"
    echo "  - No release has been published yet"
    echo "  - No binary is available for ${OS}/${ARCH}"
    echo ""
    echo "You can install from source instead:"
    echo "  git clone https://github.com/${REPO}.git"
    echo "  cd Bedrock"
    echo "  go build -o bedrock cmd/bedrock/main.go"
    echo "  sudo mv bedrock /usr/local/bin/"
    exit 1
  }

  tar -xzf "${TMPDIR}/${ARCHIVE_NAME}" -C "${TMPDIR}"

  if [ -w "${INSTALL_DIR}" ]; then
    mv "${TMPDIR}/${BINARY_NAME}" "${INSTALL_DIR}/${BINARY_NAME}"
  else
    echo "Installing to ${INSTALL_DIR} (requires sudo)..."
    sudo mv "${TMPDIR}/${BINARY_NAME}" "${INSTALL_DIR}/${BINARY_NAME}"
  fi

  chmod +x "${INSTALL_DIR}/${BINARY_NAME}"

  echo ""
  echo "Bedrock ${VERSION} installed successfully to ${INSTALL_DIR}/${BINARY_NAME}"
  echo ""
  echo "Run 'bedrock --help' to get started."
}

get_platform() {
  OS=$(uname -s | tr '[:upper:]' '[:lower:]')
  ARCH=$(uname -m)

  case "${OS}" in
    linux)  OS="linux" ;;
    darwin) OS="darwin" ;;
    *)
      echo "Error: Unsupported operating system: ${OS}"
      exit 1
      ;;
  esac

  case "${ARCH}" in
    x86_64|amd64)  ARCH="amd64" ;;
    arm64|aarch64) ARCH="arm64" ;;
    *)
      echo "Error: Unsupported architecture: ${ARCH}"
      exit 1
      ;;
  esac
}

get_latest_version() {
  VERSION=$(curl -sSfL "https://api.github.com/repos/${REPO}/releases/latest" \
    | grep '"tag_name"' \
    | head -1 \
    | sed 's/.*"tag_name": *"//;s/".*//')

  if [ -z "${VERSION}" ]; then
    echo "Error: Could not determine latest version."
    echo "Check https://github.com/${REPO}/releases for available releases."
    exit 1
  fi
}

need_cmd() {
  if ! command -v "$1" > /dev/null 2>&1; then
    echo "Error: required command '$1' not found."
    exit 1
  fi
}

main
