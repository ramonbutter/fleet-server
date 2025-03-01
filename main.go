// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License;
// you may not use this file except in compliance with the Elastic License.

//go:generate schema-generate -esdoc -s -cm "{\"Api\": \"API\", \"Id\": \"ID\"}" -o internal/pkg/model/schema.go -p model model/schema.json
//go:generate oapi-codegen --config model/oapi-cfg.yml model/openapi.yml
//go:generate go fmt internal/pkg/model/schema.go
//go:generate go fmt internal/pkg/api/openapi.gen.go

package main

import (
	"fmt"
	"os"

	"github.com/elastic/fleet-server/v7/cmd/fleet"
	"github.com/elastic/fleet-server/v7/internal/pkg/build"
	"github.com/elastic/fleet-server/v7/version"
)

var (
	Version   string = version.DefaultVersion
	Commit    string
	BuildTime string
)

func main() {
	cmd := fleet.NewCommand(build.Info{
		Version:   Version,
		Commit:    Commit,
		BuildTime: build.Time(BuildTime),
	})
	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
