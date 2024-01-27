run-prepare:
	@cd server && docker-compose up -d
	@cd server/script && go run populate_db.go
	@cd server && go mod vendor
	@cd client && npm install

run-server:
	@go run server/main.go

run-client:
	@cd client && npm run dev