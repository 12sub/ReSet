#!/bin/bash
# ReSet — 502 Error Diagnostic & Fix Script
# =========================================
# Run this to diagnose why the Python AI service is returning 502

set -e

echo "╔══════════════════════════════════════════════════════════════╗"
echo "║     ReSet 502 Error Diagnostic & Auto-Fix Tool               ║"
echo "╚══════════════════════════════════════════════════════════════╝"
echo ""

PYTHON_URL="${PYTHON_URL:-http://localhost:8001}"
API_URL="${API_URL:-http://localhost:8080}"

echo "📋 Configuration:"
echo "   Python Service URL: $PYTHON_URL"
echo "   Go API URL:         $API_URL"
echo ""

# ═══════════════════════════════════════════════════════════════
# CHECK 1: Is the Python service running?
# ═══════════════════════════════════════════════════════════════
echo "🔍 CHECK 1: Is the Python AI service responding?"
if curl -s "$PYTHON_URL/health" > /dev/null 2>&1; then
    echo "   ✅ Python service is UP at $PYTHON_URL"
    python_status=$(curl -s "$PYTHON_URL/health" | python3 -c "import sys,json; print(json.load(sys.stdin).get('status','unknown'))" 2>/dev/null || echo "unknown")
    echo "   Health status: $python_status"
else
    echo "   ❌ Python service is NOT responding at $PYTHON_URL"
    echo ""
    echo "   🔧 FIX: Start the Python service:"
    echo ""
    echo "   Option A — Direct (no Docker):"
    echo "      cd ai/"
    echo "      pip install -r requirements.txt"
    echo "      uvicorn main:app --reload --port 8001"
    echo ""
    echo "   Option B — Docker Compose:"
    echo "      docker-compose up -d python"
    echo ""
    echo "   Option C — Check if it's running on a different port:"
    echo "      lsof -i :8001    # macOS"
    echo "      ss -tlnp | grep 8001   # Linux"
    echo ""
    exit 1
fi
echo ""

# ═══════════════════════════════════════════════════════════════
# CHECK 2: Can the Go backend reach the Python service?
# ═══════════════════════════════════════════════════════════════
echo "🔍 CHECK 2: Can the Go backend classify a transaction?"
test_payload='{"text":"NETFLIX.COM 4500.00 - Monthly subscription"}'

if curl -s -X POST "$API_URL/subscriptions/detect"     -H "Content-Type: application/x-www-form-urlencoded"     -d "email=demo@reset.ng"     -d "text=NETFLIX.COM 4500.00" > /dev/null 2>&1; then
    echo "   ✅ Go API is reachable at $API_URL"
else
    echo "   ❌ Go API is NOT responding at $API_URL"
    echo ""
    echo "   🔧 FIX: Start the Go backend:"
    echo "      go build -o reset-api ."
    echo "      ./reset-api"
    echo ""
fi
echo ""

# ═══════════════════════════════════════════════════════════════
# CHECK 3: Direct Python classification test
# ═══════════════════════════════════════════════════════════════
echo "🔍 CHECK 3: Direct classification test to Python service..."
response=$(curl -s -X POST "$PYTHON_URL/classify-transaction"     -H "Content-Type: application/json"     -d "$test_payload" 2>/dev/null || echo "FAILED")

if echo "$response" | grep -q "merchant"; then
    echo "   ✅ Python classifier works directly!"
    echo "   Response preview: $(echo "$response" | python3 -c "import sys,json; d=json.load(sys.stdin); print(f"{d['merchant']} | {d['category']} | {d['confidence']}")" 2>/dev/null || echo "(parse error)")"
else
    echo "   ❌ Python classifier failed directly"
    echo "   Response: $response"
fi
echo ""

