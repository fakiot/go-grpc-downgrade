package to6

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"os"
	"strings"
)

const importStr = `	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
`

const grpcStr1 = `
// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6
`

const grpcStr2 = "// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream."

// states state machine to walk through each line pb file
type states int32

const (
	stateInitial = iota
	stateImportStart
	stateImportEnd
	stateNormal
	stateInitial2
	statePart1
	statePart2start
	statePart2end
	statePart3
	statePart4start
	statePart4end
	statePart5start
	statePart5end
	statePart6
)

// Downgrade downgrades pb.go generated files from version 7 to 6
func Downgrade(pbfile, pbgfile, outfile string) error {
	var bpb []byte
	var st states = stateInitial

	b, err := ioutil.ReadFile(pbfile)
	if err != nil {
		return err
	}
	b2, err := ioutil.ReadFile(pbgfile)
	if err != nil {
		return err
	}
	scan := bufio.NewScanner(bytes.NewReader(b))
	scan.Split(bufio.ScanLines)
	scan2 := bufio.NewScanner(bytes.NewReader(b2))
	scan2.Split(bufio.ScanLines)

	b2 = b
	bufpb := bytes.NewBuffer(bpb)
	for scan.Scan() {
		switch st {
		case stateInitial:
			if stLookForImport(scan, bufpb) {
				st = stateImportStart
			}

		case stateImportStart:
			addLine(bufpb, scan.Bytes())
			if strings.HasPrefix(scan.Text(), ")") {
				st = stateNormal
			}
		case stateNormal:
			addLine(bufpb, scan.Bytes())
		}
	}
	st = stateInitial2
	bufpb.Write([]byte(grpcStr1))
	for scan2.Scan() {
		switch st {
		case stateInitial2:
			if bytes.Contains(scan2.Bytes(), []byte("SupportPackageIsVersion")) {
				st = statePart1
			}
		case statePart1:
			if stPart1(scan2, bufpb) {
				st = statePart2start
			}
		case statePart2start:
			if stPart2(scan2, bufpb) {
				st = statePart2end
			}
		case statePart2end:
			if scan2.Text() == "// for forward compatibility" {
				st = statePart3
			}
		case statePart3:
			if stPart3(scan2, bufpb) {
				st = statePart4start
			}
		case statePart4start:
			if stPart4start(scan2, bufpb) {
				st = statePart4end
			}
		case statePart4end:
			if stPart4end(scan2, bufpb) {
				st = statePart5start
			}

		case statePart5start:
			if stPart5start(scan2, bufpb) {
				st = statePart5end
			}

		case statePart5end:
			if stPart5end(scan2, bufpb) {
				st = statePart6
			}

		case statePart6:
			addLine(bufpb, scan2.Bytes())

		}
	}
	os.Remove(pbgfile)
	ioutil.WriteFile(outfile, bufpb.Bytes(), 0644)
	return nil
}

func stPart5end(scan *bufio.Scanner, buf *bytes.Buffer) (found bool) {
	if bytes.Contains(scan.Bytes(), []byte("ServiceDesc")) {
		b := bytes.Replace(scan.Bytes(), []byte("var "), []byte("var _"), 1)
		b = bytes.Replace(b, []byte("_ServiceDesc"), []byte("_serviceDesc"), 1)
		addLine(buf, b)
		return true
	}
	return false
}

func stPart5start(scan *bufio.Scanner, buf *bytes.Buffer) (found bool) {
	if bytes.Contains(scan.Bytes(), []byte("is the grpc.ServiceDesc for")) {
		return true
	}
	addLine(buf, scan.Bytes())
	return false
}

func stPart4end(scan *bufio.Scanner, buf *bytes.Buffer) (found bool) {
	if bytes.HasPrefix(scan.Bytes(), []byte("func Register")) {
		b := bytes.Replace(scan.Bytes(), []byte("grpc.ServiceRegistrar"), []byte("*grpc.Server"), 1)
		addLine(buf, b)
		return false
	} else if bytes.Contains(scan.Bytes(), []byte("s.RegisterService(&")) {
		b := bytes.Replace(scan.Bytes(), []byte("s.RegisterService(&"), []byte("s.RegisterService(&_"), 1)
		b = bytes.Replace(b, []byte("ServiceDesc"), []byte("serviceDesc"), 1)
		addLine(buf, b)
		return true
	}

	return false
}

func stPart4start(scan *bufio.Scanner, buf *bytes.Buffer) (found bool) {
	if bytes.HasPrefix(scan.Bytes(), []byte("func (")) {
		b := bytes.Replace(scan.Bytes(), []byte("func ("), []byte("func (*"), 1)
		addLine(buf, b)
		return false
	} else if bytes.HasPrefix(scan.Bytes(), []byte("// Unsafe")) {
		return true
	}
	addLine(buf, scan.Bytes())
	return false
}

func stPart3(scan *bufio.Scanner, buf *bytes.Buffer) (found bool) {
	if bytes.HasPrefix(scan.Bytes(), []byte("// Unimplemented")) {
		b := bytes.Replace(scan.Bytes(), []byte("should"), []byte("can"), 1)
		addLine(buf, b)
		return true
	}
	addLine(buf, scan.Bytes())
	return false
}

func stPart2(scan *bufio.Scanner, buf *bytes.Buffer) (found bool) {
	if bytes.HasPrefix(scan.Bytes(), []byte("// All implementations should embed")) {
		return true
	}
	addLine(buf, scan.Bytes())
	return false
}

func stPart1(scan *bufio.Scanner, buf *bytes.Buffer) (found bool) {
	if bytes.HasPrefix(scan.Bytes(), []byte("// For semantics around ctx use")) {
		addLine(buf, []byte(grpcStr2))
		return true
	}
	addLine(buf, scan.Bytes())
	return false
}

func stLookForImport(scan *bufio.Scanner, buf *bytes.Buffer) (found bool) {
	t := scan.Bytes()
	if bytes.HasPrefix(t, []byte("import")) {
		addLine(buf, t)
		buf.Write([]byte(importStr))
		return true
	}
	addLine(buf, t)
	return false
}

func addLine(buf *bytes.Buffer, t []byte) {
	buf.Write(t)
	buf.Write([]byte("\n"))
}
