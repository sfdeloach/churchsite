#!/bin/bash
# Bootstrap Let's Encrypt certificates for the preview deployment.
#
# Solves the chicken-and-egg problem: nginx needs certs to start,
# certbot needs nginx for HTTP-01 challenge.
#
# Usage: ./scripts/init-letsencrypt.sh <domain> [email]
# Example: ./scripts/init-letsencrypt.sh sachapel.duckdns.org admin@sachapel.com

set -euo pipefail

DOMAIN="${1:?Usage: $0 <domain> [email]}"
EMAIL="${2:-}"
COMPOSE="docker compose -f compose.yml -f compose.prod.yml"
CERT_NAME="sachapel"
CERT_PATH="./certbot/conf/live/$CERT_NAME"

echo "==> Domain: $DOMAIN"
echo "==> Cert name: $CERT_NAME"

# 1. Create required directories
echo "==> Creating certbot directories..."
mkdir -p certbot/conf certbot/www

# 2. Generate self-signed placeholder so nginx can start
if [ ! -f "$CERT_PATH/fullchain.pem" ]; then
    echo "==> Generating self-signed placeholder certificate..."
    mkdir -p "$CERT_PATH"
    openssl req -x509 -nodes -newkey rsa:2048 -days 1 \
        -keyout "$CERT_PATH/privkey.pem" \
        -out "$CERT_PATH/fullchain.pem" \
        -subj "/CN=localhost" 2>/dev/null
else
    echo "==> Certificate already exists, skipping placeholder..."
fi

# 3. Start nginx (and dependencies) so certbot can reach port 80
echo "==> Starting nginx..."
$COMPOSE up -d nginx
echo "==> Waiting for nginx to be ready..."
sleep 5

# 4. Remove placeholder and obtain real certificate
echo "==> Removing placeholder certificate..."
rm -rf "$CERT_PATH"

CERTBOT_ARGS=(
    certonly --webroot
    -w /var/www/certbot
    -d "$DOMAIN"
    --cert-name "$CERT_NAME"
    --agree-tos
    --no-eff-email
    --force-renewal
)

if [ -n "$EMAIL" ]; then
    CERTBOT_ARGS+=(--email "$EMAIL")
else
    CERTBOT_ARGS+=(--register-unsafely-without-email)
fi

echo "==> Requesting certificate from Let's Encrypt..."
$COMPOSE run --rm --entrypoint certbot certbot "${CERTBOT_ARGS[@]}"

# 5. Reload nginx with the real certificate
echo "==> Reloading nginx..."
$COMPOSE exec nginx nginx -s reload

echo "==> Done! Certificate obtained for $DOMAIN"
echo "    Cert path: /etc/letsencrypt/live/$CERT_NAME/"
echo "    Auto-renewal handled by the certbot container."
