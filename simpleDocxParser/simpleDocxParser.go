package simpleDocxParser

import (
	"archive/zip"
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"strconv"
	"time"
)

type Document struct {
	originalPath string
	tempPath     string
}

func New(file string) (doc Document, err error) {
	md5hasher := md5.New()
	md5hasher.Write([]byte(strconv.FormatInt(time.Now().Unix(), 10)))
	md5hasher.Write([]byte(file))

	doc.originalPath = file
	doc.tempPath = path.Join("/tmp/docx2cleanhtml/", hex.EncodeToString(md5hasher.Sum(nil)))

	folderErr := os.MkdirAll(doc.tempPath, 0750)
	zipReader, zipErr := zip.OpenReader(file)

	if zipErr == nil {
		if folderErr == nil {
			for _, file := range zipReader.File {
				fmt.Println(file.Name)
				if isAcceptedFile(file.Name) {
					ofHandle, ofErr := os.OpenFile(path.Join(doc.tempPath, path.Base(file.Name)), os.O_WRONLY|os.O_CREATE, 0750)
					fdHandle, fdErr := file.Open()
					if fdErr == nil {
						if ofErr == nil {
							_, copyErr := io.Copy(ofHandle, fdHandle)
							if copyErr != nil {
								log.Fatal(copyErr.Error())
							}
							ofcErr := ofHandle.Close()
							if ofcErr != nil {
								log.Fatal(ofcErr.Error())
							}
						} else {
							err = ofErr
						}
						fccErr := fdHandle.Close()
						if fccErr != nil {
							log.Fatal(fccErr.Error())
						}
					} else {
						log.Fatal(fdErr.Error())
					}
				}
			}
		} else {
			log.Fatal(folderErr.Error())
		}
	} else {
		log.Fatal(zipErr.Error())
	}

	return doc, err
}

func isAcceptedFile(filename string) bool {
	requiredFiles := []string{
		"word/document.xml",
		"word/styles.xml",
		"word/_rels/document.xml.rels",
	}
	for _, elem := range requiredFiles {
		if filename == elem {
			return true
		}
	}

	return false
}

func (doc *Document) readRelations(relativePath string) {
	file, fileErr := os.Open(relativePath)
	fileBuffer := bytes.NewBuffer(nil)
	if fileErr == nil {
		_, copyErr := io.Copy(file, fileBuffer)
		if copyErr == nil {

		} else {
			log.Fatal(copyErr)
		}
	} else {
		log.Fatal(fileErr.Error())
	}
}

func (doc *Document) readParagraphs(relativePath string) {

}

func (doc *Document) Close() (err error) {
	return os.RemoveAll(doc.tempPath)
}

func (doc *Document) GetHTML() {}
