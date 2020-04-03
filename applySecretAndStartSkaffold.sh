#!/bin/bash

gopass sync -s caos-secrets
./scripts/0_imagepull-secret.sh | kubectl apply -f -
./scripts/1_argocd-secrets.sh | kubectl apply -f -
./scripts/2_ops-repo-read-secret.sh | kubectl apply -f -

skaffold run -f ./build/skaffold/skaffold.yaml