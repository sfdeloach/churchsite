---
name: deploy
description: Deploy to the preview environment on AWS EC2. Use when ready to push changes to the live preview site.
disable-model-invocation: true
---

# Preview Deployment

Deploy the current state to the AWS EC2 preview environment.

## Pre-flight Checks

Before deploying, verify:

1. **Clean working tree**: Run `git status` and confirm no uncommitted changes
2. **On main branch**: Confirm `git branch --show-current` returns `main`
3. **Templates generated**: Run `make generate` to ensure all `_templ.go` files are current
4. **Build succeeds locally**: Run `make build` to catch compile errors before deploying

If any check fails, stop and report the issue. Do not deploy with uncommitted changes or build failures.

## Deploy Process

The deployment is a manual SSH-based process. Provide the user with the exact commands to run:

```bash
# SSH into the EC2 instance
ssh -i <pemfile> ubuntu@<elastic-ip>

# Navigate to the project
cd /opt/sachapel

# Pull latest changes, rebuild, and migrate
sudo make preview-deploy
```

The `make preview-deploy` target runs: `git pull` + `docker compose rebuild` + `migrate up`.

## Post-Deploy Verification

After deployment, verify:

1. **Health check**: `curl -s https://sachapel.duckdns.org/health` should return `{"status": "ok"}`
2. **Readiness check**: `curl -s https://sachapel.duckdns.org/health/ready` should return `{"status": "ok", "postgres": "ok", "redis": "ok"}`
3. **Visual check**: Open the site in a browser and verify the deployed changes

## Troubleshooting

If deployment fails:

- **Build failure**: Check `make preview-logs` for Go compile errors
- **Migration failure**: Check if a migration needs to be rolled back with `make migrate-down`
- **Container won't start**: Check `docker compose -f compose.yml -f compose.prod.yml logs app`
- **Nginx errors**: Check `docker compose -f compose.yml -f compose.prod.yml logs nginx`
- **Out of memory**: The t3.micro has 1GB RAM + 2GB swap. Check `free -h` and `docker stats`

## Important

- This deploys to the **preview** environment (sachapel.duckdns.org), not production
- The preview environment uses `compose.yml` + `compose.prod.yml`
- Never run destructive database operations without confirming with the user first
