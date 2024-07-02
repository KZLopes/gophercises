package primitive

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
)

type Mode uint8

const (
	ModeCombo Mode = iota
	ModeTriangle
	ModeRectangle
	ModeEllipse
	ModeCircle
	ModeRotRectangle
	ModeBeziers
	ModeRotEllipse
	ModePolygon
)

type primitiveOpt func() ([]string, error)

func WithMode(mode Mode) primitiveOpt {
	return func() ([]string, error) {
		return []string{"-m", fmt.Sprintf("%d", mode)}, nil
	}
}

func WithBgColor(hexCode string) primitiveOpt {
	return func() ([]string, error) {
		// TODO: Validate hexCode, return nil
		return []string{"-bg", hexCode}, nil
	}
}

func WithColorAlpha(value int) primitiveOpt {
	return func() ([]string, error) {
		return []string{"-a", fmt.Sprintf("%d", value)}, nil

	}
}

func Buildcommand(inputFile, outputFile string, numShapes int, opts ...primitiveOpt) *exec.Cmd {
	args := strings.Fields(fmt.Sprintf("-i %s -o %s -n %d", inputFile, outputFile, numShapes))
	for _, opt := range opts {
		arg, err := opt()
		if err != nil {
			log.Println("error:", err.Error())
			continue
		}
		args = append(args, arg...)
	}
	cmd := exec.Command("primitive", args...)
	return cmd

}

func Transform(img io.Reader, ext string) (io.Reader, error) {
	in, err := os.CreateTemp("", fmt.Sprintf("*_in%s", ext))
	if err != nil {
		return nil, err
	}
	defer os.Remove(in.Name())

	out, err := os.CreateTemp("", fmt.Sprintf("*_in%s", ext))
	if err != nil {
		return nil, err
	}
	defer os.Remove(out.Name())

	// Read img into temp input file
	_, err = io.Copy(in, img)
	if err != nil {
		return nil, err
	}
	// Run Primitive
	cmd := Buildcommand(in.Name(), out.Name(), 100)
	_, err = cmd.CombinedOutput()
	if err != nil {
		log.Println("error:", err.Error())
		return nil, err
	}

	// Read temp output into a Reader that will return
	b := bytes.NewBuffer(nil)
	_, err = io.Copy(b, out)
	if err != nil {
		return nil, err
	}

	return b, nil
}
