# GinAdmin Monorepo 改造与 Vercel 部署

本文记录把 GinAdmin 从「Go 后端在根目录 + 前端在 web/」改造成标准 monorepo 结构的完整步骤，以及对应的 Vercel 部署方式和过程中踩过的坑。

## 1. 改造目标

- `apps/` 放可独立部署的应用：`apps/server`（Go 后端）、`apps/web`（React 前端）。
- `packages/` 放共享代码（暂时为空，用到时再建）。
- 根目录只放 workspace 配置、Docker 编排和通用脚本。
- Vercel 上拆成前端、后端两个项目，各自独立部署。

### 为什么 Vercel 要拆成两个项目

一个 Vercel 项目只识别其根目录下 `api/` 里的函数。Go 代码搬进 `apps/server/` 后，`go.mod` 也在里面，如果保持单项目，就得在仓库根目录再放一份 `api/` 和 `go.mod`，跨目录引用 `apps/server`，结构更乱。

拆成两个项目后，前端通过 rewrite 把 `/api/*` 转发给后端项目。浏览器看到的仍是同一个域名，`BASE_URL = '/api/v1'` 不用改，也没有跨域问题。

## 2. 目标结构

```
GinAdmin/
├── apps/
│   ├── server/                  # Go 后端，独立 Go module（module 名仍是 GinAdmin）
│   │   ├── api/index.go         # Vercel 函数入口
│   │   ├── app/app.go           # 入口中转包（见第 9 节）
│   │   ├── cmd/ config/ data/ internal/ SQL/
│   │   ├── main.go go.mod go.sum rbac_model.conf
│   │   ├── Dockerfile
│   │   ├── .dockerignore
│   │   └── vercel.json
│   └── web/                     # 前端
│       ├── Dockerfile
│       ├── vercel.json
│       └── src/ ...
├── packages/
├── docker-compose.yml
├── Makefile
├── package.json
├── pnpm-workspace.yaml
├── pnpm-lock.yaml
├── .dockerignore                # 前端 Docker 构建用（上下文为仓库根目录）
├── .gitignore
└── README.md  LICENSE
```

要点：

- **Go 的 import 路径不用改。** import 路径由 `go.mod` 中的 `module GinAdmin` 决定，与目录位置无关。
- **`apps/server` 里不要放 `package.json`**，否则 Vercel 会把它当 JS 项目执行安装和构建。Go 命令统一用 `go -C apps/server ...` 从根目录调用。

## 3. 迁移步骤

### 3.1 搬目录

在 `GinAdmin/` 根目录执行，用 `git mv` 保留文件历史：

```bash
# 清理误提交的文件（只从 git 移除，不删本地文件）
git rm --cached GinAdmin logs/gin-admin.sys.log

mkdir -p apps/server
git mv web apps/web
git mv api app cmd config data internal SQL \
       main.go go.mod go.sum rbac_model.conf Dockerfile \
       apps/server/

git rm vercel.json                 # 前后端各自有新的 vercel.json
git rm apps/web/.dockerignore      # 前端 Docker 上下文是根目录，这份不生效
```

### 3.2 workspace 配置

`pnpm-workspace.yaml`：

```yaml
packages:
  - apps/*
  - packages/*
```

`apps/*` 只匹配含 `package.json` 的目录，`apps/server` 不会被当成 JS 包。

根目录 `package.json` 的 `scripts`：

```json
"scripts": {
  "dev:web": "pnpm --filter gin-admin-web dev",
  "build:web": "pnpm --filter gin-admin-web build",
  "dev:server": "go -C apps/server run . server",
  "build:server": "go -C apps/server build -o bin/gin-admin ."
}
```

`go -C` 会先切换到 `apps/server` 再执行，代码里 `config/config.yaml`、`rbac_model.conf` 的相对路径依然有效。

然后执行一次 `pnpm install`，lockfile 的 `importers` 会从 `web` 自动更新为 `apps/web`，依赖版本不变。

### 3.3 .gitignore

把原来 `!web/docs/` 相关的几行替换为：

```
docs/
!apps/web/docs/
!apps/web/docs/ai/
!apps/web/docs/ai/**
```

`docs/` 匹配任意层级的 docs 目录，原来的 `!web/docs/` 只对旧路径生效，不改的话 `apps/web/docs/` 会被整个忽略。

追加编译产物：

```
apps/server/bin/
apps/server/GinAdmin
```

### 3.4 Go 开发环境（仅本地，不提交）

```bash
go work init ./apps/server
```

`go.mod` 不在仓库根目录后，IDE 打开 `GinAdmin/` 可能识别不到 module，`go.work` 用来告诉 IDE module 的位置。`.gitignore` 已忽略 `go.work`。

