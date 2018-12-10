package main

import (
	"bufio"
	"container/list"
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"log"
	"os"
	"regexp"
	"strconv"
)

type point struct {
	x         int
	y         int
	velocityX int
	velocityY int
}

func main() {
	//file, err := os.Open("day10/input.txt")
	file, err := os.Open("day10/test.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	ls := list.New()
	readAll(file, ls)

	// position=< 9,  1> velocity=< 0,  2>
	regex, err := regexp.Compile("position=<[\\s\\t]*(-?\\d+),[\\s\\t]*(-?\\d+)> velocity=<[\\s\\t]*(-?\\d+),[\\s\\t]*(-?\\d+)>")
	if err != nil {
		log.Fatal(err)
	}

	field := make([]*point, 0)
	maxX := 0
	maxY := 0
	minX := 0
	minY := 0
	for line := ls.Front(); line != nil; line = line.Next() {
		matches := regex.FindAllStringSubmatch(line.Value.(string), -1)[0]
		x, err := strconv.Atoi(matches[1])
		if err != nil {
			log.Fatal(err)
		}
		if x > maxX {
			maxX = x
		}
		if x < minX {
			minX = x
		}
		y, err := strconv.Atoi(matches[2])
		if err != nil {
			log.Fatal(err)
		}
		if y > maxY {
			maxY = y
		}
		if y < minY {
			minY = y
		}
		velocityX, err := strconv.Atoi(matches[3])
		if err != nil {

			log.Fatal(err)
		}
		velocityY, err := strconv.Atoi(matches[4])
		if err != nil {
			log.Fatal(err)
		}
		//fmt.Printf("Found x=%d, y=%d, velocityX=%d, velocityY=%d\n", x, y, velocityX, velocityY)
		field = append(field, &point{x: x, y: y, velocityY: velocityY, velocityX: velocityX})
	}
	fmt.Printf("Making size: minX=%v maxX=%v, minY=%v, maxY=%v\n", minX, maxX, minY, maxY)
	iter := 0

	var images []*image.Paletted
	var delays = []int{0}

	var palette = []color.Color{
		color.RGBA{0x00, 0x00, 0x00, 0xff},
		color.RGBA{0x00, 0x00, 0xff, 0xff},
		color.RGBA{0x00, 0xff, 0x00, 0xff},
		color.RGBA{0x00, 0xff, 0xff, 0xff},
		color.RGBA{0xff, 0x00, 0x00, 0xff},
		color.RGBA{0xff, 0x00, 0xff, 0xff},
		color.RGBA{0xff, 0xff, 0x00, 0xff},
		color.RGBA{0xff, 0xff, 0xff, 0xff},
	}

	for {
		iter++
		fmt.Println("Iteration ", iter)
		newMaxX := 0
		newMinX := 0
		newMaxY := 0
		newMinY := 0
		for _, point := range field {
			if point.x > newMaxX {
				newMaxX = point.x
			}
			if point.x < newMinX {
				newMinX = point.x
			}
			if point.y > newMaxY {
				newMaxY = point.y
			}
			if point.y < newMinY {
				newMinY = point.y
			}
		}

		if maxX-minX < newMaxX-newMinX || maxY-minY < newMaxY-newMinY {
			fmt.Println("Message shown at ", iter-2) // FIXME -2 is just a guess!
			break
		} else {
			img := image.NewPaletted(image.Rect(0, 0, int(scaledMaxX), int(scaledMaxY)), palette)
			images = []*image.Paletted{img}

			for _, point := range field {
				x := scaledMaxX * float64(point.x-newMinX) / float64(newMaxX-newMinX) / 2
				y := scaledMaxY * float64(point.y-newMinY) / float64(newMaxY-newMinY) / 2
				img.Set(int(x), int(y), color.RGBA{
					R: 0xff,
					G: 0xff,
					B: 0xff,
					A: 255,
				})
			}
			maxX = newMaxX
			maxY = newMaxY
			minX = newMinX
			minY = newMinY
		}

		for _, point := range field {
			//fmt.Printf("x=%v, y=%v, linear=%v\n", point.x, point.y, linear(point.x, minX, point.y, minY, maxX))
			point.x = point.x + point.velocityX
			point.y = point.y + point.velocityY
		}
	}
	// save to rgb.gif
	f, _ := os.OpenFile("/tmp/rgb3.gif", os.O_WRONLY|os.O_CREATE, 0600)
	defer f.Close()
	gif.EncodeAll(f, &gif.GIF{
		Image: images,
		Delay: delays,
	})
	//fmt.Printf("Size of the branch: %d:%d\n", maxX+1, maxY+1)
	//fmt.Printf("Field: %+v", field)
}

const scaledMaxX float64 = 540
const scaledMaxY float64 = 540

func linear(x int, minX int, y int, minY int, maxX int, maxY int) int {
	scaledX := scaledMaxX * float64(x-minX) / float64(maxX-minX)
	scaledY := scaledMaxY * float64(y-minY) / float64(maxY-minY)
	result := (int)(scaledX + scaledY*scaledMaxX)
	return result
}

func readAll(file *os.File, list *list.List) {
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		val := scanner.Text()
		list.PushBack(val)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
