.PHONY: build push exec_db install uninstall
# Builds the docker image of the database client and tags it
build:
	docker build . -t ghcr.io/kube-hack/sql-injection

# Pushes the docker image to a container registry
push:
	docker push ghcr.io/kube-hack/sql-injection

# Execs user into the database to test SQL commands
exec_db:
	kubectl exec -it $$(kubectl get pods | awk '/db-sql-injection/ {print $$1}') -- psql -U postgres

# Installs the helm chart
install:
	helm install sql-injection ./chart

# Uninstalls the helm chart
uninstall:
	helm uninstall sql-injection