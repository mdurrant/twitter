package twitter

import (
	"fmt"
	"log"

	"github.com/mrjones/oauth"
)

func NewServerClient(consumerKey, consumerSecret string) *ServerClient {
	//newClient := NewClient(consumerKey, consumerKey)

	newServer := new(ServerClient)

	newServer.OAuthConsumer = oauth.NewConsumer(
		consumerKey,
		consumerSecret,
		oauth.ServiceProvider{
			RequestTokenUrl:   OAUTH_REQUES_TOKEN,
			AuthorizeTokenUrl: OAUTH_AUTH_TOKEN,
			AccessTokenUrl:    OAUTH_ACCESS_TOKEN,
		},
	)

	//Enable debug info
	newServer.OAuthConsumer.Debug(true)

	// newServer.Client = *newClient
	fmt.Println("[server] init server")
	newServer.OAuthTokens = make(map[string]*oauth.RequestToken)
	return newServer
}

type ServerClient struct {
	Client
	OAuthConsumer *oauth.Consumer
	OAuthTokens   map[string]*oauth.RequestToken
}

func (s *ServerClient) GetAuthURL(tokenUrl string) string {
	fmt.Println("[server] tokenurl=", tokenUrl)
	fmt.Printf("[server] consumer=%v \n", s.OAuthConsumer)
	token, requestUrl, err := s.OAuthConsumer.GetRequestTokenAndUrl(tokenUrl)
	fmt.Println("[server] token=", token, " requestUrl=", requestUrl, " err=", err)
	if err != nil {
		log.Fatal(err)
	}
	// Make sure to save the token, we'll need it for AuthorizeToken()
	s.OAuthTokens[token.Token] = token
	return requestUrl
}

func (s *ServerClient) CompleteAuth(tokenKey, verificationCode string) (*oauth.AccessToken, error) {
	if _, ok := s.OAuthTokens[tokenKey]; ok {
		accessToken, err := s.OAuthConsumer.AuthorizeToken(s.OAuthTokens[tokenKey], verificationCode)
		
		if err != nil {
			return nil, err
		}

		s.HttpConn, err = s.OAuthConsumer.MakeHttpClient(accessToken)
		if err != nil {
			return nil, err
		}

		return accessToken, nil
	}
	
	return nil, nil
}
