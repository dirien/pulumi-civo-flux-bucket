CIVO_DIR=infrastructure/00-civo
FLUX_DIR=infrastructure/01-flux

STACK_NAME=dev

.PHONY: destroy
destroy:
	@echo "Destroying..."
	@pulumi destroy --cwd $(CIVO_DIR) -y
	@pulumi stack rm $(STACK_NAME) --cwd $(FLUX_DIR) --force -y

.PHONY: bootstrap
bootstrap:
	@echo "Bootstrapping Civo..."
	cd $(CIVO_DIR) && go mod tidy
	@pulumi up --cwd $(CIVO_DIR) --skip-preview -y
	@pulumi stack output kubeconfig --show-secrets --cwd $(CIVO_DIR) > kubeconfig.yaml

	@echo "Bootstrapping Flux..."
	cd $(FLUX_DIR) && go mod tidy
	@pulumi up --cwd $(FLUX_DIR) --skip-preview -y

.PHONY: check-bucket
check-bucket:
	@flux get sources bucket --kubeconfig=kubeconfig.yaml

.PHONY: upload-aws
upload-aws:
	$(eval AWS_ACCESS_KEY_ID = $(shell pulumi stack output accessKey --show-secrets --cwd $(CIVO_DIR)))
	$(eval AWS_SECRET_ACCESS_KEY = $(shell pulumi stack output secretKey --show-secrets --cwd $(CIVO_DIR)))
	$(eval BUCKET = $(shell pulumi stack output bucket --cwd $(CIVO_DIR)))

	@echo "Run following commands to upload deploy folder to the S3 bucket"
	@echo "export AWS_ACCESS_KEY_ID=$(AWS_ACCESS_KEY_ID)"
	@echo "export AWS_SECRET_ACCESS_KEY=$(AWS_SECRET_ACCESS_KEY)"
	@echo "aws s3 sync ./deploy/ s3://$(BUCKET)/"
