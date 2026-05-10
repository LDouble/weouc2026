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

type UserRepository interface {
	FindByID(ctx context.Context, userID string) (iamtypes.User, bool)
	FindByOpenID(ctx context.Context, openID string) (iamtypes.User, bool)
	Save(ctx context.Context, user iamtypes.User) iamtypes.User
	Update(ctx context.Context, userID string, mutate func(*iamtypes.User) error) (iamtypes.User, error)
	NextID() string
}

type SessionRepository interface {
	Save(ctx context.Context, session iamtypes.Session) iamtypes.Session
	Find(ctx context.Context, token string) (iamtypes.Session, bool)
	Delete(ctx context.Context, token string)
}

type CaptchaRepository interface {
	Save(ctx context.Context, ticket iamtypes.CaptchaTicket)
	Find(ctx context.Context, studentID string) (iamtypes.CaptchaTicket, bool)
	Delete(ctx context.Context, studentID string)
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

func (r *InMemoryUserRepository) findByID(userID string) (iamtypes.User, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, exists := r.byID[userID]
	return cloneUser(user), exists
}

func (r *InMemoryUserRepository) FindByOpenID(_ context.Context, openID string) (iamtypes.User, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	userID, exists := r.byOpenID[openID]
	if !exists {
		return iamtypes.User{}, false
	}

	user, userExists := r.byID[userID]
	return cloneUser(user), userExists
}

func (r *InMemoryUserRepository) Save(_ context.Context, user iamtypes.User) iamtypes.User {
	r.mu.Lock()
	defer r.mu.Unlock()

	user = cloneUser(user)
	r.byID[user.ID] = user
	if user.OpenID != "" {
		r.byOpenID[user.OpenID] = user.ID
	}

	return cloneUser(user)
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

func (r *InMemoryUserRepository) FindByID(ctx context.Context, userID string) (iamtypes.User, bool) {
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

func (r *InMemorySessionRepository) Save(_ context.Context, session iamtypes.Session) iamtypes.Session {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.sessions[session.Token] = session
	return session
}

func (r *InMemorySessionRepository) Find(_ context.Context, token string) (iamtypes.Session, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	session, exists := r.sessions[token]
	return session, exists
}

func (r *InMemorySessionRepository) Delete(_ context.Context, token string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.sessions, token)
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

func (r *InMemoryCaptchaRepository) Save(_ context.Context, ticket iamtypes.CaptchaTicket) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.tickets[ticket.StudentID] = ticket
}

func (r *InMemoryCaptchaRepository) Find(_ context.Context, studentID string) (iamtypes.CaptchaTicket, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	ticket, exists := r.tickets[studentID]
	return ticket, exists
}

func (r *InMemoryCaptchaRepository) Delete(_ context.Context, studentID string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.tickets, studentID)
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
