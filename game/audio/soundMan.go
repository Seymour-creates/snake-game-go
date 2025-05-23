package audio

import (
	"bytes"
	"log"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
)

const sampleRate = 44100

type SoundManager struct {
	audioContext *audio.Context
	sounds       map[string]*audio.Player
	looping      map[string]*audio.Player // for looping sounds
}

func NewSoundManager() *SoundManager {
	ctx := audio.NewContext(sampleRate)
	return &SoundManager{
		audioContext: ctx,
		sounds:       make(map[string]*audio.Player),
		looping:      make(map[string]*audio.Player),
	}
}

func (sm *SoundManager) LoadSound(name string, data []byte) error {
	stream, err := wav.DecodeWithoutResampling(bytes.NewReader(data))
	if err != nil {
		return err
	}
	player, err := sm.audioContext.NewPlayer(stream)
	if err != nil {
		return err
	}
	sm.sounds[name] = player
	return nil
}

func (sm *SoundManager) PlaySound(name string) {
	if player, ok := sm.sounds[name]; ok {
		log.Printf("PlaySound: playing %s", name)
		player.Rewind()
		player.Play()
	} else {
		log.Printf("PlaySound: sound %s not found", name)
	}
}

func (sm *SoundManager) StopSound(name string) {
	if player, ok := sm.sounds[name]; ok {
		player.Pause()
		player.Rewind()
	}
}

func (sm *SoundManager) LoadLoopingSound(name string, data []byte) error {
	stream, err := wav.DecodeWithoutResampling(bytes.NewReader(data))
	if err != nil {
		return err
	}
	loop := audio.NewInfiniteLoop(stream, stream.Length())
	player, err := sm.audioContext.NewPlayer(loop)
	if err != nil {
		return err
	}
	sm.looping[name] = player
	return nil
}

func (sm *SoundManager) PlayLoopingSound(name string) {
	if player, ok := sm.looping[name]; ok {
		if !player.IsPlaying() {
			player.Play()
		}
	}
}

func (sm *SoundManager) PauseLoopingSound(name string) {
	if player, ok := sm.looping[name]; ok {
		player.Pause()
	} else {
		log.Printf("PauseLoopingSound: looping sound %s not found", name)
	}
}

func (sm *SoundManager) Close() {
	for _, player := range sm.sounds {
		player.Close()
	}
	for _, player := range sm.looping {
		player.Close()
	}
}

// LoopingPlayer returns the *audio.Player for a looping sound and a bool if it exists.
func (sm *SoundManager) LoopingPlayer(name string) (*audio.Player, bool) {
	player, ok := sm.looping[name]
	return player, ok
}

// ClearSounds stops and removes all loaded sounds.
func (sm *SoundManager) ClearSounds() {
	for _, player := range sm.sounds {
		player.Pause()
		player.Close()
	}
	for _, player := range sm.looping {
		player.Pause()
		player.Close()
	}
	sm.sounds = make(map[string]*audio.Player)
	sm.looping = make(map[string]*audio.Player)
}
