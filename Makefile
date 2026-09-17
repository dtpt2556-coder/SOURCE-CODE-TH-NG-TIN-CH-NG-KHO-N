# ==============================================================================
# TenPoint — Makefile vận hành
# Mọi lệnh chạy từ THƯ MỤC GỐC dự án.   Xem toàn bộ target:  make help
# ==============================================================================

SHELL := /bin/bash
.DEFAULT_GOAL := help

ROOT        := $(patsubst %/,%,$(dir $(abspath $(lastword $(MAKEFILE_LIST)))))
COMPOSE_F   := $(ROOT)/deploy/docker-compose.yml
COMPOSE_DEV_F := $(ROOT)/deploy/docker-compose.dev.yml
ENV_FILE    := $(ROOT)/deploy/.env

COMPOSE     := docker compose -f $(COMPOSE_F)
COMPOSE_DEV := docker compose -f $(COMPOSE_F) -f $(COMPOSE_DEV_F)

# Nạp deploy/.env (nếu có) để dùng POSTGRES_USER / POSTGRES_DB trong psql & backup.
-include $(ENV_FILE)

POSTGRES_USER ?= tenpoint
POSTGRES_DB   ?= tenpoint
BACKUP_DIR    := $(ROOT)/backups
STAMP         := $(shell date +%Y%m%d-%H%M%S)

