#!/bin/bash

gopass sync -s caos-secrets
gopass caos-secrets/technical/boom/git-read-secret | kubectl apply -f -

skaffold run -f build/skaffold/skaffold.yaml