include .env
export $(shell sed 's/=.*//' .env)

include .override.env
export $(shell sed 's/=.*//' .override.env)

# Services
up:
	@echo ---Initializing all services---
	docker-compose up -d

down:
	docker-compose down

logs:
	docker-compose logs -f

restart:
	docker-compose restart

destroy:
	@echo "=============Cleaning up============="
	docker-compose down -v
	docker-compose rm -f -v -s
build:
	docker-compose build

run:
	go run main.go

# Code
graph:
	go mod graph | modgv | sfdp -Tpng -o graph.png

# Persistence
seed:
	docker-compose exec mongo rm -rf ./seeds
	docker cp ./seeds siempreabierto_mongo_1:./seeds
	docker-compose exec mongo mongoimport --username ${MONGO_USERNAME} --password ${MONGO_PASSWORD} --authenticationDatabase admin --db ${MONGO_DATABASE} --collection link --file ./seeds/link.json --jsonArray
	docker-compose exec mongo mongoimport --username ${MONGO_USERNAME} --password ${MONGO_PASSWORD} --authenticationDatabase admin --db ${MONGO_DATABASE} --collection video --file ./seeds/video.json --jsonArray
	docker-compose exec mongo rm -rf ./seeds

# Cache
redis-attatch:
	docker-compose exec redis bash

redis-insight:
	docker-compose -f ${COMPOSE_3_PARTY_CLI} up -d redis-insight
	@echo redis insight running in port ${REDIS_INSIGHT_PORT}
redis-cli: redis-insight

redis-cli-down:
	docker-compose -f ${COMPOSE_3_PARTY_CLI} down redis-insight
# Swagger
api-doc:
	swagger generate spec -o ./swagger.json
	swagger serve -F swagger ./swagger.json