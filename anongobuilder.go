package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"runtime"
	"time"
	"unsafe"
)

var src = rand.NewSource(time.Now().UnixNano())

const (
	goos        = "GOOS"
	goarch      = "GOARCH"
	charBytes   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	charIdxBits = 6
	charIdxMask = 1<<charIdxBits - 1
	charIdxMax  = 63 / charIdxBits
)

func randstrgen(n int) string { // https://play.golang.org/p/KcuJ_2c_NDj
	a := make([]byte, n)
	for j, cache, remain := n-1, src.Int63(), charIdxMax; j >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), charIdxMax
		}
		if idx := int(cache & charIdxMask); idx < len(charBytes) {
			a[j] = charBytes[idx]
			j--
		}
		cache >>= charIdxBits
		remain--
	}
	return *(*string)(unsafe.Pointer(&a))
}

func checkInSlice(str string, list []string) bool {
	for _, v := range list {
		if v == str {
			return true
		}
	}
	return false
}

var infile, osval, archval, outfile, randfilepath string

func anongobuild() {
	ofile := outfile
	_ = randstrgen(rand.Intn(100))         // Used to generate first random value (usually a single character)
	_ = ""                                 // and set it to nothing
	randname := randstrgen(rand.Intn(100)) // Feel free to adjust the 'Intn()' value accordingly
	if runtime.GOOS == "windows" {
		randfilepath = "C:\\Users\\Public\\" + randname + ".go" // Implemented to mask source path within binary artifacts.
	} else {
		randfilepath = "/tmp/" + randname + ".go"
	}
	from, err := os.Open(infile)
	if err != nil {
		log.Fatal(err)
	}
	defer from.Close()
	to, err := os.OpenFile(randfilepath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer to.Close()
	_, err = io.Copy(to, from)
	if err != nil {
		log.Fatal(err)
	}
	govar := "go"
	args := []string{"build", "-ldflags=-s -w", randfilepath}
	b := exec.Command(govar, args...)
	b.Env = os.Environ()
	b.Env = append(b.Env, fmt.Sprintf("%s=%s", goos, osval))
	b.Env = append(b.Env, fmt.Sprintf("%s=%s", goarch, archval))
	out, err := b.CombinedOutput()
	if err != nil {
		fmt.Println("Could not compile")
		os.Exit(0)
	}
	fmt.Printf("%s", out)
	fmt.Println("Successfully compiled")
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if osval == "windows" {
		os.Rename(pwd+"/"+randname+".exe", ofile+".exe")
	} else {
		os.Rename(pwd+"/"+randname, ofile)
	}
	err = os.Remove(fmt.Sprintf("%s", randfilepath))
	if err != nil {
		fmt.Println("Could not remove file")
	}
}

func cli() {
	flag.StringVar(&osval, "p", "", "Operating System: windows, linux, freebsd, nacl, netbsd, openbsd, plan9, solaris, dragonfly, darwin, android")
	flag.StringVar(&archval, "a", "", "Architecture: 386, amd64, amd64p32, arm, arm64, ppc64, ppc64le, mips, mipsle, mips64, mips64le, s390x, sparc64")
	flag.StringVar(&infile, "i", "", "Input filename: <whatever file you aim to compile.>")
	flag.StringVar(&outfile, "o", "", "Output filename: <anything goes>")
	flag.Parse()
	anongobuild()
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: ./anongobuilder -h to show the help menu.")
		os.Exit(1)
	} else {
		cli()
	}
}