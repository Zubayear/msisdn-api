//go:build wireinject
// +build wireinject

package di_context

import (
	"github.com/google/wire"
	"huspass/api"
	"huspass/external"
	"huspass/repo"
)

func DependencyProvider() (*api.Service, error) {
	wire.Build(repo.DatabaseImplProvider, api.ServiceProvider, external.HostProvider, repo.UserDBImplProvider,
		wire.Bind(new(repo.MsisdnRepository), new(*repo.MsisdnRepositoryImpl)),
		wire.Bind(new(repo.UserRepo), new(*repo.UserRepoImpl)))
	return &api.Service{}, nil
}
