#!/bin/bash
echo "****************************************************************"
echo "*** Check that go.mod does not point to a local gateway repo ***"
echo "****************************************************************"

REPLACE_EXISTS=`grep "replace " ../go.mod`
if [ -n "$REPLACE_EXISTS" ]; then
    echo "Using local gateway repo"
    echo "Replace set to: $REPLACE_EXISTS"
    exit 1
fi
echo "Using remote gateway repo"
