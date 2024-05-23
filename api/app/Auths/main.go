package Auths

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"simple_account_oauth/app/AccountStruct"
	"simple_account_oauth/app/Auths/Permissions"
	"strings"
	"sync"
	"time"
)

var (
	PermissionIDs   = []string{"read_username", "read_email", "read_last_edit", "read_create_time"}
	ReadPermissions = map[string]string{
		"read_username":    "user.username",
		"read_email":       "user.email",
		"read_last_edit":   "user.last_edit",
		"read_create_time": "user.create_time",
	}
)

type Config struct {
	AuthorizationCodeDuration time.Duration
}

type Manager struct {
	requests map[string]Request
	duration time.Duration
	mux      sync.Mutex
}

func New(config Config) *Manager {
	auths := &Manager{}
	auths.requests = make(map[string]Request)
	auths.duration = config.AuthorizationCodeDuration

	auths.ClearExpired()
	return auths
}

func (oauth *Manager) ClearExpired() {
	ticker := time.NewTicker(time.Minute)

	go func() {
		for range ticker.C {
			currentTime := time.Now().Unix()
			for key, data := range oauth.requests {
				if data.expriesTime > currentTime {
					oauth.Remove(key)
				}
			}
		}
	}()
}

func (oauth *Manager) Get(key string) *Request {
	data, exist := oauth.requests[key]
	if !exist {
		return nil
	}

	request := data
	return &request
}

func (oauth *Manager) Remove(key string) {
	oauth.mux.Lock()
	defer oauth.mux.Unlock()
	delete(oauth.requests, key)
}

func Random(length int) []byte {
	bytes := make([]byte, length)
	rand.Read(bytes)
	return bytes
}

/*
(secret, salt, hash)
*/
func GenerateSecret() (string, string, string) {
	secret := Random(32)
	salt := Random(16)

	secretStr := base64.StdEncoding.EncodeToString(secret)
	saltStr := base64.StdEncoding.EncodeToString(salt)
	hashStr := Hash(secret, salt)
	return secretStr, saltStr, hashStr
}

func Hash(secret []byte, salt []byte) string {
	hash := sha256.Sum256(append(secret, salt...))
	return base64.StdEncoding.EncodeToString(hash[:])
}

func (oauth *Manager) Add(user AccountStruct.User, app AccountStruct.ThirdApp, permissions int) (string, error) {
	rnd := Random(16)
	prev := fmt.Sprintf("%d.%s.%s", user.Id, user.Username, app.Id)
	keyBytes := append(rnd, []byte(prev)...)
	key := base64.StdEncoding.EncodeToString(keyBytes)
	key = strings.NewReplacer("+", "", "/", "", "=", "").Replace(key)

	currentTime := time.Now()
	oauth.mux.Lock()
	defer oauth.mux.Unlock()
	oauth.requests[key] = Request{
		expriesTime: currentTime.Add(oauth.duration).Unix(),
		App:         app,
		User:        user,
		permissions: permissions,
	}
	return key, nil
}

type Request struct {
	expriesTime int64
	App         AccountStruct.ThirdApp
	User        AccountStruct.User
	permissions int
}

func (req *Request) GetReadPermissions() string {
	permissions := "id"
	num := req.permissions

	if num&Permissions.Read_username != 0 {
		permissions += ",username"
	}

	if num&Permissions.Read_email != 0 {
		permissions += ",email"
	}

	if num&Permissions.Read_last_edit != 0 {
		permissions += ",last_edit"
	}

	if num&Permissions.Read_create_time != 0 {
		permissions += ",create_time"
	}

	return permissions
}
