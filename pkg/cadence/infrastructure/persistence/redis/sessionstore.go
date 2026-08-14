package redis

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	googleuuid "github.com/google/uuid"

	"cadence/pkg/cadence/app"
	"cadence/pkg/common/redis"
	"cadence/pkg/common/uuid"
)

const (
	sessionKeyPrefix      = "session:"
	userSessionsKeyPrefix = "user_sessions:"

	tokenByteLength = 32
)

type Config struct {
	TTL                time.Duration
	MaxSessionsPerUser int
}

func NewSessionStore(client redis.Client, cfg Config) app.SessionStore {
	return &sessionStore{client: client, cfg: cfg}
}

type sessionStore struct {
	client redis.Client
	cfg    Config
}

func (s *sessionStore) CreateSession(ctx context.Context, userID uuid.UUID) (string, error) {
	token, err := generateToken()
	if err != nil {
		return "", err
	}

	if err := s.client.Set(ctx, sessionKey(token), googleuuid.UUID(userID).String(), s.cfg.TTL); err != nil {
		return "", err
	}

	sessionsKey := userSessionsKey(userID)
	if err := s.client.ZAdd(ctx, sessionsKey, float64(time.Now().UnixMilli()), token); err != nil {
		return "", err
	}

	if err := s.client.Expire(ctx, sessionsKey, s.cfg.TTL); err != nil {
		return "", err
	}

	if err := s.enforceSessionLimit(ctx, sessionsKey); err != nil {
		return "", err
	}

	return token, nil
}

func (s *sessionStore) ValidateSession(ctx context.Context, token string) (uuid.UUID, error) {
	value, err := s.client.Get(ctx, sessionKey(token))
	if err != nil {
		if errors.Is(err, redis.ErrKeyNotFound) {
			return uuid.UUID{}, app.ErrSessionNotFound
		}
		return uuid.UUID{}, err
	}

	userID, err := googleuuid.Parse(value)
	if err != nil {
		return uuid.UUID{}, err
	}

	if err := s.client.Expire(ctx, sessionKey(token), s.cfg.TTL); err != nil {
		return uuid.UUID{}, err
	}

	return uuid.UUID(userID), nil
}

// enforceSessionLimit drops stale entries left behind by expired sessions, then evicts the
// oldest live sessions until the user is back within the configured limit.
func (s *sessionStore) enforceSessionLimit(ctx context.Context, sessionsKey string) error {
	members, err := s.client.ZRangeWithScores(ctx, sessionsKey)
	if err != nil {
		return err
	}

	live := make([]string, 0, len(members))
	stale := make([]string, 0)
	for _, m := range members {
		exists, existsErr := s.client.Exists(ctx, sessionKey(m.Member))
		if existsErr != nil {
			return existsErr
		}
		if exists {
			live = append(live, m.Member)
		} else {
			stale = append(stale, m.Member)
		}
	}

	if len(stale) > 0 {
		if err := s.client.ZRem(ctx, sessionsKey, stale...); err != nil {
			return err
		}
	}

	if s.cfg.MaxSessionsPerUser <= 0 {
		return nil
	}

	if excess := len(live) - s.cfg.MaxSessionsPerUser; excess > 0 {
		evicted := live[:excess]
		sessionKeys := make([]string, len(evicted))
		for i, token := range evicted {
			sessionKeys[i] = sessionKey(token)
		}
		if err := s.client.Delete(ctx, sessionKeys...); err != nil {
			return err
		}
		if err := s.client.ZRem(ctx, sessionsKey, evicted...); err != nil {
			return err
		}
	}

	return nil
}

func generateToken() (string, error) {
	buf := make([]byte, tokenByteLength)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return hex.EncodeToString(buf), nil
}

func sessionKey(token string) string {
	return sessionKeyPrefix + token
}

func userSessionsKey(userID uuid.UUID) string {
	return userSessionsKeyPrefix + googleuuid.UUID(userID).String()
}
