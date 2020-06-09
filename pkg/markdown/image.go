package markdown

import (
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
)

func renderItermImage(filename string, rootDir string) (string, error) {
	fullFilePath := path.Join(rootDir, filename)
	// Skip errors and directories
	if fi, err := os.Stat(fullFilePath); err != nil || fi.IsDir() {
		return "", err
	}

	f, err := os.Open(fullFilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	str, err := display(f)

	if err != nil {
		log.Fatal(err)
	}
	return str, nil
}

func display(r io.Reader) (string, error) {
	preserveAspectRatio := true
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return "", err
	}

	width, height := widthAndHeight("", "", "")

	str := fmt.Sprintf("\033]1337;")
	str += fmt.Sprintf("File=inline=1")
	if width != "" || height != "" {
		if width != "" {
			str += fmt.Sprintf(";width=%s", width)
		}
		if height != "" {
			str += fmt.Sprintf(";height=%s", height)
		}
	}
	if preserveAspectRatio {
		str += ("preserveAspectRatio=1")
	}
	str += (":")
	str += fmt.Sprintf("%s", base64.StdEncoding.EncodeToString(data))
	str += ("\a\n")

	return str, nil
}

func widthAndHeight(width, height, size string) (w, h string) {
	if width != "" {
		w = width
	}
	if height != "" {
		h = height
	}
	if size != "" {
		sp := strings.SplitN(size, ",", -1)
		if len(sp) == 2 {
			w = sp[0]
			h = sp[1]
		}
	}
	return
}
