package seamless

import (
	"math"
	"regexp"
	"strconv"
	"strings"

	log "github.com/schollz/logger"
	"github.com/schollz/seamlessloop/src/sox"
)

type AudioFile struct {
	Filename string
	Duration float64
	BPM      float64
	Beats    float64
}

func Load(filename string, bpm ...float64) (af *AudioFile, err error) {
	af = new(AudioFile)
	af.Filename = filename

	// get bpm
	r, _ := regexp.Compile(`bpm(\d+)`)
	af.BPM, err = strconv.ParseFloat(strings.TrimPrefix(r.FindString(filename), "bpm"), 64)
	if err != nil {
		return
	}
	if len(bpm) > 0 {
		af.BPM = bpm[0]
	}

	// get duration
	af.Duration, err = sox.Length(af.Filename)
	if err != nil {
		return
	}

	// get beats
	af.Beats = af.Duration / (60 / af.BPM)

	return
}

func (af *AudioFile) Process() (af2 *AudioFile, err error) {
	// truncate with silence
	fname2, err := sox.SilenceTrimDB(af.Filename, 0.05, -50)
	if err != nil {
		return
	}
	err = sox.Copy(fname2, af.Filename+"_processed.wav")
	if err != nil {
		return
	}
	af2, err = Load(af.Filename+"_processed.wav", af.BPM)
	beatNum := 16.0
	if af2.Beats < 16 {
		beatNum = 8
	}
	if af2.Beats < 8 {
		beatNum = 4
	}
	beats := math.Round(af2.Beats/beatNum)*beatNum - af2.Beats
	log.Debugf("leftover beats: %f", beats)

	if beats < 0 {
		// crossfade
		fadeTime := math.Abs(beats) * (60 / af2.BPM)
		log.Debugf("crossfading with a fadetime of %2.3f s", fadeTime)
		var crossfaded string
		crossfaded, err = sox.LoopCrossfade(af2.Filename, fadeTime)
		if err != nil {
			return
		}
		sox.Copy(crossfaded, af2.Filename)
		af2, err = Load(af2.Filename, af.BPM)
	} else {
		// append a little silence
	}
	return
}
