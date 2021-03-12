build:
	mkdir -p bin/
	cd dashboard/ && yarn build;
	rice embed-go
	env GOOS=linux go build -ldflags="-s -w" -trimpath -o bin/engauge
windows:
	mkdir -p bin/
	cd dashboard/ && yarn build;
	rice embed-go
	env GOOS=windows go build -ldflags="-s -w" -trimpath -o bin/engauge.exe