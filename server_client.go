package twitter

import (
	"log"

	"github.com/mrjones/oauth"
)

func NewServerClient(consumerKey, consumerSecret string) *ServerClient {
	newClient := NewClient(consumerKey, consumerKey)
	newServer := new(ServerClient)
	newServer.Client = *newClient
	newServer.OAuthTokens = make(map[string]*oauth.RequestToken)
	return newServer
}

type ServerClient struct {
	Client
	OAuthTokens map[string]*oauth.RequestToken
}

func (s *ServerClient) GetAuthURL(tokenUrl string) string {
	token, requestUrl, err := s.OAuthConsumer.GetRequestTokenAndUrl(tokenUrl)
	if err != nil {
		log.Fatal(err)
	}
	// Make sure to save the token, we'll need it for AuthorizeToken()
	s.OAuthTokens[token.Token] = token
	return requestUrl
}

func (s *ServerClient) CompleteAuth(tokenKey, verificationCode string) error {
	accessToken, err := s.OAuthConsumer.AuthorizeToken(s.OAuthTokens[tokenKey], verificationCode)
	if err != nil {
		log.Fatal(err)
	}

	s.HttpConn, err = s.OAuthConsumer.MakeHttpClient(accessToken)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}
