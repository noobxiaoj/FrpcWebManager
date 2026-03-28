#!/bin/sh
set -e

# 禁用颜色输出（避免 ANSI 转义码污染日志）
echo_info() {
    echo "[INFO] $1"
}

echo_success() {
    echo "[SUCCESS] $1"
}

echo_warn() {
    echo "[WARN] $1"
}

echo_error() {
    echo "[ERROR] $1"
}

# 数据目录
DATA_DIR="/app/data"
SETTINGS_FILE="$DATA_DIR/settings.json"

# 确保数据目录结构存在（静默执行）
mkdir -p "$DATA_DIR/configs"
mkdir -p "$DATA_DIR/logs"

# 检查 settings.json 是否存在，不存在则创建默认配置
if [ ! -f "$SETTINGS_FILE" ]; then
    cat > "$SETTINGS_FILE" <<EOF
{
  "frontendPort": 4500,
  "showServerPort": false,
  "refreshInterval": 1,
  "showRefreshTime": true,
  "showServerName": true
}
EOF
fi

# 读取前端端口配置
FRONTEND_PORT=$(jq -r '.frontendPort // 4500' "$SETTINGS_FILE" 2>/dev/null || echo "4500")

# 验证端口是否有效（静默修正）
if ! [ "$FRONTEND_PORT" -eq "$FRONTEND_PORT" ] 2>/dev/null || [ "$FRONTEND_PORT" -lt 1024 ] || [ "$FRONTEND_PORT" -gt 65535 ]; then
    FRONTEND_PORT=4500
fi

# 启动后端服务（托管前端静态文件）
# 通过环境变量 PORT 传递端口配置
export PORT="$FRONTEND_PORT"
exec /app/frpc_webmanager
