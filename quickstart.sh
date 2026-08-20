#!/bin/bash
# ReSet Quick Start & Test Script
# ================================
# This script starts all services and runs the test suite.

set -e

echo "╔══════════════════════════════════════════════════════════════╗"
echo "║           ReSet — Quick Start & Test Script                  ║"
echo "╚══════════════════════════════════════════════════════════════╝"
echo ""

# Check dependencies
command -v docker-compose >/dev/null 2>&1 || { echo "docker-compose is required but not installed. Aborting." >&2; exit 1; }
command -v python3 >/dev/null 2>&1 || { echo "python3 is required but not installed. Aborting." >&2; exit 1; }

# Step 1: Start infrastructure
echo "📦 Step 1: Starting PostgreSQL, Redis, and Python services..."
docker-compose up -d db redis python

# Wait for services to be healthy
echo "⏳ Waiting for services to be healthy..."
sleep 10

# Step 2: Install Python test dependencies
echo "🐍 Step 2: Installing test dependencies..."
pip3 install requests --quiet 2>/dev/null || pip install requests --quiet 2>/dev/null

# Step 3: Run classifier tests
echo "🧪 Step 3: Running classifier tests..."
python3 test/test_runner.py --mode classifier --python-url http://localhost:8001

# Step 4: Start Go API (if built)
if [ -f "./reset-api" ] || [ -f "./main" ]; then
    echo "🚀 Step 4: Starting Go API..."
    # In a real scenario, you'd start the Go binary here
    # ./reset-api &
    echo "   (Skipping Go API start — build and run manually if needed)"
else
    echo "⚠️  Step 4: Go API binary not found. Build with: go build -o reset-api cmd/api/main.go"
fi

echo ""
echo "✅ Setup complete!"
echo ""
echo "📋 Next steps:"
echo "   1. Open http://localhost:8080 in your browser"
echo "   2. Try the Detect page with sample transactions"
echo "   3. Check the Health page at http://localhost:8080/health"
echo ""
echo "🧪 To run tests again:"
echo "   python3 test/test_runner.py --mode full --api-url http://localhost:8080"
echo ""