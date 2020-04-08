build: 
	  CGO_ENABLED=0 GOOS=linux go build -o go4translator vcap/main.go
docker:
	  docker build -t qibobo/go4translator_vcap:latest -f Dockerfile.vcap . && docker push qibobo/go4translator_vcap:latest
