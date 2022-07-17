# seamlessloop

this project can take a piece of audio and convert it to a *seamless loop* and can optionally be quantized to the nearest beat. "*seamless*" means that zero clicks/pops from discontinuities should occur with this transformation.


the process is shown here using audacity:

![seamless loop](https://user-images.githubusercontent.com/6550035/177219531-2efca0a8-07c7-4055-8fd0-b9b66060799a.gif)

in the case that the audio is too short (say 31 beats) then the audio is rounded to 32 beats and just appended with 1 beat of silence and the begining/end are given 5 ms fades to prevent clips.

## install

the easiest way to install is to [download the latest release](https://github.com/schollz/seamlessloop/releases/latest).

first install Go, then `go install github.com/schollz/seamlessloop@latest`.

also you need to install Sox: https://sourceforge.net/projects/sox/

## usage

you can specify a folder of files and specify a folder to output the resulting loops.

```
seamlessloop -in-folder INPUTFOLDER -out-folder OUTPUTFOLDER
```

the `OUTPUTFOLDER` will be created if it does not exist.

for example, you can make seamless quantized loops out of the files in this repo:

```
$ seamlessloop --in-folder src --out-folder quantized
wrote 'quantized/136/amenbreak_bpm136_beats8.wav'
wrote 'quantized/174/loop1_bpm174_beats16.wav'
wrote 'quantized/120/pad_bpm120_beats64.wav'
```


## quantized loops

quantizing only works if you include `bpmX` in the filename of the sample being quantized. for example, this program will not quantize a file named `sample.wav` but *will work* if the filename is `sample_bpm120.wav` or `bpm138_blahblah.wav`, etc. the `bpmX` has to be in the filename for this program to work.

quantized seamless loops are made by first figuring out the closest 4/8/16-multiple beat and then implementing a crossfade between the extra end with the beginning. for example, if you have a 35-beat piece of audio it will crop it to 32 beats. the continuous piece of audio is made by taking the X+1 beat (in example, the 33rd beat) and fading it out, and then mixing it into the beginning which has been faded in.




## thanks

thanks to Frederik Olofsson for [the crossfading graphic / explaination](https://fredrikolofsson.com/f0blog/buffer-xfader/).