package audio

import (
	"bytes"
	"io"

	ffmpeg "github.com/u2takey/ffmpeg-go"
)

func TranscodeOggToMp3(oggfile io.Reader) (io.Reader, error) {
	var out bytes.Buffer

	err := ffmpeg.
		Input("pipe:0").
		Output("pipe:1", ffmpeg.KwArgs{
			"format": "mp3",
			"acodec": "libmp3lame",
		}).
		WithInput(oggfile).
		WithOutput(&out).
		Run()

	if err != nil {
		return nil, err
	}

	return &out, nil
}
