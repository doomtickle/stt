package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Word struct {
	Confidence float64 `json:"confidence"`
	End        int     `json:"end"`
	Start      int     `json:"start"`
	Text       string  `json:"text"`
}

type SentimentAnalysisResult struct {
	Text       string      `json:"text"`
	Start      int         `json:"start"`
	End        int         `json:"end"`
	Sentiment  string      `json:"sentiment"`
	Confidence float64     `json:"confidence"`
	Speaker    interface{} `json:"speaker"`
}

type Entity struct {
	EntityType string `json:"entity_type"`
	Text       string `json:"text"`
	Start      int    `json:"start"`
	End        int    `json:"end"`
}

type IabCategoriesResult struct {
	Status  string `json:"status"`
	Results []struct {
		Text   string `json:"text"`
		Labels []struct {
			Relevance float64 `json:"relevance"`
			Label     string  `json:"label"`
		} `json:"labels"`
		Timestamp struct {
			Start int `json:"start"`
			End   int `json:"end"`
		} `json:"timestamp"`
	} `json:"results"`
	Summary map[string]float64 `json:"summary"`
}

type NLPResult struct {
	AcousticModel            string                    `json:"acoustic_model"`
	AudioDuration            float64                   `json:"audio_duration"`
	AudioURL                 string                    `json:"audio_url"`
	Confidence               float64                   `json:"confidence"`
	DualChannel              interface{}               `json:"dual_channel"`
	FormatText               bool                      `json:"format_text"`
	ID                       string                    `json:"id"`
	LanguageModel            string                    `json:"language_model"`
	Punctuate                bool                      `json:"punctuate"`
	Status                   string                    `json:"status"`
	Text                     string                    `json:"text"`
	Utterances               interface{}               `json:"utterances"`
	WebhookStatusCode        interface{}               `json:"webhook_status_code"`
	WebhookURL               string                    `json:"webhook_url"`
	Entities                 []Entity                  `json:"entities"`
	Words                    []Word                    `json:"words"`
	SentimentAnalysisResults []SentimentAnalysisResult `json:"sentiment_analysis_results"`
	IabCategoriesResult      `json:"iab_categories_result"`
}

const UPLOAD_URL = "https://api.assemblyai.com/v2/upload"
const TRANSCRIPT_URL = "https://api.assemblyai.com/v2/transcript"

func main() {
	loadEnv()

	u := uploadAudio()

	fmt.Printf("Upload URL: %s\n", u)

	id := startTranscription(u)

	fmt.Printf("Transcript ID: %s\n", id)

	getTranscription(id)
}

func writeJSONResponse(b []byte) {
	f, err := os.Create("out.json")

	if err != nil {
		log.Fatalln(err)
	}

	defer f.Close()

	w := bufio.NewWriter(f)
	_, err = w.Write(b)
	if err != nil {
		log.Fatalln(err)
	}

	w.Flush()
}

func loadEnv() {
	err := godotenv.Load()

	if err != nil {
		log.Fatalln(err)
	}
}

func getTranscription(id string) {
	POLLING_URL := TRANSCRIPT_URL + "/" + id
  apiKey := os.Getenv("ASSEMBLY_AI_KEY")

	// Send GET request
	client := &http.Client{}
	req, _ := http.NewRequest("GET", POLLING_URL, nil)
	req.Header.Set("content-type", "application/json")
	req.Header.Set("authorization", apiKey)
	res, err := client.Do(req)

	if err != nil {
		log.Fatalln(err)
	}

	defer res.Body.Close()

	var result NLPResult
	json.NewDecoder(res.Body).Decode(&result)

	fmt.Println("Checking status...")
	status := result.Status

	// Check status and print the transcribed text
	if status == "completed" {
		fmt.Print("Complete!\n")
		fmt.Printf("Text: %s\n", fmt.Sprint(result.Text))
		fmt.Printf("Words: %s\n", fmt.Sprint(result.Words))
		fmt.Printf("Sentiment Analysis Results: %s\n", fmt.Sprint(result.SentimentAnalysisResults))
		fmt.Printf("IAB Cateory Results: %s\n", fmt.Sprint(result.IabCategoriesResult))
		fmt.Printf("Entities: %s\n", fmt.Sprint(result.Entities))

		out, err := json.Marshal(result)
		if err != nil {
			log.Fatalln(err)
		}

		writeJSONResponse(out)
	}

	if status != "completed" {
		fmt.Printf("Current Status: %s\n", status)
		time.Sleep(5 * time.Second)
		getTranscription(id)
	}
}

func startTranscription(url string) string {
  apiKey := os.Getenv("ASSEMBLY_AI_KEY")
	fmt.Print("Starting transcription process...")
	// Prepare json data
	values := map[string]interface{}{"audio_url": url, "iab_categories": true, "entity_detection": true, "sentiment_analysis": true, "auto_chapters": true}
	jsonData, err := json.Marshal(values)

	if err != nil {
		log.Fatalln(err)
	}

	// Setup HTTP client and set header
	client := &http.Client{}
	req, _ := http.NewRequest("POST", TRANSCRIPT_URL, bytes.NewBuffer(jsonData))
	req.Header.Set("content-type", "application/json")
	req.Header.Set("authorization", apiKey)
	res, err := client.Do(req)

	if err != nil {
		log.Fatalln(err)
	}

	defer res.Body.Close()

	// Decode json and store it in a map
	var result map[string]interface{}
	json.NewDecoder(res.Body).Decode(&result)

	fmt.Print("Done.\n")

	return fmt.Sprint(result["id"])
}

func uploadAudio() string {
  apiKey := os.Getenv("ASSEMBLY_AI_KEY")
	file := os.Args[1]
	// Load file
	data, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Print("Uploading file...")
	// Setup HTTP client and set header
	client := &http.Client{}
	req, _ := http.NewRequest("POST", UPLOAD_URL, bytes.NewBuffer(data))
	req.Header.Set("authorization", apiKey)
	res, err := client.Do(req)

	if err != nil {
		log.Fatalln(err)
	}

	defer res.Body.Close()

	// Decode json and store it in a map
	var result map[string]interface{}
	json.NewDecoder(res.Body).Decode(&result)

	fmt.Print("Done.\n")

	return fmt.Sprint(result["upload_url"])
}
