# Ubuntu Server deployment

This is Step 3 of the DevOps hardening project. The goal is to deploy the same Dockerized application manually on an Ubuntu Server before introducing Kubernetes.

## 1. Prepare the server

```bash
sudo apt update
sudo apt upgrade -y
sudo apt install -y ca-certificates curl git
```

Install Docker using Docker's official Ubuntu instructions, then verify:

```bash
docker --version
docker compose version
```

## 2. Clone the repository

```bash
git clone https://github.com/XU4NLOC/go_http-rest-api.git
cd go_http-rest-api
```

Copy the DevOps files from this project into the repository if you have not committed them yet.

## 3. Create local environment variables

```bash
cat > .env <<'EOF_ENV'
DB_USER=budgetuser
DB_PASSWORD=change-this-password
DB_NAME=budgettracker
EOF_ENV
```

Do not commit `.env`.

## 4. Build and start

```bash
docker compose up -d --build
```

Check:

```bash
docker compose ps
docker compose logs -f api
```

Verify the application:

```bash
curl http://localhost:8080/health
curl http://localhost:8080/ready
```

Expected:

```json
{"status":"ok"}
```

and:

```json
{"status":"ready"}
```

## 5. Basic Linux troubleshooting

```bash
docker ps
docker logs <container>
ss -lntp
df -h
free -m
```

The point of this step is to learn how to diagnose the application at the Linux host + container level before adding Kubernetes.
