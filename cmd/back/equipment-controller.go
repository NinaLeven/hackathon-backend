package main

import (
	"context"
	"time"

	"github.com/go-pg/pg/v9"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/WantsToFress/hackathon-backend/internal/model"
	resequip "github.com/WantsToFress/hackathon-backend/pkg"
)

const (
	approvalPrice = 100000
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
		query.Where(model.Columns.Equipment.Name+" ilike concat('%', ?::text, '%')", r.GetSearch().GetValue())
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
		CreatedAt: timeToTimestamp(equipment.CreatedAt),
		Deadline:  timeToTimestamp(equipment.Deadline),
		PersonId:  equipment.PersonID,
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
		query.Where(model.Columns.Equipment.Name+" ilike concat('%', ?::text, '%')", r.GetSearch().GetValue())
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

func (is *IncidentService) createEquipmentAssignment(ctx context.Context, tx *pg.Tx, incidentId, personId string) error {
	log := loggerFromContext(ctx)

	ai := &model.EquipmentIncident{}

	err := tx.ModelContext(ctx, ai).
		Where(model.Columns.EquipmentIncident.IncidentID+" = ?", incidentId).
		Select()
	if err != nil {
		log.WithError(err).Error("unable to select equipment incident")
		return status.Error(codes.Internal, "unable to select equipment incident")
	}

	ea := &model.EquipmentAssignment{
		ID:          model.GenStringUUID(),
		PersonID:    personId,
		EquipmentID: ai.EquipmentID,
		Deadline:    ai.Deadline,
		CreatedAt:   time.Now(),
	}
	_, err = tx.ModelContext(ctx, ea).
		Insert()
	if err != nil {
		log.WithError(err).Error("unable to create equipment assignment")
		return status.Error(codes.Internal, "unable to create equipment assignment")
	}

	return nil
}

func (is *IncidentService) CreateEquipment(ctx context.Context, r *resequip.EquipmentCreate) (*empty.Empty, error) {
	log := loggerFromContext(ctx)

	equipment := &model.Equipment{
		ID:          model.GenStringUUID(),
		Name:        r.GetName().GetValue(),
		Description: r.GetDescription().GetValue(),
		Price:       int(r.GetPrice().GetValue()),
	}
	_, err := is.db.ModelContext(ctx, equipment).
		Insert()
	if err != nil {
		log.WithError(err).Error("unable to create equipment")
		return nil, status.Error(codes.Internal, "unable to create equipment")
	}

	return &empty.Empty{}, nil
}

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
