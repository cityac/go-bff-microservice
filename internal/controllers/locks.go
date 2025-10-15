package controllers

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"

	"gitlab.HIDDEN.com/HIDDEN/sms-core/concurrency/redlock"
)

type LockEntry struct {
	Entity   string `json:"entity"`
	EntityId string `json:"entityId"`
	OrgId    string `json:"orgId"`
}

type LockRequest struct {
	UserId  string            `json:"userId"`
	Action  string            `json:"action"`
	Token   string            `json:"token"`
	Entries []redlock.LockKey `json:"entries"`
	TTL     int               `json:"ttl"`
}

type UnlockRequest struct {
	Token   string            `json:"token"`
	Entries []redlock.LockKey `json:"entries"`
}

func (rc *RestController) Lock(ctx *gin.Context, req LockRequest) (bool, error) {
	client := rc.redis.GetClient()

	ok, conflicts, _, err := redlock.AcquireMulti(
		client,
		req.Entries,
		time.Duration(req.TTL),
		redlock.LockPayload{
			UserID: req.UserId,
			Action: req.Action,
			Token:  req.Token,
		},
	)
	if err != nil {
		return false, err
	}

	if !ok {
		msg := "Resource is locked"
		if len(conflicts) > 0 {
			if key, err := redlock.ParseLockKey(conflicts[0].Key); err == nil {
				msg = fmt.Sprintf("Resource '%s #%s' is being updated at the moment, please try again later", key.Entity, key.EntityId)
			}
		}

		ctx.AbortWithStatusJSON(423, gin.H{"message": msg})
		return false, nil
	}

	return true, nil
}

func (rc *RestController) Unlock(_ *gin.Context, req UnlockRequest) (bool, error) {
	client := rc.redis.GetClient()

	if err := redlock.ReleaseMulti(client, req.Entries, req.Token); err != nil {
		return false, err
	}

	return true, nil
}
