
run_apps: run_app_1 run_app_2 run_app_3

run_app_1:
	@echo "Setting CONFIG_PATH and starting application..."
	CONFIG_PATH="config/local1.yaml" go run internal/cmd/main.go &

run_app_2:
	@echo "Setting CONFIG_PATH2 and starting application..."
	CONFIG_PATH="config/local2.yaml" go run internal/cmd/main.go &

run_app_3:
	@echo "Setting CONFIG_PATH3 and starting application..."
	CONFIG_PATH="config/local3.yaml" go run internal/cmd/main.go &



clear:
	@docker-compose down; \
	sudo pkill nginx || true; \
	echo "Clean up done."
