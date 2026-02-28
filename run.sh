#!/bin/bash

# Configuration
SERVER_DIR="server"
WEB_DIR="web"
SERVER_PORT=8080
WEB_PORT=5173
SERVER_LOG="server.log"
WEB_LOG="web.log"

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Helper to check if a port is in use
check_port() {
    lsof -i:$1 >/dev/null 2>&1
    return $?
}

# Helper to kill process by port
kill_by_port() {
    local port=$1
    local name=$2
    local pids=$(lsof -t -i:$port)
    
    if [ -n "$pids" ]; then
        echo -e "${YELLOW}Stopping $name on port $port (PIDs: $pids)...${NC}"
        kill $pids
        # Wait for process to exit
        sleep 1
        # Force kill if still running
        if check_port $port; then
            echo -e "${RED}Force killing $name...${NC}"
            kill -9 $pids 2>/dev/null
        fi
        echo -e "${GREEN}$name stopped.${NC}"
    else
        echo -e "${YELLOW}$name is not running on port $port.${NC}"
    fi
}

# Helper to wait for a port to be ready
wait_for_port() {
    local port=$1
    local timeout=15
    local count=0
    
    echo -n "Waiting for port $port..."
    while ! check_port $port; do
        sleep 1
        count=$((count+1))
        echo -n "."
        if [ $count -ge $timeout ]; then
            echo ""
            return 1
        fi
    done
    echo ""
    return 0
}

start() {
    echo -e "${GREEN}=== Starting TimeTrace ===${NC}"

    # Check dependencies
    if ! command -v go &> /dev/null; then
        echo -e "${RED}Error: Go is not installed.${NC}"
        exit 1
    fi
    if ! command -v npm &> /dev/null; then
        echo -e "${RED}Error: npm is not installed.${NC}"
        exit 1
    fi

    # Stop existing processes first
    if check_port $SERVER_PORT; then
        echo -e "${YELLOW}Port $SERVER_PORT is already in use. Stopping backend...${NC}"
        kill_by_port $SERVER_PORT "Backend"
    fi
    if check_port $WEB_PORT; then
        echo -e "${YELLOW}Port $WEB_PORT is already in use. Stopping frontend...${NC}"
        kill_by_port $WEB_PORT "Frontend"
    fi

    # Start Backend
    echo -e "Starting Backend Server..."
    cd $SERVER_DIR
    nohup go run main.go > ../$SERVER_LOG 2>&1 &
    cd ..
    
    # Wait for backend to be ready
    if wait_for_port $SERVER_PORT; then
        echo -e "${GREEN}Backend Server started successfully on port $SERVER_PORT.${NC}"
    else
        echo -e "${RED}Failed to start Backend Server. Check $SERVER_LOG for details.${NC}"
        exit 1
    fi

    # Start Frontend
    echo -e "Starting Frontend Web..."
    cd $WEB_DIR
    nohup npm run dev > ../$WEB_LOG 2>&1 &
    cd ..

    echo -e "${GREEN}Frontend Web starting... logs at $WEB_LOG${NC}"
    echo -e "${GREEN}=== TimeTrace is running! ===${NC}"
    echo -e "Backend: http://localhost:$SERVER_PORT"
    echo -e "Frontend: http://localhost:$WEB_PORT"
}

stop() {
    echo -e "${GREEN}=== Stopping TimeTrace ===${NC}"
    kill_by_port $SERVER_PORT "Backend"
    kill_by_port $WEB_PORT "Frontend"
    echo -e "${GREEN}All services stopped.${NC}"
}

case "$1" in
    start)
        start
        ;;
    stop)
        stop
        ;;
    restart)
        stop
        sleep 1
        start
        ;;
    *)
        echo "Usage: $0 {start|stop|restart}"
        exit 1
esac