## 4. Vercel 配置

### 4.1 后端：apps/server/vercel.json

```json
{
  "$schema": "https://openapi.vercel.sh/vercel.json",
  "framework": null,
  "buildCommand": "mkdir -p public",
  "outputDirectory": "public",
  "functions": {
    "api/index.go": {
      "maxDuration": 60,
      "includeFiles": "{config/config.yaml,rbac_model.conf}"
    }
  },
  "rewrites": [
    { "source": "/(.*)", "destination": "/api/index" }
  ]
}
```

- `buildCommand` + `outputDirectory` 指向空的 `public` 目录：没有构建产物时，Vercel 会把项目根目录当静态文件对外提供，`config/config.yaml` 会被公开下载。
- `includeFiles`：配置文件和 Casbin 模型文件不会自动打包进函数，必须显式包含。

### 4.2 前端：apps/web/vercel.json

```json
{
  "$schema": "https://openapi.vercel.sh/vercel.json",
  "rewrites": [
    { "source": "/api/:path*", "destination": "https://<后端域名>.vercel.app/api/:path*" },
    { "source": "/(.*)", "destination": "/index.html" }
  ]
}
```

前端不需要写 `installCommand` / `buildCommand`，Vercel 会自动识别 Vite 和上层的 pnpm workspace。最后一条 rewrite 用于单页路由，否则刷新子页面会 404。

### 4.3 控制台设置

1. **现有项目改为前端项目**：Root Directory 设为 `apps/web`；Framework Preset 设为 Vite；确认 "Include files outside the root directory" 已开启（需要读取根目录 lockfile）；数据库相关环境变量可删除。
2. **新建后端项目**：导入同一仓库，Root Directory 设为 `apps/server`，配置环境变量（见 4.4）。
3. **先部署后端**，访问 `https://<后端域名>/ping` 确认返回 `pong`，再把域名填进前端 rewrite，部署前端。
4. 可选：两个项目都开启 "Skip deployments when there are no changes to the root directory"，只改一端时另一端不会重新部署。

### 4.4 后端环境变量

| 变量 | 值 |
|---|---|
| `MYSQL_HOST` | TiDB 控制台提供的 host |
| `MYSQL_PORT` | `4000` |
| `MYSQL_USER` | 形如 `xxxx.root`，带前缀 |
| `MYSQL_PASSWORD` | TiDB 密码 |
| `MYSQL_DATABASE` | 库名 |
| `MYSQL_TLS` | `true` |
| `JWT_SECRET` | 随机长字符串 |
| `LOG_OUTPUT` | `stderr` |
| `GIN_MODE` | `release` |

数据库使用 TiDB Cloud Serverless（兼容 MySQL，现有 gorm 代码和 `UNIX_TIMESTAMP()` 等 SQL 无需修改，强制 TLS）。

## 5. Docker 相关

### 5.1 apps/web/Dockerfile（第一阶段）

```dockerfile
FROM node:22-alpine AS builder

WORKDIR /app
RUN corepack enable

COPY package.json pnpm-workspace.yaml pnpm-lock.yaml ./
COPY apps/web/package.json apps/web/
RUN pnpm install --frozen-lockfile --filter gin-admin-web...

COPY apps/web/ apps/web/
RUN pnpm --filter gin-admin-web build
```

第二阶段：

```dockerfile
COPY --from=builder /app/apps/web/dist /usr/share/nginx/html
```

注意删掉旧文件中 workspace 构建之后残留的 `pnpm install` / `COPY . .` / `pnpm build` 三行。

### 5.2 .dockerignore

根目录 `.dockerignore` 追加 `apps/server`（前端镜像不需要后端代码）。

新建 `apps/server/.dockerignore`：

```
logs
uploadFiles
bin
GinAdmin
.env
```

`apps/server/Dockerfile` 内容不变，构建上下文就是 `apps/server`。

### 5.3 docker-compose.yml

```yaml
  mysql:
    volumes:
      - mysql-data:/var/lib/mysql
      - ./apps/server/SQL:/docker-entrypoint-initdb.d

  backend:
    build:
      context: ./apps/server
      dockerfile: Dockerfile

  frontend:
    build:
      context: .
      dockerfile: apps/web/Dockerfile
```

### 5.4 Makefile

```makefile
dev-frontend-local:
	pnpm dev:web
dev-backend-local:
	pnpm dev:server
```

## 6. 验证清单

