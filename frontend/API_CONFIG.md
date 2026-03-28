# API 配置说明

## 开发环境配置

当其他电脑需要访问时,请修改 `.env.development` 文件:

```bash
# 将 localhost 改为后端服务器的实际 IP 地址
VITE_API_URL=http://192.168.1.100:8080
```

### 步骤:

1. 查看后端服务器的 IP 地址:
   - Windows: 打开 cmd,输入 `ipconfig`
   - Mac/Linux: 打开终端,输入 `ifconfig` 或 `ip addr`

2. 编辑 `frontend/.env.development` 文件
   - 将 `localhost` 替换为实际的 IP 地址

3. 重启前端开发服务器:
   ```bash
   cd frontend
   npm run dev
   ```

## 生产环境配置

部署时修改 `.env.production` 文件中的 `VITE_API_URL` 为实际的后端地址。

## 示例

假设后端服务器 IP 是 `192.168.1.100`:

```bash
# 前端访问地址: http://192.168.1.100:5173
# 后端 API 地址: http://192.168.1.100:8080

# .env.development 配置:
VITE_API_URL=http://192.168.1.100:8080
```

## 注意事项

- 修改 `.env` 文件后需要重启前端服务才能生效
- 确保后端服务器允许外部访问(Gin 默认监听 0.0.0.0)
- 确保防火墙允许相应端口的访问
