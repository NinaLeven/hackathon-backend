package main

import (
	"context"
	resequip "github.com/WantsToFress/hackathon-backend/pkg"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (es *IncidentService) ListEquipment(ctx context.Context, r *resequip.EquipmentFilter) (*resequip.EquipmentList, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (es *IncidentService) ListEquipmentForPerson(ctx context.Context, r *resequip.EquipmentFilter) (*resequip.AssignedEquipmentList, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}
