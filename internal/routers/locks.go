package routers

import (
	"bytes"
	"encoding/json"
	"io"
	"time"

	"bff/internal/controllers"
	"bff/internal/middlewares/debug"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"

	"gitlab.HIDDEN.com/HIDDEN/sms-core/concurrency/redlock"
	HIDDENcore "gitlab.HIDDEN.com/HIDDEN/sms-core/HIDDEN"
	ginutils "gitlab.HIDDEN.com/HIDDEN/sms-core/utils/gin"
)

const (
	lockRuleEntity       = "rule"
	lockGlobalRuleEntity = "global-rule"
	lockDaEntity         = "da"
	lockVGEntity         = "vendor-group"
)

type masterRuleIds []string

func (p *masterRuleIds) UnmarshalJSON(data []byte) (err error) {
	var op struct {
		Ids []string `json:"masterRuleIds"`
	}
	err = json.Unmarshal(data, &op)
	if err == nil {
		*p = op.Ids
	}
	return
}

type ids []string

func (p *ids) UnmarshalJSON(data []byte) (err error) {
	var op struct {
		Ids []string `json:"ids"`
	}
	err = json.Unmarshal(data, &op)
	if err == nil {
		*p = op.Ids
	}
	return
}

func (r *BffRouter) lockByParam(param, entity string, ttl time.Duration) gin.HandlerFunc {
	return redlock.AcquireHandler(r.redis.GetClient(), func(ctx *gin.Context) ([]redlock.LockKey, bool) {
		return []redlock.LockKey{{Entity: entity, EntityId: ctx.Param(param), OrgId: ctx.GetString("orgId")}}, true
	}, ttl)
}

func (r *BffRouter) lockById(entity string, ttl time.Duration) gin.HandlerFunc {
	return r.lockByParam("id", entity, ttl)
}

func lockByPayload[Ids ~[]string](client *redis.Client, entity string, ttl time.Duration) gin.HandlerFunc {
	return redlock.AcquireHandler(client, func(ctx *gin.Context) ([]redlock.LockKey, bool) {
		var keys []redlock.LockKey

		orgId := ctx.GetString("orgId")

		b, err := io.ReadAll(ctx.Request.Body)
		if err != nil {
			log.Warn().Err(err).Msg("Cannot read request body for lock")
			return nil, false
		}

		ctx.Request.Body = io.NopCloser(bytes.NewBuffer(b))

		var ids Ids
		if err := json.Unmarshal(b, &ids); err != nil {
			log.Warn().Err(err).Msg("Cannot parse payload for lock")
			return nil, false
		}

		for _, id := range ids {
			if id != "" {
				keys = append(keys, redlock.LockKey{Entity: entity, EntityId: id, OrgId: orgId})
			}
		}

		return keys, true
	}, ttl)
}

func (r *BffRouter) lockRule(ruleType string, ttl time.Duration) gin.HandlerFunc {
	return redlock.AcquireHandler(r.redis.GetClient(), func(ctx *gin.Context) ([]redlock.LockKey, bool) {
		var keys []redlock.LockKey

		orgId := ctx.GetString("orgId")
		id := ctx.Param("id")

		keys = append(keys, redlock.LockKey{Entity: ruleType, EntityId: id, OrgId: orgId})

		b, err := io.ReadAll(ctx.Request.Body)
		if err != nil {
			log.Warn().Err(err).Msg("Cannot read request body for rule lock")
			return nil, false
		}

		ctx.Request.Body = io.NopCloser(bytes.NewBuffer(b))

		var rule HIDDENcore.Rule
		if err := json.Unmarshal(b, &rule); err != nil {
			log.Warn().Err(err).Msg("Cannot parse rule for lock")
			return nil, false
		}

		for _, vg := range rule.VendorGroups {
			if vg.Ref != "" {
				keys = append(keys, redlock.LockKey{Entity: lockVGEntity, EntityId: id, OrgId: orgId})
			}
		}

		return keys, true
	}, ttl)
}

func (r *BffRouter) unlockFromCtx() gin.HandlerFunc {
	return redlock.ReleaseHandler(r.redis.GetClient())
}

func proxyLockTokenHeader(ctx *gin.Context) {
	if token, ok := ctx.Get(redlock.LockTokenKey); ok {
		ctx.Request.Header.Set("X-Lock-Token", token.(string))
	}
}

func (r *BffRouter) MountLocksRouter() {
	locks := r.router.Group("/api/locks")
	locks.PUT("/", ginutils.Wrap(controllers.Controller.Lock))
	locks.DELETE("/", debug.LockTimeout, ginutils.Wrap(controllers.Controller.Unlock))
}
