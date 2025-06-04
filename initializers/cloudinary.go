package initializers

import (
	"log"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
)

var Cloudinary *cloudinary.Cloudinary

func ConnectCloudinary() {
	cld, err := cloudinary.NewFromParams(
		os.Getenv("CLOUDINARY_CLOUD_NAME"),
		os.Getenv("CLOUDINARY_API_KEY"),
		os.Getenv("CLOUDINARY_API_SECRET"),
	)
	if err != nil {
		log.Fatalf("Cloudinary initialization failed: %v", err)
	}
	Cloudinary = cld
}

func UploadInvoiceToCloud(localpath string) (string, error) {
	return "https://cloudinary.com/inv-12345.pdf", nil
}