# Màu cho `make help`
C_RESET := \033[0m
C_BOLD  := \033[1m
C_CYAN  := \033[36m
C_GREEN := \033[32m
C_DIM   := \033[2m

.PHONY: help up down dev logs ps build rebuild migrate ingest psql \
        test-be test-fe lint backup-db restore-db clean env-check

# ------------------------------------------------------------------------------
##@ Trợ giúp
# ------------------------------------------------------------------------------

help: ## Hiện danh sách lệnh (mặc định)
	@printf "\n$(C_BOLD)TenPoint$(C_RESET) — nền tảng tổng hợp tin chứng khoán VN\n"
	@printf "$(C_DIM)Cách dùng: make <target>$(C_RESET)\n"
	@awk 'BEGIN {FS = ":.*##"} \
		/^##@/ { printf "\n  $(C_BOLD)%s$(C_RESET)\n", substr($$0, 5); next } \
		/^[a-zA-Z_-]+:.*?##/ { printf "    $(C_CYAN)%-14s$(C_RESET) %s\n", $$1, $$2 }' \
		$(MAKEFILE_LIST)
	@printf "\n  $(C_BOLD)Chạy lần đầu$(C_RESET)\n"
	@printf "    $(C_GREEN)cp deploy/.env.example deploy/.env  &&  make build  &&  make dev$(C_RESET)\n\n"

env-check: ## Kiểm tra deploy/.env đã tồn tại chưa
	@if [ ! -f "$(ENV_FILE)" ]; then \
		printf "  \033[33m!\033[0m Chưa có deploy/.env — đang tạo từ deploy/.env.example\n"; \
		cp "$(ROOT)/deploy/.env.example" "$(ENV_FILE)"; \
		printf "  \033[33m!\033[0m ĐỌC LẠI deploy/.env và đổi POSTGRES_PASSWORD + ADMIN_TOKEN trước khi lên production.\n"; \
	fi

# ------------------------------------------------------------------------------
##@ Vòng đời hệ thống
# ------------------------------------------------------------------------------

up: env-check ## Khởi động toàn bộ (cấu hình production)
	$(COMPOSE) up -d
	@$(MAKE) --no-print-directory ps

dev: env-check ## Khởi động chế độ dev (mở cổng 5432/8080/3000, tắt cron, log verbose)
	$(COMPOSE_DEV) up -d
	@$(MAKE) --no-print-directory ps
	@printf "\n  web  → http://localhost:3000\n  api  → http://localhost:8080/healthz\n  caddy→ http://localhost\n\n"

down: ## Dừng và xoá container (GIỮ NGUYÊN dữ liệu trong volume)
	$(COMPOSE_DEV) down --remove-orphans

ps: ## Xem trạng thái các service
	@$(COMPOSE) ps

logs: ## Xem log realtime (make logs S=api để lọc 1 service)
	$(COMPOSE) logs -f --tail=200 $(S)

# ------------------------------------------------------------------------------
##@ Build
# ------------------------------------------------------------------------------

build: env-check ## Build image api + web (có dùng cache)
	$(COMPOSE) build

rebuild: env-check ## Build lại từ đầu, bỏ qua toàn bộ cache
	$(COMPOSE) build --no-cache --pull

# ------------------------------------------------------------------------------
##@ Dữ liệu & pipeline
# ------------------------------------------------------------------------------

migrate: env-check ## Chạy migration (embed trong binary → khởi động lại api là chạy)
	$(COMPOSE) up -d db
	$(COMPOSE) up -d --force-recreate api
	@sleep 3
	$(COMPOSE) logs --tail=60 api

ingest: env-check ## Chạy pipeline ingest THỦ CÔNG một lần (không đợi cron)
	$(COMPOSE) run --rm -T --entrypoint /usr/local/bin/ingest api

psql: ## Mở psql vào database
	$(COMPOSE) exec db psql -U "$(POSTGRES_USER)" -d "$(POSTGRES_DB)"

backup-db: ## Sao lưu DB ra backups/tenpoint-<timestamp>.dump
	@mkdir -p "$(BACKUP_DIR)"
	$(COMPOSE) exec -T db pg_dump -U "$(POSTGRES_USER)" -d "$(POSTGRES_DB)" \
		--format=custom --no-owner --no-privileges \
		> "$(BACKUP_DIR)/tenpoint-$(STAMP).dump"
	@printf "  ✓ Đã lưu: backups/tenpoint-$(STAMP).dump ($$(du -h "$(BACKUP_DIR)/tenpoint-$(STAMP).dump" | cut -f1))\n"

restore-db: ## Khôi phục DB.  Bắt buộc: make restore-db FILE=backups/xxx.dump
	@if [ -z "$(FILE)" ]; then \
		printf "  \033[31m✗\033[0m Thiếu tham số FILE. Ví dụ: make restore-db FILE=backups/tenpoint-20260916-070000.dump\n"; \
		exit 1; \
	fi
	@if [ ! -f "$(FILE)" ]; then printf "  \033[31m✗\033[0m Không thấy file: $(FILE)\n"; exit 1; fi
	@printf "  \033[33m!\033[0m Thao tác này GHI ĐÈ dữ liệu hiện tại của '$(POSTGRES_DB)'. Ctrl-C trong 5s để huỷ.\n"
	@sleep 5
	$(COMPOSE) exec -T db pg_restore -U "$(POSTGRES_USER)" -d "$(POSTGRES_DB)" \
		--clean --if-exists --no-owner --no-privileges < "$(FILE)"
	@printf "  ✓ Khôi phục xong từ $(FILE)\n"

# ------------------------------------------------------------------------------
##@ Chất lượng
# ------------------------------------------------------------------------------

test-be: ## Chạy unit test backend (Go)
	@if command -v go >/dev/null 2>&1; then \
		cd "$(ROOT)/backend" && go test ./... ; \
	else \
		docker run --rm -v "$(ROOT)/backend":/src -w /src -e CGO_ENABLED=0 \
			golang:1.26-alpine go test ./... ; \
	fi

test-fe: ## Chạy test frontend (bỏ qua nếu chưa có script test)
	@if command -v npm >/dev/null 2>&1; then \
		cd "$(ROOT)/frontend" && npm test --if-present ; \
	else \
		docker run --rm -v "$(ROOT)/frontend":/app -w /app node:22-alpine \
			sh -c "npm ci --no-audit --no-fund && npm test --if-present" ; \
	fi

lint: ## Lint cả backend (go vet + gofmt) và frontend (eslint)
	@printf "\n$(C_BOLD)→ backend$(C_RESET)\n"
	@if command -v go >/dev/null 2>&1; then \
		cd "$(ROOT)/backend" && go vet ./... && \
		out=$$(gofmt -l .); if [ -n "$$out" ]; then printf "  gofmt cần chạy lại trên:\n%s\n" "$$out"; exit 1; fi ; \
	else printf "  (bỏ qua — chưa cài Go)\n"; fi
	@printf "\n$(C_BOLD)→ frontend$(C_RESET)\n"
	@if command -v npm >/dev/null 2>&1; then \
		cd "$(ROOT)/frontend" && npm run lint ; \
	else printf "  (bỏ qua — chưa cài Node)\n"; fi

# ------------------------------------------------------------------------------
##@ Dọn dẹp
# ------------------------------------------------------------------------------

clean: ## ⚠️ Xoá container + network + VOLUME (MẤT TOÀN BỘ DỮ LIỆU DB)
	@printf "  \033[31m!\033[0m Sắp XOÁ volume tenpoint_pgdata (mất sạch dữ liệu). Ctrl-C trong 5s để huỷ.\n"
	@sleep 5
	$(COMPOSE_DEV) down -v --remove-orphans
	docker image prune -f --filter "label=com.docker.compose.project=tenpoint"
	@printf "  ✓ Đã dọn sạch.\n"
