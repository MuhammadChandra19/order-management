run-prepare:
	@cd server && docker-compose up -d
	@cd server && go mod vendor
	@cd client && npm install
	@cd server/script && go run populate_db.go

run-server:
	@go run server/main.go

run-client:
	@cd client && npm run dev