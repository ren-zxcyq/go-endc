// A command line utility which encrypts files provided and can decrypt files encrypted.
package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
)

var debug bool

type cmdline struct {
	srcf string
	dstf string
	pass string
	mode bool //	true for enc false for dec
}

func main() {
	// stringdemo()
	// filedemo()
	act()

}

// Function setupflags sets up parses command line options
// and populates struct cmdline
func setupflags() *cmdline {
	// get needed sys info
	cwd, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	// fmt.Println("[*]\t", cwd)  // for example /home/user

	// set up flags
	// var cmdlflags = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	var srcfpointer = flag.String("s", "", "File to encrypt. Absolute path.")
	var dstfpointer = flag.String("d", "", "Save output file to. Absolute path.")
	var passpointer = flag.String("p", "", "Password.")
	var decrypt = flag.Bool("decrypt", false, "Mode Decrypt")
	var encrypt = flag.Bool("encrypt", false, "Mode Encrypt")
	var debugpointer = flag.Bool("debug", false, "Enable debugging.")

	// parse flags
	flag.Parse()
	var c cmdline = cmdline{}
	c.srcf = *srcfpointer
	c.dstf = *dstfpointer
	c.pass = *passpointer
	if (*encrypt && *decrypt) || ((*encrypt == false) && (*decrypt == false)) {
		// fmt.Fprintf(flag.Output(), "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults() //	Print Usage and Exit
		os.Exit(0)
	} else if *encrypt {
		c.mode = true
	} else if *decrypt {
		c.mode = false
	}

	debug = *debugpointer

	if debug {
		fmt.Println("[*]\t", cwd)
	}
	return &c
}

// Function returnfilecontents returns contents of the provided ABSPATH
// in form []byte
func returnfilecontents(abspath string) []byte {
	b, err := ioutil.ReadFile(abspath) // just pass the file name
	if err != nil {
		fmt.Println("[*]\tRead File:", err)
	}

	//fmt.Println("[*]\t", b) // print the content as 'bytes'
	//str := string(b) // convert content to a 'string'
	// fmt.Println("[*]\t", str) // print the content as a 'string'

	return b //	fmt.Sprintln(str)
}

// Function savebytestofile saves provided byte array to a file.
func savebytestofile(abspath string, newcontents []byte) error {
	//	Open file
	f, err := os.Create(abspath)
	if err != nil {
		fmt.Println("[*]\t", err)
	}
	defer f.Close()

	//	Write
	n, err := f.Write(newcontents)
	if err != nil {
		fmt.Println("[*]\t", err)
	}

	//	Report
	if debug {
		fmt.Printf("[*]\tWrote %d bytes\r\n", n)
	}
	f.Sync() //	Flush
	return err
}

// Function act performs the core operation defined by Boolean "cmdline.mode"
func act() {
	cl := setupflags()

	fmt.Println("[*]\tSrc:\t", cl.srcf)
	fmt.Println("[*]\tDst:\t", cl.dstf)
	var r []byte
	b := returnfilecontents(cl.srcf)

	if cl.mode {
		fmt.Println("[*]\tEncrypting file contents")
		r = enc([]byte(cl.pass), b)
		savebytestofile(cl.dstf, r)
	} else {
		fmt.Println("[*]\tDecrypting file contents")
		r = dec([]byte(cl.pass), b)
		savebytestofile(cl.dstf, r)
	}
	// fmt.Println("[*]\tUsing pass:", cl.pass)
}

// Encrypts the given plaintext[]byte using the given key []byte
// Using AES in CTR mode of operation
// Nonce size = 12
// Returns [nonce+ciphertext] []byte
func enc(key []byte, plaintext []byte) (ciphertext []byte) {

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}

	// Never use more than 2^32 random nonces with a given key because of the risk of a repeat
	nonce := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}

	ciphertext = nonce
	// fmt.Println("[*]\tNonce: ", nonce)

	// fmt.Println("[*]\tNonce: ", ciphertext)

	ciphertext = aesgcm.Seal(ciphertext, nonce, plaintext, nil)

	return //	Basically returns nonce:ciphertext
}

// Func dec decrypts ciphertext []byte for the given key []byte
// ciphertext contents are [nonce+len(ciphertext)]
// Returns plaintext []byte
func dec(key []byte, ciphertext []byte) (plaintext []byte) {

	// nonce, _ := hex.DecodeString("64a9433eae7ccceee2fc0eda")
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}

	nonceSize := aesgcm.NonceSize()
	if len(ciphertext) < nonceSize {
		fmt.Println("[*]\tciphertext < nonceSize:", err)
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err = aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		panic(err.Error())
	}
	return
}

// Function filedemo contains a demo version of the mechanism.
// It encrypts the given files contents, saves them to a file
// and then shows its decrypted contents.
func filedemo() {

	cl := setupflags()

	fmt.Println("[*]\tsrc: ", cl.srcf)
	fmt.Println("[*]\tdst: ", cl.dstf)
	fmt.Println("[*]\tpass: ", cl.pass)

	b := returnfilecontents(cl.srcf)
	fmt.Printf("Init File Contents %x\n", b)

	eb := enc([]byte(cl.pass), b)
	fmt.Printf("Enc'd File's Contents %x\n", eb)
	savebytestofile(cl.dstf, eb)

	db := dec([]byte(cl.pass), eb)
	fmt.Printf("Resulting Dec'd Contents %x\n", db)

	returnfilecontents(cl.dstf)
}

// Function demo contains a demo version of the mechanism.
// This function encrypts a string, decrypts it and shows contents of the 3 files.
func demo() {
	// Load key
	// If you want to convert a passphrase to a key, use a suitable
	// package like bcrypt or scrypt.
	// decoded key:	16 bytes (AES-128) or 32 (AES-256)
	key, _ := hex.DecodeString("6368616e676520746869732070617373776f726420746f206120736563726574")
	plaintext := []byte("exampleplaintext")

	fmt.Printf("Plaintext - %x\n", plaintext)

	ciphertext := enc(key, plaintext)
	fmt.Printf("Ciphertext - %x\n", ciphertext)

	retrievedtext := dec(key, ciphertext)
	fmt.Printf("Plaintext - %x\n", retrievedtext)
}
