package maps

import (
	"context"
	"fmt"
	"strconv"

	"go.formulabun.club/metadatadb"

	"go.formulabun.club/srb2kart/addons"
	"go.formulabun.club/srb2kart/lump/soc"
)

type MapExtractor struct {
}

func (e MapExtractor) Extract(file string, a addons.Addon, c *metadatadb.Client, ctx context.Context) error {
	socs, err := a.Socs()
	if err != nil {
		return err
	}

	for _, soc := range socs {
		for _, block := range soc {
			if block.IsLevel() {
				c.AddMap(file, blockToMapData(block), ctx)
			}
		}
	}
	return nil
}

func blockToMapData(s soc.Block) metadatadb.MapData {
	return metadatadb.MapData{
		fmt.Sprintf("%02s", s.Header.Name)[:2],
		s.Properties["LEVELNAME"],
		s.Properties["ACT"],
		s.Properties["SUBTITLE"],
		s.Properties["ZONETITLE"],
		parseBool(s, "NOZONE", false),
		metadatadb.LevelType(s.Properties["TYPEOFLEVEL"]),
		parseInt(s, "PALETTE", 0),
		parseInt(s, "SKY", 0),
		parseInt(s, "NUMLAPS", 3),
		s.Properties["MUSIC"],
	}
}

func parseInt(s soc.Block, prop string, fallBack int) int {
	prop, ok := s.Properties[prop]
	if ok {
		i, err := strconv.Atoi(prop)
		if err == nil {
			return i
		}
	}
	return fallBack
}

func parseBool(b soc.Block, prop string, fallBack bool) bool {
	prop, ok := b.Properties[prop]
	if ok {
		b, err := strconv.ParseBool(prop)
		if err == nil {
			return b
		}
	}
	return fallBack
}
