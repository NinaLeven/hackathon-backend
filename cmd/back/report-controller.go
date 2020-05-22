package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"

	resequip "github.com/WantsToFress/hackathon-backend/pkg"
)

func (is *IncidentService) GetReport(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := loggerFromContext(ctx)

	incidents, err := is.ListIncidents(ctx, &resequip.IncidentFilter{
		CreatedAt: &resequip.TimestampSelector{
			LowerBound: timeToTimestamp(time.Now().Add(time.Duration(-1)*time.Hour*24*30)),
		},
	})
	if err != nil {
		log.WithError(err).Error("unable to get incidents")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	f := excelize.NewFile()

	f.SetSheetName("Sheet1", "main")

	f.SetCellValue("main", "A1", "incident_id")
	f.SetCellValue("main", "B1", "ordinal")
	f.SetCellValue("main", "C1", "type")
	f.SetCellValue("main", "D1", "status")
	f.SetCellValue("main", "E1", "created_at")
	f.SetCellValue("main", "F1", "resolved_at")
	f.SetCellValue("main", "G1", "deadline")
	f.SetCellValue("main", "H1", "creator_login")
	f.SetCellValue("main", "I1", "assignee_login")
	f.SetCellValue("main", "J1", "priority")
	f.SetCellValue("main", "K1", "equipment_name")
	f.SetCellValue("main", "L1", "equipment_requires_approval")
	f.SetCellValue("main", "M1", "equipment_approved")

	for i := range incidents.GetIncidents() {
		f.SetCellValue("main", fmt.Sprintf("A%d", i+2), incidents.GetIncidents()[i].GetId())
		f.SetCellValue("main", fmt.Sprintf("B%d", i+2), incidents.GetIncidents()[i].GetOrdinal())
		f.SetCellValue("main", fmt.Sprintf("C%d", i+2), incidents.GetIncidents()[i].GetType().String())
		f.SetCellValue("main", fmt.Sprintf("D%d", i+2), incidents.GetIncidents()[i].GetStatus().String())
		f.SetCellValue("main", fmt.Sprintf("E%d", i+2), timestampToTime(incidents.GetIncidents()[i].GetCreatedAt()))
		f.SetCellValue("main", fmt.Sprintf("F%d", i+2), timestampToTime(incidents.GetIncidents()[i].GetResolvedAt()))
		f.SetCellValue("main", fmt.Sprintf("G%d", i+2), timestampToTime(incidents.GetIncidents()[i].GetDeadline()))
		f.SetCellValue("main", fmt.Sprintf("H%d", i+2), incidents.GetIncidents()[i].GetCreator().GetLogin())
		f.SetCellValue("main", fmt.Sprintf("I%d", i+2), incidents.GetIncidents()[i].GetAssignee().GetLogin())
		f.SetCellValue("main", fmt.Sprintf("J%d", i+2), incidents.GetIncidents()[i].GetPriority().String())
		f.SetCellValue("main", fmt.Sprintf("K%d", i+2), incidents.GetIncidents()[i].GetEquipmentIncident().GetEquipment().GetName())
		f.SetCellValue("main", fmt.Sprintf("L%d", i+2), incidents.GetIncidents()[i].GetEquipmentIncident().GetRequiresApproval())
		f.SetCellValue("main", fmt.Sprintf("M%d", i+2), incidents.GetIncidents()[i].GetEquipmentIncident().GetApproved())
	}

	err = f.Write(w)
	if err != nil {
		log.WithError(err).Error("unable to write xlsx")
		return
	}
}
