package main

import (
	"context"

	"github.com/go-pg/pg/v9"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/WantsToFress/hackathon-backend/internal/model"
	resequip "github.com/WantsToFress/hackathon-backend/pkg"
)

func modelToEquipment(equipment *model.Equipment) *resequip.Equipment {
	res := &resequip.Equipment{
		Id:          equipment.ID,
		Name:        equipment.Name,
		Description: equipment.Description,
		Price:       int64(equipment.Price),
	}
	return res
}

func (is *IncidentService) ListEquipment(ctx context.Context, r *resequip.EquipmentFilter) (*resequip.EquipmentList, error) {
	log := loggerFromContext(ctx)

	equipment := []*model.Equipment{}

	query := is.db.ModelContext(ctx, &equipment)

	if r.GetSearch() != nil {
		query.Where(model.Columns.Equipment.Name + " ilike concat('%', ?::text, '%')", r.GetSearch().GetValue())
	}

	err := query.Select()
	if err != nil {
		log.WithError(err).Error("unable to select equipment")
		return nil, status.Error(codes.Internal, "unable to select equipment")
	}

	res := make([]*resequip.Equipment, 0, len(equipment))

	for i := range equipment {
		res = append(res, modelToEquipment(equipment[i]))
	}

	return &resequip.EquipmentList{Equipment: res}, nil
}

func modelToAssignedEquipment(equipment *model.EquipmentAssignment) *resequip.AssignedEquipment {
	res := &resequip.AssignedEquipment{
		CreatedAt:            timeToTimestamp(equipment.CreatedAt),
		Deadline:             timeToTimestamp(equipment.Deadline),
		PersonId:             equipment.PersonID,
	}
	if equipment.Equipment != nil {
		res.Equipment = modelToEquipment(equipment.Equipment)
	}
	return res
}

func (is *IncidentService) ListEquipmentForPerson(ctx context.Context, r *resequip.EquipmentFilter) (*resequip.AssignedEquipmentList, error) {
	log := loggerFromContext(ctx)

	equipment := []*model.EquipmentAssignment{}

	query := is.db.ModelContext(ctx, &equipment).
		Relation(model.Columns.EquipmentAssignment.Equipment)

	if r.GetSearch() != nil {
		query.Where(model.Columns.Equipment.Name + " ilike concat('%', ?::text, '%')", r.GetSearch().GetValue())
	}

	err := query.Select()
	if err != nil {
		log.WithError(err).Error("unable to select equipment")
		return nil, status.Error(codes.Internal, "unable to select equipment")
	}

	res := make([]*resequip.AssignedEquipment, 0, len(equipment))

	for i := range equipment {
		res = append(res, modelToAssignedEquipment(equipment[i]))
	}

	return &resequip.AssignedEquipmentList{Equipment: res}, nil
}

const (
	approvalPrice = 100000
)

func (is *IncidentService) equipmentRequiresApproval(ctx context.Context, tx *pg.Tx, equipmentId string) (bool, error) {
	needApproval := false
	err := tx.ModelContext(ctx, (*model.Equipment)(nil)).
		ColumnExpr(model.Columns.Equipment.Price+" >= ? as need_approval", approvalPrice).
		Where(model.Columns.Equipment.ID+" = ?", equipmentId).
		Select(&needApproval)
	if err != nil {
		return false, err
	}
	return needApproval, nil
}
