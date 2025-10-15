package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/cgrates/cgrates/engine"
	"github.com/gin-gonic/gin"
	"gitlab.HIDDEN.com/HIDDEN/sms-core/fp"
	"gitlab.HIDDEN.com/HIDDEN/sms-core/logger"
	"gitlab.HIDDEN.com/HIDDEN/sms-core/types"
	"gitlab.HIDDEN.com/HIDDEN/sms-core/utils"
	"io"
	"net/http"
	"strconv"
)

func (rc *RestController) Accounts(c *gin.Context) {
	tenantName := c.Param("tenantName")
	logger.Info().Msgf("Accounts FOR %v", tenantName)
	tenantName = fp.Ternary(tenantName == "", "HIDDEN", tenantName)
	list, h, err := rc.getAccounts(c, tenantName)

	if err != nil {
		logger.Error().Err(err).Msg("error getting accounts")
		utils.SendErrorResponse(c, err)
		return
	}

	if list == nil {
		utils.SendOkSliceResponse[engine.Account](c, nil, nil)
		return
	}

	utils.SendOkSliceResponse(c, h, &list)
	return

}

func (rc *RestController) getAccounts(c *gin.Context, tenantName string) ([]engine.Account, *types.ReplyMessageHeaders, error) {
	requestURL := fmt.Sprintf("%s/api/ping/%s", rc.config.HIDDENProviderHost, tenantName)
	res, headers, err := utils.HttpGet(c, requestURL)

	if err != nil {
		logger.Error().Err(err).Msg("HttpGet")
		return nil, nil, err
	}

	defer res.Body.Close()

	if res.StatusCode == http.StatusOK {
		bodyBytes, err := io.ReadAll(res.Body)
		if err != nil {
			logger.Error().Err(err).Msg("ReadAll")
			return nil, nil, err
		}

		var accounts []engine.Account
		err = json.Unmarshal(bodyBytes, &accounts)

		if err != nil {
			logger.Error().Err(err).Msg("Unmarshal")
			return nil, nil, err
		}

		headers.Total = strconv.Itoa(len(accounts))

		return accounts, headers, err
	}

	return nil, headers, nil
}
