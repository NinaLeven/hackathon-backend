package main


import (
	"context"
	"github.com/WantsToFress/hackathon-backend/internal/model"
	resequip "github.com/WantsToFress/hackathon-backend/pkg"
	"github.com/mailgun/mailgun-go/v4"
	"github.com/mailgun/mailgun-go/v4/events"
	"strings"
	"time"
)

type MailConfig struct {
	Address    string `yaml:"address"`
	Password   string `yaml:"password"`
	Domain     string `yaml:"domain"`
	PrivateKey string `yaml:"private_key"`
	BaseAPI    string `yaml:"base_api"`
}

func NewMailClient(ctx context.Context, config *MailConfig) mailgun.Mailgun {
	client := mailgun.NewMailgun(config.Domain, config.PrivateKey)
	client.SetAPIBase(mailgun.APIBaseUS)
	return client
}

func (is *IncidentService) watchMail(ctx context.Context) {
	log := loggerFromContext(ctx)

	eventList := is.mailClient.PollEvents(&mailgun.ListEventOptions{
		Begin:           time.Now().Add(time.Duration(-1) * time.Hour * 24),
		PollInterval:    time.Second * 5,
	})

	err := eventList.Err()
	if err != nil {
		log.WithError(err).Error("unable to list eventList")
		return
	}

	page := make([]mailgun.Event, 1000, 1000)
	for eventList.Poll(ctx, &page) {
		for _, e := range page {
			if es, ok := e.(*events.Stored); ok {
				msg, err := is.mailClient.GetStoredMessage(ctx, es.Storage.URL)
				if err != nil {
					log.WithError(err).Error("unable to get message")
					continue
				}
				log.WithField("msg", msg.StrippedText).Info()
				is.processMail(ctx, msg.Sender, msg.StrippedText)
			}
		}
	}
	err = eventList.Err()
	if err != nil {
		log.WithError(err).Error("unable to list eventList")
		return
	}

	log.Info("stop watching mail")
}

func (is *IncidentService) processMail(ctx context.Context, sender string, msg string) {
	log := loggerFromContext(ctx)

	person, err := is.getPersonByEmail(ctx, sender)
	if err != nil {
		log.WithError(err).Error("unable to select person")
		return
	}

	equipmentId := ""
	description := ""
	deadline := time.Time{}

	lines := strings.Split(msg, "\n")
	for i := range lines {
		lines[i] = strings.Trim(lines[i], "\r")
		switch {
		case strings.HasPrefix(lines[i], "название оборудования:"):

			equipmentName := strings.Trim(strings.TrimPrefix(lines[i], "название оборудования:"), " ")
			err = is.db.ModelContext(ctx, (*model.Equipment)(nil)).
				ColumnExpr(model.Columns.Equipment.ID).
				Where(model.Columns.Equipment.Name+" ilike ?", equipmentName).
				Select(&equipmentId)
			if err != nil {
				log.WithError(err).Warning()
				continue
			}

		case strings.HasPrefix(lines[i], "выдать до:"):
			deadline, err = time.Parse(time.RFC3339, strings.Trim(strings.TrimPrefix(lines[i], "выдать до:"), " "))
			if err != nil {
				log.WithError(err).Warning()
				continue
			}

		case strings.HasPrefix(lines[i], "описание:"):
			description = strings.Trim(strings.TrimPrefix(lines[i], "описание:"), " ")
		}
	}

	if equipmentId == "" || deadline.Unix() == 0 || description == "" {
		return
	}

	_, err = is.CreateEquipmentIncident(userToContext(ctx, person), &resequip.EquipmentIncidentCreate{
		Incident: &resequip.MaintenanceIncidentCreate{
			Description: stringToStringWrapper(description),
			Deadline:    timeToTimestamp(time.Now().Add(time.Hour * 24 * 7)),
			Priority:    1,
		},
		EquipmentId: equipmentId,
		Deadline:    timeToTimestamp(deadline),
	})
	if err != nil {
		log.WithError(err).Error()
	}
}