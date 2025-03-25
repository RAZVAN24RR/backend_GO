package controllers

import (
	"context"
	"net/http"
	"time"

	"example.com/go-mongo-auth/config"
	"example.com/go-mongo-auth/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
)

// jwtSecret este cheia folosită pentru semnarea token-urilor JWT.
// În producție, aceasta trebuie gestionată prin variabile de mediu sau un secret manager.
var jwtSecret = []byte("secret-key")

// Register gestionează înregistrarea unui nou utilizator.
func Register(c *gin.Context) {
	var creds models.Credentials
	if err := c.ShouldBindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Date JSON invalide"})
		return
	}

	newUser := models.User{
		Username: creds.Username,
		Password: creds.Password, // Pentru demonstrație; în producție, criptează parola!
	}

	collection := config.MongoClient.Database("testdb").Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := collection.InsertOne(ctx, newUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Nu s-a putut crea utilizatorul"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Utilizator creat cu succes"})
}

// Login validează credențialele și generează un token JWT la autentificare.
func Login(c *gin.Context) {
	var creds models.Credentials
	if err := c.ShouldBindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Date JSON invalide"})
		return
	}

	var user models.User
	collection := config.MongoClient.Database("testdb").Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := collection.FindOne(ctx, bson.M{"username": creds.Username}).Decode(&user)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Utilizator inexistent"})
		return
	}

	if user.Password != creds.Password {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Parolă incorectă"})
		return
	}

	// Creăm token-ul JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": creds.Username,
		"exp":      time.Now().Add(72 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Nu s-a putut genera token-ul"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

// Protected este un handler pentru ruta protejată.
func Protected(c *gin.Context) {
	claims := c.MustGet("user").(jwt.MapClaims)
	c.JSON(http.StatusOK, gin.H{
		"message": "Ai acces la ruta protejată!",
		"user":    claims,
	})
}