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