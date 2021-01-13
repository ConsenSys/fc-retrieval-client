#!/bin/bash
# set -e

echo "**************************************************************************************"
echo "*** Update go.mod and go.sum to point to the latest gateway packages on github.com ***"
echo "**************************************************************************************"

# Remove any local references
sed '/replace .*/d' ../go.mod > ../go.mod.new
rm ../go.mod
mv ../go.mod.new ../go.mod


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
GATEWAY_BRANCH=`git rev-parse --abbrev-ref HEAD`
echo "Gateway repo branch: $GATEWAY_BRANCH"
BRANCH_EXISTS_ON_GATEWAY=`git branch --list | grep $CLIENT_BRANCH`
#echo "BRANCH_EXISTS_ON_GATEWAY: $BRANCH_EXISTS_ON_GATEWAY"

cd $CLIENT_DIR
if [ $CLIENT_BRANCH != $GATEWAY_BRANCH ]; then 
    echo "Client and Gateway branch do not match"
    if [ -n "$BRANCH_EXISTS_ON_GATEWAY" ]; then
        echo "ERROR: Branch $CLIENT_BRANCH exists on the gateway repo, but the gateway repo is currently using branch $GATEWAY_BRANCH"
        exit 1
    fi

    echo "Calling go get to use main on fc-retrieval-gateway"
    go get github.com/ConsenSys/fc-retrieval-gateway@main
else
    echo Client and Gateway branch match
    echo "Calling go get to use $GATEWAY_BRANCH on fc-retrieval-gateway"
    go get github.com/ConsenSys/fc-retrieval-gateway@$GATEWAY_BRANCH
fi

