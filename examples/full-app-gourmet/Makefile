dev:
	air -- -debug

prepare-db: backup
	cp gourmet.bak.db recipe.db

db-upload:
	operations/upload-db.sh

preview: build
	./bin/gourmet

prettier:
	npx prettier --write .

css:
	tailwindcss -i ./tailwind.css -o ./static/tailwind.min.css --minify

css-watch:
	tailwindcss -i ./tailwind.css -o ./static/tailwind.min.css --minify --watch 


build: css
	go generate ./...
	go build -v -ldflags="-s -w" -o gourmet-app

# Build for Prod
docker-build:
	go generate ./...
	docker build --platform linux/amd64 -t ewenquim/gourmet .

docker-push:
	docker push ewenquim/gourmet

docker-build-and-push: docker-build docker-push

docker-preview:
	docker run --env-file .env --rm -p 8083:8083 ewenquim/gourmet

deploy:
	make backup & GOARCH=amd64 GOOS=linux make build
	operations/deploy.sh
	open https://gourmet.quimerch.com
	make logs

logs:
	operations/logs.sh

backup:
	operations/backup.sh
