#!/bin/bash
set -e

echo "=== Afritech Online Startup ==="

# Ensure Node.js is in PATH for Git Bash
export PATH="/c/Program Files/nodejs:$PATH"

# Check Go
if ! command -v go &> /dev/null; then
    echo "ERROR: Go is not installed or not in PATH"
    exit 1
fi
echo "[OK] Go found: $(go version)"

# Check Node.js
if ! command -v node &> /dev/null; then
    echo "ERROR: Node.js is not installed or not in PATH"
    exit 1
fi
echo "[OK] Node.js found: $(node -v)"

# Check npm
if ! command -v npm &> /dev/null; then
    echo "ERROR: npm is not installed or not in PATH"
    exit 1
fi
echo "[OK] npm found: $(npm -v)"

# Clean and install frontend dependencies
echo "[INFO] Installing frontend dependencies..."
cd frontend
rm -rf node_modules .next
npm install
cd ..

# Start backend in background with retry for Go module downloads
echo "[START] Starting backend on http://localhost:8080 ..."
cd backend
for i in 1 2 3; do
    echo "[BACKEND] Attempt $i: starting..."
    go run ./cmd/server &
    BACKEND_PID=$!
    sleep 3
    if kill -0 $BACKEND_PID 2>/dev/null; then
        echo "[BACKEND] Started successfully"
        break
    else
        echo "[BACKEND] Attempt $i failed, retrying..."
        kill $BACKEND_PID 2>/dev/null || true
        sleep 2
    fi
done
cd ..

# Start frontend in background
echo "[START] Starting frontend on http://localhost:3000 ..."
cd frontend
npm run dev &
FRONTEND_PID=$!
cd ..

# Wait for services to start
echo "[WAIT] Waiting for services to start..."
sleep 8

# Check if frontend is running
if kill -0 $FRONTEND_PID 2>/dev/null; then
    echo "[OPEN] Opening frontend in browser..."
    start http://localhost:3000
else
    echo "[WARN] Frontend failed to start. Check output above."
fi

echo ""
echo "=== Afritech Online is running ==="
echo "  Backend:  http://localhost:8080"
echo "  Frontend: http://localhost:3000"
echo "  API Docs: See docs/api.md"
echo ""
echo "Press Ctrl+C to stop all services"

# Trap Ctrl+C and kill both processes
trap "echo ''; echo '[STOP] Shutting down...'; kill $BACKEND_PID $FRONTEND_PID 2>/dev/null; exit 0" INT TERM

# Wait for processes
wait
