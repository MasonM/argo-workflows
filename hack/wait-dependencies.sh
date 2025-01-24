#!/usr/bin/env bash
set -eu -o pipefail

apps=$(kubectl get deploy -o=jsonpath='{.items[?(@.spec.replicas>0)].metadata.name}')
echo "waiting for deployments: $apps"
kubectl wait --timeout 2m --for=condition=Available deploy $apps

ports=$(kubectl get svc -o=jsonpath='{.items[?(@.spec.type=="LoadBalancer")]..port}')
echo "waiting for ports: $ports"
./hack/wait-ports.sh "$ports"
