
run:
	go run . -debug

watch:
	air -- -debug

deploy: backup
	operations/deploy.sh && open https://gourmet.quimerch.com

backup:
	operations/backup.sh
