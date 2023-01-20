package blender

import (
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	BLENDER          = "/snap/bin/blender"
	DEFAULTBLENDFILE = "/root/threehead/head.blend"
	VOXELSBLENDFILE  = "/root/threehead/head_voxels.blend"
	PYSCRIPT         = "/root/threehead/render.py"
	DEBUG            = true
)

func Run(input string, output string, blendfile string) bool {
	if DEBUG == true {
		in, _ := os.Open(input)
		out, _ := os.OpenFile(output, os.O_CREATE|os.O_RDWR, 0755)
		defer in.Close()
		defer out.Close()
		io.Copy(out, in)
		return true
	}
	i, _ := filepath.Abs(input)
	o, _ := filepath.Abs(output)
	p := exec.Command(BLENDER, "-b", blendfile, "-P", PYSCRIPT, "--", i, o)
	log.Println(strings.Join([]string{BLENDER, "-b", blendfile, "-P", PYSCRIPT, "--", i, o}, " "))
	err := p.Run()
	return err == nil
}
