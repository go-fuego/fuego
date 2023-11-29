dev:
	npx tailwindcss -i ./tailwind.css -o ./static/tailwind.min.css --watch &
	air -- -debug

preview: build
	./bin/gourmet

build:
	npx tailwindcss -i ./tailwind.css -o ./static/tailwind.min.css --minify
	go build -o bin/gourmet -v .

deploy: build backup
	operations/deploy.sh && open https://gourmet.quimerch.com

backup:
	operations/backup.sh
