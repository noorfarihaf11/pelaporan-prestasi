package utils 
 
import ( 
    "pelaporan-prestasi/app/model"
    "time" 
 
    "github.com/golang-jwt/jwt/v5" 
) 
 
var jwtSecret = []byte("your-secret-key-min-32-characters-long") 
 
func GenerateToken(user model.User) (string, error) { 
    claims := model.JWTClaims{ 
        UserID:   user.ID, 
        Username: user.Username, 
        RegisteredClaims: jwt.RegisteredClaims{ 
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), 
            IssuedAt:  jwt.NewNumericDate(time.Now()), 
        }, 
    } 
 
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims) 
    return token.SignedString(jwtSecret) 
} 

func GenerateRefreshToken(user model.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID.String(),
		"exp":     time.Now().Add(7 * 24 * time.Hour).Unix(), // 7 hari
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

 
func ValidateToken(tokenString string) (*model.JWTClaims, error) { 
    token, err := jwt.ParseWithClaims(tokenString, &model.JWTClaims{}, 
func(token *jwt.Token) (interface{}, error) { 
        return jwtSecret, nil 
    }) 
 
    if err != nil { 
        return nil, err 
    } 
 
    if claims, ok := token.Claims.(*model.JWTClaims); ok && token.Valid { 
        return claims, nil 
    } 
 
    return nil, jwt.ErrInvalidKey 
}