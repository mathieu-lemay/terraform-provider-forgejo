// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"context"
	"log"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mariadb"
)

const (
	networkName        = "forgejo"
	forgejoUserName    = "tfadmin"
	forgejoUserEmail   = "tfadmin@localhost"
	forgejoTokenScopes = "write:organization,write:repository,write:user,write:admin"
)

func main() {
	ctx := context.Background()

	getDBTestContainer(ctx)

	// listDockerContainers(ctx)
}

func getDBTestContainer(ctx context.Context) {
	containerName := "forgejo_db"

	containerRequest := testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Name:     containerName,
			Networks: []string{networkName},
			ConfigModifier: func(config *container.Config) {
				config.Hostname = containerName
			},
		},
		Reuse: true,
	}

	c, err := mariadb.Run(
		ctx,
		"mariadb:lts",
		mariadb.WithDatabase("forgejo"),
		mariadb.WithUsername("forgejo"),
		mariadb.WithPassword("password"),
		testcontainers.CustomizeRequest(containerRequest),
		testcontainers.WithCmd(
			"--transaction-isolation=READ-COMMITTED",
			"--binlog-format=ROW",
		),
	)
	if err != nil {
		log.Fatalf("error starting db container: %v", err)
	}

	log.Printf("Successfully started container: %v", c)

}

func listDockerContainers(ctx context.Context) {
	docker, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		log.Fatalf("error creating docker client: %v", err)
	}
	defer docker.Close()

	containers, err := docker.ContainerList(ctx, container.ListOptions{
		All: true,
	})

	if err != nil {
		log.Fatalf("error getting docker containers: %v", err)
	}

	log.Printf("Found %d containers", len(containers))
}
