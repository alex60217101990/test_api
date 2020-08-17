include makefiles/Makefile.colors

compile: 
	@printf "${COLOR_LIGHT_CYAN}Compile service...${COLOR_NC}\n"

DOT := $(shell command -v helm -help 2> /dev/null)
NAMESPACE := "default"
POSTGRES_NAME := "postgres"

install-helm:
ifndef DOT
	@printf "${COLOR_GREEN}Install Helm...${COLOR_NC}\n"
	curl -fsSL -o get_helm.sh https://raw.githubusercontent.com/helm/helm/master/scripts/get-helm-3
	chmod 700 get_helm.sh
	./get_helm.sh
	rm get_helm.sh
else
	@printf "${COLOR_YELLOW}Helm already installed [%s]${COLOR_NC}\n" $(shell helm version)
endif

deploy-postgres:
ifeq ($(strip $(shell helm repo list | grep bitnami 2> /dev/null)),)
	@printf "${COLOR_GREEN}Add Bitnami to Helm charts list...${COLOR_NC}\n"
	helm repo add bitnami https://charts.bitnami.com/bitnami
else
	@printf "${COLOR_YELLOW}Helm Bitnami chart already exists.${COLOR_NC}\n"
endif
ifeq ($(strip $(shell helm list | grep postgres 2> /dev/null)),)
	@printf "${COLOR_GREEN}Deploy PostgreSQL chart...${COLOR_NC}\n"
	helm repo add bitnami https://charts.bitnami.com/bitnami
	#kubectl create configmap init-db-script -n ${NAMESPACE} --from-file=./sql-scripts/init-script.sql -o yaml --dry-run | kubectl apply -f -
	@printf "${COLOR_YELLOW}$(shell kubectl apply -f ./sql-scripts/init.yaml)${COLOR_NC}\n"
	helm upgrade \
		--set initdbScriptsConfigMap=sql-init \
      --namespace ${NAMESPACE} \
      --install \
      --wait \
      --timeout 200s \
      --atomic \
      ${POSTGRES_NAME} bitnami/postgresql
else
	@printf "${COLOR_YELLOW}Helm Bitnami chart already exists.${COLOR_NC}\n"
endif
	
loop:
	@printf "${COLOR_YELLOW}Update configs file.${COLOR_NC}\n"
	while [ $(shell kubectl get secret --namespace ${NAMESPACE} postgres-postgresql -o jsonpath="{.data.postgresql-password}" | base64 --decode) == '']; \
    do \
        printf "\t${COLOR_GRAY}watch update secret state...${COLOR_NC}\n"; \
        sleep 1; \
    done
	sed -i "" "/^\([[:space:]]*db-password: \).*/s//\1$(shell kubectl get secret --namespace ${NAMESPACE} postgres-postgresql -o jsonpath="{.data.postgresql-password}" | base64 --decode)/" ./configs/application.yaml

deploy: deploy-postgres loop

delete-postgres:
	@printf "${COLOR_YELLOW}Uninstall PostgreSQL chart...${COLOR_NC}\n"
	helm delete ${POSTGRES_NAME}

go_lint:
	@printf "${COLOR_YELLOW}Validate golang code...${COLOR_NC}\n"
	docker build -t alex6021710/go-linter:v0.0.1 -f ./dockerfiles/Dockerfile.golint ../ && \
	docker run --rm -it alex6021710/go-linter:v0.0.1 golangci-lint run /app/cmd/app 

MYDIR = ./dockerfiles
docker_lint: $(MYDIR)/Dockerfile.*
	@printf "${COLOR_YELLOW}Validate dockerfiles...${COLOR_NC}\n"
	for file in $^ ; do \
		docker run --rm -i hadolint/hadolint < $${file} ; \
	done

lint: docker_lint go_lint

go_push:
	@printf "${COLOR_YELLOW}Build app image...${COLOR_NC}\n"
	docker build -t alex6021710/test-api:v0.0.8 -f ./dockerfiles/Dockerfile.golang ../ && \
	docker login && \
	docker push alex6021710/test-api:v0.0.8  


# helm template -f ./test-api/values.yaml ./test-api/





