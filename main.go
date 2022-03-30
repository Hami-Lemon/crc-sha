package main

import (
	"bufio"
	"crypto"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strings"
	"sync"

	_ "crypto/md5"
	_ "crypto/sha1"
	_ "crypto/sha256"
	_ "crypto/sha512"
	_ "golang.org/x/crypto/blake2b"
	_ "golang.org/x/crypto/blake2s"
	_ "golang.org/x/crypto/md4"
	_ "golang.org/x/crypto/ripemd160"
	_ "golang.org/x/crypto/sha3"
)

var (
	compare string
	sha     string
	//校验的文件或文件夹，如果是文件夹，校验其中的所有文件,默认处理当前文件夹
	target string
	//支持的算法
	support = map[string]crypto.Hash{
		"MD4":         crypto.MD4,
		"MD5":         crypto.MD5,
		"SHA1":        crypto.SHA1,
		"SHA224":      crypto.SHA224,
		"SHA256":      crypto.SHA256,
		"SHA384":      crypto.SHA384,
		"SHA512":      crypto.SHA512,
		"RIPEMD160":   crypto.RIPEMD160,
		"SHA3_224":    crypto.SHA3_224,
		"SHA3_256":    crypto.SHA3_256,
		"SHA3_384":    crypto.SHA3_384,
		"SHA3_512":    crypto.SHA3_512,
		"SHA512_224":  crypto.SHA512_224,
		"SHA512_256":  crypto.SHA512_256,
		"BLAKE2s_256": crypto.BLAKE2s_256,
		"BLAKE2b_256": crypto.BLAKE2b_256,
		"BLAKE2b_384": crypto.BLAKE2b_384,
		"BLAKE2b_512": crypto.BLAKE2b_512,
	}
)

func pareFlag() {
	flag.StringVar(&sha, "s", "sha256", "指定hash算法。")
	flag.StringVar(&compare, "c", "", "该值不为空时，会将校验码和该值进行比较。")
	flag.Usage = func() {
		fmt.Println("用法：sha [-s hash] [-c value] file")
		flag.PrintDefaults()
		fmt.Println("支持的hash算法：")
		s := make([]string, 0, len(support))
		for k, _ := range support {
			s = append(s, k)
		}
		sort.Slice(s, func(i, j int) bool {
			return strings.Compare(s[i], s[j]) < 0
		})
		for _, a := range s {
			fmt.Printf("  %s\n", a)
		}
		fmt.Println("")
	}
	flag.Parse()
	sha = strings.ToUpper(sha)
	_, ok := support[sha]
	if !ok {
		log.Printf("不支持：%s\n", sha)
		flag.Usage()
		return
	}
	more := flag.Args()
	if len(more) == 0 {
		var err error
		if target, err = os.Getwd(); err != nil {
			log.Fatalln(err)
		}
	} else {
		target = more[0]
	}
}

func main() {
	pareFlag()
	cyc(target)
}

func output(s, n, f string) {
	fmt.Printf("%s: %s\t%s\n", s, n, f)
}

func cyc(path string) {
	state, err := os.Stat(path)
	if err != nil {
		log.Fatalln(err)
	}
	if state.IsDir() {
		cycDir(path)
		return
	}

	file, _ := os.Open(path)
	fr := bufio.NewReader(file)
	num := sum(fr)
	output(sha, num, path)
	if compare != "" {
		equal := strings.Compare(num, compare)
		fmt.Printf("equal: %v\n", equal == 0)
	}
}

func cycDir(dir string) {
	d, err := os.Open(dir)
	if err != nil {
		log.Fatalln(err)
	}

	files, err := d.ReadDir(-1)
	if err != nil {
		log.Fatalln(err)
	}
	wait := sync.WaitGroup{}
	//改变工作目录
	_ = os.Chdir(dir)
	for _, file := range files {
		if !file.IsDir() {
			//这里获取名称只有文件名，不包含路径
			f := file.Name()
			wait.Add(1)
			go func() {
				fp, _ := os.Open(f)
				defer fp.Close()
				fr := bufio.NewReader(fp)
				num := sum(fr)
				output(sha, num, f)
				wait.Done()
			}()
		}
	}
	wait.Wait()
}

func sum(r io.Reader) string {
	h := support[sha].New()

	_, err := io.Copy(h, r)
	if err != nil {
		println(err)
		log.Fatalln(err)
	}
	return hex(h.Sum(nil), false)
}

func hex(binary []byte, upper bool) string {
	//码表
	tableLower := []byte{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', 'a', 'b', 'c', 'd', 'e', 'f'}
	tableUpper := []byte{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', 'A', 'B', 'C', 'D', 'E', 'F'}
	var table []byte
	if upper {
		table = tableUpper
	} else {
		table = tableLower
	}
	//长度为binary的两倍
	h := make([]byte, len(binary)<<1)
	i := 0
	for _, b := range binary {
		h[i] = table[b>>4&0x0f]
		h[i+1] = table[b&0x0f]
		i += 2
	}
	return string(h)
}
