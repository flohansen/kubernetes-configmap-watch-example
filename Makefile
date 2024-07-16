build-image-local:
	docker build -t importer:latest -f build/importer.dockerfile .
	minikube image load importer:latest

generate-configmap:
	python scripts/generate-configmap.py
