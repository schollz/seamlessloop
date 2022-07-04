package seamless

import (
	"encoding/json"
	"fmt"
	"testing"

	log "github.com/schollz/logger"
	"github.com/stretchr/testify/assert"
)

func TestAudio(t *testing.T) {
	log.SetLevel("debug")
	var b []byte
	var err error
	var af *AudioFile
	for _, fname := range []string{"loop1_bpm174.wav", "amenbreak_bpm136.wav", "pad_bpm120.wav"} {
		fmt.Println(fname)
		af, err = Load(fname)
		assert.Nil(t, err)
		af, err = af.Process()
		b, _ = json.MarshalIndent(af, "", " ")
		fmt.Println(string(b))
	}
}
