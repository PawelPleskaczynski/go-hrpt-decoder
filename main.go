package main

import (
	"bytes"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"log"
	"math"
	"os"
)

var width, height int
var filename *string

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func processFrame(channel, width, height int, data []uint16, stretch, southbound bool) *image.Gray16 {
	log.Printf("Processing channel %d...", channel)
	data = data[750:]

	j, count10240, count850 := 0, 0, 0
	array512 := make([]uint16, len(data))

	var maxVal uint16

	for i := channel; i < len(data); i++ {
		if count10240 < 10240 {
			array512[j] = data[i]
			j++
			count10240 += 5
			i += 4
		} else if count850 < 850-1 && count10240 == 10240 {
			count850++
		} else if count850 == 850-1 {
			count10240 = 0
			count850 = 0
		}
	}
	data = array512

	if stretch {
		for i := 0; i < len(array512); i++ {
			if array512[i] > maxVal {
				maxVal = array512[i]
			}
		}
	}

	frame := image.NewGray16(image.Rect(0, 0, width, height))

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if southbound {
				if stretch {
					frame.SetGray16(width-x-1, height-y-1, color.Gray16{uint16(float64(data[y*width+x]) / (float64(maxVal) / float64(math.MaxUint16)))})
				} else {
					frame.SetGray16(width-x-1, height-y-1, color.Gray16{data[y*width+x]})
				}
			} else {
				if stretch {
					frame.SetGray16(x, y, color.Gray16{uint16(float64(data[y*width+x]) / (float64(maxVal) / float64(math.MaxUint16)))})
				} else {
					frame.SetGray16(x, y, color.Gray16{data[y*width+x]})
				}
			}
		}
	}

	return frame
}

func saveImageRGB(frameR, frameG, frameB *image.Gray16, ch1, ch2, ch3 int) {
	out, err := os.Create(fmt.Sprintf("%s_%d%d%d.png", (*filename), ch1, ch2, ch3))
	check(err)
	defer out.Close()
	frame := image.NewRGBA64(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			frame.SetRGBA64(x, y, color.RGBA64{frameR.Gray16At(x, y).Y, frameG.Gray16At(x, y).Y, frameB.Gray16At(x, y).Y, math.MaxUint16})
		}
	}
	png.Encode(out, frame)
	log.Printf("Saved RGB %d%d%d image", ch1, ch2, ch3)
}

func saveImageGray(frame *image.Gray16, channel int) {
	out, err := os.Create(fmt.Sprintf("%s_channel_%d.png", (*filename), channel))
	check(err)
	defer out.Close()
	png.Encode(out, frame)
	log.Printf("Saved channel %d", channel)
}

func main() {
	filename = flag.String("input", "", "input file")
	channel := flag.Int("channel", -1, "channel of the image, -1 returns all the channels")
	stretch := flag.Bool("stretch", false, "stretch histogram")
	southbound := flag.Bool("south", false, "enable this option for southbound passes")
	rChannel := flag.Int("r", -1, "image channel to be represented by red channel (for NOAA images ranges from 0-4)")
	gChannel := flag.Int("g", -1, "image channel to be represented by green channel (for NOAA images ranges from 0-4)")
	bChannel := flag.Int("b", -1, "image channel to be represented by blue channel (for NOAA images ranges from 0-4)")
	flag.Parse()

	if *filename == "" {
		fmt.Fprintln(os.Stderr, "Usage of go-hrpt-decoder:")
		flag.PrintDefaults()
		os.Exit(1)
	}

	buf := bytes.NewBuffer(nil)
	file, err := os.Open(*filename)
	check(err)
	io.Copy(buf, file)
	fileslice := []byte(buf.Bytes())
	file.Close()

	width = 2048
	height = len(fileslice) / (10240 + 850) / 2

	slice16bits := make([]uint16, len(fileslice))

	j := 0

	for i := 0; i < len(fileslice); i++ {
		if i%2 == 0 {
			slice16bits[j] = (uint16(fileslice[i+1])<<8<<0 | (uint16(fileslice[i+0]) << 0))
			j++
		}
	}

	if *rChannel != -1 && *gChannel != -1 && *bChannel != -1 {
		imageR := processFrame(*rChannel, width, height, slice16bits, *stretch, *southbound)
		imageG := processFrame(*gChannel, width, height, slice16bits, *stretch, *southbound)
		imageB := processFrame(*bChannel, width, height, slice16bits, *stretch, *southbound)
		saveImageRGB(imageR, imageG, imageB, *rChannel, *gChannel, *bChannel)
	} else {
		if *channel == -1 {
			for i := 0; i < 5; i++ {
				saveImageGray(processFrame(i, width, height, slice16bits, *stretch, *southbound), i)
			}
		} else {
			saveImageGray(processFrame(*channel, width, height, slice16bits, *stretch, *southbound), *channel)
		}
	}
}
