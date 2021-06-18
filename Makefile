build-app:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o "build/url-collector" -ldflags '-w'
