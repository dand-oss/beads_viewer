#!/bin/bash
set -e

# Build and install bv via Makefile
cd "$(dirname "$0")"
make install
