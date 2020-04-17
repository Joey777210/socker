package container

import (
	"Socker/overlay2"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"
)

type image struct {
	Name string
	ModTime string
	Size string
}

func ImageLs() error {
	files, err := ioutil.ReadDir(overlay2.ROOT)
	if err != nil {
		log.Errorf("Open dir %s error %v", overlay2.ROOT, err)
		return err
	}
	//get all image
	var images []image
	for _, f := range files {
		strs := strings.Split(f.Name(), ".")
		if len(strs) == 2 && strs[1] == "tar"{
			image := getImage(f)
			images = append(images, image)
		}
	}

	w := tabwriter.NewWriter(os.Stdout, 12, 1, 3, ' ', 0)
	fmt.Fprint(w, "Name\tModTime\tSize\n")
	for _, item := range images {
		fmt.Fprintf(w,"%s\t%s\t%s\n",
			item.Name,
			item.ModTime,
			item.Size)
	}

	//flush stdout stream buffer, print the container list
	if err := w.Flush(); err != nil {
		log.Errorf("Flush error %v", err)
		return err
	}
	return nil
}

func getImage(f os.FileInfo) image {
	var image image
	names := strings.Split(f.Name(), ".")
	modTime := f.ModTime().String()
	times := strings.Split(modTime, " ")
	hms := strings.Split(times[1], ".")

	image.Name = names[0]
	image.ModTime = times[0] + " " + hms[0]
	image.Size = strconv.Itoa(int(f.Size()/1024/1024)) + "MB"
	return image
}