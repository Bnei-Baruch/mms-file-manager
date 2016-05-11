package models

import (
	"database/sql/driver"
	"encoding/json"
)

type Exif struct {
	FileSize  string
	FileType  string
	ImageSize string

	FileTypeExtension string
	MIMEType          string
	ImageWidth        int
	ImageHeight       int
	AspectRatio       string
	FrameRate         string
	VideoBitrate      string
	MPEGAudioVersion  int
	AudioLayer        int
	AudioBitrate      string
	SampleRate        int
	ChannelMode       string
	ModeExtension     string
	OriginalMedia     bool
	Emphasis          string
	Duration          string
	MegaPixels        float64 `json:"Megapixels"`
}

func (exif Exif) Value() (driver.Value, error) {
	valueString, err := json.Marshal(exif)
	return string(valueString), err
}

func (exif *Exif) Scan(value interface{}) error {
	if err := json.Unmarshal(value.([]byte), &exif); err != nil {
		return err
	}
	return nil
}
