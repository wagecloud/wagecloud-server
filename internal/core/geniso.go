package core

import (
	"io"
	"log"

	"github.com/kdomanski/iso9660"
)

type GenISOParams struct {
	Userdata  io.Reader
	Metadata  io.Reader
	ResultIso io.Writer
}

// GenISO generate cloud init iso based on user data and metadata
func GenISO(params GenISOParams) {
	writer, err := iso9660.NewWriter()
	if err != nil {
		log.Fatalf("failed to create writer: %s", err)
	}
	defer writer.Cleanup()

	err = writer.AddFile(params.Userdata, "user-data")
	if err != nil {
		log.Fatalf("failed to add user-data: %s", err)
	}

	err = writer.AddFile(params.Metadata, "meta-data")
	if err != nil {
		log.Fatalf("failed to add meta-data: %s", err)
	}

	err = writer.WriteTo(params.ResultIso, "cidata")
	if err != nil {
		log.Fatalf("failed to write ISO image: %s", err)
	}

	log.Printf("ISO image created successfully")
}
