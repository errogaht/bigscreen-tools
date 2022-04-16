// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/errogaht/bigscreen-tools/repo"
)

import (
	_ "github.com/joho/godotenv/autoload"
)

// Injectors from wire.go:

func InitializeRoomRepo() *repo.RoomRepo {
	pgxConn := NewConn()
	accountProfileRepo := InitializeAccountProfileRepo()
	oculusProfileRepo := repo.NewOculusProfileRepo(pgxConn)
	steamProfileRepo := repo.NewSteamProfileRepo(pgxConn)
	roomUsersRepo := repo.NewRoomUsersRepo(pgxConn)
	settingsRepo := repo.NewSettingsRepo(pgxConn)
	roomRepo := repo.NewRoomRepo(pgxConn, accountProfileRepo, oculusProfileRepo, steamProfileRepo, roomUsersRepo, settingsRepo)
	return roomRepo
}

func InitializeAccountProfileRepo() *repo.AccountProfileRepo {
	pgxConn := NewConn()
	oculusProfileRepo := repo.NewOculusProfileRepo(pgxConn)
	steamProfileRepo := repo.NewSteamProfileRepo(pgxConn)
	accountProfileRepo := repo.NewAccountProfileRepo(pgxConn, oculusProfileRepo, steamProfileRepo)
	return accountProfileRepo
}

func InitializeSteamProfileRepo() *repo.SteamProfileRepo {
	pgxConn := NewConn()
	steamProfileRepo := repo.NewSteamProfileRepo(pgxConn)
	return steamProfileRepo
}

func InitializeRoomUsersRepo() *repo.RoomUsersRepo {
	pgxConn := NewConn()
	roomUsersRepo := repo.NewRoomUsersRepo(pgxConn)
	return roomUsersRepo
}

func InitializeOculusProfileRepo() *repo.OculusProfileRepo {
	pgxConn := NewConn()
	oculusProfileRepo := repo.NewOculusProfileRepo(pgxConn)
	return oculusProfileRepo
}

func InitializeSettingsRepo() *repo.SettingsRepo {
	pgxConn := NewConn()
	settingsRepo := repo.NewSettingsRepo(pgxConn)
	return settingsRepo
}
