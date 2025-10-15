package controllers

import (
	"context"
	"github.com/gin-gonic/gin"
	"gitlab.HIDDEN.com/HIDDEN/sms-core/constants"
	"gitlab.HIDDEN.com/HIDDEN/sms-core/customErrors"
	"gitlab.HIDDEN.com/HIDDEN/sms-core/featureFlags"
	"gitlab.HIDDEN.com/HIDDEN/sms-core/fp"
	"gitlab.HIDDEN.com/HIDDEN/sms-core/logger"
	"gitlab.HIDDEN.com/HIDDEN/sms-core/utils"
	"time"
)

func (rc *RestController) GetFeatureFlags(c *gin.Context) {
	logger.Info().Msg("GetFeatureFlags")
	orgId := c.GetString(constants.OrgId)
	userId := c.GetString(constants.UserId)

	flags, err := featureFlags.ScanFeatureFlagsByUser(c, rc.redis.GetClient(), orgId, userId)
	if err != nil {
		logger.Error().Err(err).Msg("featureFlags.ScanFeatureFlagsByUser")
		utils.SendErrorResponse(c, err)
		return
	}

	flags = fp.Map(flags, featureFlags.ExtractFlagNameFromRedisKey)

	logger.Info().Interface("flags", flags).
		Str("userId", userId).
		Str("orgId", orgId).
		Msg("GetFeatureFlags got flags")

	utils.SendOkSliceResponse(c, nil, &flags)
}

func (rc *RestController) ReplaceAllFeatureFlagUsers(c *gin.Context) {
	orgId := c.Param("orgId")
	logger.Info().Str("orgId", orgId).Msg("ReplaceAllFeatureFlagUsers")

	payload, err := utils.ParseJSON[featureFlags.ReplaceAllFeatureFlagsUsersPayload](c)
	if err != nil {
		logger.Error().Err(err).Msg("utils.ParseJSON[featureFlags.UpdateFeatureFlagsPayload]")
		utils.SendErrorResponse(c, err, customErrors.JSONUnmarshal)
		return
	}

	err = payload.Validate()
	if err != nil {
		logger.Error().Err(err).Msg("payload.Validate")
		utils.SendErrorResponse(c, err, customErrors.BadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Second*5)
	defer cancel()

	err = featureFlags.ReplaceAllFeatureFlagUsers(ctx, rc.redis.GetClient(), payload.Flag, orgId, payload.UserIds)
	if err != nil {
		logger.Error().Err(err).Msg("featureFlags.ReplaceAllFeatureFlagUsers")
		utils.SendErrorResponse(c, err, customErrors.JSONUnmarshal)
		return
	}

	logger.Info().Str("orgId", orgId).
		Str("flag", payload.Flag).
		Interface("userIds", payload.UserIds).
		Msg("Replaced all users for flag")

	utils.SendNoContentResponse(c)
}

func (rc *RestController) UpdateFeatureFlag(c *gin.Context) {
	userId := c.Param("userId")
	orgId := c.Param("orgId")
	logger.Info().Str("orgId", orgId).Str("userId", userId).Msg("UpdateFeatureFlag")

	payload, err := utils.ParseJSON[featureFlags.UpdateSingleFeatureFlagPayload](c)
	if err != nil {
		logger.Error().Err(err).Msg("utils.ParseJSON[featureFlags.UpdateSingleFeatureFlagPayload]")
		utils.SendErrorResponse(c, err, customErrors.JSONUnmarshal)
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Second*5)
	defer cancel()

	update := fp.Ternary(payload.Enable, featureFlags.SetFeatureFlagUser, featureFlags.DeleteFeatureFlagUser)

	err = update(ctx, rc.redis.GetClient(), payload.Flag, orgId, userId)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to toggle feature flag")
		utils.SendErrorResponse(c, err, customErrors.JSONUnmarshal)
		return
	}

	logger.Info().Str("orgId", orgId).
		Str("userId", userId).
		Bool("enabled", payload.Enable).
		Str("flag", payload.Flag).
		Msg("Updated feature flag")

	utils.SendNoContentResponse(c)
}
