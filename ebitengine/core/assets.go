package core

import (
	"bytes"
	"embed"
	"fmt"
	"io"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const sampleRate = 44100

var (
	assetFS      embed.FS
	audioContext *audio.Context
)

var (
	imageCache   = map[string]*ebiten.Image{}
	imageCacheMu sync.Mutex

	soundCache   = map[string][]byte{}
	soundCacheMu sync.Mutex
)

// InitAssets sets the embedded filesystem and creates the audio context.
// Must be called before any Load/Play functions.
func InitAssets(fs embed.FS) {
	assetFS = fs
	audioContext = audio.NewContext(sampleRate)
}

// AudioContext returns the global audio context.
func AudioContext() *audio.Context {
	return audioContext
}

// LoadImage loads (or returns cached) an image from the embedded assets.
// name is relative to assets/ without the .png extension.
func LoadImage(name string) *ebiten.Image {
	imageCacheMu.Lock()
	defer imageCacheMu.Unlock()

	if img, ok := imageCache[name]; ok {
		return img
	}

	path := "assets/" + name + ".png"
	f, err := assetFS.Open(path)
	if err != nil {
		panic(fmt.Sprintf("assets: failed to open image %s: %v", path, err))
	}
	defer f.Close()

	img, _, err := ebitenutil.NewImageFromReader(f)
	if err != nil {
		panic(fmt.Sprintf("assets: failed to decode image %s: %v", path, err))
	}

	imageCache[name] = img
	return img
}

// loadSoundBytes loads and caches the raw decoded PCM bytes for an OGG file.
func loadSoundBytes(name string) []byte {
	soundCacheMu.Lock()
	defer soundCacheMu.Unlock()

	if data, ok := soundCache[name]; ok {
		return data
	}

	path := "assets/" + name + ".ogg"
	f, err := assetFS.Open(path)
	if err != nil {
		panic(fmt.Sprintf("assets: failed to open sound %s: %v", path, err))
	}
	defer f.Close()

	stream, err := vorbis.DecodeF32(f)
	if err != nil {
		panic(fmt.Sprintf("assets: failed to decode sound %s: %v", path, err))
	}

	data, err := io.ReadAll(stream)
	if err != nil {
		panic(fmt.Sprintf("assets: failed to read sound %s: %v", path, err))
	}

	soundCache[name] = data
	return data
}

// PlaySFX plays a sound effect. Each call creates a new player so overlapping
// plays work (like the Lua clone approach).
func PlaySFX(name string) {
	data := loadSoundBytes(name)
	player, err := audioContext.NewPlayerF32(bytes.NewReader(data))
	if err != nil {
		return
	}
	player.Play()
}

// ClearAssets flushes all caches.
func ClearAssets() {
	imageCacheMu.Lock()
	imageCache = map[string]*ebiten.Image{}
	imageCacheMu.Unlock()

	soundCacheMu.Lock()
	soundCache = map[string][]byte{}
	soundCacheMu.Unlock()
}
