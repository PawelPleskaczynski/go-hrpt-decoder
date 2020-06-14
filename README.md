# go-hrpt-decoder
WIP software for decoding HRPT images from demodulated data. It's been tested on HRPT images obtained with GNU Radio from NOAA satellites.

### Features
- processing individual channels
- combining channels to color RGB images
- automatic histogram stretching
- option to process northbound/southbound passes

### Usage
```
Usage of go-hrpt-decoder:
  -input string
        input file
  -channel int
        channel of the image, -1 returns all the channels (default -1)
  -r int
        image channel to be represented by red channel (for NOAA images ranges from 0-4) (default -1)
  -b int
        image channel to be represented by blue channel (for NOAA images ranges from 0-4) (default -1)
  -g int
        image channel to be represented by green channel (for NOAA images ranges from 0-4) (default -1)
  -south
        enable this option for southbound passes
  -stretch
        stretch histogram
```

### Useful links
- https://www.sat.dundee.ac.uk/hrptformat.html - HRPT format documentation
- https://tynet.eu/hrpt/hrpt-decoder - GNU Radio-based HRPT demodulator
- http://www.alblas.demon.nl/wsat/satinfo.html - CHRPT format documentation
- http://www.sat.cc.ua/page5.html - MetFY3x software for unpacking HRPT frames from MetOP, FengYun and Meteor satellites

### To-do list
- test on files from different demodulators
- add support for FengYun, Meteor and MetOP HRPT images