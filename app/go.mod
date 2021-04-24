module github.com/maito1201/go-vue-openapi-fullstack-webapp/app

go 1.16

replace github.com/maito1201/go-vue-openapi-fullstack-webapp/server => ../server

require (
	github.com/go-openapi/loads v0.20.2
	github.com/jessevdk/go-flags v1.5.0
	github.com/maito1201/go-vue-openapi-fullstack-webapp/server v0.0.0-00010101000000-000000000000
)
