package mockdata

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/MrKyomoto/ieat/internal/auth"
	"github.com/MrKyomoto/ieat/internal/catalog"
	"golang.org/x/crypto/bcrypt"
)

type mockUser struct {
	user         auth.User
	passwordHash string
}

type mockSession struct {
	userID    string
	expiresAt time.Time
}

type Store struct {
	mu           sync.Mutex
	usersByEmail map[string]mockUser
	usersByID    map[string]mockUser
	sessions     map[string]mockSession
	canteens     []catalog.Canteen
}

func New(password string) (*Store, error) {
	if password == "" {
		return nil, fmt.Errorf("DEV_SEED_PASSWORD is required for mock data")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hash mock password: %w", err)
	}
	users := []auth.User{
		{ID: "00000000-0000-0000-0000-000000000001", Email: "student@mail.ustc.edu.cn", Nickname: "开发同学", Role: "member"},
		{ID: "00000000-0000-0000-0000-000000000002", Email: "manager@ustc.edu.cn", Nickname: "开发管理人员", Role: "manager"},
		{ID: "00000000-0000-0000-0000-000000000003", Email: "admin@ustc.edu.cn", Nickname: "开发平台管理员", Role: "admin"},
	}
	store := &Store{
		usersByEmail: make(map[string]mockUser, len(users)),
		usersByID:    make(map[string]mockUser, len(users)),
		sessions:     make(map[string]mockSession),
		canteens:     sampleCanteens(),
	}
	for _, user := range users {
		item := mockUser{user: user, passwordHash: string(hash)}
		store.usersByEmail[user.Email] = item
		store.usersByID[user.ID] = item
	}
	return store, nil
}

func (s *Store) FindUserByEmail(_ context.Context, email string) (auth.User, string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	item, ok := s.usersByEmail[email]
	if !ok {
		return auth.User{}, "", auth.ErrNotFound
	}
	return item.user, item.passwordHash, nil
}

func (s *Store) CreateSession(_ context.Context, tokenHash []byte, userID string, expiresAt time.Time) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.usersByID[userID]; !ok {
		return auth.ErrNotFound
	}
	s.sessions[string(tokenHash)] = mockSession{userID: userID, expiresAt: expiresAt}
	return nil
}

func (s *Store) FindUserBySession(_ context.Context, tokenHash []byte) (auth.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	key := string(tokenHash)
	session, ok := s.sessions[key]
	if !ok || !session.expiresAt.After(time.Now()) {
		delete(s.sessions, key)
		return auth.User{}, auth.ErrNotFound
	}
	item, ok := s.usersByID[session.userID]
	if !ok {
		return auth.User{}, auth.ErrNotFound
	}
	return item.user, nil
}

func (s *Store) DeleteSession(_ context.Context, tokenHash []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.sessions, string(tokenHash))
	return nil
}

func (s *Store) List(_ context.Context) ([]catalog.Canteen, error) {
	return s.canteens, nil
}

func sampleCanteens() []catalog.Canteen {
	return []catalog.Canteen{{
		ID:   "10000000-0000-0000-0000-000000000001",
		Name: "示例食堂",
		Floors: []catalog.Floor{
			{
				ID:   "20000000-0000-0000-0000-000000000001",
				Name: "一层",
				Windows: []catalog.Window{{
					ID:            "30000000-0000-0000-0000-000000000001",
					ExternalCode:  "DEMO-01",
					Name:          "示例窗口一",
					Description:   "用于开发食堂目录和评价功能",
					BusinessHours: "10:30-13:30 / 16:30-19:30",
				}},
			},
			{
				ID:   "20000000-0000-0000-0000-000000000002",
				Name: "二层",
				Windows: []catalog.Window{{
					ID:            "30000000-0000-0000-0000-000000000002",
					ExternalCode:  "DEMO-02",
					Name:          "示例窗口二",
					Description:   "用于开发管理范围和流水导入",
					BusinessHours: "10:30-13:30 / 16:30-19:30",
				}},
			},
		},
	}}
}
