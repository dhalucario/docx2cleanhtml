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
	linkRelations  map[string]string
	styles         map[string]string
}

var htmlElementAliases = map[string]string{
	"Title":     "<h1>%s</h1>",
	"Heading 1": "<h2>%s</h2>",
	"Heading 2": "<h3>%s</h3>",
	"Heading 3": "<h4>%s</h4>",
	"Heading 4": "<h5>%s</h5>",
}

func New(file string) (doc Document, err error) {
	doc.styles = make(map[string]string)
	doc.linkRelations = make(map[string]string)

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
	doc.readStyles()
	doc.getLinkRelations()
	//doc.close()
}

func (doc *Document) readStyles() {
	file, fileErr := os.Open(path.Join(doc.tempPath, "styles.xml"))
	var parsedStyles xmlStyles
	if fileErr == nil {
		readAllContent, readAllErr := ioutil.ReadAll(file)
		if readAllErr == nil {
			parseErr := xml.Unmarshal(readAllContent, &parsedStyles)
			if parseErr != nil {
				log.Fatal(parseErr)
			}

			for _, style := range parsedStyles.Xstyles {
				if style.XstyleId != "" {
					doc.styles[style.XstyleId] = style.Xname.Xval
				}
			}
		}
	}
}

func (doc *Document) getLinkRelations() {
	file, fileErr := os.Open(path.Join(doc.tempPath, "document.xml.rels"))
	var relationships xmlRelationships
	if fileErr == nil {
		byteContent, readAllErr := ioutil.ReadAll(file)
		if readAllErr == nil {
			parseErr := xml.Unmarshal(byteContent, &relationships)
			if parseErr != nil {
				log.Fatal(parseErr)
			}

			for _, rl := range relationships.Xrelationships {
				doc.linkRelations[rl.Xid] = rl.Xtarget
			}
		} else {
			log.Fatal(readAllErr)
		}
	} else {
		log.Fatal(fileErr.Error())
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

func (doc *Document) close() {
	os.RemoveAll(doc.tempPath)
}

func (doc *Document) PrintHTML() {
	bufferPara := ""
	for _, p := range doc.parsedDocument.Xbody.Xparagraphs {
		for _, r := range p.paragraph.run {
			if r.urlId != "" {
				bufferPara += fmt.Sprintf("<a href=\"%s\">%s</a>", doc.linkRelations[r.urlId], r.text)
			} else {
				bufferPara += r.text
			}
		}

		if _, exists := htmlElementAliases[doc.styles[p.paragraph.style]]; exists {
			fmt.Printf(htmlElementAliases[doc.styles[p.paragraph.style]], bufferPara)
			fmt.Print("\n")
		} else {
			fmt.Printf("<p>%s<p>\n", bufferPara)
		}
		bufferPara = ""
	}
}
