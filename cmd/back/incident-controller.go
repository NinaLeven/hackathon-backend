package main

import (
	"context"
	"github.com/WantsToFress/hackathon-backend/internal/model"
	resequip "github.com/WantsToFress/hackathon-backend/pkg"
	"github.com/go-pg/pg/v9"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

func (es *IncidentService) createIncident(ctx context.Context, tx *pg.Tx, r *resequip.MaintenanceIncidentCreate, creatorId string, incidentType resequip.IncidentType) error {
	if r.GetDescription() == nil {
		return status.Error(codes.InvalidArgument, "description is mandatory")
	}

	if r.GetPriority() == resequip.IncidentPriority_none_priority {
		return status.Error(codes.InvalidArgument, "priority is mandatory")
	}

	if r.GetDeadline() == nil {
		return status.Error(codes.InvalidArgument, "deadline is mandatory")
	}

	incident := model.Incident{
		ID:          model.GenStringUUID(),
		Description: r.GetDescription().GetValue(),
		CreatedAt:   time.Now(),
		Deadline:    timestampToTime(r.GetDeadline()),
		CreatorID:   creatorId,
		Status:      resequip.IncidentStatus_created.String(),
		Comment:     stringWrapperToPtr(r.GetComment()),
		Type:        incidentType.String(),
		Priority:    int(r.GetPriority()),
	}

	_, err := tx.ModelContext(ctx, incident).Insert()
	if err != nil {
		return err
	}

	return nil
}

func (es *IncidentService) CreateMaintenanceIncident(ctx context.Context, r *resequip.MaintenanceIncidentCreate) (*empty.Empty, error) {
	log := loggerFromContext(ctx)

	user, err := userFromContext(ctx)
	if err != nil {
		return nil, status.Error(codes.PermissionDenied, err.Error())
	}

	tx, err := es.db.Begin()
	if err != nil {
		log.WithError(err).Error("unable to begin transaction")
		return nil, status.Error(codes.Internal, "unable to begin transaction")
	}

	err = es.createIncident(ctx, tx, r, user.Id, resequip.IncidentType_maintenance)
	if err != nil {
		log.WithError(err).Error("unable to create incident")
		terr := tx.Rollback()
		if terr != nil {
			log.WithError(terr).Error("unable to rollback transaction")
		}
		return nil, status.Error(codes.Internal, "unable to create incident")
	}

	err = tx.Commit()
	if err != nil {
		log.WithError(err).Error("unable to commit transaction")
		return nil, status.Error(codes.Internal, "unable to commit transaction")
	}

	return &empty.Empty{}, nil
}

func (es *IncidentService) CreateEquipmentIncident(ctx context.Context, r *resequip.EquipmentIncidentCreate) (*empty.Empty, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (es *IncidentService) ListIncidents(ctx context.Context, r *resequip.IncidentFilter) (*resequip.IncidentList, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (es *IncidentService) AssignIncident(ctx context.Context, r *resequip.AssignmentRequest) (*empty.Empty, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (es *IncidentService) ChangeIncidentStatus(ctx context.Context, r *resequip.IncidentStatusRequest) (*empty.Empty, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (es *IncidentService) CommentOnIncident(ctx context.Context, r *resequip.IncidentCommentRequest) (*empty.Empty, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (es *IncidentService) ApproveEquipmentIncident(ctx context.Context, r *resequip.IncidentApprovalRequest) (*empty.Empty, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}
