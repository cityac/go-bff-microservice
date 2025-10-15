package controllers

import (
	"github.com/gin-gonic/gin"
	"gitlab.HIDDEN.com/HIDDEN/sms-core/utils"
)

func (rc *RestController) TestViewPartners(c *gin.Context) {
	utils.SendOkResponse(c, nil, "READ WORKS")
}

func (rc *RestController) TestViewDraftPartners(c *gin.Context) {
	utils.SendOkResponse(c, nil, "CREATE WORKS")
}

func (rc *RestController) TestCreateAndEditPartners(c *gin.Context) {
	utils.SendOkResponse(c, nil, "EDIT WORKS")
}

//func (rc *RestController) TestPostMutex(c *gin.Context) {
//	var clientPayload redis.RedisPayload
//	ttlParam := c.DefaultQuery("ttl", strconv.Itoa(30)) // default 30s
//
//	ttl, err := strconv.ParseInt(ttlParam, 10, 64)
//	if err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{
//			"message": "Invalid TTL parameter, must be a number",
//		})
//		return
//	}
//
//	if err := c.ShouldBindJSON(&clientPayload); err != nil {
//		utils.SendErrorResponse(c, err, customErrors.JSONMarshal)
//		return
//	}
//
//	redisKey := redis.LockKey{
//		Action:     clientPayload.Action,
//		Collection: clientPayload.Collection,
//	}
//
//	// Try to acquire the lock
//	isLockSet, err := rc.redis.AcquireLock(redisKey, clientPayload, ttl)
//	if err != nil {
//		fmt.Println("Error acquiring lock:", err)
//
//		c.JSON(http.StatusInternalServerError, gin.H{
//			"message": "Internal Server Error",
//		})
//		return
//	}
//
//	// Check if key is already in Redis store
//	if !isLockSet {
//		valueInRedis, err := rc.redis.GetValue(redisKey)
//		if err != nil {
//			fmt.Println("Error during getting value from Redis:", err)
//		}
//
//		// Only lock owner can replace current lock. For other users must be locked
//		if valueInRedis.UserID == clientPayload.UserID {
//			c.JSON(http.StatusOK, gin.H{
//				"message": fmt.Sprintf("Lock sucessfully created and will be automatically released after %v seconds!", ttl),
//			})
//			return
//		}
//
//		// Get TTL for current lock
//		ttlInSeconds, err := rc.redis.CheckTTL(redisKey)
//		if err != nil {
//			fmt.Println(err)
//
//			c.JSON(http.StatusInternalServerError, gin.H{
//				"message": "Internal Server Error",
//			})
//		}
//
//		msg := fmt.Sprintf("The request can't be proccesed! User with ID: %v already working with this resourse. Try after %v seconds!", valueInRedis.UserID, ttlInSeconds)
//
//		c.JSON(http.StatusConflict, gin.H{
//			"message": msg,
//		})
//		return
//	}
//
//	c.JSON(http.StatusOK, gin.H{
//		"message": fmt.Sprintf("Lock sucessfully created and will be automatically released after %v seconds!", ttl),
//	})
//}
