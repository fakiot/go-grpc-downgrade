package main

import (
	"flag"

	"github.com/fakiot/go-grpc-downgrade/version7/to6"
)

// main main function to downgrade
func main() {
	pbfile := flag.String("pb", "xyz.pb.go", "file xyz.pb.go")
	pbgrpcfile := flag.String("grpc", "xyz_grpc.pb.go", "file xyz_grpc.pb.go")
	outfile := flag.String("o", "xyz.pb.go", "output file xyz.pb.go")
	verfile := flag.String("ver", "7to6", "used versions: [7to6]")
	flag.Parse()

	switch *verfile {
	case "7to6":
		if err := to6.Downgrade(*pbfile, *pbgrpcfile, *outfile); err != nil {
			println("failed to downgrade: " + err.Error())
		}
	default:
		println("Downgrade versions not found")
	}
}
