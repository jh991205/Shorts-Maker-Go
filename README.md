
# Reddit Shorts Video Automation in Go

This project is a Go application that automates fetching Reddit posts, generating audio using Google Cloud Text-to-Speech, editing videos with FFmpeg, transcribing video subtitles with AssemblyAI, and uploading the final videos to YouTube.

## Features

- Fetches the latest post from a subreddit.
- Converts the Reddit post text to speech using Google Cloud Text-to-Speech API.
- Edits videos using FFmpeg by synchronizing the generated audio with a video file.
- Generates subtitles using AssemblyAI and burns them onto the video.
- Uploads the final video to YouTube using the YouTube Data API.

## Prerequisites

1. **Go**: Make sure Go is installed on your system. You can download it from [here](https://golang.org/dl/).
   
2. **FFmpeg**: This tool is used for video processing. Install FFmpeg using the following commands:

    - Ubuntu/Linux:
      ```bash
      sudo apt-get install ffmpeg
      ```
    - macOS:
      ```bash
      brew install ffmpeg
      ```

3. **Google Cloud API**:
   - Enable the Google Cloud Text-to-Speech and YouTube APIs on your Google Cloud project.
   - Download your `google-credentials.json` file from Google Cloud and place it in the correct path.

4. **AssemblyAI**: Get an API key from [AssemblyAI](https://www.assemblyai.com/) for video transcription.

5. **Reddit API**: Create a Reddit App [here](https://www.reddit.com/prefs/apps) to get the `client_id`, `client_secret`, `username`, and `password` for API access.

## Environment Setup

Create a `.env` file in the root of the project with the following content:

```bash
# Reddit API Credentials
REDDIT_CLIENT_ID=your_reddit_client_id
REDDIT_CLIENT_SECRET=your_reddit_client_secret
REDDIT_USERNAME=your_reddit_username
REDDIT_PASSWORD=your_reddit_password

# AssemblyAI API Key
ASSEMBLYAI_API_KEY=your_assemblyai_api_key

# Google Cloud Credentials
# Ensure the path points to the JSON file with your Google Cloud credentials
GOOGLE_APPLICATION_CREDENTIALS=path/to/your-google-credentials.json

# YouTube API OAuth2 Credentials
# This path points to the credentials.json file used for OAuth2
YOUTUBE_CREDENTIALS_PATH=path/to/credentials.json
```

## Installation

1. **Clone the repository**:

   ```bash
   git clone https://github.com/your-repo/reddit-shorts-go.git
   cd reddit-shorts-go
   ```

2. **Install dependencies**:

   Install required Go packages:

   ```bash
   go get github.com/AssemblyAI/assemblyai-go-sdk
   go get github.com/vartanbeno/go-reddit/v2/reddit
   go get cloud.google.com/go/texttospeech/apiv1
   go get google.golang.org/api/youtube/v3
   ```

3. **Install FFmpeg**:
   
   Install FFmpeg if you haven't already (instructions in Prerequisites).

4. **Prepare OAuth for YouTube**:

   - Follow [this guide](https://developers.google.com/youtube/v3/guides/auth/server-side-web-apps) to set up OAuth for the YouTube Data API.
   - Save the `credentials.json` file into the location you specified in your `.env` file.

## Usage

1. **Start the Go HTTP server**:

   Run the following command to start the server:

   ```bash
   go run main.go
   ```

   The server will start at `http://localhost:8000`.

2. **Test the `/process_video` endpoint**:

   Use any API testing tool (like Postman or Curl) to send a `POST` request to `http://localhost:8000/process_video` to trigger the video processing.
