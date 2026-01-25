# Deployment Guide

This project supports multiple deployment targets.

## Vercel

1. Install Vercel CLI: `npm i -g vercel`
2. Run `vercel` in the project root.
3. The project is configured to use the Go runtime via `vercel.json` and `api/index.go`.

## Railway

1. Connect your GitHub repository to Railway.
2. Railway will automatically detect the `railway.toml` and use the `Dockerfile` to build the service.

## Zeabur

1. Connect your GitHub repository to Zeabur.
2. Zeabur will automatically detect the `zeabur.toml` and use the `Dockerfile`.

## Docker (Local / Other Clouds)

Build the image:
```bash
docker build -t category-api .
```

Run the container:
```bash
docker run -p 8080:8080 category-api
```

The server will be available at `http://localhost:8080`.