```bash
pnpm install
pnpm build:web            # 生成 apps/web/dist
pnpm build:server         # 生成 apps/server/bin/gin-admin
go -C apps/server vet ./...
pnpm dev:server           # 另开终端 pnpm dev:web，确认能登录
docker compose build      # 如使用 Docker
git status
```

- `git status` 中应大量出现 `renamed:`，说明历史保留。
- 确认没有 `node_modules`、`bin/`、`GinAdmin` 被纳入版本管理。
- Vercel 构建日志顶部的 Commit 号应为最新提交。

## 7. 日常开发约定

所有命令在仓库根目录执行：

| 要做的事 | 命令 |
|---|---|
| 安装依赖 | `pnpm install` |
| 启动前端 | `pnpm dev:web` |
| 启动后端 | `pnpm dev:server` |
| 前端加依赖 | `pnpm --filter gin-admin-web add <包名>` |
| 前端加开发依赖 | `pnpm --filter gin-admin-web add -D <包名>` |
| 后端加依赖 | `go -C apps/server get <包名>` |
| 构建前端 | `pnpm build:web` |
| 构建后端 | `pnpm build:server` |

不要 `cd apps/web` 后执行 `pnpm install`。

## 8. 后续扩展

- **共享 TS 代码**（如接口类型）：建 `packages/shared/`，自动纳入 workspace，在 `apps/web/package.json` 中用 `"@gin-admin/shared": "workspace:*"` 引用。
- **构建编排**：应用和共享包变多、需要构建缓存和增量构建时，再引入 Turborepo。

## 9. 踩坑记录

### 9.1 `Project framework is set to "services", but no services are declared`

**原因**：首次导入项目时 `vercel.json` 带有 `services` 块，Vercel 把 Framework Preset 保存成了 Services，之后删掉 `services` 两边对不上。

**解决**：`vercel.json` 中写 `"framework": null`（或控制台改为 Other / 对应框架）。

### 9.2 `ERR_INVALID_THIS` / `Value of "this" must be of type URLSearchParams`

**原因**：执行安装的是 Vercel 预装的旧版 pnpm（6.x 一类），在 Node 20+ 上会报这个错。`npm install -g pnpm@x` 装的新版本在 `PATH` 中排在预装版本之后，没有生效。

**根本原因**：Vercel 根据**项目根目录**的 lockfile 判断 pnpm 版本，而当时 lockfile 在 `web/` 下，根目录没有。

### 9.3 `Ignoring not compatible lockfile` / `Headless installation requires a pnpm-lock.yaml`

**原因**：同 9.2，旧版 pnpm 读不了 `lockfileVersion: '9.0'`。`corepack enable` 也没生效，因为 Vercel 只读根目录 `package.json` 的 `packageManager`。

**解决**：改为 pnpm workspace，lockfile 放到根目录，Vercel 根据 lockfile 版本自动选择 pnpm 9。

### 9.4 `build:web: command not found`

**原因**：`buildCommand` 写成了 `"build:web"`，会被 shell 当作可执行文件。

**解决**：写成 `"pnpm build:web"`。

### 9.5 `vite: command not found` / `node_modules missing`

**原因**：漏建 `pnpm-workspace.yaml`，pnpm 没识别到子包，本地 `pnpm install` 生成了空 lockfile（`importers: .: {}`）并覆盖了原来的 lockfile。

**解决**：补上 `pnpm-workspace.yaml`，从 git 历史恢复原 lockfile（`git show <commit>:web/pnpm-lock.yaml > pnpm-lock.yaml`），再 `pnpm install`。提交时确认 workspace 文件已包含。

### 9.6 `use of internal package GinAdmin/internal/routers not allowed`

**原因**：Vercel 的 Go 构建器会把 `api/index.go` 改写成 `handler/api` 包导入（日志中的 `imports handler/api`）。Go 规定 `GinAdmin/internal/...` 只能被 `GinAdmin/` 下的代码导入，`handler/api` 不满足。本地能编译是因为本地路径是 `GinAdmin/api`。

**解决**：新增非 internal 的中转包 `app/app.go`，入口 `api/index.go` 只导入 `GinAdmin/app`。

## 10. 已知限制：大文件上传

Vercel 函数没有持久磁盘（只有实例私有的 `/tmp`），请求体上限约 4.5MB，且有执行时长限制。当前「分片写本地磁盘再合并」的方案在 Vercel 上不可靠，只能保证接口不崩溃。

后续计划：改为对象存储预签名直传（推荐 Cloudflare R2，S3 兼容）。分片不再经过后端函数，后端只负责 `init`（CreateMultipartUpload）、签名、`status`（ListParts）、`complete`（CompleteMultipartUpload），哈希秒传、并发、重试、暂停等前端逻辑保持不变。