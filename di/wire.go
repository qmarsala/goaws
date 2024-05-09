//go:build wireinject

package di

import (
	goaws "goaws/internal"

	"github.com/google/wire"
)

func InitializeDatabase() (goaws.DatabaseConnection, error) {
	wire.Build(goaws.ProvideDatabase)
	return goaws.DatabaseConnection{}, nil
}
