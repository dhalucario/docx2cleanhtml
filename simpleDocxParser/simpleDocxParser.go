package simpleDocxParser

import (
	"archive/zip"
	"crypto/md5"
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strconv"
	"time"
)

type Document struct {
	originalPath string
	tempPath     string

	parsedDocument xmlDocument
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

func (doc *Document) ReadRelations() {
	doc.readDocuments()
	fmt.Printf("%v", doc.parsedDocument)
}

func (doc *Document) readStyles() {
	// TODO: Fix this up. Basicly same as readDocuments.
	file, fileErr := os.Open(path.Join(doc.tempPath, "styles.xml"))
	if fileErr == nil {
		readAllContent, readAllErr := ioutil.ReadAll(file)
		if readAllErr == nil {
			parseErr := xml.Unmarshal()
		}
	}
}

func (doc *Document) readDocuments() {
	file, fileErr := os.Open(path.Join(doc.tempPath, "document.xml"))
	if fileErr == nil {
		byteContent, readAllErr := ioutil.ReadAll(file)
		if readAllErr == nil {
			parseErr := xml.Unmarshal(byteContent, &doc.parsedDocument)
			if parseErr != nil {
				log.Fatal(parseErr)
			}
		} else {
			log.Fatal(readAllErr)
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
