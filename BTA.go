package main

import (
	"image"
	"image/color"

	"golang.org/x/exp/shiny/driver"
	"golang.org/x/exp/shiny/screen"
	"golang.org/x/mobile/event/lifecycle"
)

func load_config() {

}

func main() {
	driver.Main(func(s screen.Screen) {
		windowopts := screen.NewWindowOptions{Title: "BTA"}
		bgcolor := color.RGBA{30, 30, 30, 1}
		bgrectangle := image.Rectangle{image.Point{0, 0}, image.Point{1280, 720}}
		windowcolour := color.Color(bgcolor)
		w, err := s.NewWindow(&windowopts)
		if err != nil {

			return
		}
		defer w.Release()
		w.Fill(bgrectangle, windowcolour, 0)
		w.Publish()
		for {
			switch e := w.NextEvent().(type) {
			case lifecycle.Event:
				if e.To == lifecycle.StageDead {
					return
				}
			}
		}
	})

}
