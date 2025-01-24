#!/usr/bin/env bash
set -eu -o pipefail

for port in $@; do
  # The /dev/tcp/localhost/$port pseudo-device is a Bashism that opens a connection to localhost:$port.
  # Docs: https://tldp.org/LDP/abs/html/devref1.html
  #
  # We're not using "lsof -i :$port" because that only checks if something is
  # listening on the port, not whether it can accept connections. Also, when
  # using services of type "LoadBalancer", the load balancer uses iptables rules
  # to forward ports, so they won't show up in "lsof".
  until (: < "/dev/tcp/localhost/$port") 2>/dev/null; do
    echo "waiting for port $port to be accessible"
    sleep 1
  done
  echo "port $port accessible"
done
