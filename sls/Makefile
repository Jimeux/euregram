.PHONY: clean build build-go-% build-js-% deploy destroy deploy-% log-% destroy
SHELL := /bin/bash



UPLOAD_BUCKET_NAME := $(UPLOAD_BUCKET_NAME)
IMAGE_BUCKET_NAME := $(IMAGE_BUCKET_NAME)
FRONTEND_BUCKET_NAME := $(FRONTEND_BUCKET_NAME)
LOG_BUCKET_NAME := $(LOG_BUCKET_NAME)
DISTRIBUTION_ID := $(DISTRIBUTION_ID)



#############################################
############## SET UP ENV VARS ##############
#############################################
gen-env:
	@printf "%s\n" \
"# Set these AFTER running make deploy and BEFORE running make deploy-frontend" \
"FRONTEND_BUCKET_NAME=See Outputs.FrontendBucket" \
"IMAGE_BUCKET_NAME=See Outputs.ImageBucket" \
"LOG_BUCKET_NAME=See Outputs.LogBucket" \
"UPLOAD_BUCKET_NAME=See Outputs.UploadBucket" \
"DISTRIBUTION_ID=See Outputs.Distribution" \
> .env
	@echo -e "\nSuccessfully wrote to file .env"
	@printf "%s\n" \
"# Set these AFTER running make deploy and BEFORE running make deploy-frontend" \
"VITE_BASE_URL=Should match SSM param euregram-dev-domain" \
> ../client/.env
	@echo -e "\nSuccessfully wrote to file ../client/env\n"



##################################################################################
################################## SERVERLESS ####################################
##################################################################################

build: build-go build-js

build-go:
	@echo "Building all Go handlers:"
	@for fn in authenticator authorizer lister persister presigner; do \
		echo -n "	- building $${fn}..."; \
		env GOARCH=amd64 GOOS=linux go1.18rc1 build -ldflags="-s -w" -o bin/$${fn}/$${fn} handlers/$${fn}/$${fn}.go || exit 1; \
		echo "	done"; \
	done

build-js:
	@echo "Building all JS handlers:"
	@echo -n "	- building validator...";
	@cp -r handlers/validator/ bin/validator/
	@echo "	done"
	@echo -n "	- building resizer...";
	@echo "	done"
	@cp -r handlers/resizer/ bin/resizer/

build-go-*:
	env GOARCH=amd64 GOOS=linux go1.18rc1 build -ldflags="-s -w" -o bin/$*/$* handlers/$*/$*.go

deploy-go-*:
	env GOARCH=amd64 GOOS=linux go1.18rc1 build -ldflags="-s -w" -o bin/$*/$* handlers/$*/$*.go
	sls deploy -f $*

build-js-*:
	cp -r handlers/$*/ bin/$*/

log-%:
	sls logs -f $* --startTime 2m

clean:
	@echo -n "Cleaning bin directory... "
	@rm -rf ./bin
	@echo "done"

deploy: clean build
	@echo "Starting deployment:"
	@sls deploy --verbose

destroy: delete-buckets
	sls remove



##################################################################################
################################### FRONTEND #####################################
##################################################################################

deploy-frontend:
	@cd ../client && npm run build && cd -
	aws s3 sync --region us-east-1 ../client/dist s3://$(FRONTEND_BUCKET_NAME) --
	aws cloudfront create-invalidation --distribution-id $(DISTRIBUTION_ID) --paths "/index.html"

delete-buckets:
	aws s3 rm s3://$(FRONTEND_BUCKET_NAME) --recursive
	aws s3 rm s3://$(IMAGE_BUCKET_NAME) --recursive
	aws s3 rm s3://$(UPLOAD_BUCKET_NAME) --recursive
	aws s3 rm s3://$(LOG_BUCKET_NAME) --recursive
