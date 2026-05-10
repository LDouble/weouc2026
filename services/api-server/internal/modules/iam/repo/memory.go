package repo

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"

	iamtypes "github.com/liangluo/weouc2026/services/api-server/internal/modules/iam/types"
)

var ErrUserNotFound = errors.New("user not found")
var ErrSessionNotFound = errors.New("session not found")
var ErrCaptchaNotFound = errors.New("captcha not found")

type UserRepository interface {
	FindByID(ctx context.Context, userID string) (iamtypes.User, error)
	FindByOpenID(ctx context.Context, openID string) (iamtypes.User, error)
	Save(ctx context.Context, user iamtypes.User) (iamtypes.User, error)
	Update(ctx context.Context, userID string, mutate func(*iamtypes.User) error) (iamtypes.User, error)
	NextID() string
}

type SessionRepository interface {
	Save(ctx context.Context, session iamtypes.Session) (iamtypes.Session, error)
	Find(ctx context.Context, token string) (iamtypes.Session, error)
	Delete(ctx context.Context, token string) error
}

type CaptchaRepository interface {
	Save(ctx context.Context, ticket iamtypes.CaptchaTicket) error
	Find(ctx context.Context, studentID string) (iamtypes.CaptchaTicket, error)
	Delete(ctx context.Context, studentID string) error
}

type InMemoryUserRepository struct {
	mu       sync.RWMutex
	byID     map[string]iamtypes.User
	byOpenID map[string]string
	seq      uint64
}

func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{
		byID:     make(map[string]iamtypes.User),
		byOpenID: make(map[string]string),
	}
}

func (r *InMemoryUserRepository) findByID(userID string) (iamtypes.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, exists := r.byID[userID]
	if !exists {
		return iamtypes.User{}, ErrUserNotFound
	}

	return cloneUser(user), nil
}

func (r *InMemoryUserRepository) FindByOpenID(_ context.Context, openID string) (iamtypes.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	userID, exists := r.byOpenID[openID]
	if !exists {
		return iamtypes.User{}, ErrUserNotFound
	}

	user, userExists := r.byID[userID]
	if !userExists {
		return iamtypes.User{}, ErrUserNotFound
	}

	return cloneUser(user), nil
}

func (r *InMemoryUserRepository) Save(_ context.Context, user iamtypes.User) (iamtypes.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	user = cloneUser(user)
	r.byID[user.ID] = user
	if user.OpenID != "" {
		r.byOpenID[user.OpenID] = user.ID
	}

	return cloneUser(user), nil
}

func (r *InMemoryUserRepository) Update(_ context.Context, userID string, mutate func(*iamtypes.User) error) (iamtypes.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	user, exists := r.byID[userID]
	if !exists {
		return iamtypes.User{}, ErrUserNotFound
	}

	updated := cloneUser(user)
	if err := mutate(&updated); err != nil {
		return iamtypes.User{}, err
	}

	r.byID[userID] = updated
	if updated.OpenID != "" {
		r.byOpenID[updated.OpenID] = updated.ID
	}

	return cloneUser(updated), nil
}

func (r *InMemoryUserRepository) NextID() string {
	value := atomic.AddUint64(&r.seq, 1)
	return fmt.Sprintf("user-%06d", value)
}

func (r *InMemoryUserRepository) FindByID(ctx context.Context, userID string) (iamtypes.User, error) {
	return r.findByID(userID)
}

type InMemorySessionRepository struct {
	mu       sync.RWMutex
	sessions map[string]iamtypes.Session
}

func NewInMemorySessionRepository() *InMemorySessionRepository {
	return &InMemorySessionRepository{
		sessions: make(map[string]iamtypes.Session),
	}
}

func (r *InMemorySessionRepository) Save(_ context.Context, session iamtypes.Session) (iamtypes.Session, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.sessions[session.Token] = session
	return session, nil
}

func (r *InMemorySessionRepository) Find(_ context.Context, token string) (iamtypes.Session, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	session, exists := r.sessions[token]
	if !exists {
		return iamtypes.Session{}, ErrSessionNotFound
	}

	return session, nil
}

func (r *InMemorySessionRepository) Delete(_ context.Context, token string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.sessions, token)
	return nil
}

type InMemoryCaptchaRepository struct {
	mu      sync.RWMutex
	tickets map[string]iamtypes.CaptchaTicket
}

func NewInMemoryCaptchaRepository() *InMemoryCaptchaRepository {
	return &InMemoryCaptchaRepository{
		tickets: make(map[string]iamtypes.CaptchaTicket),
	}
}

func (r *InMemoryCaptchaRepository) Save(_ context.Context, ticket iamtypes.CaptchaTicket) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.tickets[ticket.StudentID] = ticket
	return nil
}

func (r *InMemoryCaptchaRepository) Find(_ context.Context, studentID string) (iamtypes.CaptchaTicket, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	ticket, exists := r.tickets[studentID]
	if !exists {
		return iamtypes.CaptchaTicket{}, ErrCaptchaNotFound
	}

	return ticket, nil
}

func (r *InMemoryCaptchaRepository) Delete(_ context.Context, studentID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.tickets, studentID)
	return nil
}

func cloneUser(user iamtypes.User) iamtypes.User {
	cloned := user
	cloned.Roles = append([]string(nil), user.Roles...)
	cloned.Permissions = append([]string(nil), user.Permissions...)
	if user.StudentProfile != nil {
		profile := *user.StudentProfile
		cloned.StudentProfile = &profile
	}

	return cloned
}
