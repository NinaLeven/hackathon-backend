package main

import (
	"context"
	resequip "github.com/WantsToFress/hackathon-backend/pkg"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (es *IncidentService) getPersonByLogin(ctx context.Context, login string) (*resequip.Person, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (es *IncidentService) GetPerson(ctx context.Context, r *resequip.Id) (*resequip.Person, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (es *IncidentService) WhoAmI(ctx context.Context, r *empty.Empty) (*resequip.Person, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}
