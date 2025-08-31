run-b:
	cd backend && docker compose up --build

stop-b:
	cd backend && docker compose down -v
