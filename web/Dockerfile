FROM node:22-alpine AS builder


WORKDIR /app
# 阶段一
# 启用 pnpm
RUN corepack enable && corepack prepare pnpm@latest --activate

# 先复制依赖文件，尽量利用 docker 缓存
COPY package.json pnpm-local.yaml ./

RUN pnpm install --frozen-lockfile

# 复制源码并构建
COPY . .
RUN pnpm build

# 阶段二
# 运行
FROM nginx:1.27-alpine
# 把构建产物放在 nginx 的默认目录下
COPY --from=builder /app/dist /usr/share/nginx/html

# 写入 nginx 配置
RUN printf 'server {\n\
listen 80;\n\
server_name localhost;\n\
root /usr/share/nginx/html;\n\
index index.html;\n\
location / {\n\
    try_files $uri $uri/ /index.html;\n\
}\n\
location /api {\n\
    proxy_pass http://backend:8080;\n\
    proxy_set_header Host $host;\n\
    proxy_set_header X-Real-IP $remote_addr;\n\
}\n\
}\n' > /etc/nginx/conf.d/default.conf


EXPOSE 80