# Installation instructions
1) [Register github app]('https://docs.github.com/en/apps/creating-github-apps/registering-a-github-app/registering-a-github-app')
2) Run local / host on AWS

# Run local instructions
### Back
1) edit `.env.dev` file, you can take variable names from `.env.example`
2) `sh run-dev.sh`
3) Backend is working

### Front
1) edit `front/.env`, you can take variable names from `.env.example`
2) `cd front && npm i && npm run dev`
3) Frontend is working

### Set up nginx config
```
upstream back {
    server 127.0.0.1:3000;
}

upstream front {
    server localhost:5173;
}

server {
    listen 80;
    listen [::]:80;
    server_name localhost;
    client_max_body_size 50M;

    location / {
        proxy_pass http://front;
    }
    location /repos {
        proxy_pass http://front;
    }
    location /\w/new {
        proxy_pass http://front;
    }
    location /api/status {
        proxy_pass http://back;
        proxy_redirect     off;
        proxy_read_timeout 300s;
    }
    location /redirect {
        proxy_pass http://back;
        proxy_redirect     off;
        proxy_read_timeout 300s;
    }
    location /callback {
        proxy_pass http://back;
        proxy_redirect     off;
        proxy_read_timeout 300s;
    }
    location /api/ {
        proxy_pass http://back;
        proxy_redirect     off;
        proxy_read_timeout 300s;
    }
}
```
2) sudo systemctl restart nginx
