package main

import (
	"context"

	"go.formulabun.club/metadatadb"

	"go.formulabun.club/srb2kart/addons"
)

type Extractor interface {
	Extract(fileName string, addon addons.Addon, databaseClient *metadatadb.Client, ctx context.Context) error
}
