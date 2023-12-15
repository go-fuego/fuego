dev:
	npx tailwindcss -i ./tailwind.css -o ./static/tailwind.min.css --watch &
	air -- -debug

prepare-db: backup
	cp gourmet.bak.db recipe.db

db-upload:
	operations/upload-db.sh

preview: build
	./bin/gourmet

prettier:
	npx prettier --write .

build: 
	npx tailwindcss -i ./tailwind.css -o ./static/tailwind.min.css --minify
	go build -v -ldflags="-s -w" -o gourmet-app

deploy:
	make backup & GOARCH=amd64 GOOS=linux make build
	operations/deploy.sh
	open https://gourmet.quimerch.com
	make logs

logs:
	operations/logs.sh

backup:
	operations/backup.sh
