package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"local-storage/storage"
)

var localStorage storage.Storage
var dirPath string

func init() {
	dirPath = filepath.Join(os.TempDir(), "local-storage")
	localStorage = storage.NewLocalStorage(dirPath)
	fmt.Printf("local storage : %s\n", dirPath)
}

func generateUploadName(src string) (string, error) {
	f, err := os.Open(src)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		log.Fatal(err)
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	f, err := ioutil.TempFile(os.TempDir(), "*.received.tmp")
	if err != nil {
		fmt.Println(err)
		return
	}

	n, err := io.Copy(f, r.Body)
	if err != nil {
		panic(err)
	}

	err = f.Close()
	if err != nil {
		panic(err)
	}

	dst, err := generateUploadName(f.Name())
	if err != nil {
		panic(err)
	}

	err = localStorage.Move(f.Name(), dst)
	if err != nil {
		panic(err)
	}
	w.Write([]byte(fmt.Sprintf("%d bytes are recieved. \nchecksum: %s ", n, dst)))
}

func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	checksum := r.URL.Query().Get("checksum")
	exist, err := localStorage.Exist(checksum)
	if err != nil {
		panic(err)
	}
	if exist {
		w.Header().Set("Content-Type", "application/octet-stream")
		// w.Header().Set("Content-Transfer-Encoding", "binary")
		// w.Header().Set("Expires", "0")
		w.WriteHeader(http.StatusOK)
		dst := filepath.Join(dirPath, checksum)
		f, err := os.Open(dst)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		io.Copy(w, f)
	} else {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(fmt.Sprintf("checksum: %s no exist.\n", checksum)))
	}
}

func main() {
	http.HandleFunc("/upload", UploadHandler)
	http.HandleFunc("/download", DownloadHandler)
	log.Fatal(http.ListenAndServe(":10000", nil))
}
