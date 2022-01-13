package main

import (
	"backend/models"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/pascaldekloe/jwt"
	"golang.org/x/crypto/bcrypt"
)

// ログインユーザ(モック)
var validUser = models.User {
	ID: 10, 
	Email: "me@here.com",
	Password: "$2a$12$TBZJBBs0TfWdXHeujpGBn.TTwJq5V7Ra4yu.w9VV/Xgp9R3XS2YCq",
}

type Credentials struct {
	Username string `json:"email"`
	Password string `json:"password"`
}

func (app *application) Signin(w http.ResponseWriter, r *http.Request) {
	var creds Credentials

	// リクエストのJSONを構造体credsに変換する
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		app.errorJSON(w, errors.New("unauthorized"))
		return
	}

	// ハッシュ化されたパスワード
	hashedPassword := validUser.Password

	// 入力したパスワードを照合する
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(creds.Password))
	if err != nil {
		app.errorJSON(w, errors.New("unauthorized"))
		return
	}
	
	// JWTトークンのclaimsを生成する
	var claims jwt.Claims
	claims.Subject = fmt.Sprint(validUser.ID) // JWTのタイトル
	claims.Issued = jwt.NewNumericTime(time.Now()) // JWTが発行された日時
	claims.NotBefore = jwt.NewNumericTime(time.Now()) // JWTが有効になる日時
	claims.Expires = jwt.NewNumericTime(time.Now().Add(24 * time.Hour)) // JWTが失効する日時
	claims.Issuer = "mydomain.com" // JWTの発行者
	claims.Audiences = []string{"mydomain.com"} // JWTの想定利用者

	// ハッシュ関数(SHA-256)と秘密鍵から署名を作成する
	jwtBytes, err := claims.HMACSign(jwt.HS256, []byte(app.config.jwt.secret))
	if err != nil {
		app.errorJSON(w, errors.New("error signing"))
		return
	}

	// 署名(JWTトークン)をレスポンスとして返す
	app.writeJSON(w, http.StatusOK, string(jwtBytes), "response")
}