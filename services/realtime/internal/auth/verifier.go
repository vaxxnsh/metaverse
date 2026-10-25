package auth

import (
	"github.com/vaxxnsh/metaverse/realtime/internal/ws"
	sharedauth "github.com/vaxxnsh/shared/auth"
)

type JWTVerifier struct{ secret string }

func NewJWTVerifier(secret string) *JWTVerifier { return &JWTVerifier{secret: secret} }

func (v *JWTVerifier) Verify(token string) (ws.Claims, error) {
	claims, err := sharedauth.VerifyJWT(token, v.secret)
	return ws.Claims{Subject: claims.Subject}, err
}
