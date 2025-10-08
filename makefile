run-b:
	cd backend && docker compose up --build

stop-b:
	cd backend && docker compose down -v

run-gestion:
	cd backend && docker compose up --build gestion-db gestion-service

run-doc:
	cd backend && docker compose up --build documentacion-db documentacion-service

run-front:
	cd frontend && npm run dev

SERVICES := \
	documentacion-db \
	gestion-db \
	notification-db \
	oauth2DB \
	mongu \
	influxdb \
	emqx \
	iot-service \
	notification-service \
	gestion-service \
	documentacion-service \
	middleware_mqtt \
	data-sync-init \
	database1 \
	oauth2

# Levanta todo menos apigateway
run-b-wa:
	cd backend && docker compose up --build $(SERVICES)
# Baja todo menos apigateway
stop-b-wa:
	cd backend && docker compose down -v $(SERVICES)

run-b-ag:
	cd backend && docker compose up --build apigateway
	
stop-b-ag:
	cd backend && docker compose stop apigateway || true
	cd backend && docker compose rm -f -s apigateway || true

NGROK_URL := anita-submicroscopic-overgently.ngrok-free.app

run-b-ngrok:
	cd backend && docker compose up --build -d && ngrok http --domain=$(NGROK_URL) 3500

run-ngrok:
	ngrok http --domain=$(NGROK_URL) 3500

err-f:
	cd frontend && npx eslint . --ext .js,.ts,.tsx

ml:
	cd backend/scripts/ML && python3 main.py