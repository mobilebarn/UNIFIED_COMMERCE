package graphql

import (
	"context"
	"errors"
	"strings"

	"github.com/99designs/gqlgen/plugin/federation/fedruntime"
)

func (ec *executionContext) __resolve__service(ctx context.Context) (fedruntime.Service, error) {
	if ec.DisableIntrospection {
		return fedruntime.Service{}, errors.New("federated introspection disabled")
	}

	var sdl []string

	for _, src := range sources {
		if src.BuiltIn {
			continue
		}
		sdl = append(sdl, src.Input)
	}

	return fedruntime.Service{
		SDL: strings.Join(sdl, "\n"),
	}, nil
}

func (ec *executionContext) __resolve_entities(ctx context.Context, representations []map[string]interface{}) []fedruntime.Entity {
	// For now, return empty slice to avoid compilation errors
	// This can be expanded later when entity resolution is needed
	return []fedruntime.Entity{}
}