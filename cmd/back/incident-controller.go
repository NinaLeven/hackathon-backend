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

func (is *IncidentService) getNextAssignee(ctx context.Context, tx *pg.Tx) (string, error) {
	supportId := ""
	_, err := tx.QueryOneContext(ctx, &supportId, model.GetNextAssignee)
	if err != nil {
		return "", err
	}
	return supportId, nil
}

func (is *IncidentService) createIncident(ctx context.Context, tx *pg.Tx, r *resequip.MaintenanceIncidentCreate, creatorId string, incidentType resequip.IncidentType) (string, error) {
	log := loggerFromContext(ctx)

	if r.GetDescription() == nil {
		return "", status.Error(codes.InvalidArgument, "description is mandatory")
	}

	if r.GetPriority() == resequip.IncidentPriority_none_priority {
		return "", status.Error(codes.InvalidArgument, "priority is mandatory")
	}

	if r.GetDeadline() == nil {
		return "", status.Error(codes.InvalidArgument, "deadline is mandatory")
	}

	incident := model.Incident{
		ID:          model.GenStringUUID(),
		Description: r.GetDescription().GetValue(),
		CreatedAt:   time.Now(),
		Deadline:    timestampToTime(r.GetDeadline()),
		CreatorID:   creatorId,
		Status:      resequip.IncidentStatus_assigned.String(),
		Comment:     stringWrapperToPtr(r.GetComment()),
		Type:        incidentType.String(),
		Priority:    int(r.GetPriority()),
	}

	_, err := tx.ModelContext(ctx, incident).Insert()
	if err != nil {
		log.WithError(err).Error("unable to insert new incident")
		return "", err
	}

	assigneeId, err := is.getNextAssignee(ctx, tx)
	if err != nil {
		log.WithError(err).Error("unable to get new assignee")
		return "", status.Error(codes.Internal, "unable to get new assignee")
	}

	_, err = tx.ModelContext(ctx, (*model.Incident)(nil)).
		Set(model.Columns.Incident.AssigneeID+" = ?", assigneeId).
		Where(model.Columns.Incident.ID+" = ?", incident.ID).
		Update()
	if err != nil {
		log.WithError(err).Error("unable to set assignee")
		return "", status.Error(codes.Internal, "unable to set assignee")
	}

	return incident.ID, nil
}

