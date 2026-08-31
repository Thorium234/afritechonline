#!/bin/bash
set -e

echo "=== Afritech Online Startup ==="

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

# Install frontend dependencies if needed
if [ ! -d "frontend/node_modules" ]; then
    echo "[INFO] Installing frontend dependencies..."
    cd frontend
    npm install
    cd ..
fi

# Start backend in background
echo "[START] Starting backend on http://localhost:8080 ..."
cd backend
go run ./cmd/server &
BACKEND_PID=$!
cd ..

# Start frontend in background
echo "[START] Starting frontend on http://localhost:3000 ..."
cd frontend
npm run dev &
FRONTEND_PID=$!
cd ..

# Wait for services to start
echo "[WAIT] Waiting for services to start..."
sleep 5

# Open frontend in default browser
echo "[OPEN] Opening frontend in browser..."
start http://localhost:3000

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
