package handlers

import (
	"context"
	"log"
	"net/http"
	"notification-service/internal/db"
	"notification-service/internal/models"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

func CreateNotificationHandler(c *gin.Context) {
	var notification models.Notification
	if err := c.ShouldBindJSON(&notification); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	notification.CreatedAt = time.Now()
	notification.IsRead = false

	_, err := db.NotificationCollection.InsertOne(context.Background(), notification)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create notification"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Notification created successfully"})
}

func MarkAsReadHandler(c *gin.Context) {
	userID := c.Param("userID")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"user_id": userID, "is_read": false}
	update := bson.M{"$set": bson.M{"is_read": true}}

	opts := options.Update().SetMulti(true)
	_, err := db.NotificationCollection.UpdateMany(ctx, filter, update, opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to mark notifications as read"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Notifications marked as read"})
}

func GetNotificationsHandler(c *gin.Context) {
	userID := c.Param("userID")
	var notifications []models.Notification

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := db.NotificationCollection.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch notifications"})
		return
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var notification models.Notification
		if err := cursor.Decode(&notification); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode notification"})
			return
		}
		notifications = append(notifications, notification)
	}

	c.JSON(http.StatusOK, notifications)
}

const batchSize = 100

func MarkAsReadHandlerBatch(c *gin.Context) {
	userID := c.Param("userID")

	var wg sync.WaitGroup
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := db.NotificationCollection.Find(ctx, bson.M{"user_id": userID, "is_read": false})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch notifications"})
		return
	}
	defer cursor.Close(ctx)

	// Process notifications in batches
	batch := make([]interface{}, 0, batchSize)
	for cursor.Next(ctx) {
		var notification bson.M
		if err := cursor.Decode(&notification); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode notification"})
			return
		}
		batch = append(batch, notification["_id"])
		if len(batch) >= batchSize {
			wg.Add(1)
			go markBatchAsRead(batch, &wg)
			batch = make([]interface{}, 0, batchSize)
		}
	}

	// Process the remaining batch
	if len(batch) > 0 {
		wg.Add(1)
		go markBatchAsRead(batch, &wg)
	}

	wg.Wait()
	c.JSON(http.StatusOK, gin.H{"message": "Notifications marked as read"})
}

func markBatchAsRead(batch []interface{}, wg *sync.WaitGroup) {
	defer wg.Done()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := db.NotificationCollection.UpdateMany(ctx, bson.M{"_id": bson.M{"$in": batch}}, bson.M{"$set": bson.M{"is_read": true}})
	if err != nil {
		log.Printf("Failed to mark batch as read: %v", err)
	}
}
