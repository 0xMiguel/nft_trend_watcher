package monitor

import (
	"encoding/json"
	"github.com/DisgoOrg/disgohook"
	"github.com/DisgoOrg/disgohook/api"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type TokenData struct {
	AveragePriceInEth float64   `json:"averagePriceInEth"`
	Count             int       `json:"count"`
	DailyVolumes      []float64 `json:"dailyVolumes"`
	MaxPriceInEth     float64   `json:"maxPriceInEth"`
	MinPriceInEth     float64   `json:"minPriceInEth"`
	VolumeInEth       float64   `json:"volumeInEth"`
	DeltaStats        struct {
		AveragePriceInEth float64 `json:"averagePriceInEth"`
		Count             int     `json:"count"`
		MaxPriceInEth     float64 `json:"maxPriceInEth"`
		MinPriceInEth     float64 `json:"minPriceInEth"`
		VolumeInEth       float64 `json:"volumeInEth"`
	} `json:"deltaStats"`
	Address     string `json:"address"`
	ExternalURL string `json:"externalUrl"`
	ImageURL    string `json:"imageUrl"`
	Name        string `json:"name"`
	Symbol      string `json:"symbol"`
	UUID        string `json:"uuid"`
}

type IcyToolsResponse struct {
	Data []TokenData `json:"data"`
	Total int `json:"total"`
}

type monitor struct {
	CurrentTokens map[string]TokenData
}

var (
	httpClient = http.Client{Timeout: 10 * time.Second}
	cookieValue = "token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZGRyZXNzIjoiMHhiZjlhZjI1YTE0ZTdhYTQ1NTk5ZmNkZWYzYTg5OWIzOGY0N2VjY2I3IiwiaWF0IjoxNjMyODc5NzQyLCJhdWQiOiJhY2NvdW50OnJlZ2lzdGVyZWQifQ.r0hw2ZDIxYFOX4lJewjw__IOEXzxMwuwNRZOBrfeOgc"
	apiUrl = "https://icy.tools/api/collections/trending?period=15m"
	Task   = new(monitor)
)


func (m *monitor) StartMonitor()  {

	m.CurrentTokens = map[string]TokenData{}
	for {
		Task.GetLatest()
		time.Sleep(2 * time.Minute)
	}
}

func (m *monitor) readAndClose(r io.ReadCloser) ([]byte, error) {
	readBytes, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return readBytes, r.Close()
}

func (m *monitor) GetLatest() {
	log.Println("Getting latest trending...")

	r, err := http.NewRequest("GET", apiUrl, nil)
	if err != nil {
		log.Println("error creating new request", err)
	}
	r.Header.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/94.0.4606.61 Safari/537.36")
	r.Header.Set("cookie", cookieValue)

	resp, err := httpClient.Do(r)
	if err != nil {
		log.Println("error executing request", err)
	}

	var response IcyToolsResponse
	body, err := m.readAndClose(resp.Body)
	if err != nil {
		log.Println("error reading response body")
	}

	if err := json.Unmarshal(body, &response); err != nil{
		log.Println("error parsing response body")
	}

	for _, t := range response.Data {
		_,dataExists :=  m.CurrentTokens[t.Name]

		if !dataExists {
			log.Println("Found new trending token", t.Name)
			m.CurrentTokens[t.Name] = t
			m.SendHook(t.Name, t.Address, t.ImageURL, t.ExternalURL)
		}

	}
}

func (m *monitor) SendHook(tokenName string, tokenContract string, image string, website string)  {
	//https://discord.com/api/webhooks/892767666933211216/A1qcumNl8LjL2Gw8z-WttPLyu4cEnwuDXX8Ff1gFmEJAaCleTo7F9myP2HfV9mgQlc13
	hook, err := disgohook.NewWebhookClientByToken(nil, nil, "892767666933211216/A1qcumNl8LjL2Gw8z-WttPLyu4cEnwuDXX8Ff1gFmEJAaCleTo7F9myP2HfV9mgQlc13")
	if err != nil {
		log.Println("error creating webhook client", err)
	}

	embed := api.NewEmbedBuilder()

	embed.SetAuthorIcon("https://cdn.shopify.com/s/files/1/1061/1924/products/Money_Bag_Emoji_large.png?v=1571606064")
	embed.SetAuthorName("New trend found")

	embed.AddField("Token name", tokenName, false)
	embed.AddField("Contract", tokenContract, false)
	embed.AddField("Website", website, false)


	_, err = hook.SendEmbeds(embed.Build())

	if err != nil {
		log.Println("error sending webhook", err)
	}

	log.Println("Sent webhook")
}
