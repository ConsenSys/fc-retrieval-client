#!/bin/bash
echo "*************************************************"
echo "*** Check that go.mod is configured correctly ***"
echo "*************************************************"

# Check that it isn't set to using a local checkout.
REPLACE_EXISTS=`grep "replace " ../go.mod`
if [ -n "$REPLACE_EXISTS" ]; then
    echo "ERROR: Using local gateway repo"
    echo "Replace set to: $REPLACE_EXISTS"
    exit 1
fi
echo "Using remote gateway repo"

GATEWAY_DIR="../../fc-retrieval-gateway"
CLIENT_DIR="../fc-retrieval-client"
if [ ! -d "$GATEWAY_DIR" ]; then
    echo "ERROR: Gateway directory $GATEWAY_DIR does not exist"
    exit 1
fi
echo "Found Gateway repo: $GATEWAY_DIR"

CLIENT_BRANCH=`git rev-parse --abbrev-ref HEAD`
echo "Client repo branch: $CLIENT_BRANCH"
cd $GATEWAY_DIR
BRANCH_EXISTS_ON_GATEWAY=`git branch -r --list origin/$CLIENT_BRANCH`
#echo "BRANCH_EXISTS_ON_GATEWAY: $BRANCH_EXISTS_ON_GATEWAY"

echo HERE1
pwd
echo HERE2
git branch -r --list origin/$CLIENT_BRANCH
echo HERE3

cd $CLIENT_DIR
GATEWAY_BRANCH_TO_USE=$CLIENT_BRANCH
if [ -z "$BRANCH_EXISTS_ON_GATEWAY" ]; then
    echo "Branch $CLIENT_BRANCH does not exist on gateway"
    GATEWAY_BRANCH_TO_USE=main
else
    echo "Branch $CLIENT_BRANCH does exist on gateway"
fi

cp go.mod go.mod.temp
go get github.com/ConsenSys/fc-retrieval-gateway@$GATEWAY_BRANCH_TO_USE
go mod tidy
if cmp -s go.mod go.mod.temp; then
    echo go.mod file correctly configured.
else 
    echo ERROR: go.mod file not configured correctly. 
    echo Please execute:
    echo   go get github.com/ConsenSys/fc-retrieval-gateway@$GATEWAY_BRANCH_TO_USE
    echo   go mod tidy
    exit 1
fi

