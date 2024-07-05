package cache

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.etcd.io/bbolt"
	"golang.org/x/crypto/bcrypt"
)

const bucketName = "CacheBucket"

type PasswordRequest struct {
	Password string `json:"password" binding:"required"`
}

func ClearCache(db *bbolt.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		var req PasswordRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := bcrypt.CompareHashAndPassword([]byte("$2a$10$cvD77jlhjwTXKLkX.KeI0Ool7Kp5HCjovhJ.mX01P18qxt4R/CbIu"), []byte(req.Password))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": errors.New("password doesn't match").Error()})
			return
		}

		err = db.Update(func(tx *bbolt.Tx) error {
			b := tx.Bucket([]byte(bucketName))
			if b == nil {
				return nil
			}
			return b.ForEach(func(k, v []byte) error {
				return b.Delete(k)
			})
		})

		if err != nil {
			log.Printf("Failed to clear cache: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "error clearing cache",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "success",
		})
	}
}
