package main

import (
	"context"
	"github.com/WantsToFress/hackathon-backend/internal/model"
	resequip "github.com/WantsToFress/hackathon-backend/pkg"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func modelToPerson(person *model.Person) *resequip.Person {
	res := &resequip.Person{
		Id:        person.ID,
		Login:     person.Login,
		Email:     person.Email,
		FullName:  person.FullName,
		Role:      resequip.Role(resequip.Role_value[person.Role]),
		ManagerId: ptrToStringWrapper(person.ManagerID).GetValue(),
	}
	return res
}

func (is *IncidentService) getPersonByLogin(ctx context.Context, login string) (*resequip.Person, error) {
	log := loggerFromContext(ctx)

	person := &model.Person{}

	err := is.db.ModelContext(ctx, person).
		Where(model.Columns.Person.Login+" = ?", login).
		Select()
	if err != nil {
		log.WithError(err).Error("unable to select person")
		return nil, status.Error(codes.Internal, "unable to select person")
	}

	return modelToPerson(person), nil
}

func (is *IncidentService) getPersonByEmail(ctx context.Context, email string) (*resequip.Person, error) {
	log := loggerFromContext(ctx)

	person := &model.Person{}

	err := is.db.ModelContext(ctx, person).
		Where(model.Columns.Person.Email+" = ?", email).
		Select()
	if err != nil {
		log.WithError(err).Error("unable to select person")
		return nil, status.Error(codes.Internal, "unable to select person")
	}

	return modelToPerson(person), nil
}

func (is *IncidentService) GetPerson(ctx context.Context, r *resequip.Id) (*resequip.Person, error) {
	log := loggerFromContext(ctx)

	person := &model.Person{}

	err := is.db.ModelContext(ctx, person).
		Where(model.Columns.Person.ID+" = ?", r.GetId()).
		Select()
	if err != nil {
		log.WithError(err).Error("unable to select person")
		return nil, status.Error(codes.Internal, "unable to select person")
	}

	return modelToPerson(person), nil
}

func (is *IncidentService) WhoAmI(ctx context.Context, r *empty.Empty) (*resequip.Person, error) {
	user, err := userFromContext(ctx)
	if err != nil {
		return nil, status.Error(codes.PermissionDenied, err.Error())
	}
	return is.GetPerson(ctx, &resequip.Id{Id: user.Id})
}
