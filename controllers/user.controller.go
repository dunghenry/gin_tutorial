package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
	"trandung/gin_tutorial/configs"
	"trandung/gin_tutorial/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

func hashPassword(password string) string {
	fmt.Printf("Value: %v\n", password)
	rounds, _ := strconv.Atoi(os.Getenv("ROUNDS"))
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), rounds)
	return string(hash)
}
func comparePassword(hashedPassword string, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return err
	}
	return nil
}

var userCollection = configs.GetCollection(configs.DB, "users")

func GetUsers(c *gin.Context) {
	var rs = c.Request.Header["Token"][0]
	if rs == "Bearer" || rs == "" || rs == "Bearer " {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "failure",
			"message": "You're not authenticated",
		})
		return
	} else {
		tokenString := strings.Split(c.Request.Header["Token"][0], " ")[1]
		if len(tokenString) == 195 {
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				return []byte(os.Getenv("ACCESS_TOKEN_SECRET")), nil
			})
			if token.Valid {
				cur, err := userCollection.Find(context.Background(), bson.D{})
				if err != nil {
					log.Fatal(err)
				}
				defer cur.Close(context.Background())

				var results []models.User
				if err = cur.All(context.Background(), &results); err != nil {
					log.Fatal(err)
				}
				c.JSON(http.StatusOK, gin.H{
					"status": "success",
					"users":  results,
				})
				return
			} else if ve, ok := err.(*jwt.ValidationError); ok {
				if ve.Errors&jwt.ValidationErrorMalformed != 0 {
					c.JSON(http.StatusBadRequest, gin.H{
						"status":  "failure",
						"message": "Token is not valid",
					})
					return
				} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
					// Token is either expired or not active yet
					c.JSON(http.StatusBadRequest, gin.H{
						"status":  "failure",
						"message": "Token is expired",
					})
					return
				} else {
					c.JSON(http.StatusBadRequest, gin.H{
						"status":  "failure",
						"message": "Couldn't handle this token",
					})
					return
				}
			} else {
				c.JSON(http.StatusBadRequest, gin.H{
					"status":  "failure",
					"message": "Couldn't handle this token",
				})
				return
			}
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "failure",
				"message": "Token is not valid",
			})
			return
		}
	}

}
func Register(c *gin.Context) {
	var newUser models.User
	newUser.Id = primitive.NewObjectID()
	if err := c.BindJSON(&newUser); err != nil {
		return
	}
	hashedPassword := hashPassword(newUser.Password)
	newUser.Password = hashedPassword
	res, err := userCollection.InsertOne(context.Background(), newUser)
	if err != nil {
		log.Fatal(err)
	}
	id := res.InsertedID
	fmt.Println(res)
	c.JSON(http.StatusCreated, gin.H{
		"status": "success",
		"_id":    id,
	})
}
func Login(c *gin.Context) {
	var user models.Login
	if err := c.BindJSON(&user); err != nil {
		return
	}
	var u models.User
	userCollection.FindOne(context.TODO(), bson.D{{"email", user.Email}}).Decode(&u)

	if u.Username == "" {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "failure",
			"message": "Email not found!",
		})
	} else {
		e := comparePassword(u.Password, user.Password)
		if e != nil {
			c.JSON(http.StatusOK, gin.H{
				"status":  "failure",
				"message": "Password is not valid!",
			})
		} else {
			userId := string(u.Id.Hex())
			claims := &jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
				NotBefore: jwt.NewNumericDate(time.Now()),
				ID:        userId,
			}
			mysecret := []byte(os.Getenv("ACCESS_TOKEN_SECRET"))
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
			tokenString, _ := token.SignedString(mysecret)
			var result models.LoginSuccess
			result.Id = u.Id
			result.Username = u.Username
			result.Email = u.Email
			result.Token = tokenString
			c.JSON(http.StatusOK, gin.H{
				"status": "success",
				"data":   result,
			})
		}

	}
}

func GetUserById(c *gin.Context) {
	var id = c.Param("id")
	todoIdObject, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  "failure",
			"message": "Id invalid",
		})
		return
	}
	var result models.User
	err = userCollection.FindOne(context.TODO(), bson.D{{"_id", todoIdObject}}).Decode(&result)
	if result.Username == "" && result.Email == "" {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "failure",
			"message": "Todo not found!",
		})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status": "success",
			"user":   result,
		})
		return
	}
}
