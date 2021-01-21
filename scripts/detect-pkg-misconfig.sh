#!/bin/bash
echo "*************************************************"
echo "*** Check that go.mod is configured correctly ***"
echo "*************************************************"

cp ../go.mod ../go.mod.temp
bash ./use-remote-repos.sh
go mod tidy

if cmp -s ../go.mod ../go.mod.temp; then
    echo go.mod file correctly configured.
else 
    echo ERROR: go.mod file not configured correctly. 
    echo go.mod is:
    cat ../go.mod.temp
    echo go mod should be
    cat ../go.mod
    echo
    echo Please execute:
    echo   bash use-remote-repos.sh
    exit 1
fi
