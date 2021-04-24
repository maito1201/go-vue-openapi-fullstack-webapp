package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"

	"github.com/go-openapi/loads"
	flags "github.com/jessevdk/go-flags"

	"regexp"
	"github.com/maito1201/go-vue-openapi-fullstack-webapp/server/gen/restapi"
	"github.com/maito1201/go-vue-openapi-fullstack-webapp/server/gen/restapi/factory"
)

var proxyRegexp = regexp.MustCompile(`^/api`)

//go:embed frontend/dist/*
var static embed.FS

func main() {

	swaggerSpec, err := loads.Embedded(restapi.SwaggerJSON, restapi.FlatSwaggerJSON)
	if err != nil {
		log.Fatalln(err)
	}

	api := factory.NewFactoryAPI(swaggerSpec)
	server := restapi.NewServer(api)
	defer server.Shutdown()

	parser := flags.NewParser(server, flags.Default)
	parser.ShortDescription = "Greeting Server"
	parser.LongDescription = swaggerSpec.Spec().Info.Description
	server.ConfigureFlags()
	for _, optsGroup := range api.CommandLineOptionsGroups {
		_, err := parser.AddGroup(optsGroup.ShortDescription, optsGroup.LongDescription, optsGroup.Options)
		if err != nil {
			log.Fatalln(err)
		}
	}

	if _, err := parser.Parse(); err != nil {
		code := 1
		if fe, ok := err.(*flags.Error); ok {
			if fe.Type == flags.ErrHelp {
				code = 0
			}
		}
		os.Exit(code)
	}

	// serve swagger api server.
	server.ConfigureAPI()
	http.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) {
		// remove "/api" fron api path for swagger api
		r.URL.Path = proxyRegexp.ReplaceAllString(r.URL.Path, "")
		server.GetHandler().ServeHTTP(w, r)
	})

	// serve frontend HTML.
	public, err := fs.Sub(static, "frontend/dist")
	if err != nil {
		panic(err)
	}
	http.Handle("/", http.FileServer(http.FS(public)))

	log.Println("listening on localhost:3000...")
	// NOTE: if you want to use another port, you also have to modify app\frontend\src\client-axios\base.ts
	log.Fatal(http.ListenAndServe(":3000", nil))
}
