package main

import (
	"errors"
	"image"
	"net/url"
	"os"
	"sync"
	"threehead/blender"
	"threehead/mojang"
	"time"

	"image/png"
)

type Request struct {
	Data      string
	Params    string
	WaitGroup *sync.WaitGroup
	Skin      string
	Head      string
	Error     error
}

var order = []*Request{}

func RunOrderThread() {
	go func() {
		for {
			if len(order) == 0 {
				time.Sleep(time.Second)
				continue
			}
			request := order[0]
			if len(order) > 1 {
				order = order[1:]
			} else {
				order = []*Request{}
			}

			request.Error = request.Execute()

			request.WaitGroup.Done()
		}
	}()
}

var ErrInvalidSource = errors.New("invalid source")
var ErrRenderError = errors.New("render error")

func (r *Request) Execute() error {
	skinUrl := ""
	_, err := url.ParseRequestURI(r.Data)
	if err != nil {
		uuid, err := mojang.GetUUID(r.Data)
		if err != nil {
			return ErrInvalidSource
		}
		textures, err := mojang.GetTexture(uuid)
		if err != nil || textures.Textures.Skin == nil {
			return ErrInvalidSource
		}
		skinUrl = textures.Textures.Skin.URL
	} else {
		skinUrl = r.Data
	}
	skin, err := DownloadFile("skins/", skinUrl)
	if err != nil {
		return ErrInvalidSource
	}
	err = TryFormat("skins/" + skin + ".png")
	if err != nil {
		os.Remove("skins/" + skin + ".png") // Само очистка
		return err
	}
	_, err = os.Stat("heads/" + skin + "-" + r.Params + ".png")
	if err == nil {
		r.Head = "heads/" + skin + "-" + r.Params + ".png"
		r.Skin = "skins/" + skin + ".png"
		return nil
	}
	blendfile := blender.DEFAULTBLENDFILE
	if r.Params == "voxel" {
		blendfile = blender.VOXELSBLENDFILE
	}
	if blender.Run("skins/"+skin+".png", "heads/"+skin+"-"+r.Params+".png", blendfile) {
		r.Head = "heads/" + skin + "-" + r.Params + ".png"
		r.Skin = "skins/" + skin + ".png"
		return nil
	}
	return ErrRenderError
}

var ErrInvalidFormat = errors.New("invalid format")

func TryFormat(input string) error {
	f, err := os.Open(input)
	if err != nil {
		return err
	}
	img, _, err := image.Decode(f)
	if err != nil {
		return ErrInvalidFormat
	}
	f.Close()
	size := img.Bounds().Max
	if size.X == size.Y {
		return nil
	}
	resized := image.NewRGBA(image.Rect(0, 0, size.X, size.X))
	for y := 0; y < size.Y; y++ {
		for x := 0; x < size.X; x++ {
			resized.Set(x, y, img.At(x, y))
		}
	}
	f, err = os.OpenFile(input, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	defer f.Close()
	png.Encode(f, resized)
	return nil
}

func AddRequest(data string, params string) (string, string, error) {
	request := &Request{
		Data:      data,
		Params:    params,
		WaitGroup: &sync.WaitGroup{},
	}
	request.WaitGroup.Add(1)
	order = append(order, request)
	request.WaitGroup.Wait()
	return request.Skin, request.Head, request.Error
}
