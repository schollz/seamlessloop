package sox

import (
	"math"
	"os"
	"strings"
	"testing"

	log "github.com/schollz/logger"
	"github.com/stretchr/testify/assert"
)

var s = New("sss")

func TestRun(t *testing.T) {
	log.SetLevel("trace")
	stdout, stderr, err := run("sox", "--help")
	assert.Nil(t, err)
	assert.True(t, strings.Contains(stdout, "SoX"))
	assert.Empty(t, stderr)
}

func TestLength(t *testing.T) {
	length, err := s.Length("sample.wav")
	assert.Nil(t, err)
	assert.Equal(t, 1.133354, length)
}

func TestInfo(t *testing.T) {
	samplerate, channnels, err := s.Info("sample.wav")
	assert.Nil(t, err)
	assert.Equal(t, 48000, samplerate)
	assert.Equal(t, 2, channnels)
}

func TestSilence(t *testing.T) {
	fname2, err := s.SilenceAppend("sample.wav", 1)
	assert.Nil(t, err)
	length1, _ := s.Length("sample.wav")
	length2, _ := s.Length(fname2)
	assert.Less(t, math.Abs(length2-length1-1), 0.00001)

	fname2, err = s.SilencePrepend("sample.wav", 1)
	assert.Nil(t, err)
	length1, _ = s.Length("sample.wav")
	length2, _ = s.Length(fname2)
	assert.Less(t, math.Abs(length2-length1-1), 0.00001)

	fname3 := MustString(s.SilenceTrim(fname2))
	length3 := MustFloat(s.Length(fname3))
	assert.Greater(t, length2-length3, 1.0)

	os.Rename(fname3, "test.wav")
}

func TestTrim(t *testing.T) {
	var fname2 string
	var err error
	fname2, err = s.Trim("sample.wav", 0.5, 0.5)
	assert.Nil(t, err)
	assert.Equal(t, 0.5, MustFloat(s.Length(fname2)))
	fname2, err = s.Trim("sample.wav", 0.5)
	assert.Nil(t, err)
	assert.Equal(t, 0.633354, MustFloat(s.Length(fname2)))
}

func TestReverse(t *testing.T) {
	var fname2 string
	var err error
	fname2, err = s.Reverse("sample.wav")
	assert.Nil(t, err)
	assert.Equal(t, MustFloat(s.Length("sample.wav")), MustFloat(s.Length(fname2)))
}

func TestPitch(t *testing.T) {
	var fname2 string
	var err error
	fname2, err = s.Pitch("sample.wav", 3)
	assert.Nil(t, err)
	assert.Equal(t, MustFloat(s.Length("sample.wav")), MustFloat(s.Length(fname2)))
	os.Rename(fname2, "test.wav")
}

func TestJoin(t *testing.T) {
	var fname2 string
	var err error
	fname2, err = s.Join("sample.wav", "sample.wav", "sample.wav")
	assert.Nil(t, err)
	assert.LessOrEqual(t, math.Abs(MustFloat(s.Length(fname2))-3*MustFloat(s.Length("sample.wav"))), 0.001)
	os.Rename(fname2, "test.wav")
}

func TestRepeat(t *testing.T) {
	var fname2 string
	var err error
	fname2, err = s.Repeat("sample.wav", 2)
	assert.Nil(t, err)
	assert.LessOrEqual(t, math.Abs(MustFloat(s.Length(fname2))-3*MustFloat(s.Length("sample.wav"))), 0.001)
	os.Rename(fname2, "test.wav")
}

func TestRetempoSpeed(t *testing.T) {
	var fname2 string
	var err error
	fname2, err = s.RetempoSpeed("sample.wav", 60, 120)
	assert.Nil(t, err)
	assert.LessOrEqual(t, math.Abs(MustFloat(s.Length("sample.wav"))/2-MustFloat(s.Length(fname2))), 0.001)
}
func TestRetempoStretch(t *testing.T) {
	var fname2 string
	var err error
	fname2, err = s.RetempoStretch("sample.wav", 60, 120)
	assert.Nil(t, err)
	assert.LessOrEqual(t, math.Abs(MustFloat(s.Length("sample.wav"))/2-MustFloat(s.Length(fname2))), 0.001)
}

func TestCopyPaste(t *testing.T) {
	var fname2 string
	var err error
	fname2, err = s.CopyPaste("sample.wav", 0.14, 0.27, 0.57, 0.02)
	assert.Nil(t, err)
	assert.Equal(t, MustFloat(s.Length("sample.wav")), MustFloat(s.Length(fname2)))
	os.Rename(fname2, "test.wav")
}

func TestPaste(t *testing.T) {
	var fname2 string
	var err error
	crossfade := 0.04
	piece := MustString(s.Trim("sample.wav", 0.14-crossfade, 0.27+crossfade))
	fname2, err = s.Paste("sample.wav", piece, 0.57, crossfade)
	assert.Nil(t, err)
	assert.Equal(t, MustFloat(s.Length("sample.wav")), MustFloat(s.Length(fname2)))
	os.Rename(fname2, "test.wav")
}

func TestGain(t *testing.T) {
	var fname2 string
	var err error
	fname2, err = s.Gain("sample.wav", 6)
	assert.Nil(t, err)
	assert.Equal(t, MustFloat(s.Length("sample.wav")), MustFloat(s.Length(fname2)))
	os.Rename(fname2, "test.wav")
}

func TestSampleRate(t *testing.T) {
	var fname2 string
	var err error
	fname2, err = s.SampleRate("sample.wav", 8000)
	assert.Nil(t, err)
	assert.Less(t, math.Floor(MustFloat(s.Length("sample.wav"))-MustFloat(s.Length(fname2))), 0.001)
	os.Rename(fname2, "test.wav")
}

func TestStretch(t *testing.T) {
	var fname2 string
	var err error
	fname2, err = s.Stretch("sample.wav", 2)
	assert.Nil(t, err)
	assert.Less(t, math.Abs(MustFloat(s.Length("sample.wav"))*2-MustFloat(s.Length(fname2))), 0.01)
	os.Rename(fname2, "test.wav")
}

func TestStutter(t *testing.T) {
	var fname2 string
	var err error
	fname2, err = s.Stutter("sample.wav", 60.0/160/4, 0.5, 4, 0.005)
	assert.Nil(t, err)
	if fname2 != "sample.wav" {
		os.Rename(fname2, "test.wav")
	}
}

func TestFade(t *testing.T) {
	fname2, err := s.Fade("sample.wav", 0.3, 0.5)
	assert.Nil(t, err)
	os.Rename(fname2, "test.wav")
}

func TestLoopCrossfade(t *testing.T) {
	fname2, err := s.LoopCrossfade("sample.wav", 0.2)
	assert.Nil(t, err)
	os.Rename(fname2, "test.wav")
}

// keep this last
func TestClean(t *testing.T) {
	assert.Nil(t, s.Clean())
}
