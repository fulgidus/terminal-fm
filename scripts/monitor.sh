#!/bin/bash

# Terminal-Radio Monitoring Script
# Real-time dashboard for monitoring the service

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

APP_NAME="terminal-radio"
APP_PORT=22

# Clear screen and hide cursor
clear
tput civis

# Trap to show cursor on exit
trap 'tput cnorm; exit' INT TERM

while true; do
    tput cup 0 0
    
    echo -e "${CYAN}=================================================="
    echo -e "  Terminal-Radio Monitoring Dashboard"
    echo -e "==================================================${NC}"
    echo ""
    
    # Service Status
    echo -e "${BLUE}📊 Service Status${NC}"
    echo "─────────────────────────────────────────────────"
    if systemctl is-active --quiet $APP_NAME; then
        echo -e "Status: ${GREEN}● RUNNING${NC}"
    else
        echo -e "Status: ${RED}● STOPPED${NC}"
    fi
    
    UPTIME=$(systemctl show $APP_NAME --property=ActiveEnterTimestamp --value)
    if [ -n "$UPTIME" ] && [ "$UPTIME" != "n/a" ]; then
        echo "Uptime: $(systemctl show $APP_NAME --property=ActiveEnterTimestamp --value | xargs -I {} date -d {} '+%Y-%m-%d %H:%M:%S')"
    fi
    
    echo ""
    
    # Active Connections
    echo -e "${BLUE}🔌 Active Connections${NC}"
    echo "─────────────────────────────────────────────────"
    CONN_COUNT=$(ss -tnp 2>/dev/null | grep ":$APP_PORT" | grep ESTAB | wc -l)
    if [ "$CONN_COUNT" -gt 0 ]; then
        echo -e "Connected users: ${GREEN}$CONN_COUNT${NC}"
    else
        echo -e "Connected users: ${YELLOW}0${NC}"
    fi
    echo ""
    
    # Resource Usage
    echo -e "${BLUE}💻 Resource Usage${NC}"
    echo "─────────────────────────────────────────────────"
    
    # CPU & Memory
    if systemctl is-active --quiet $APP_NAME; then
        PID=$(systemctl show $APP_NAME --property=MainPID --value)
        if [ "$PID" != "0" ] && [ -n "$PID" ]; then
            CPU=$(ps -p $PID -o %cpu= 2>/dev/null | xargs)
            MEM=$(ps -p $PID -o %mem= 2>/dev/null | xargs)
            RSS=$(ps -p $PID -o rss= 2>/dev/null | xargs)
            RSS_MB=$((RSS / 1024))
            
            echo "CPU: ${CPU}%"
            echo "Memory: ${MEM}% (${RSS_MB}MB)"
        else
            echo "N/A (process not found)"
        fi
    else
        echo "N/A (service not running)"
    fi
    echo ""
    
    # System Resources
    echo -e "${BLUE}🖥️  System Resources${NC}"
    echo "─────────────────────────────────────────────────"
    
    # Total CPU
    TOTAL_CPU=$(top -bn1 | grep "Cpu(s)" | sed "s/.*, *\([0-9.]*\)%* id.*/\1/" | awk '{print 100 - $1}')
    echo "System CPU: ${TOTAL_CPU}%"
    
    # Memory
    MEM_INFO=$(free -m | awk 'NR==2{printf "%.1f%%", $3*100/$2}')
    MEM_USED=$(free -m | awk 'NR==2{printf "%dMB", $3}')
    MEM_TOTAL=$(free -m | awk 'NR==2{printf "%dMB", $2}')
    echo "System Memory: ${MEM_INFO} (${MEM_USED}/${MEM_TOTAL})"
    
    # Disk
    DISK_INFO=$(df -h / | awk 'NR==2{print $5}')
    DISK_USED=$(df -h / | awk 'NR==2{print $3}')
    DISK_TOTAL=$(df -h / | awk 'NR==2{print $2}')
    echo "Disk Usage: ${DISK_INFO} (${DISK_USED}/${DISK_TOTAL})"
    
    echo ""
    
    # Recent Logs
    echo -e "${BLUE}📝 Recent Logs (last 10)${NC}"
    echo "─────────────────────────────────────────────────"
    journalctl -u $APP_NAME -n 10 --no-pager --output=short-iso 2>/dev/null | tail -10 || echo "No logs available"
    
    echo ""
    echo "─────────────────────────────────────────────────"
    echo -e "${YELLOW}Press Ctrl+C to exit${NC} | Refreshing every 5s..."
    
    # Wait 5 seconds before refresh
    sleep 5
done
