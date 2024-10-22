run_app:
	@echo "Setting CONFIG_PATH and starting application..."
	export CONFIG_PATH="config/local1.yaml"; \
    docker-compose up -d; \
	go run internal/cmd/main.go; \
	sudo nginx -t -c nginx/nginx.conf
    
clear:
	@docker-compose down; \
	sudo pkill nginx || true; \
	echo "Clean up done."
