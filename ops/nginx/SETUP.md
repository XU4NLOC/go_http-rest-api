# Nginx reverse proxy

After the Dockerized API is running on Ubuntu at `127.0.0.1:8080`, put Nginx in front of it.

## Install

```bash
sudo apt update
sudo apt install -y nginx
```

Copy `budget-api.conf` to:

```bash
sudo cp budget-api.conf /etc/nginx/sites-available/budget-api
sudo ln -s /etc/nginx/sites-available/budget-api /etc/nginx/sites-enabled/budget-api
```

Remove the default site if necessary:

```bash
sudo rm -f /etc/nginx/sites-enabled/default
```

Test and reload:

```bash
sudo nginx -t
sudo systemctl reload nginx
```

Now the request path becomes:

```text
Client -> Nginx :80 -> Docker container :8080 -> Go API
```

Verify:

```bash
curl http://SERVER_IP/health
curl http://SERVER_IP/ready
```

Later, when a domain and TLS are available, this same reverse-proxy layer is where HTTPS termination can be added.
