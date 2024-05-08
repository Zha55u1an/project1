 // PATH: go-auth/utils/ParseToken.go

 package utils

 import (
    "go_project/internal/models"

     "github.com/dgrijalva/jwt-go"
 )

 func ParseToken(tokenString string) (claims *models.Claims, err error) {
     token, err := jwt.ParseWithClaims(tokenString, &models.Claims{}, func(token *jwt.Token) (interface{}, error) {
         return []byte("jwtkey_go_project"), nil
     })

     if err != nil {
         return nil, err
     }

     claims, ok := token.Claims.(*models.Claims)

     if !ok {
         return nil, err
     }

     return claims, nil
 }
