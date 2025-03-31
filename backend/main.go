package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/net/context"
)

// Models
type Calculation struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Expression string             `bson:"expression" json:"expression"`
	Result     string             `bson:"result" json:"result"`
	Timestamp  time.Time          `bson:"timestamp" json:"timestamp"`
}

type CurrencyRate struct {
	Base  string             `bson:"base" json:"base"`
	Rates map[string]float64 `bson:"rates" json:"rates"`
	Date  time.Time          `bson:"date" json:"date"`
}

var calculationCollection *mongo.Collection
var currencyCollection *mongo.Collection
var ctx context.Context

func main() {
	// Set up MongoDB connection
	ctx = context.Background()
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal(err)
	}

	// Initialize collections
	calculationCollection = client.Database("calculator").Collection("calculations")
	currencyCollection = client.Database("calculator").Collection("currency_rates")

	// Set up Gin router
	router := gin.Default()

	// Configure CORS
	router.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowMethods:    []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:    []string{"Origin", "Content-Type"},
	}))

	// Routes
	router.POST("/api/calculate", calculateExpression)
	router.GET("/api/history", getCalculationHistory)
	router.DELETE("/api/history/:id", deleteCalculation)
	router.GET("/api/currency", getCurrencyRates)

	// Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	router.Run(":" + port)
}

func calculateExpression(c *gin.Context) {
	var input struct {
		Expression string `json:"expression"`
	}

	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// This is a simplified example - in reality, you'd want to use a proper expression parser
	// with security checks. Here we're just demonstrating the structure.
	result := "42.0" // Placeholder for actual calculation

	// Save calculation to database
	calculation := Calculation{
		Expression: input.Expression,
		Result:     result,
		Timestamp:  time.Now(),
	}

	_, err := calculationCollection.InsertOne(ctx, calculation)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save calculation"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": result})
}

func getCalculationHistory(c *gin.Context) {
	limit := 50 // Default limit
	if limitStr := c.Query("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil {
			limit = parsedLimit
		}
	}

	findOptions := options.Find()
	findOptions.SetSort(bson.D{{Key: "timestamp", Value: -1}})
	findOptions.SetLimit(int64(limit))

	cursor, err := calculationCollection.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch calculation history"})
		return
	}
	defer cursor.Close(ctx)

	var calculations []Calculation
	if err := cursor.All(ctx, &calculations); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode calculation history"})
		return
	}

	c.JSON(http.StatusOK, calculations)
}

func deleteCalculation(c *gin.Context) {
	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	result, err := calculationCollection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete calculation"})
		return
	}

	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Calculation not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Calculation deleted successfully"})
}

func getCurrencyRates(c *gin.Context) {
	base := c.DefaultQuery("base", "USD")

	// Check if we have recent rates in the database
	var currencyRate CurrencyRate
	err := currencyCollection.FindOne(
		ctx,
		bson.M{
			"base": base,
			"date": bson.M{"$gte": time.Now().Add(-24 * time.Hour)},
		},
	).Decode(&currencyRate)

	if err == nil {
		// We have recent rates in the database
		c.JSON(http.StatusOK, currencyRate)
		return
	}

	// If not found or outdated, fetch new rates (in a real app, this would call an external API)
	rates := map[string]float64{
		"USD": 1.0,
		"EUR": 0.93,
		"GBP": 0.79,
		"JPY": 151.15,
		"CNY": 7.24,
		"RUB": 92.5,
	}

	// Adjust rates based on the selected base currency
	if base != "USD" && rates[base] != 0 {
		baseRate := rates[base]
		for currency, rate := range rates {
			rates[currency] = rate / baseRate
		}
	}

	currencyRate = CurrencyRate{
		Base:  base,
		Rates: rates,
		Date:  time.Now(),
	}

	// Store in database
	_, err = currencyCollection.UpdateOne(
		ctx,
		bson.M{"base": base},
		bson.M{"$set": currencyRate},
		options.Update().SetUpsert(true),
	)
	if err != nil {
		log.Printf("Failed to update currency rates: %v", err)
	}

	c.JSON(http.StatusOK, currencyRate)
} 