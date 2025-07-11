package api

import (
	"elitecode/internal/storage"
)

type SignupRequest struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	UsernameOrEmail string `json:"usernameOrEmail"`
	Password        string `json:"password"`
}

type AuthResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

type User struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Stats    Stats  `json:"stats"`
}

type Stats struct {
	ProblemsSolved int     `json:"problemsSolved"`
	Accuracy       float64 `json:"accuracy"`
	CurrentStreak  int     `json:"currentStreak"`
	MaxStreak      int     `json:"maxStreak"`
}

func (c *Client) Signup(req SignupRequest) (*AuthResponse, error) {
	var response AuthResponse
	err := c.Post("/auth/signup", req, &response)
	if err != nil {
		return nil, err
	}

	// Store the token
	config := storage.GetConfig()
	config.AuthToken = response.Token
	config.User = storage.UserConfig{
		ID:       response.User.ID,
		Name:     response.User.Name,
		Username: response.User.Username,
		Email:    response.User.Email,
	}
	storage.SaveConfig(config)

	return &response, nil
}

func (c *Client) Login(req LoginRequest) (*AuthResponse, error) {
	var response AuthResponse
	err := c.Post("/auth/login", req, &response)
	if err != nil {
		return nil, err
	}

	// Store the token
	config := storage.GetConfig()
	config.AuthToken = response.Token
	config.User = storage.UserConfig{
		ID:       response.User.ID,
		Name:     response.User.Name,
		Username: response.User.Username,
		Email:    response.User.Email,
	}
	storage.SaveConfig(config)

	// Set token for future requests
	c.SetAuthToken(response.Token)

	return &response, nil
}

func (c *Client) Logout() error {
	err := c.Post("/auth/logout", nil, nil)
	if err != nil {
		return err
	}

	// Clear local storage
	config := storage.GetConfig()
	config.AuthToken = ""
	config.User = storage.UserConfig{}
	storage.SaveConfig(config)

	return nil
}

func (c *Client) GetCurrentUser() (*User, error) {
	var user User
	err := c.Get("/auth/me", &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (c *Client) RefreshToken() (*AuthResponse, error) {
	var response AuthResponse
	err := c.Post("/auth/refresh", nil, &response)
	if err != nil {
		return nil, err
	}

	// Update stored token
	config := storage.GetConfig()
	config.AuthToken = response.Token
	storage.SaveConfig(config)

	// Set token for future requests
	c.SetAuthToken(response.Token)

	return &response, nil
}

func (c *Client) IsAuthenticated() bool {
	config := storage.GetConfig()
	return config.AuthToken != ""
}