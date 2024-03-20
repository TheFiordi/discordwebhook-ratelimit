package discordwebhook

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"net/url"
	"time"
)

func SendMessage(webhook string, message Message, proxy string) error {
	// Validate parameters
	if webhook == "" {
		return errors.New("empty URL")
	}

	// Prepare the HTTP client
	var httpClient *http.Client
	if proxy != "" {
		proxyURL, err := url.Parse(proxy)
		if err != nil {
			return err
		}
		httpClient = &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(proxyURL),
			},
		}
	} else {
		httpClient = http.DefaultClient
	}

	for {
		payload := new(bytes.Buffer)

		err := json.NewEncoder(payload).Encode(message)
		if err != nil {
			return err
		}

		// Make the HTTP request
		resp, err := httpClient.Post(webhook, "application/json", payload)

		if err != nil {
			log.Printf("HTTP request failed: %v", err)
			return err
		}

		switch resp.StatusCode {
		case http.StatusOK, http.StatusNoContent:
			// Success
			err := resp.Body.Close()
			if err != nil {
				return err
			}
			return nil
		case http.StatusTooManyRequests:
			// Rate limit exceeded, retry after backoff duration
			var response DiscordResponse
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}
			err = json.Unmarshal(body, &response)
			if err != nil {
				return err
			}

			/*
				Calculate the time until reset and add it to the current local time.
				Some extra time of 750ms is added because without it I still encountered 429s.
			*/

			if response.RetryAfter != 0 {

				whole, frac := math.Modf(response.RetryAfter)
				resetAt := time.Now().Add(time.Duration(whole) * time.Second).Add(time.Duration(frac*1000) * time.Millisecond).Add(750 * time.Millisecond)
				time.Sleep(time.Until(resetAt))
			} else {
				time.Sleep(5 * time.Second)
			}

			err = resp.Body.Close()
			if err != nil {
				return err
			}
		default:
			// Handle other HTTP status codes
			err := resp.Body.Close()
			if err != nil {
				return err
			}
			responseBody, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			return fmt.Errorf("HTTP request failed with status %d, body: \n %s", resp.StatusCode, responseBody)
		}
	}
}
