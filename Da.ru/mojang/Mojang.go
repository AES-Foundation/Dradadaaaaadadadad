package mojang

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"net/http"
)

const (
	USERNAME2UUID = "https://api.mojang.com/users/profiles/minecraft/"
	UUID2SKIN     = "https://sessionserver.mojang.com/session/minecraft/profile/"
)

type Profile struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

type MinecraftProfile struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Properties []struct {
		Name      string `json:"name"`
		Value     string `json:"value"`
		Signature string `json:"signature"`
	} `json:"properties"`
}

func GetUUID(username string) (string, error) {
	r, err := http.Get(USERNAME2UUID + username)
	if err != nil {
		return "", err
	}
	defer r.Body.Close()
	profile := Profile{}
	err = json.NewDecoder(r.Body).Decode(&profile)
	if err != nil {
		return "", err
	}
	return profile.ID, nil
}

func GetProfile(uuid string) (*MinecraftProfile, error) {
	r, err := http.Get(UUID2SKIN + uuid)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	profile := MinecraftProfile{}
	err = json.NewDecoder(r.Body).Decode(&profile)
	if err != nil {
		return nil, err
	}
	return &profile, nil
}

type Textures struct {
	Timestamp   int64  `json:"timestamp"`
	ProfileID   string `json:"profileId"`
	ProfileName string `json:"profileName"`
	Textures    struct {
		Skin *struct {
			URL string `json:"url"`
		} `json:"SKIN"`
		Cape *struct {
			URL string `json:"url"`
		} `json:"CAPE"`
	} `json:"textures"`
}

func GetTexture(uuid string) (*Textures, error) {
	profile, err := GetProfile(uuid)
	if err != nil {
		return nil, err
	}
	textureBase64 := ""
	for _, v := range profile.Properties {
		if v.Name == "textures" {
			textureBase64 = v.Value
		}
	}
	textureBase64Reader := base64.NewDecoder(base64.StdEncoding, bytes.NewBufferString(textureBase64))
	textures := Textures{}
	err = json.NewDecoder(textureBase64Reader).Decode(&textures)
	if err != nil {
		return nil, err
	}
	return &textures, nil
}
