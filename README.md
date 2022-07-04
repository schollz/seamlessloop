# seamlessloop

this project can take a piece of audio with a bpm (in the filename) and convert it to a *seamless loop* quantized to the nearest 4/8/16-multiple beat.

it works by first figuring out the closest 4/8/16-multiple beat and then implementing a crossfade. for example, if you have a 35-beat piece of audio it will crop it to 32 beats. the continuous piece of audio is made by taking the X+1 beat (in example, the 33rd beat) and fading it out, and then mixing it into the beginning which has been faded in.

the process is shown here using audacity:

![seamless loop](https://user-images.githubusercontent.com/6550035/177219531-2efca0a8-07c7-4055-8fd0-b9b66060799a.gif)

in the case that the audio is too short (say 31 beats) then the audio is rounded to 32 beats and just appended with 1 beat of silence.

## install

first install Go, then `go install github.com/schollz/seamlessloop@latest`.

also you need to install Sox: https://sourceforge.net/projects/sox/

## usage

you can specify a folder of files and specify a folder to output the resulting loops.

```
./seamlessloop --in INPUTFOLDER --out OUTPUTFOLDER
```

the `OUTPUTFOLDER` will be created if it does not exist.


## important! you must have properly named audio files

all the files are assumed to have `bpmX` in their filename! this is very improtant, as this program *does not guess the BPM*. for example, this program will not work on a file named `sample.wav` but *will work* if the filename is `sample_bpm120.wav` or `bpm138_blahblah.wav`, etc. the `bpmX` has to be in the filename for this program to work.
