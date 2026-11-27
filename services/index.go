package services

import (
	"github.com/mbarolo/webp_backend/config"
	"github.com/mbarolo/webp_backend/repository"
)

var sqliteConnection = config.NewSQLiteConnection()

var (
	userRepo = repository.NewUserRepository(sqliteConnection.DB)
)
