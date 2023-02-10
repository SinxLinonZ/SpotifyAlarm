package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2"
	"gopkg.in/yaml.v3"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
)

type AppSetting struct {
	RedirectURL  string
	ClientId     string
	ClientSecret string
}

func getClient() *spotify.Client {
	var client *spotify.Client

	var appSetting AppSetting
	b, err := os.ReadFile("config.yaml")
	if err != nil {
		log.Fatal(err)
	}
	err = yaml.Unmarshal(b, &appSetting)
	if err != nil {
		log.Fatal(err)
	}

	auth := spotifyauth.New(
		spotifyauth.WithRedirectURL(appSetting.RedirectURL),
		spotifyauth.WithScopes(
			spotifyauth.ScopeStreaming,
			spotifyauth.ScopeUserReadPrivate,
			spotifyauth.ScopeUserReadPlaybackState,
			spotifyauth.ScopeUserModifyPlaybackState,
			spotifyauth.ScopeUserReadCurrentlyPlaying,
		),
		spotifyauth.WithClientID(appSetting.ClientId),
		spotifyauth.WithClientSecret(appSetting.ClientSecret),
	)

	if _, err := os.Stat("token.json"); os.IsNotExist(err) {
		// file does not exist

		token := getTokenFromOAuth(auth)
		//log.Println("Auth successful!")

		// save token to file
		saveTokenToFile(token)
		client = spotify.New(auth.Client(context.Background(), token))
	} else {
		// file exists, read token from file
		token := getTokenFromFile()
		client = spotify.New(auth.Client(context.Background(), token))
	}

	return client
}

func getTokenFromOAuth(auth *spotifyauth.Authenticator) *oauth2.Token {
	state := createState()
	url := auth.AuthURL(state)

	//log.Println("Please log in to Spotify by visiting the following page in your browser:")
	//log.Println(url)
	openBrowser(url)
	token := waitForAuthToken(auth, state)

	return token
}

func getTokenFromFile() *oauth2.Token {
	var token oauth2.Token
	file, err := os.Open("token.json")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&token)
	if err != nil {
		log.Fatal(err)
	}

	return &token
}

func waitForAuthToken(auth *spotifyauth.Authenticator, state string) *oauth2.Token {
	c := make(chan *oauth2.Token, 255)
	m := http.NewServeMux()
	s := http.Server{Addr: ":13333", Handler: m}
	m.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		token, err := auth.Token(r.Context(), state, r)
		if err != nil {
			http.Error(w, "Couldn't get token", http.StatusNotFound)
			return
		}
		c <- token

		flusher, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "Server does not support Flusher!",
				http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "text/html")
		_, _ = w.Write([]byte(`<html><head><title>Auth successful</title></head><body><h1>Auth successful</h1><p>You can close this window now.</p><script>setTimeout(function(){window.close()}, 1000);</script></body></html>`))

		flusher.Flush()

		if err := s.Shutdown(context.Background()); err != nil {
			log.Fatal(err)
		}
	})
	if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}

	client := <-c
	return client
}

func createState() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}

func openBrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Fatal(err)
	}
}

func saveTokenToFile(token *oauth2.Token) {
	b, err := json.Marshal(token)
	if err != nil {
		log.Fatal(err)
	}
	tokenFile, err := os.Create("token.json")
	if err != nil {
		log.Fatal(err)
	}
	defer func(tokenFile *os.File) {
		err := tokenFile.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(tokenFile)
	_, err = tokenFile.Write(b)
	if err != nil {
		log.Fatal(err)
	}
}
