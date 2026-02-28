// Package message implements the business logic for sending and retrieving
// direct messages between users, including image attachments.
package message

import (
	"github.com/CackSocial/cack-backend/internal/domain"
	"github.com/CackSocial/cack-backend/internal/dto"
	"github.com/CackSocial/cack-backend/internal/infrastructure/storage"
	"github.com/CackSocial/cack-backend/internal/repository"
	ucerrors "github.com/CackSocial/cack-backend/internal/usecase/errors"
)

// MessageUseCase encapsulates all message-related business logic including
// sending messages, retrieving conversations, and listing conversations.
type MessageUseCase struct {
	messageRepo repository.MessageRepository
	userRepo    repository.UserRepository
	storage     storage.Storage
}

// NewMessageUseCase creates a new MessageUseCase with the given dependencies.
func NewMessageUseCase(messageRepo repository.MessageRepository, userRepo repository.UserRepository, storage storage.Storage) *MessageUseCase {
	return &MessageUseCase{
		messageRepo: messageRepo,
		userRepo:    userRepo,
		storage:     storage,
	}
}

// Send delivers a message from the authenticated user to the receiver
// identified by username. If the request includes an image, it is uploaded
// via the storage backend.
func (uc *MessageUseCase) Send(senderID string, receiverUsername string, req *dto.SendMessageRequest) (*dto.MessageResponse, error) {
	receiver, err := uc.userRepo.GetByUsername(receiverUsername)
	if err != nil || receiver == nil {
		return nil, ucerrors.ErrUserNotFound
	}

	var imageURL string
	if req.Image != nil {
		url, err := uc.storage.Upload(req.Image)
		if err != nil {
			return nil, err
		}
		imageURL = url
	}

	msg := &domain.Message{
		SenderID:   senderID,
		ReceiverID: receiver.ID,
		Content:    req.Content,
		ImageURL:   imageURL,
	}

	if err := uc.messageRepo.Create(msg); err != nil {
		return nil, err
	}

	return toMessageResponse(msg), nil
}

// GetConversation returns a paginated conversation between the authenticated
// user and the partner identified by username. Messages sent to the current
// user are marked as read.
func (uc *MessageUseCase) GetConversation(userID string, partnerUsername string, page, limit int) ([]dto.MessageResponse, int64, error) {
	partner, err := uc.userRepo.GetByUsername(partnerUsername)
	if err != nil || partner == nil {
		return nil, 0, ucerrors.ErrUserNotFound
	}

	messages, total, err := uc.messageRepo.GetConversation(userID, partner.ID, page, limit)
	if err != nil {
		return nil, 0, err
	}

	// Mark messages addressed to the current user as read.
	_ = uc.messageRepo.MarkAsRead(userID, partner.ID)

	responses := make([]dto.MessageResponse, 0, len(messages))
	for i := range messages {
		responses = append(responses, *toMessageResponse(&messages[i]))
	}

	return responses, total, nil
}

// GetConversations returns a paginated list of the authenticated user's
// conversations, each showing the partner, last message, and unread count.
func (uc *MessageUseCase) GetConversations(userID string, page, limit int) ([]dto.ConversationListResponse, int64, error) {
	previews, total, err := uc.messageRepo.GetConversations(userID, page, limit)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]dto.ConversationListResponse, 0, len(previews))
	for _, p := range previews {
		responses = append(responses, dto.ConversationListResponse{
			User: dto.UserProfile{
				ID:          p.User.ID,
				Username:    p.User.Username,
				DisplayName: p.User.DisplayName,
				Bio:         p.User.Bio,
				AvatarURL:   p.User.AvatarURL,
			},
			LastMessage: *toMessageResponse(&p.LastMessage),
			UnreadCount: p.UnreadCount,
		})
	}

	return responses, total, nil
}

// SendFromWS creates a message directly from WebSocket data where the
// receiver ID and image URL are already resolved. Returns the persisted
// domain message.
func (uc *MessageUseCase) SendFromWS(senderID, receiverID, content, imageURL string) (*domain.Message, error) {
	msg := &domain.Message{
		SenderID:   senderID,
		ReceiverID: receiverID,
		Content:    content,
		ImageURL:   imageURL,
	}

	if err := uc.messageRepo.Create(msg); err != nil {
		return nil, err
	}

	return msg, nil
}

// toMessageResponse converts a domain Message into a MessageResponse DTO.
func toMessageResponse(m *domain.Message) *dto.MessageResponse {
	return &dto.MessageResponse{
		ID:         m.ID,
		SenderID:   m.SenderID,
		ReceiverID: m.ReceiverID,
		Content:    m.Content,
		ImageURL:   m.ImageURL,
		ReadAt:     m.ReadAt,
		CreatedAt:  m.CreatedAt,
	}
}
