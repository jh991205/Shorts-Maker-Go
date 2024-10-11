package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	aai "github.com/AssemblyAI/assemblyai-go-sdk"
	"cloud.google.com/go/texttospeech/apiv1"
	texttospeechpb "google.golang.org/genproto/googleapis/cloud/texttospeech/v1"
	"github.com/joho/godotenv"
	"github.com/vartanbeno/go-reddit/v2/reddit"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func getCurrentDatetime() string {
	return time.Now().Format("20060102_150405")
}

func fetchRedditPost() (string, string, error) {
	client, err := reddit.NewClient(reddit.Credentials{
		ID:       os.Getenv("REDDIT_CLIENT_ID"),
		Secret:   os.Getenv("REDDIT_CLIENT_SECRET"),
		Username: os.Getenv("REDDIT_USERNAME"),
		Password: os.Getenv("REDDIT_PASSWORD"),
	})
	if err != nil {
		return "", "", err
	}

	posts, _, err := client.Subreddit.NewPosts(context.Background(), "AmItheAsshole", &reddit.ListOptions{Limit: 1})
	if err != nil {
		return "", "", err
	}

	if len(posts) > 0 {
		return posts[0].Title, posts[0].SelfText, nil
	}

	return "", "", fmt.Errorf("No posts found")
}

func generateSpeech(content, outputPath string) (string, error) {
	ctx := context.Background()
	client, err := texttospeech.NewClient(ctx, option.WithCredentialsFile("path/to/your-google-credentials.json"))
	if err != nil {
		return "", err
	}

	req := &texttospeechpb.SynthesizeSpeechRequest{
		Input: &texttospeechpb.SynthesisInput{
			InputSource: &texttospeechpb.SynthesisInput_Text{Text: content},
		},
		Voice: &texttospeechpb.VoiceSelectionParams{
			LanguageCode: "en-US",
			SsmlGender:   texttospeechpb.SsmlVoiceGender_NEUTRAL,
		},
		AudioConfig: &texttospeechpb.AudioConfig{
			AudioEncoding: texttospeechpb.AudioEncoding_MP3,
		},
	}

	resp, err := client.SynthesizeSpeech(ctx, req)
	if err != nil {
		return "", err
	}

	err = ioutil.WriteFile(outputPath, resp.AudioContent, 0644)
	if err != nil {
		return "", err
	}

	return outputPath, nil
}

func editVideo(audioPath, videoPath, outputPath string) error {
	start := rand.Intn(480)
	ffmpegCmd := exec.Command("ffmpeg", "-i", videoPath, "-i", audioPath,
		"-ss", fmt.Sprintf("%d", start), "-t", "58", "-vf", "scale=1080:1920",
		"-c:v", "libx264", "-c:a", "aac", outputPath)
	return ffmpegCmd.Run()
}

func transcribeVideo(videoPath, srtOutputPath string) error {
	client := aai.NewClient(os.Getenv("ASSEMBLYAI_API_KEY"))

	// Open video file to be transcribed
	f, err := os.Open(videoPath)
	if err != nil {
		return err
	}
	defer f.Close()

	// Transcribe the video using AssemblyAI
	transcript, err := client.Transcripts.TranscribeFromReader(context.TODO(), f, nil)
	if err != nil {
		return err
	}

	// Write the SRT file
	err = ioutil.WriteFile(srtOutputPath, []byte(*transcript.Text), 0644)
	if err != nil {
		return err
	}

	return nil
}

func addSubtitlesToVideo(videoPath, srtPath, outputPath string) error {
	// Add subtitles to the video using ffmpeg
	ffmpegCmd := exec.Command("ffmpeg", "-i", videoPath, "-vf", fmt.Sprintf("subtitles=%s", srtPath), outputPath)
	return ffmpegCmd.Run()
}

func uploadToYouTube(videoPath, title, description string) error {
	ctx := context.Background()

	// Use OAuth2 for Google YouTube API
	service, err := youtube.NewService(ctx, option.WithCredentialsFile("path/to/credentials.json"))
	if err != nil {
		return err
	}

	// Create YouTube video metadata
	snippet := &youtube.VideoSnippet{
		Title:       title,
		Description: description,
		CategoryId:  "22",
	}
	status := &youtube.VideoStatus{PrivacyStatus: "public"}

	video := &youtube.Video{Snippet: snippet, Status: status}
	file, err := os.Open(videoPath)
	if err != nil {
		return err
	}
	defer file.Close()

	request := service.Videos.Insert([]string{"snippet", "status"}, video)
	request.Media(file)
	_, err = request.Do()
	return err
}

func handleProcessVideo(w http.ResponseWriter, r *http.Request) {
	title, content, err := fetchRedditPost()
	if err != nil {
		log.Println("Error fetching Reddit post:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	speechPath := filepath.Join("generated", "speech_"+getCurrentDatetime()+".mp3")
	videoPath := filepath.Join("generated", "video_"+getCurrentDatetime()+".mp4")
	srtPath := filepath.Join("generated", getCurrentDatetime()+".srt")
	finalVideoPath := filepath.Join("generated", "final_video_"+getCurrentDatetime()+".mp4")

	_, err = generateSpeech(content, speechPath)
	if err != nil {
		log.Println("Error generating speech:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = editVideo(speechPath, "minecraft.mp4", videoPath)
	if err != nil {
		log.Println("Error editing video:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = transcribeVideo(videoPath, srtPath)
	if err != nil {
		log.Println("Error transcribing video:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = addSubtitlesToVideo(videoPath, srtPath, finalVideoPath)
	if err != nil {
		log.Println("Error adding subtitles to video:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = uploadToYouTube(finalVideoPath, title+" #shorts", content+" #aitah #aita #shorts")
	if err != nil {
		log.Println("Error uploading to YouTube:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Video processing completed successfully"))
}

func main() {
	loadEnv()

	http.HandleFunc("/process_video", handleProcessVideo)
	log.Fatal(http.ListenAndServe(":8000", nil))
}
