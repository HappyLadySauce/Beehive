#!/bin/sh
set -e

# Install etcdctl if not present
if ! command -v etcdctl >/dev/null 2>&1; then
  echo "Installing etcdctl..."
  ETCD_VERSION="3.5.9"
  ARCH="amd64"
  wget -q "https://github.com/etcd-io/etcd/releases/download/v${ETCD_VERSION}/etcd-v${ETCD_VERSION}-linux-${ARCH}.tar.gz" -O /tmp/etcd.tar.gz
  tar -xzf /tmp/etcd.tar.gz -C /tmp
  cp "/tmp/etcd-v${ETCD_VERSION}-linux-${ARCH}/etcdctl" /usr/local/bin/etcdctl
  chmod +x /usr/local/bin/etcdctl
  rm -rf /tmp/etcd*
fi

echo "Waiting for etcd cluster to be ready..."
sleep 10

ENDPOINTS="${ETCD_ENDPOINTS:-http://etcd-0:2379,http://etcd-1:2379,http://etcd-2:2379}"
USER="${ETCD_USER:-Beehive}"
PASSWORD="${ETCD_PASSWORD:-Beehive}"

echo "Checking if authentication is already enabled..."
if etcdctl --endpoints="${ENDPOINTS}" auth status 2>/dev/null | grep -q 'Authentication Status: true'; then
  echo "Authentication already enabled, skipping initialization."
  exit 0
fi

# etcd requires a 'root' user to exist before enabling authentication
# Create root user first (use the same password as the main user)
echo "Creating root user (required for enabling auth)..."
if etcdctl --endpoints="${ENDPOINTS}" user add "root:${PASSWORD}" 2>/dev/null; then
  echo "Root user created successfully"
else
  echo "Root user already exists, skipping..."
fi

echo "Granting root role to root user..."
if etcdctl --endpoints="${ENDPOINTS}" user grant-role root root 2>/dev/null; then
  echo "Root role granted successfully"
else
  echo "Root role already granted, skipping..."
fi

# Enable authentication (requires root user to exist)
echo "Enabling authentication..."
if etcdctl --endpoints="${ENDPOINTS}" auth enable 2>/dev/null; then
  echo "Authentication enabled successfully"
else
  echo "Authentication already enabled, skipping..."
fi

# If the specified user is not 'root', create it after enabling auth
if [ "${USER}" != "root" ]; then
  echo "Creating user: ${USER} (after enabling auth)..."
  if etcdctl --endpoints="${ENDPOINTS}" --user="root:${PASSWORD}" user add "${USER}:${PASSWORD}" 2>/dev/null; then
    echo "User ${USER} created successfully"
  else
    echo "User ${USER} already exists, skipping..."
  fi
  
  echo "Granting root role to user: ${USER}..."
  if etcdctl --endpoints="${ENDPOINTS}" --user="root:${PASSWORD}" user grant-role "${USER}" root 2>/dev/null; then
    echo "Root role granted to ${USER} successfully"
  else
    echo "Root role already granted to ${USER}, skipping..."
  fi
fi

echo "Authentication setup completed!"
echo "Username: ${USER}"
echo "Password: ${PASSWORD}"
echo "Note: Root user also exists with the same password for administrative purposes."
