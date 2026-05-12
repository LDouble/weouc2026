.PHONY: help check check-backend check-admin-web check-miniapp ci-check ci-check-backend ci-check-admin-web ci-check-miniapp generate-sdk

help:
	@echo "可用命令:"
	@echo "  make check              # 使用本地依赖执行最小校验"
	@echo "  make ci-check           # 使用 CI 同步方式执行最小校验"
	@echo "  make generate-sdk       # 基于 OpenAPI 生成 JS / Dart SDK"

check: check-backend check-admin-web check-miniapp

check-backend:
	cd services/api-server && go test ./...

check-admin-web:
	cd apps/admin-web && npm run build

check-miniapp:
	cd apps/miniapp-wechat && npm run check:syntax

ci-check: ci-check-backend ci-check-admin-web ci-check-miniapp

ci-check-backend:
	cd services/api-server && go test ./...

ci-check-admin-web:
	cd apps/admin-web && npm ci && npm run build

ci-check-miniapp:
	cd apps/miniapp-wechat && npm ci && npm run check:syntax

generate-sdk:
	bash packages/contracts/scripts/generate-sdks.sh