# ═══════════════════════════════════════════════════════════════
# CHECK 4: Docker-specific networking (if applicable)
# ═══════════════════════════════════════════════════════════════
if command -v docker-compose &> /dev/null || command -v docker &> /dev/null; then
    echo "🔍 CHECK 4: Docker container status..."

    if docker ps --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}" 2>/dev/null | grep -q "reset"; then
        echo "   Docker containers found:"
        docker ps --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}" | grep "reset" || true
        echo ""

        # Check if python container is healthy
        if docker ps | grep -q "reset-python"; then
            echo "   ✅ reset-python container is running"
            echo "   Checking container health..."
            docker exec reset-python python -c "import urllib.request; urllib.request.urlopen('http://localhost:8001/health')" 2>/dev/null && echo "   ✅ Container healthcheck passes" || echo "   ⚠️ Container healthcheck failed"
        else
            echo "   ❌ reset-python container NOT found"
            echo "   🔧 FIX: docker-compose up -d python"
        fi

        # Check Go container
        if docker ps | grep -q "reset-api"; then
            echo "   ✅ reset-api container is running"
        else
            echo "   ⚠️ reset-api container NOT found (may be running locally)"
        fi
    else
        echo "   No ReSet Docker containers found (running locally?)"
    fi
    echo ""
fi

# ═══════════════════════════════════════════════════════════════
# CHECK 5: Environment variables
# ═══════════════════════════════════════════════════════════════
echo "🔍 CHECK 5: Environment variables..."
echo "   PYTHON_URL=${PYTHON_URL:-(not set, using default)}"
echo "   PAYSTACK_SECRET_KEY=${PAYSTACK_SECRET_KEY:0:15}... (truncated)"
echo "   DATABASE_URL=${DATABASE_URL:-(not set)}"
echo "   REDIS_ADDR=${REDIS_ADDR:-(not set)}"
echo ""

# ═══════════════════════════════════════════════════════════════
# CHECK 6: Python dependencies
# ═══════════════════════════════════════════════════════════════
echo "🔍 CHECK 6: Python dependencies..."
if python3 -c "import fastapi, uvicorn, rapidfuzz" 2>/dev/null; then
    echo "   ✅ All Python dependencies installed"
else
    echo "   ❌ Missing Python dependencies"
    echo "   🔧 FIX: pip install -r ai/requirements.txt"
fi
echo ""

# ═══════════════════════════════════════════════════════════════
# QUICK FIXES
# ═══════════════════════════════════════════════════════════════
echo "═══════════════════════════════════════════════════════════════"
echo "  🔧 COMMON FIXES FOR 502 ERROR"
echo "═══════════════════════════════════════════════════════════════"
echo ""
echo "  PROBLEM 1: Python service not running"
echo "  ───────────────────────────────────────"
echo "  Fix: cd ai/ && uvicorn main:app --reload --port 8001"
echo ""
echo "  PROBLEM 2: Wrong PYTHON_URL in Go backend"
echo "  ───────────────────────────────────────────"
echo "  If running locally (no Docker):"
echo "     PYTHON_URL=http://localhost:8001"
echo "  If running in Docker:"
echo "     PYTHON_URL=http://python:8001"
echo ""
echo "  Fix: Set in .env or export PYTHON_URL=http://localhost:8001"
echo ""
echo "  PROBLEM 3: Port conflict on 8001"
echo "  ─────────────────────────────────"
echo "  Fix: Use a different port"
echo "       uvicorn main:app --reload --port 8002"
echo "       export PYTHON_URL=http://localhost:8002"
echo ""
echo "  PROBLEM 4: Docker networking issue"
echo "  ───────────────────────────────────"
echo "  Fix: Restart all services"
echo "       docker-compose down && docker-compose up -d"
echo ""
echo "  PROBLEM 5: Python service crashed"
echo "  ─────────────────────────────────"
echo "  Fix: Check Python logs"
echo "       docker logs reset-python   # Docker"
echo "       # or check terminal output # Local"
echo ""
echo "═══════════════════════════════════════════════════════════════"