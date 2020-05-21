package main

import (
	"context"
	resequip "github.com/WantsToFress/hackathon-backend/pkg"
	"github.com/centrifugal/centrifuge-go"
)

func (es *IncidentService) GetChatToken(ctx context.Context, r *resequip.Id) (*resequip.ChatToken, error) {
	//log := loggerFromContext(ctx)
	//
	//if !model.IsValidUUID(r.GetId()) {
	//	return nil, status.Error(codes.InvalidArgument, "invalid id")
	//}
	//
	//user, err := userFromContext(ctx)
	//if err != nil {
	//	log.WithError(err).Error("unable to get user from context")
	//	return nil, err
	//}
	//
	//token := jwt.New(jwt.SigningMethodHS256)
	//token.Claims = jwt.MapClaims{
	//	"sub": user.Login,
	//	"exp": time.Now().Add(time.Hour * 24),
	//}
	//tokenRaw, err := token.SignedString(es.hmacSecret)
	//if err != nil {
	//	return nil, err
	//}
	//
	//return &resequip.ChatToken{Token: tokenRaw}, nil
	return &resequip.ChatToken{}, nil
}

type Message struct {
	Id       string `json:"id"`
	FullName string `json:"full_name"`
	UID      string `json:"uid"`
	Login    string `json:"login"`
	EventId  string `json:"event_id"`
	Time     int64  `json:"time"`
	Message  string `json:"message"`
}

//
//func modelToMessage(m *model.Message) *resequip.Message {
//	res := &resequip.Message{}
//
//	res.Id = m.ID
//	res.Message = m.Message
//	res.Time = m.Time.Unix()
//	res.Login = m.Login
//	res.EventId = m.EventID
//	res.FullName = m.FullName
//	res.Uid = m.PersonID
//
//	return res
//}

func (es *IncidentService) GetChatHistory(ctx context.Context, r *resequip.Id) (*resequip.ChatHistory, error) {
	//log := loggerFromContext(ctx)
	//
	//messages := []*model.Message{}
	//err := es.db.ModelContext(ctx, &messages).
	//	Where(model.Columns.Message.EventID+" = ?", r.GetId()).
	//	Order(model.Columns.Message.Time + " ASC").
	//	Select()
	//if err != nil {
	//	log.WithError(err).Error("unable to select chat")
	//	return nil, status.Error(codes.Internal, err.Error())
	//}
	//
	//res := &resequip.ChatHistory{
	//	Messages: make([]*resequip.Message, 0, len(messages)),
	//}
	//
	//for i := range messages {
	//	res.Messages = append(res.Messages, modelToMessage(messages[i]))
	//}
	//
	//return res, nil
	return &resequip.ChatHistory{}, nil
}

func (es *IncidentService) OnPublish(sub *centrifuge.Subscription, e centrifuge.PublishEvent) {
	//data, err := e.Data.MarshalJSON()
	//if err != nil {
	//	log.Error(err)
	//	return
	//}
	//msg := &Message{}
	//err = json.Unmarshal(data, msg)
	//if err != nil {
	//	log.Error(msg)
	//	return
	//}
	//
	//message := &model.Message{
	//	ID:       msg.Id,
	//	PersonID: msg.UID,
	//	EventID:  msg.EventId,
	//	Login:    msg.Login,
	//	FullName: msg.FullName,
	//	Time:     time.Unix(msg.Time/1000, 0),
	//	Message:  msg.Message,
	//}
	//
	//_, err = es.db.Model(message).
	//	OnConflict("do nothing").
	//	Insert()
	//if err != nil {
	//	log.Error(err)
	//}
}

func (es *IncidentService) WatchChat(ctx context.Context) error {
	sub, err := es.cent.NewSubscription("all")
	if err != nil {
		return err
	}

	sub.OnPublish(es)

	err = sub.Subscribe()
	if err != nil {
		return nil
	}

	return nil
}
