package audio_engine

type Engine interface {
	LoadSound(folderPath string)
	PlaySound(audioFile string)
}

type AudioPlayer struct {
	Engine
}

func NewAudioPlayer() *AudioPlayer {
	return &AudioPlayer{
		Engine: NewEngineBeep(),
	}
}
