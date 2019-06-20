package handler

import (
	"secret/app"
	"secret/handler/request"
	"secret/log"
	"secret/model"
	"secret/random"
	"time"

	"github.com/go-redis/redis"
)

// Default handles all game requests
func Default(c *app.Ctx) error {
	if c.Path.Next() == "v1" {
		return c.Next(v1)
	}
	return c.NotFound()
}

func v1(c *app.Ctx) error {
	if c.Path.Next() == "secret" {
		return c.Next(secret)

	}
	return c.NotFound()
}

func secret(c *app.Ctx) error {
	if c.Path.Next() == "" && c.Req.Method == "POST" {
		return c.Call(addSecret)
	}
	if c.Path.Next() != "" && c.Req.Method == "GET" {
		return c.Next(getSecretByHash)
	}
	return c.NotFound()
}

func addSecret(c *app.Ctx) error {
	err := c.Req.ParseForm()
	if err != nil {
		return c.InternalServerError(err)
	}
	params, err := request.NewAddSecret(c.Req.FormValue("secret"), c.Req.FormValue("expireAfterViews"), c.Req.FormValue("expireAfter"))
	if err != nil {
		return c.BadRequest(err)
	}
	err = params.Validate()
	if err != nil {
		return c.BadRequest(err)
	}
	hash := random.String(10)
	secret := model.Secret{
		Hash:           hash,
		SecretText:     params.Secret,
		CreatedAt:      time.Now(),
		ExpiresAt:      time.Now().Local().Add(params.ExpireAfter),
		RemainingViews: params.ExpireAfterViews,
	}
	log.With("params", params).Info("params")
	err = c.App.Redis(func(rc *redis.Client) error {
		return rc.Set(hash, secret, params.ExpireAfter).Err()
	})
	if err != nil {
		return err
	}
	return c.Success(secret)
}

// getSecretByHash - Find a secret by hash
func getSecretByHash(c *app.Ctx) error {
	if c.Path.Next() != "" {
		return c.NotFound()
	}
	hash := c.Path.Current()
	log.With("hash", hash).Info("hash")
	var secret model.Secret
	err := c.App.Redis(func(rc *redis.Client) error {
		data, err := rc.Get(hash).Bytes()
		if err != nil {
			return c.NotFound()
		}
		log.With("data", string(data)).Warn("data")
		err = secret.UnmarshalBinary(data)
		if err != nil {
			return err
		}
		if secret.RemainingViews == 0 {
			return c.NotFound()
		}
		secret.RemainingViews--
		expireAfter := 0 * time.Minute
		if secret.CreatedAt != secret.ExpiresAt {
			expireAfter = secret.ExpiresAt.Sub(time.Now())
		}
		return rc.Set(hash, secret, expireAfter).Err()
	})
	if err != nil {
		return err
	}
	return c.Success(secret)
}
