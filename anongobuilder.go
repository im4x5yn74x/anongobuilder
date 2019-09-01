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

var infile, osval, archval, outfile, randfilepath string
var src = rand.NewSource(time.Now().UTC().UnixNano())

var installvar bool

const (
	goos        = "GOOS"
	goarch      = "GOARCH"
	charBytes   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	charIdxBits = 6                  // 6 bits to represent a letter index
	charIdxMask = 1<<charIdxBits - 1 // All 1-bits, as many as letterIdxBits
	charIdxMax  = 63 / charIdxBits   // # of letter indices fitting in 63 bits
)

func randstrgen(n int) string { // NOW with more R@ND0M! Credit goes to Go Playground found here: https://play.golang.org/p/KcuJ_2c_NDj
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

func randInt(min, max int) int {
	return min + rand.Intn(max-min)
}

func checkInSlice(str string, list []string) bool {
	for _, v := range list {
		if v == str {
			return true
		}
	}
	return false
}

func anongobuild() {
	ofile := outfile
	randname := randstrgen(randInt(rand.Intn(len(charBytes)), randInt(43, 252)))
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
	if installvar {
		gobin := os.Getenv("GOBIN")
		if gobin == "" {
			fmt.Println("Installation failed.")
			fmt.Println("GOBIN has null value within your environment.")
			fmt.Println("Please set a filepath value:")
			fmt.Println("Windows: set GOBIN=$GOPATH/bin")
			fmt.Println("Linux: export GOBIN=$GOPATH/bin")
			os.Exit(1)
		}
		goinstall := "go"
		args := []string{"install", "-ldflags=-s -w", randfilepath}
		b := exec.Command(goinstall, args...)
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
		if osval == "windows" {
			os.Rename(gobin+"/"+randname+".exe", gobin+"/"+ofile+".exe")
		} else {
			os.Rename(gobin+"/"+randname, gobin+"/"+ofile)
		}
	} else {
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
	flag.BoolVar(&installvar, "install", false, "Install your compiled binary.")
	flag.Parse()
	if installvar == true {
		fmt.Println("Installing...")
	}
	anongobuild()
}

func main() {
	rand.Seed(time.Now().UnixNano())
	if len(os.Args) < 2 {
		fmt.Println("Usage: ./anongobuilder -h to show the help menu.")
		os.Exit(1)
	} else {
		cli()
	}
}
