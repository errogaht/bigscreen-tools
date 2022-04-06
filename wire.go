//go:build wireinject

package main

import (
	"github.com/errogaht/bigscreen-tools/repo"
	"github.com/google/wire"
)

func InitializeRoomRepo() *repo.RoomRepo {
	wire.Build(
		NewConn,
		repo.NewRoomRepo,
		repo.NewSteamProfileRepo,
		repo.NewOculusProfileRepo,
		repo.NewSettingsRepo,
		InitializeAccountProfileRepo,
	)
	return &repo.RoomRepo{}
}
func InitializeAccountProfileRepo() *repo.AccountProfileRepo {
	wire.Build(
		NewConn,
		repo.NewSteamProfileRepo,
		repo.NewOculusProfileRepo,
		repo.NewAccountProfileRepo,
	)
	return &repo.AccountProfileRepo{}
}

func InitializeSteamProfileRepo() *repo.SteamProfileRepo {
	wire.Build(
		NewConn,
		repo.NewSteamProfileRepo,
	)
	return &repo.SteamProfileRepo{}
}

func InitializeOculusProfileRepo() *repo.OculusProfileRepo {
	wire.Build(
		NewConn,
		repo.NewOculusProfileRepo,
	)
	return &repo.OculusProfileRepo{}
}

func InitializeSettingsRepo() *repo.SettingsRepo {
	wire.Build(
		NewConn,
		repo.NewSettingsRepo,
	)
	return &repo.SettingsRepo{}
}
