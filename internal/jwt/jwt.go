package jwt

import (
	"bytes"
	"crypto/rsa"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/bignyap/verifyjwt/internal/utils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/patrickmn/go-cache"
)

var certCache *cache.Cache

func init() {
	// Initialize the cache with a default expiration time of 1 hour, and cleanup interval of 30 minutes
	certCache = cache.New(1*time.Hour, 30*time.Minute)
}

func getJWKSet(issuer string) (jwk.Set, error) {
	if jwks, found := certCache.Get(issuer); found {
		return jwks.(jwk.Set), nil
	}

	// Create the form input
	form := url.Values{}

	// Fetch the public keys from the certificate endpoint
	certEndpoint, err := url.JoinPath(issuer, "protocol/openid-connect/certs")
	if err != nil {
		return nil, fmt.Errorf("error getting the public certificate: %v", err)
	}

	// Create the request
	req, err := http.NewRequest("GET", certEndpoint, bytes.NewBufferString(form.Encode()))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	body, err := utils.AppHTTPClient(req)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	jwks, err := jwk.Parse(body)
	if err != nil {
		return nil, fmt.Errorf("failed to decode JWKS: %w", err)
	}

	// Cache the jwks for future use
	certCache.Set(issuer, jwks, cache.DefaultExpiration)
	return jwks, nil
}

func extractRealmFromPath(path string) (string, error) {

	segments := strings.Split(path, "/")
	if len(segments) < 3 {
		return "", fmt.Errorf("issuer URL path does not contain enough segments to extract realm")
	}

	return segments[3], nil // The realm is the fourth segment in the path
}

func ParseAndVerifyJWT(tokenString string) (jwt.MapClaims, error) {
	// Step 1: Parse the JWT token without verifying the signature
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return jwt.MapClaims{}, fmt.Errorf("failed to parse token: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return jwt.MapClaims{}, fmt.Errorf("failed to parse claims")
	}

	// Step 2: Extract the issuer and then the realm from it
	issuer, ok := claims["iss"].(string)
	if !ok {
		return jwt.MapClaims{}, fmt.Errorf("issuer (iss) not found in the token")
	}
	issuer = strings.Replace(issuer, "localhost", "host.docker.internal", 1)

	parsedUrl, err := url.Parse(issuer)
	if err != nil {
		return jwt.MapClaims{}, fmt.Errorf("issuer (iss) not found in the token")
	}
	baseUrlHost := parsedUrl.Host
	parsedOrgUrl, err := url.Parse(os.Getenv("AUTH_URL"))
	if err != nil {
		return jwt.MapClaims{}, fmt.Errorf("wrong base url")
	}
	parsedOrgUrlHost := parsedOrgUrl.Host
	if baseUrlHost != parsedOrgUrlHost {
		return jwt.MapClaims{}, fmt.Errorf("token not issued by %s", parsedOrgUrlHost)
	}

	realm, err := extractRealmFromPath(parsedUrl.Path)
	if err != nil {
		return jwt.MapClaims{}, fmt.Errorf("realm not found in the issuer")
	}
	claims["realm"] = realm

	// Step 3: Extract the issuer and then the realm from it
	// audience, ok := claims["aud"].(string)
	// if !ok {
	// 	return jwt.MapClaims{}, fmt.Errorf("audience (aud) not found in the token")
	// }

	// Step 4: Get or fetch the JWK set
	jwks, err := getJWKSet(issuer)
	if err != nil {
		return jwt.MapClaims{}, fmt.Errorf("failed to get JWK set: %w", err)
	}

	// Get the kid from the unverified header
	unverifiedHeader := token.Header
	kid, ok := unverifiedHeader["kid"].(string)
	if !ok {
		return jwt.MapClaims{}, fmt.Errorf("kid not found in the token header")
	}

	// Find the corresponding key
	key, found := jwks.LookupKeyID(kid)
	if !found {
		return jwt.MapClaims{}, fmt.Errorf("key ID not found in the certificate endpoint")
	}

	// Extract the public key
	var pubKey rsa.PublicKey
	if err := key.Raw(&pubKey); err != nil {
		return jwt.MapClaims{}, fmt.Errorf("failed to parse public key: %w", err)
	}

	// Step 4: Verify the signature
	parsedToken, err := jwt.ParseWithClaims(
		tokenString,
		jwt.MapClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return &pubKey, nil
		},
		// jwt.WithAudience(audience),
	)
	if err != nil {
		return jwt.MapClaims{}, fmt.Errorf("failed to verify token: %w", err)
	}

	if !parsedToken.Valid {
		return jwt.MapClaims{}, fmt.Errorf("invalid token")
	}

	return claims, nil
}

func ExtractToken(r *http.Request) (string, error) {

	// Extract the token from header
	authHeader := strings.Trim(r.Header.Get("Authorization"), ";")
	if authHeader == "" {
		return "", fmt.Errorf("authorization header is missing")
	}

	// Split the token to get the token
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", fmt.Errorf("authorization header format must be Bearer {token}")
	}

	return parts[1], nil
}
