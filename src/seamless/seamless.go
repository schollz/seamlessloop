package seamless

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"

	log "github.com/schollz/logger"
	"github.com/schollz/seamlessloop/src/sox"
)

type AudioFile struct {
	Filename   string
	Duration   float64
	SampleRate int
	Channels   int
	Samples    int64
	BPM        float64
	Beats      float64
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

	// get samples
	af.Samples, err = sox.Samples(af.Filename)
	if err != nil {
		return
	}

	// get channels and sample rate
	af.SampleRate, af.Channels, err = sox.Info(af.Filename)

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
	fmt.Printf("%+v\n", af2)
	beatNum := 16.0
	if af2.Beats < 32 {
		beatNum = 8
	}
	if af2.Beats < 8 {
		beatNum = 4
	}
	targetBeats := math.Round(af2.Beats/beatNum) * beatNum
	beats := targetBeats - af2.Beats
	beatSamples := int64(math.Round(targetBeats * 60 / af2.BPM * float64(af2.SampleRate)))
	log.Debugf("beatnum: %2.1f", beatNum)
	log.Debugf("target beats: %f", targetBeats)
	log.Debugf("target samples: %d", beatSamples)
	log.Debugf("leftover beats: %f", beats)

	if beats < 0 {
		// crossfade
		if beats < -1.0 {
			beats = -1.0
			var fname2 string
			fname2, err = sox.TrimSamples(af2.Filename, 0, int64((targetBeats-1)*60/af2.BPM)*af2.Samples)
			if err != nil {
				return
			}
			sox.Copy(fname2, af2.Filename)
			af2, err = Load(af2.Filename)
			if err != nil {
				return
			}
		}
		fadeTime := math.Abs(beats) * (60 / af2.BPM)
		log.Debugf("crossfading %2.3f beats with a fadetime of %2.3f s", beats, fadeTime)
		log.Debugf("%d samples -> %d samples", af2.Samples, beatSamples)
		var crossfaded string
		crossfaded, err = sox.LoopCrossfadeSamples(af2.Filename, af2.Samples-beatSamples)
		if err != nil {
			return
		}
		err = sox.Copy(crossfaded, af2.Filename)
	} else {
		// append silence
		secondsAddSilence := float64(beatSamples-af2.Samples) / float64(af2.SampleRate)
		log.Debugf("adding %2.6f seconds of silence", secondsAddSilence)
		var silenceAdded string
		silenceAdded, err = sox.SilenceAppend(af2.Filename, secondsAddSilence)
		if err != nil {
			return
		}
		err = sox.Copy(silenceAdded, af2.Filename)
	}
	if err != nil {
		return
	}
	af2, err = Load(af2.Filename, af.BPM)
	return
}
