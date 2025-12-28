# NewsApp-Go

A minimal news website with a Go backend and a TypeScript (React + Vite) frontend.

## Structure

- `backend/` — Go HTTP server exposing `/api/articles` and serving static files from `./static`.
- `frontend/` — Vite + React TypeScript app that fetches `/api/articles`.

## Requirements

- Go 1.20+
- Node 18+ and npm or yarn

## Development

1. Run backend:

```bash
cd backend
go run .
```

Server runs on `:8080` by default.

2. Run frontend (in separate terminal):

```bash
cd frontend
npm install
npm run dev
```

Frontend dev server runs on `:5173` and proxies `/api` to the Go backend.

## Production build

1. Build frontend:

```bash
cd frontend
npm run build
```

2. Copy `frontend/dist` to `backend/static`, then run the Go server:

```bash
rm -rf backend/static
cp -R frontend/dist backend/static
cd backend
go run .
```

Now the Go server will serve the static site and API on the same port.
