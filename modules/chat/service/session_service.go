package service

import (
	"context"
	db "shofy/db/sqlc"
	"shofy/modules/chat/model"
)

func (s *ChatService) GetOrCreateSession(ctx context.Context, chatSession model.ChatSession) (db.Session, error) {
	session, err := s.Queries.GetCurrentSessions(ctx, db.GetCurrentSessionsParams{
		ChannelID: int32(chatSession.ChannelID),
		UserID:    int32(chatSession.UserID),
	})
	if err == nil {
		return session, nil
	}

	newSession, err := s.Queries.CreateSession(ctx, db.CreateSessionParams{
		ChannelID: int32(chatSession.ChannelID),
		UserID:    int32(chatSession.UserID),
	})

	if err != nil {
		return db.Session{}, err
	}

	return newSession, nil
}
