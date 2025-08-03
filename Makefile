# Путь до ClickHouse миграции
CLICKHOUSE_SQL=./event-worker/migrations/init_clickhouse.sql

# Название контейнера ClickHouse
CLICKHOUSE_CONTAINER=taskflow-clickhouse

# ==== Docker ====
up:
	docker compose up --build

down:
	docker compose down -v

restart: down up

logs:
	docker compose logs -f

logs-%:
	docker compose logs -f $*

# ==== ClickHouse ====
migrate-clickhouse:
	docker exec -i $(CLICKHOUSE_CONTAINER) clickhouse-client < $(CLICKHOUSE_SQL)

clickhouse-shell:
	docker exec -it $(CLICKHOUSE_CONTAINER) clickhouse-client

# ==== PostgreSQL (task-service) ====
psql-task:
	docker exec -it taskflow-task-service psql -U postgres -d taskflow

migrate-task:
	docker exec -i taskflow-task-service psql -U postgres -d taskflow < ./task-service/migrations/init.sql

# ==== Утилиты ====
chmod-entrypoint:
	chmod +x ./event-worker/migrations/clickhouse-entrypoint.sh
