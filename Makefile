SSO_PATH=./cmd/sso/main.go

sso-run:
	go run $(SSO_PATH)

# to use this you need to install https://github.com/michurin/human-readable-json-logging
pp-sso-run:
	pplog go run $(SSO_PATH)
