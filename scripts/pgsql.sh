#!/bin/bash

# PostgreSQL Installation Script for Ubuntu
# This script installs PostgreSQL and configures it to be accessible from a Kubernetes cluster

set -e  # Exit on error

# Color output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Configuration variables (can be overridden via environment variables)
PG_VERSION=${PG_VERSION:-16}
PG_DB_NAME=${PG_DB_NAME:-zewi}
PG_USER=${PG_USER:-zewi}
PG_PASSWORD=${PG_PASSWORD:-zewi_pwd_ww3i3y}
PG_PORT=${PG_PORT:-5432}

# Functions
log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if running as root
if [ "$EUID" -ne 0 ]; then
    log_error "Please run as root (use sudo)"
    exit 1
fi

log_info "Starting PostgreSQL installation and configuration..."

# Update package list
log_info "Updating package list..."
apt-get update -qq

# Install PostgreSQL
log_info "Installing PostgreSQL ${PG_VERSION}..."
if ! command -v psql &> /dev/null; then
    apt-get install -y postgresql-${PG_VERSION} postgresql-contrib-${PG_VERSION}
else
    log_warn "PostgreSQL is already installed"
fi

# Start and enable PostgreSQL service
log_info "Starting PostgreSQL service..."
systemctl start postgresql
systemctl enable postgresql

# Get PostgreSQL data directory
PG_DATA_DIR=$(sudo -u postgres psql -t -P format=unaligned -c 'SHOW data_directory;' 2>/dev/null || echo "/var/lib/postgresql/${PG_VERSION}/main")
PG_CONFIG_DIR="/etc/postgresql/${PG_VERSION}/main"

log_info "PostgreSQL data directory: ${PG_DATA_DIR}"
log_info "PostgreSQL config directory: ${PG_CONFIG_DIR}"

# Configure PostgreSQL to listen on all interfaces
log_info "Configuring PostgreSQL to accept remote connections..."

# Backup original postgresql.conf
if [ ! -f "${PG_CONFIG_DIR}/postgresql.conf.backup" ]; then
    cp "${PG_CONFIG_DIR}/postgresql.conf" "${PG_CONFIG_DIR}/postgresql.conf.backup"
fi

# Update postgresql.conf
sed -i "s/#listen_addresses = 'localhost'/listen_addresses = '*'/" "${PG_CONFIG_DIR}/postgresql.conf"
sed -i "s/#port = 5432/port = ${PG_PORT}/" "${PG_CONFIG_DIR}/postgresql.conf"

# Configure pg_hba.conf to allow connections from Kubernetes cluster
log_info "Configuring pg_hba.conf for remote access..."

# Backup original pg_hba.conf
if [ ! -f "${PG_CONFIG_DIR}/pg_hba.conf.backup" ]; then
    cp "${PG_CONFIG_DIR}/pg_hba.conf" "${PG_CONFIG_DIR}/pg_hba.conf.backup"
fi

# Add entries to pg_hba.conf
cat >> "${PG_CONFIG_DIR}/pg_hba.conf" << EOF

# Allow remote connections (configure firewall rules for security)
host    all             all             0.0.0.0/0               md5
EOF

# Create database and user
log_info "Creating database and user..."

sudo -u postgres psql << EOF
-- Create user
DO \$\$
BEGIN
    IF NOT EXISTS (SELECT FROM pg_catalog.pg_roles WHERE rolname = '${PG_USER}') THEN
        CREATE USER ${PG_USER} WITH PASSWORD '${PG_PASSWORD}';
    ELSE
        ALTER USER ${PG_USER} WITH PASSWORD '${PG_PASSWORD}';
    END IF;
END
\$\$;

-- Create database
SELECT 'CREATE DATABASE ${PG_DB_NAME} OWNER ${PG_USER}'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = '${PG_DB_NAME}')\gexec

-- Grant privileges
GRANT ALL PRIVILEGES ON DATABASE ${PG_DB_NAME} TO ${PG_USER};
\c ${PG_DB_NAME}
GRANT ALL ON SCHEMA public TO ${PG_USER};
EOF

# Configure firewall (UFW)
log_info "Configuring firewall rules..."

if command -v ufw &> /dev/null; then
    # Check if UFW is active
    if ufw status | grep -q "Status: active"; then
        log_info "UFW is active, adding PostgreSQL port rule..."
        ufw allow ${PG_PORT}/tcp comment 'PostgreSQL'
    else
        log_warn "UFW is not active. Please configure firewall manually if needed."
    fi
else
    log_warn "UFW is not installed. Please configure firewall manually if needed."
fi

# Restart PostgreSQL to apply changes
log_info "Restarting PostgreSQL to apply configuration changes..."
systemctl restart postgresql

# Wait for PostgreSQL to be ready
log_info "Waiting for PostgreSQL to be ready..."
sleep 3

# Test connection
log_info "Testing PostgreSQL connection..."
if sudo -u postgres psql -c "SELECT version();" > /dev/null 2>&1; then
    log_info "PostgreSQL is running successfully!"
else
    log_error "PostgreSQL connection test failed!"
    exit 1
fi

# Display connection information
log_info "=========================================="
log_info "PostgreSQL Installation Complete!"
log_info "=========================================="
echo ""
log_info "Connection Details:"
echo "  Host: $(hostname -I | awk '{print $1}')"
echo "  Port: ${PG_PORT}"
echo "  Database: ${PG_DB_NAME}"
echo "  Username: ${PG_USER}"
echo "  Password: ${PG_PASSWORD}"
echo ""
log_info "Connection String:"
echo "  postgresql://${PG_USER}:${PG_PASSWORD}@$(hostname -I | awk '{print $1}'):${PG_PORT}/${PG_DB_NAME}"
echo ""
log_warn "IMPORTANT: Save the password above securely!"
log_warn "Update your Kubernetes ConfigMaps/Secrets with these connection details."
echo ""
log_info "To connect from Kubernetes, use these environment variables:"
echo "  ZEWI_DATABASE_DRIVER=postgres"
echo "  ZEWI_DATABASE_HOST=$(hostname -I | awk '{print $1}')"
echo "  ZEWI_DATABASE_PORT=${PG_PORT}"
echo "  ZEWI_DATABASE_USERNAME=${PG_USER}"
echo "  ZEWI_DATABASE_PASSWORD=${PG_PASSWORD}"
echo "  ZEWI_DATABASE_NAME=${PG_DB_NAME}"
echo ""
log_info "For Kubernetes Secret, you can create it with:"
echo "  kubectl create secret generic zewi-db-secret \\"
echo "    --from-literal=host=$(hostname -I | awk '{print $1}') \\"
echo "    --from-literal=port=${PG_PORT} \\"
echo "    --from-literal=username=${PG_USER} \\"
echo "    --from-literal=password=${PG_PASSWORD} \\"
echo "    --from-literal=database=${PG_DB_NAME} \\"
echo "    --namespace=zewi"
echo ""
