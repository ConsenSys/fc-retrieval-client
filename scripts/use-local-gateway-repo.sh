#!/bin/bash
set -e

echo "********************************************************"
echo "*** Update go.mod to point to the local gateway repo ***"
echo "********************************************************"

REPLACE_TEXT="replace github.com/ConsenSys/fc-retrieval-gateway => ../fc-retrieval-gateway"

sed '/replace .*/d' ../go.mod > ../go.mod.new
rm ../go.mod
mv ../go.mod.new ../go.mod
echo $REPLACE_TEXT >> ../go.mod