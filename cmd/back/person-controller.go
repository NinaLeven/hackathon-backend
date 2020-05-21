package main

import (
	"context"
	resequip "github.com/WantsToFress/hackathon-backend/pkg"
	"github.com/golang/protobuf/ptypes/empty"
)

func (is *IncidentService) getPersonByLogin(ctx context.Context, login string) (*resequip.Person, error) {
	return is.GetPerson(ctx, nil)
}

func (is *IncidentService) GetPerson(ctx context.Context, r *resequip.Id) (*resequip.Person, error) {
	return &resequip.Person{
		Id:       "34b3af0d-6203-4603-ba2e-74c3c1edd19e",
		Login:    "vpbukhi",
		Email:    "vpbukhti@sas.kek",
		FullName: "Бухйтичук Владимир Кекович",
		Role:     1,
	}, nil
}

func (is *IncidentService) WhoAmI(ctx context.Context, r *empty.Empty) (*resequip.Person, error) {
	return is.GetPerson(ctx, nil)
}
