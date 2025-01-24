#!/usr/bin/env sh

# Install helper and minimal init script for k3s, the Kubernetes distribution we
# use for CI builds (via GitHub Actions) and the devcontainer.
#
# Not meant to be run directly on a developer's host machine.

PIDFILE=/var/run/k3s.pid

set -eux

case "$1" in
  install)
    if ! echo "${INSTALL_K3S_VERSION:-}" | egrep '^v[0-9]+\.[0-9]+\.[0-9]+\+k3s1$'; then
      export INSTALL_K3S_VERSION=v1.31.2+k3s1
    fi
    curl -sfL https://get.k3s.io | INSTALL_K3S_SKIP_ENABLE=true sh -
    ;;
  start)
    # The k3s installer automatically creates systemd unit files, but the devcontainer doesn't have systemd,
    # so we can't use "systemctl start k3s". Instead, we can use
    # start-stop-daemon (which comes standard with Debian/Ubuntu) to start in the same way as the unit files.
    if sudo start-stop-daemon \
      --start \
      --verbose \
      --background \
      --notify-await \
      --pidfile "$PIDFILE" \
      --make-pidfile \
      --startas /usr/local/bin/k3s \
      -- server --docker --log /var/log/k3s.log --write-kubeconfig-mode=644 "${INSTALL_K3S_EXEC:-}"; then
      until kubectl --kubeconfig=/etc/rancher/k3s/k3s.yaml cluster-info ; do sleep 1s ; done
      [ -f "$KUBECONFIG" ] || cp /etc/rancher/k3s/k3s.yaml "$KUBECONFIG"
    fi
    ;;
  status|stop)
    sudo start-stop-daemon --"$1" --pidfile "$PIDFILE"
    ;;
esac
