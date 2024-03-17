package audio_engine

import (
	"github.com/gopxl/beep"
	"github.com/gopxl/beep/speaker"
	"github.com/gopxl/beep/wav"
	"log"
	"os"
	"path/filepath"
)

type EngineBeep struct {
	buffers map[string]*beep.Buffer
}

func NewEngineBeep() *EngineBeep {
	return &EngineBeep{
		buffers: make(map[string]*beep.Buffer),
	}
}

func (e *EngineBeep) LoadSound(folderPath string) {
	log.Printf("receiving audio files list from the path: %s\n", folderPath)
	files, err := filepath.Glob(filepath.Join(folderPath, "*"))
	if err != nil || len(files) == 0 {
		log.Printf("error occurred while receiving audio files list from the path: %s %v\n", folderPath, err)
		if len(files) == 0 {
			log.Println("no files received")
		}
	}
	log.Printf("received filenames: %v\n", files)
	for _, file := range files {
		log.Printf("opening audio file: %s\n", file)
		f, err := os.Open(file)
		if err != nil {
			log.Fatal(err)
		}
		streamer, format, err := wav.Decode(f)
		if err != nil {
			log.Fatalf("error occurred while decoding file %s %v", file, err)
		}
		buffer := beep.NewBuffer(format)
		buffer.Append(streamer)
		log.Printf("adding the buffer for %s to buffers map\n", filepath.Base(file))
		e.buffers[filepath.Base(file)] = buffer
		streamer.Close()

	}
}

func (e *EngineBeep) PlaySound(audioFile string) {
	log.Printf("receiving StreamSeeker for file: %s\n", audioFile)
	log.Println(e.buffers)
	fx := e.buffers[audioFile].Streamer(0, e.buffers[audioFile].Len())
	speaker.Play(fx)
}
