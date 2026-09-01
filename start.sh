#!/bin/bash
set -e

echo "=== Afritech Online Startup ==="

# Ensure Go and Node.js are in PATH for Git Bash
export PATH="/c/Program Files/Go/bin:/c/Program Files/nodejs:$PATH"

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

# Install frontend dependencies only if missing
if [ ! -d "frontend/node_modules" ]; then
    echo "[INFO] Installing frontend dependencies..."
    cd frontend
    npm install
    cd ..
else
    echo "[OK] Frontend dependencies already installed"
fi

# Start backend in background
echo "[START] Starting backend on http://localhost:8080 ..."
cd backend
go run ./cmd/server &
BACKEND_PID=$!
cd ..

# Wait for backend to start and apply migrations
echo "[WAIT] Waiting for backend to start..."
sleep 5

# Check if backend is running
if ! kill -0 $BACKEND_PID 2>/dev/null; then
    echo "[ERROR] Backend failed to start. Check the output above."
    exit 1
fi

echo "[OK] Backend started. Migrations applied."

# Start frontend in background
echo "[START] Starting frontend on http://localhost:3000 ..."
cd frontend
npm run dev &
FRONTEND_PID=$!
cd ..

# Wait for frontend to start
echo "[WAIT] Waiting for frontend to start..."
sleep 8

# Check if frontend is running
if kill -0 $FRONTEND_PID 2>/dev/null; then
    echo "[OPEN] Opening frontend in browser..."
    start http://localhost:3000
else
    echo "[WARN] Frontend failed to start. Check the output above."
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
