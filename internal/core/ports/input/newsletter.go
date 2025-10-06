package input

import (
	"context"
	"zpwoot/internal/core/application/dto"
)


type NewsletterService interface {

	ListNewsletters(ctx context.Context, sessionID string) (*dto.ListNewslettersResponse, error)
	GetNewsletterInfo(ctx context.Context, sessionID string, newsletterJID string) (*dto.NewsletterInfo, error)
	GetNewsletterInfoWithInvite(ctx context.Context, sessionID string, req *dto.NewsletterInfoWithInviteRequest) (*dto.NewsletterInfo, error)


	CreateNewsletter(ctx context.Context, sessionID string, req *dto.CreateNewsletterRequest) (*dto.NewsletterInfo, error)


	FollowNewsletter(ctx context.Context, sessionID string, req *dto.FollowNewsletterRequest) error
	UnfollowNewsletter(ctx context.Context, sessionID string, newsletterJID string) error


	GetMessages(ctx context.Context, sessionID string, newsletterJID string, req *dto.GetNewsletterMessagesRequest) (*dto.ListNewsletterMessagesResponse, error)
	MarkViewed(ctx context.Context, sessionID string, newsletterJID string, req *dto.NewsletterMarkViewedRequest) error


	SendReaction(ctx context.Context, sessionID string, newsletterJID string, req *dto.NewsletterReactionRequest) error
	ToggleMute(ctx context.Context, sessionID string, newsletterJID string, req *dto.NewsletterMuteRequest) error
}
