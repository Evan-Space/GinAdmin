help:
	@echo "可用命令:"
	@echo "  make up         启动所有服务"
	@echo "  make down       停止并删除所有容器"
	@echo "  make build      重新构建镜像（自动清理旧镜像）"
	@echo "  make restart    重启所有容器"
	@echo "  make status     查看容器运行状态"
	@echo "  make logs     实时查看所有日志"
	@echo "  make logs-back  实时查看后端日志"
	@echo "  make dev-frontend-local  启动本地前端项目"
	@echo "  make dev-backend-local  启动本地后端项目"

up:
	docker compose up -d
down:
	docker compose down
build:
	docker rmi $$(docker images -q ginadmin-backend ginadmin-frontend) 2>/dev/null; \
	docker compose build --no-cache && \
	docker image prune -f
# 	docker compose build --no-cache && docker image prune -f
restart:
	docker compose restart
status:
	docker compose ps
logs:
	docker compose logs -f
logs-back:
	docker compose logs -f backend
dev-frontend-local:
	pnpm dev:web
dev-backend-local:
	pnpm dev:server