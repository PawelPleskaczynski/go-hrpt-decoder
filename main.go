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
	"os"
)

var width, height int
var filename *string

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func saveImage(channel, width, height int, data []uint16) {
	log.Printf("Processing channel %d...", channel)
	data = data[750:]

	j, count10240, count850 := 0, 0, 0
	array512 := make([]uint16, len(data))

	for i := channel * 2; i < len(data); i++ {
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

	image := image.NewGray16(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			image.SetGray16(x, y, color.Gray16{data[y*width+x]})
		}
	}
	out, err := os.Create(fmt.Sprintf("%s_channel_%d.png", (*filename), channel))
	check(err)
	defer out.Close()
	png.Encode(out, image)
	log.Printf("Saved channel %d", channel)
}

func main() {
	filename = flag.String("input", "", "input file")
	channel := flag.Int("channel", -1, "channel of the image, -1 returns all the channels")
	flag.Parse()

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

	if *channel == -1 {
		for i := 0; i < 5; i++ {
			saveImage(i, width, height, slice16bits)
		}
	} else {
		saveImage(*channel, width, height, slice16bits)
	}
}
