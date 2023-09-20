package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.formulabun.club/extractor/maps"
	"go.formulabun.club/metadatadb"
	"go.formulabun.club/storage"

	"go.formulabun.club/srb2kart/addons"
)

const parallel = 50

var extractors []Extractor = []Extractor{
	maps.MapExtractor{},
}

func extract(files *[]metadatadb.File) {
	ctx := context.Background()
	connectCtx, _ := context.WithTimeout(ctx, time.Second*5)
	c, err := metadatadb.NewClient(connectCtx)
	if err != nil {
		log.Println(err)
	}

	s := make(chan struct{}, parallel)
	for _, f := range *files {
		file := f
		s <- struct{}{}
		go func() {
			ctx, cancel := context.WithTimeout(ctx, time.Second*2)
			defer cancel()
			ExtractFile(file, c, ctx)
			<-s
		}()
	}
}

func ExtractFile(file metadatadb.File, c *metadatadb.Client, ctx context.Context) {
	log.Printf("Starting file %s\n", file.Filename)
	defer log.Printf("Stopped file %s\n", file.Filename)

	fileR, err := storage.Get(file)

	if err != nil {
		log.Printf("Could not read the file %v: %s", file.Filename, err)
		return
	}

	defer fileR.Close()

	addon, err := addons.Read(fileR)
	for _, e := range extractors {
		if err != nil {
			continue
		}
		err = e.Extract(file.Filename, addon, c, ctx)
		if err != nil {
			fmt.Println(err)
		}
	}
}
