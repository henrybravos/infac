#!/bin/bash

# Kill any process using port 8885
echo "🔥 Stopping any process using port 8885..."
lsof -ti:8885 | xargs kill -9 2>/dev/null || true

# Wait a moment for the port to be free
sleep 1

# Start development server with Air
echo "🚀 Starting development server with Air..."
air