# ============================================
# 多阶段构建 Dockerfile
# 同时构建前后端，打包到单个容器
# ============================================

# --------------------------------------------
# 阶段1: 构建前端
# --------------------------------------------
FROM node:20-alpine AS frontend-builder

WORKDIR /frontend

# 复制前端依赖文件
COPY frontend/package*.json ./
RUN npm ci

# 复制前端源码并构建
COPY frontend/ ./
RUN npm run build

# --------------------------------------------
# 阶段2: 构建后端
# --------------------------------------------
FROM golang:alpine AS backend-builder

WORKDIR /backend

# 安装必要的构建工具
RUN apk add --no-cache git

# 复制依赖文件
COPY backend/go.mod backend/go.sum ./
RUN go mod download

# 复制源码并构建（包括 bin 目录）
COPY backend/ ./

# 构建时指定数据目录为 /app/data（自动检测架构）
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-w -s" \
    -o frpc_webmanager main.go

# --------------------------------------------
# 阶段3: 运行镜像
# --------------------------------------------
FROM alpine:latest

# 安装运行时依赖
RUN apk --no-cache add ca-certificates jq wget

WORKDIR /app

# 创建数据目录
RUN mkdir -p /app/data/configs /app/data/logs

# 从后端构建阶段复制二进制文件
COPY --from=backend-builder /backend/frpc_webmanager /app/

# 自动下载对应架构的 frpc 二进制文件
# 支持的架构: amd64, arm, arm64
RUN mkdir -p /app/bin && \
    ARCH=$(uname -m) && \
    case "$ARCH" in \
        x86_64) FRPC_ARCH="amd64" ;; \
        aarch64) FRPC_ARCH="arm64" ;; \
        armv7l) FRPC_ARCH="arm" ;; \
        *) FRPC_ARCH="amd64" ;; \
    esac && \
    FRPC_VERSION="0.58.1" && \
    wget -O /tmp/frpc.tar.gz \
        "https://github.com/fatedier/frp/releases/download/v${FRPC_VERSION}/frp_${FRPC_VERSION}_linux_${FRPC_ARCH}.tar.gz" && \
    tar -xzf /tmp/frpc.tar.gz -C /tmp && \
    find /tmp -name "frpc" -type f -exec mv {} /app/bin/frpc \; && \
    rm -rf /tmp/frpc.tar.gz /tmp/frp_* && \
    chmod +x /app/bin/frpc

# 复制文档目录（用于 about 页面的更新日志）
COPY --from=backend-builder /backend/docs /app/docs/

# 从前端构建阶段复制静态文件
COPY --from=frontend-builder /frontend/dist /app/frontend/dist/

# 复制启动脚本
COPY docker/entrypoint.sh /app/entrypoint.sh
RUN chmod +x /app/entrypoint.sh

# 暴露默认端口（可通过 settings.json 修改）
EXPOSE 4500

# 数据卷
VOLUME ["/app/data"]

# 健康检查(使用轻量级health端点,避免极空间等NAS系统显示非正常状态)
# 注意: wget --spider 模式有bug,使用 -O /dev/null 代替
HEALTHCHECK --interval=30s --timeout=3s --start-period=10s --retries=3 \
    CMD wget -q -O /dev/null http://localhost:4500/api/health || exit 1

# 启动命令
ENTRYPOINT ["/app/entrypoint.sh"]
