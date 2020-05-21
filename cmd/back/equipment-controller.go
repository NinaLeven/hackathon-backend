package main

import (
	"context"
	"github.com/WantsToFress/hackathon-backend/internal/model"
	resequip "github.com/WantsToFress/hackathon-backend/pkg"
	"github.com/go-pg/pg/v9"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (is *IncidentService) ListEquipment(ctx context.Context, r *resequip.EquipmentFilter) (*resequip.EquipmentList, error) {
	return &resequip.EquipmentList{
		Equipment: []*resequip.Equipment{
			{
				Id:          "3baf4d2b-0ed6-4ce8-af72-498cb1d3f5c8",
				Name:        "name",
				Description: "desctiption",
				Price:       1000000,
			},
			{
				Id:          "b86a1864-90f0-45b7-b1c4-60a14ed518f7",
				Name:        "110 montauk",
				Description: "*REDACTED*",
				Price:       1000000,
			},
		},
	}, nil
}

func (is *IncidentService) ListEquipmentForPerson(ctx context.Context, r *resequip.EquipmentFilter) (*resequip.AssignedEquipmentList, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
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
