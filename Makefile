build: 
	  CGO_ENABLED=0 GOOS=linux go build -o go4translator vcap/main.go
