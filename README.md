<h2>Anon-Go-Builder</h2>

Anongobuilder was created mid flight from the South East to the Pacific North West coasts across the United States. 
The idea came about while building new features into <a href="https://github.com/im4x5yn74x/dropper2">Dropper2</a>. As explained during that project, the Golang environment compiles both its runtime libraries and dependencies into every binary produced. Along with these artifacts, Golang brings local system information and builds it into each binary compiled within an individual's Go environment; meaning if you have compiled any ".go" file within your GOPATH, the path to said Go file will appear as an artifact for every platform you compile for. With all of this said, I have created the Anon-Go-Builder.

It's a simple tool that takes some of the compiling wizardry from Dropper2 and allows an individual to compile their Golang files more anonymously. Unfortunately we cannot do much about the necessary Golang related pseudo environment brought along with each binary created, however we've done what we can for masking origin and local environmental details in regards.

With that said, the execution is straight foreward.<br>
<code>go get github.com/im4x5yn74x/anongobuilder</code><br>
Build and Install the Anongobuilder binary.<br>
<code>go build anongobuilder.go</code><br>
<code>./anongobuilder -h</code>
<pre>Usage of ./anongobuilder:
  -a string
        Architecture: 386, amd64, amd64p32, arm, arm64, ppc64, ppc64le, mips, mipsle, mips64, mips64le, s390x, sparc64
  -i string
        Input filename: <whatever file you aim to compile.>
  -o string
        Output filename: <anything goes>
  -p string
        Operating System: windows, linux, freebsd, nacl, netbsd, openbsd, plan9, solaris, dragonfly, darwin, android
</pre>
If you're feeling froggy, you can even compile Anongobuilder with itself!<br>
<code>cd $GOPATH/src/github.com/im4x5yn74x/anongobuilder</code><br>
<code>go run anongobuilder.go -a amd64 -p linux -i anongobuilder.go -o anongobuilder</code><br>
OR:<br>
<code>go run anongobuilder.go -a 386 -p windows -i anongobuilder.go -o anongobuilder</code>