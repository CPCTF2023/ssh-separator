// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/google/wire"
	"github.com/mazrean/separated-webshell/api"
	"github.com/mazrean/separated-webshell/repository"
	"github.com/mazrean/separated-webshell/repository/badger"
	"github.com/mazrean/separated-webshell/service"
	"github.com/mazrean/separated-webshell/ssh"
	"github.com/mazrean/separated-webshell/store"
	"github.com/mazrean/separated-webshell/store/gomap"
	"github.com/mazrean/separated-webshell/workspace"
	"github.com/mazrean/separated-webshell/workspace/docker"
)

// Injectors from wire.go:

func InjectServer() (*Server, func(), error) {
	workspace, err := docker.NewWorkspace()
	if err != nil {
		return nil, nil, err
	}
	gomapWorkspace := gomap.NewWorkspace()
	db, cleanup, err := badger.NewDB()
	if err != nil {
		return nil, nil, err
	}
	transaction := badger.NewTransaction(db)
	user := badger.NewUser(db)
	setup := service.NewSetup(workspace, gomapWorkspace, transaction, user)
	serviceUser := service.NewUser(workspace, gomapWorkspace, user, transaction)
	apiUser := api.NewUser(serviceUser)
	apiAPI := api.NewAPI(apiUser)
	workspaceConnection := docker.NewWorkspaceConnection()
	pipe := service.NewPipe(gomapWorkspace, workspaceConnection, workspace)
	sshSSH := ssh.NewSSH(serviceUser, pipe)
	server, err := NewServer(setup, apiAPI, sshSSH)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	return server, func() {
		cleanup()
	}, nil
}

// wire.go:

var (
	transactionBind         = wire.Bind(new(repository.ITransaction), new(*badger.Transaction))
	storeWorkspaceBind      = wire.Bind(new(store.IWorkspace), new(*gomap.Workspace))
	repositoryUserBind      = wire.Bind(new(repository.IUser), new(*badger.User))
	workspaceBind           = wire.Bind(new(workspace.IWorkspace), new(*docker.Workspace))
	workspaceConnectionBind = wire.Bind(new(workspace.IWorkspaceConnection), new(*docker.WorkspaceConnection))
	serviceUserBind         = wire.Bind(new(service.IUser), new(*service.User))
	servicePipeBind         = wire.Bind(new(service.IPipe), new(*service.Pipe))
)

type Server struct {
	*service.Setup
	*api.API
	*ssh.SSH
}

func NewServer(setup *service.Setup, a *api.API, s *ssh.SSH) (*Server, error) {
	return &Server{
		Setup: setup,
		API:   a,
		SSH:   s,
	}, nil
}
