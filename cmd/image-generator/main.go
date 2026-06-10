package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const (
	falAPIEndpoint = "https://queue.fal.run/fal-ai/flux/dev"
)

type FalRequest struct {
	Prompt       string `json:"prompt"`
	OutputFormat string `json:"output_format,omitempty"`
}

type FalResponse struct {
	Images []struct {
		URL string `json:"url"`
	} `json:"images"`
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: image-generator <path-to-image>.prompt.md")
		os.Exit(1)
	}

	promptFile := os.Args[1]
	if err := run(promptFile); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}

func run(promptFile string) error {
	if !strings.HasSuffix(promptFile, ".prompt.md") {
		return fmt.Errorf("prompt file must have .prompt.md extension")
	}

	apiKey := os.Getenv("FAL_KEY")
	if apiKey == "" {
		return fmt.Errorf("FAL_KEY environment variable is not set")
	}

	promptBytes, err := os.ReadFile(promptFile)
	if err != nil {
		return fmt.Errorf("failed to read prompt file: %w", err)
	}

	prompt := strings.TrimSpace(string(promptBytes))
	if prompt == "" {
		return fmt.Errorf("prompt file is empty")
	}

	outputPath := strings.TrimSuffix(promptFile, ".prompt.md") + ".png"

	imageURL, err := generateImage(apiKey, prompt)
	if err != nil {
		return fmt.Errorf("failed to generate image: %w", err)
	}

	if err := downloadImage(imageURL, outputPath); err != nil {
		return fmt.Errorf("failed to save image: %w", err)
	}

	fmt.Printf("Image saved to: %s\n", outputPath)
	return nil
}

func generateImage(apiKey, prompt string) (string, error) {
	reqBody := FalRequest{
		Prompt:       prompt,
		OutputFormat: "png",
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", falAPIEndpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Key "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var falResp FalResponse
	if err := json.Unmarshal(body, &falResp); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	if len(falResp.Images) == 0 {
		return "", fmt.Errorf("no images returned from API")
	}

	return falResp.Images[0].URL, nil
}

func downloadImage(url, outputPath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download image: status %d", resp.StatusCode)
	}

	dir := filepath.Dir(outputPath)
	if dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}

	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	return err
}
