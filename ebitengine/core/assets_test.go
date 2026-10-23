package core

import (
	"bytes"
	"errors"
	"image"
	"image/png"
	"io/fs"
	"os"
	"testing"
	"testing/fstest"
)

type closeErrorFile struct{ fs.File }

func (closeErrorFile) Close() error { return errors.New("close failed") }

type closeErrorFS struct{ fs.FS }

func (c closeErrorFS) Open(name string) (fs.File, error) {
	f, err := c.FS.Open(name)
	if err != nil {
		return nil, err
	}
	return closeErrorFile{f}, nil
}

func useAssets(t *testing.T, files fstest.MapFS) {
	t.Helper()
	previous := assetFS
	assetFS = closeErrorFS{files}
	ClearAssets()
	t.Cleanup(func() {
		assetFS = previous
		ClearAssets()
	})
}

func pngBytes(t *testing.T) []byte {
	t.Helper()
	var buf bytes.Buffer
	if err := png.Encode(&buf, image.NewNRGBA(image.Rect(0, 0, 2, 3))); err != nil {
		t.Fatal(err)
	}
	return buf.Bytes()
}

func TestLoadImageIgnoresCloseError(t *testing.T) {
	useAssets(t, fstest.MapFS{"assets/test.png": {Data: pngBytes(t)}})

	img := LoadImage("test")

	if w, h := img.Bounds().Dx(), img.Bounds().Dy(); w != 2 || h != 3 {
		t.Fatalf("image size = %dx%d, want 2x3", w, h)
	}
	if LoadImage("test") != img {
		t.Fatal("second load did not return the cached image")
	}
}

func TestLoadImagePanicsOnMissingFile(t *testing.T) {
	useAssets(t, fstest.MapFS{})
	defer func() {
		if recover() == nil {
			t.Fatal("expected a panic for a missing image")
		}
	}()
	LoadImage("missing")
}

func TestLoadSoundBytesIgnoresCloseError(t *testing.T) {
	ogg, err := os.ReadFile("../assets/audio/ui/click2.ogg")
	if err != nil {
		t.Fatal(err)
	}
	useAssets(t, fstest.MapFS{"assets/click.ogg": {Data: ogg}})

	data := loadSoundBytes("click")

	if len(data) == 0 {
		t.Fatal("decoded sound is empty")
	}
}