func (is *IncidentService) CreateMaintenanceIncident(ctx context.Context, r *resequip.MaintenanceIncidentCreate) (*empty.Empty, error) {
	log := loggerFromContext(ctx)

	user, err := userFromContext(ctx)
	if err != nil {
		return nil, status.Error(codes.PermissionDenied, err.Error())
	}

	tx, err := is.db.Begin()
	if err != nil {
		log.WithError(err).Error("unable to begin transaction")
		return nil, status.Error(codes.Internal, "unable to begin transaction")
	}

	_, err = is.createIncident(ctx, tx, r, user.Id, resequip.IncidentType_maintenance)
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

func (is *IncidentService) CreateEquipmentIncident(ctx context.Context, r *resequip.EquipmentIncidentCreate) (*empty.Empty, error) {
	log := loggerFromContext(ctx)

	user, err := userFromContext(ctx)
	if err != nil {
		return nil, status.Error(codes.PermissionDenied, err.Error())
	}

	tx, err := is.db.Begin()
	if err != nil {
		log.WithError(err).Error("unable to begin transaction")
		return nil, status.Error(codes.Internal, "unable to begin transaction")
	}

	incidentId, err := is.createIncident(ctx, tx, r.GetIncident(), user.Id, resequip.IncidentType_equipment)
	if err != nil {
		log.WithError(err).Error("unable to create incident")
		terr := tx.Rollback()
		if terr != nil {
			log.WithError(terr).Error("unable to rollback transaction")
		}
		return nil, status.Error(codes.Internal, "unable to create incident")
	}

	requiresApproval, err := is.equipmentRequiresApproval(ctx, tx, r.GetEquipmentId())
	if err != nil {
		log.WithError(err).Error("unable to determine equipment approval")
		terr := tx.Rollback()
		if terr != nil {
			log.WithError(terr).Error("unable to rollback transaction")
		}
		return nil, status.Error(codes.Internal, "unable to determine equipment approval")
	}

	equipInc := model.EquipmentIncident{
		ID:           model.GenStringUUID(),
		EquipmentID:  r.GetEquipmentId(),
		IncidentID:   incidentId,
		Deadline:     timestampToTime(r.GetDeadline()),
		NeedApproval: requiresApproval,
	}
	_, err = tx.ModelContext(ctx, equipInc).
		Insert()
	if err != nil {
		log.WithError(err).Error("unable to insert new equipment incident")
		terr := tx.Rollback()
		if terr != nil {
			log.WithError(terr).Error("unable to rollback transaction")
		}
		return nil, status.Error(codes.Internal, "unable to insert new equipment incident")
	}

	err = tx.Commit()
	if err != nil {
		log.WithError(err).Error("unable to commit transaction")
		return nil, status.Error(codes.Internal, "unable to commit transaction")
	}

	return &empty.Empty{}, nil
}

type incidentWithEquipment struct {
	ID          string     `sql:"id,pk,type:uuid"`
	Ordinal     int        `sql:"ordinal"`
	Description string     `sql:"description,notnull"`
	CreatedAt   time.Time  `sql:"created_at,notnull"`
	ResolvedAt  *time.Time `sql:"resolved_at"`
	Deadline    time.Time  `sql:"deadline,notnull"`
	Status      string     `sql:"status,notnull"`
	Comment     *string    `sql:"comment"`
	Type        string     `sql:"type,notnull"`
	Priority    int        `sql:"priority,notnull"`

	EquipmentIncidentID *string    `sql:"equipment_incident_id,type:uuid"`
	EquipmentDeadline   *time.Time `sql:"equipment_deadline"`
	NeedApproval        *bool      `sql:"need_approval"`
	Approved            *bool      `sql:"approved"`

	EquipmentID          *string `sql:"equipment_id,type:uuid"`
	EquipmentName        *string `sql:"equipment_name"`
	EquipmentDescription *string `sql:"equipment_description"`
	EquipmentPrice       *int    `sql:"equipment_price"`

	AssigneeID       *string `sql:"assignee_id,type:uuid"`
	AssigneeLogin    *string `sql:"assignee_id,type:uuid"`
	AssigneeFullName *string `sql:"assignee_full_name,type:uuid"`

	CreatorID       *string `sql:"creator_id,type:uuid"`
	CreatorLogin    *string `sql:"creator_id,type:uuid"`
	CreatorFullName *string `sql:"creator_full_name,type:uuid"`
}

func modelToIncident(inc *incidentWithEquipment) *resequip.Incident {
	res := &resequip.Incident{
		Id:          inc.ID,
		Ordinal:     int64(inc.Ordinal),
		Description: inc.Description,
		Priority:    resequip.IncidentPriority(inc.Priority),
		Deadline:    timeToTimestamp(inc.Deadline),
		Comment:     ptrToStringWrapper(inc.Comment),
		Status:      resequip.IncidentStatus(resequip.IncidentStatus_value[inc.Status]),
		Type:        resequip.IncidentType(resequip.IncidentType_value[inc.Type]),
		CreatedAt:   timeToTimestamp(inc.CreatedAt),
		ResolvedAt:  prtToTimestamp(inc.ResolvedAt),
	}

	if inc.EquipmentIncidentID != nil {
		res.EquipmentIncident = &resequip.EquipmentIncident{}

		if inc.EquipmentDeadline != nil {
			res.EquipmentIncident.Deadline = timeToTimestamp(*inc.EquipmentDeadline)
		}
		if inc.NeedApproval != nil {
			res.EquipmentIncident.RequiresApproval = *inc.NeedApproval
		}
		if inc.Approved != nil {
			res.EquipmentIncident.Approved = *inc.Approved
		}

		if inc.EquipmentID != nil {
			res.EquipmentIncident.Equipment = &resequip.Equipment{}

			if inc.EquipmentName != nil {
				res.EquipmentIncident.Equipment.Name = *inc.EquipmentName
			}
			if inc.EquipmentDescription != nil {
				res.EquipmentIncident.Equipment.Description = *inc.EquipmentDescription
			}
			if inc.EquipmentPrice != nil {
				res.EquipmentIncident.Equipment.Price = int64(*inc.EquipmentPrice)
			}
		}
	}

	if inc.AssigneeID != nil {
		res.Assignee = &resequip.Person{
			Id: *inc.AssigneeID,
		}

		if inc.AssigneeLogin != nil {
			res.Assignee.Login = *inc.AssigneeLogin
		}
		if inc.AssigneeFullName != nil {
			res.Assignee.FullName = *inc.AssigneeFullName
		}
	}

	if inc.CreatorID != nil {
		res.Creator = &resequip.Person{
			Id: *inc.CreatorID,
		}

		if inc.CreatorLogin != nil {
			res.Creator.Login = *inc.CreatorLogin
		}
		if inc.CreatorFullName != nil {
			res.Creator.FullName = *inc.CreatorFullName
		}
	}

	return res
}

func (is *IncidentService) ListIncidents(ctx context.Context, r *resequip.IncidentFilter) (*resequip.IncidentList, error) {
	log := loggerFromContext(ctx)

	incidents := []*incidentWithEquipment{}

	query := is.db.ModelContext(ctx, (*model.Incident)(nil)).
		ColumnExpr("t."+model.Columns.Incident.ID+" as id").
		ColumnExpr("t."+model.Columns.Incident.Ordinal+" as ordinal").
		ColumnExpr("t."+model.Columns.Incident.Description+" as description").
		ColumnExpr("t."+model.Columns.Incident.CreatedAt+" as created_at").
		ColumnExpr("t."+model.Columns.Incident.ResolvedAt+" as resolved_at").
		ColumnExpr("t."+model.Columns.Incident.Deadline+" as deadline").
		ColumnExpr("t."+model.Columns.Incident.Status+" as status").
		ColumnExpr("t."+model.Columns.Incident.Comment+" as comment").
		ColumnExpr("t."+model.Columns.Incident.Type+" as type").
		ColumnExpr("t."+model.Columns.Incident.Priority+" as priority").
		ColumnExpr("ei."+model.Columns.EquipmentIncident.ID+" as equipment_incident_id").
		ColumnExpr("ei."+model.Columns.EquipmentIncident.Deadline+" as equipment_deadline").
		ColumnExpr("ei."+model.Columns.EquipmentIncident.NeedApproval+" as need_approval").
		ColumnExpr("ei."+model.Columns.EquipmentIncident.Approved+" as approved").
		ColumnExpr("e."+model.Columns.Equipment.ID+" as equipment_id").
		ColumnExpr("e."+model.Columns.Equipment.Name+" as equipment_name").
		ColumnExpr("e."+model.Columns.Equipment.Description+" as equipment_description").
		ColumnExpr("e."+model.Columns.Equipment.Price+" as equipment_price").
		ColumnExpr("ass."+model.Columns.Person.ID+" as assignee_id").
		ColumnExpr("ass."+model.Columns.Person.Login+" as assignee_login").
		ColumnExpr("ass."+model.Columns.Person.FullName+" as assignee_full_name").
		ColumnExpr("cr."+model.Columns.Person.ID+" as creator_id").
		ColumnExpr("cr."+model.Columns.Person.Login+" as creator_login").
		ColumnExpr("cr."+model.Columns.Person.FullName+" as creator_full_name").
		Join("left join "+model.Tables.EquipmentIncident.Name+" as ei").
		JoinOn("t."+model.Columns.Incident.ID+" = "+"ei."+model.Columns.EquipmentIncident.IncidentID).
		Join("left join "+model.Tables.Equipment.Name+" as e").
		JoinOn("ei."+model.Columns.EquipmentIncident.EquipmentID+" = "+"e."+model.Columns.Equipment.ID).
		Join("inner join "+model.Tables.Person.Name+" as cr").
		JoinOn("t."+model.Columns.Incident.CreatorID+" = "+"cr."+model.Columns.Person.ID).
		Join("left join "+model.Tables.Support.Name+" as s").
		JoinOn("t."+model.Columns.Incident.AssigneeID+" = "+"s."+model.Columns.Support.ID).
		Join("left join "+model.Tables.Person.Name+" as ass").
		JoinOn("s."+model.Columns.Support.PersonID+" = "+"ass."+model.Columns.Person.ID).
		Order("t."+model.Columns.Incident.Priority+" asc", "t."+model.Columns.Incident.Deadline+" asc")

	if r.GetSearch() != nil {
		query.Where("t."+model.Columns.Incident.Description+" ilike concat('%', ?::text, '%')", r.GetSearch().GetValue())
	}

	if r.GetDeadline() != nil {
		if r.GetDeadline().GetLowerBound() != nil {
			query.Where("t."+model.Columns.Incident.Deadline+" >= ?", timestampToTime(r.GetDeadline().GetLowerBound()))
		}
		if r.GetDeadline().GetUpperBound() != nil {
			query.Where("t."+model.Columns.Incident.Deadline+" <= ?", timestampToTime(r.GetDeadline().GetUpperBound()))
		}
	}

	if r.GetCreatedAt() != nil {
		if r.GetCreatedAt().GetLowerBound() != nil {
			query.Where("t."+model.Columns.Incident.Deadline+" >= ?", timestampToTime(r.GetCreatedAt().GetLowerBound()))
		}
		if r.GetCreatedAt().GetUpperBound() != nil {
			query.Where("t."+model.Columns.Incident.Deadline+" <= ?", timestampToTime(r.GetCreatedAt().GetUpperBound()))
		}
	}

	if r.GetCreatorId() != nil {
		query.Where("t."+model.Columns.Incident.CreatorID+" = ?", r.GetCreatorId().GetValue())
	}

	if r.GetType() != resequip.IncidentType_none_type {
		query.Where("t."+model.Columns.Incident.Type+" = ?", r.GetType().String())
	}

	if r.GetStatus() != resequip.IncidentStatus_none_status {
		query.Where("t."+model.Columns.Incident.Status+" = ?", r.GetStatus().String())
	}

	if r.GetOrdinal() != nil {
		query.Where("t."+model.Columns.Incident.Ordinal+" = ?", r.GetOrdinal().GetValue())
	}

	if r.GetAssigneeId() != nil {
		query.Join("inner join "+model.Tables.Support.Name+" as s").
			JoinOn("t."+model.Columns.Incident.AssigneeID+" = "+"s."+model.Columns.Support.ID).
			Where("s."+model.Columns.Support.PersonID+" = ?", r.GetAssigneeId().GetValue())
	}

	_, err := query.SelectAndCount(&incidents)
	if err != nil {
		log.WithError(err).Error("unable to list incidents")
		return nil, status.Error(codes.Internal, "unable to list incidents")
	}

	res := make([]*resequip.Incident, 0, len(incidents))

	for i := range incidents {
		res = append(res, modelToIncident(incidents[i]))
	}

	return &resequip.IncidentList{Incidents: res}, nil
}

func (is *IncidentService) AssignIncident(ctx context.Context, r *resequip.AssignmentRequest) (*empty.Empty, error) {
	log := loggerFromContext(ctx)

	user, err := userFromContext(ctx)
	if err != nil {
		return nil, status.Error(codes.PermissionDenied, err.Error())
	}

	if user.Role != resequip.Role_team_leader {
		return nil, status.Error(codes.PermissionDenied, "only team leader is allowed")
	}

	_, err = is.db.ModelContext(ctx, (*model.Incident)(nil)).
		Set(model.Columns.Incident.AssigneeID+" = ?",
			is.db.ModelContext(ctx, (*model.Support)(nil)).
				ColumnExpr(model.Columns.Support.ID).
				Where(model.Columns.Support.PersonID+" = ?", r.GetPersonId()),
		).
		Set(model.Columns.Incident.Status+" = ?", resequip.IncidentStatus_assigned.String()).
		Where(model.Columns.Incident.ID+" = ?", r.GetIncidentId()).
		Update()
	if err != nil {
		log.WithError(err).Error("unable to set assignee")
		return nil, status.Error(codes.Internal, "unable to set assignee")
	}

	return &empty.Empty{}, nil
}

func (is *IncidentService) incidentCouldBeResolved(ctx context.Context, tx *pg.Tx, incident *model.Incident) (bool, error) {
	if incident.Type != resequip.IncidentType_equipment.String() {
		return true, nil
	}

	ei := &model.EquipmentIncident{}

	requiresApproval := false
	err := tx.ModelContext(ctx, ei).
		Where(model.Columns.EquipmentIncident.IncidentID+" = ?", incident.ID).
		Select(&requiresApproval)
	if err != nil {
		return false, err
	}

	return !ei.NeedApproval || ei.NeedApproval && ei.Approved, nil
}

func (is *IncidentService) ChangeIncidentStatus(ctx context.Context, r *resequip.IncidentStatusRequest) (*empty.Empty, error) {
	log := loggerFromContext(ctx)

	user, err := userFromContext(ctx)
	if err != nil {
		return nil, status.Error(codes.PermissionDenied, err.Error())
	}

	if user.Role != resequip.Role_team_leader && user.Role != resequip.Role_support {
		return nil, status.Error(codes.PermissionDenied, "only support team allowed")
	}

	tx, err := is.db.Begin()
	if err != nil {
		log.WithError(err).Error("unable to begin transaction")
		return nil, status.Error(codes.Internal, "unable to begin transaction")
	}

	incident := &model.Incident{}
	err = tx.ModelContext(ctx, incident).
		Where(model.Columns.Incident.ID+" = ?", r.GetIncidentId()).
		For("update").
		Select()
	if err != nil {
		log.WithError(err).Error("unable to select incident")
		terr := tx.Rollback()
		if terr != nil {
			log.WithError(terr).Error("unable to rollback transaction")
		}
		return nil, status.Error(codes.Internal, "unable to select incident")
	}

	switch resequip.IncidentStatus(resequip.IncidentStatus_value[incident.Status]) {
	case resequip.IncidentStatus_created:
		{
			if r.GetStatus() != resequip.IncidentStatus_assigned {
				terr := tx.Rollback()
				if terr != nil {
					log.WithError(terr).Error("unable to rollback transaction")
				}
				return nil, status.Error(codes.PermissionDenied, "not allowed")
			}
			if incident.AssigneeID == nil {
				terr := tx.Rollback()
				if terr != nil {
					log.WithError(terr).Error("unable to rollback transaction")
				}
				return nil, status.Error(codes.PermissionDenied, "no assignee is set")
			}
		}
	case resequip.IncidentStatus_assigned:
		{
			switch r.GetStatus() {
			case resequip.IncidentStatus_resolved:
				ok, err := is.incidentCouldBeResolved(ctx, tx, incident)
				if err != nil {
					log.WithError(err).Error("unable to determine if the incident could be resolved")
					terr := tx.Rollback()
					if terr != nil {
						log.WithError(terr).Error("unable to rollback transaction")
					}
					return nil, status.Error(codes.Internal, "unable to determine if the incident could be resolved")
				}
				if !ok {
					terr := tx.Rollback()
					if terr != nil {
						log.WithError(terr).Error("unable to rollback transaction")
					}
					return nil, status.Error(codes.PermissionDenied, "incident must be approved first")
				}
				err = is.createEquipmentAssignment(ctx, tx, incident.ID, incident.CreatorID)
				if err != nil {
					log.WithError(err).Error("unable to assign equipment")
					terr := tx.Rollback()
					if terr != nil {
						log.WithError(terr).Error("unable to rollback transaction")
					}
					return nil, status.Error(codes.Internal, "unable to assign equipment")
				}
			case resequip.IncidentStatus_dismissed:
				break
			default:
				terr := tx.Rollback()
				if terr != nil {
					log.WithError(terr).Error("unable to rollback transaction")
				}
				return nil, status.Error(codes.PermissionDenied, "not allowed")
			}
		}
	default:
		terr := tx.Rollback()
		if terr != nil {
			log.WithError(terr).Error("unable to rollback transaction")
		}
		return nil, status.Error(codes.PermissionDenied, "not allowed")
	}

	_, err = tx.ModelContext(ctx, (*model.Incident)(nil)).
		Where(model.Columns.Incident.ID+" = ?", r.GetIncidentId()).
		Set(model.Columns.Incident.Status+" = ?", r.GetStatus().String()).
		Update()
	if err != nil {
		log.WithError(err).Error("unable to set new status")
		terr := tx.Rollback()
		if terr != nil {
			log.WithError(terr).Error("unable to rollback transaction")
		}
		return nil, status.Error(codes.Internal, "unable to set new status")
	}

	err = tx.Commit()
	if err != nil {
		log.WithError(err).Error("unable to commit transaction")
		return nil, status.Error(codes.Internal, "unable to commit transaction")
	}

	return &empty.Empty{}, nil
}

func (is *IncidentService) CommentOnIncident(ctx context.Context, r *resequip.IncidentCommentRequest) (*empty.Empty, error) {
	log := loggerFromContext(ctx)

	user, err := userFromContext(ctx)
	if err != nil {
		return nil, status.Error(codes.PermissionDenied, err.Error())
	}

	if user.Role != resequip.Role_team_leader && user.Role != resequip.Role_support {
		return nil, status.Error(codes.PermissionDenied, "only support team allowed")
	}

	_, err = is.db.ModelContext(ctx, (*model.Incident)(nil)).
		Where(model.Columns.Incident.ID+" = ?", r.GetIncidentId()).
		Set(model.Columns.Incident.Comment+" = ?", r.GetComment().GetValue()).
		Update()
	if err != nil {
		log.WithError(err).Error("unable to set new comment")
		return nil, status.Error(codes.Internal, "unable to set new comment")
	}

	return &empty.Empty{}, nil
}

func (is *IncidentService) ApproveEquipmentIncident(ctx context.Context, r *resequip.IncidentApprovalRequest) (*empty.Empty, error) {
	log := loggerFromContext(ctx)

	user, err := userFromContext(ctx)
	if err != nil {
		return nil, status.Error(codes.PermissionDenied, err.Error())
	}

	incident := &model.Incident{}
	err = is.db.ModelContext(ctx, incident).
		Where(model.Columns.Incident.ID+" = ?", r.GetIncidentId()).
		Select()
	if err != nil {
		log.WithError(err).Error("unable to select incident")
		return nil, status.Error(codes.Internal, "unable to select incident")
	}

	if incident.Type != resequip.IncidentType_equipment.String() {
		return nil, status.Error(codes.InvalidArgument, "invalid incident type")
	}

	isManager := false
	err = is.db.ModelContext(ctx, (*model.Person)(nil)).
		ColumnExpr(model.Columns.Person.ManagerID+" = ? as is_manager", user.GetId()).
		Where(model.Columns.Person.ID+" = ?", incident.CreatorID).
		Select(&isManager)
	if err != nil {
		log.WithError(err).Error("unable to determine persons manager")
		return nil, status.Error(codes.Internal, "unable to determine persons manager")
	}

	if !isManager {
		return nil, status.Error(codes.PermissionDenied, "not allowed")
	}

	_, err = is.db.ModelContext(ctx, (*model.EquipmentIncident)(nil)).
		Set(model.Columns.EquipmentIncident.Approved+" = true").
		Where(model.Columns.EquipmentIncident.IncidentID+" = ?", r.GetIncidentId()).
		Update()
	if err != nil {
		log.WithError(err).Error("unable to approve incident")
		return nil, status.Error(codes.Internal, "unable to approve incident")
	}

	return &empty.Empty{}, nil
}
