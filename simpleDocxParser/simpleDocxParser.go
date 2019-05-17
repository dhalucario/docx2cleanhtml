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

func New(file string) (doc *Document, err error) {
	newDoc := Document{}
	newDoc.styles = make(map[string]string)
	newDoc.linkRelations = make(map[string]string)

	md5hasher := md5.New()
	md5hasher.Write([]byte(strconv.FormatInt(time.Now().Unix(), 10)))
	md5hasher.Write([]byte(file))

	newDoc.originalPath = file
	newDoc.tempPath = path.Join("/tmp/docx2cleanhtml/", hex.EncodeToString(md5hasher.Sum(nil)))


	folderErr := os.MkdirAll(newDoc.tempPath, 0750)
	zipReader, zipErr := zip.OpenReader(file)

	if zipErr != nil {
		log.Fatal(zipErr.Error())
	}

	if folderErr != nil {
		log.Fatal(folderErr.Error())
	}

	for _, file := range zipReader.File {

		if !isAcceptedFile(file.Name) {
			continue
		}

		ofHandle, err := os.OpenFile(path.Join(newDoc.tempPath, path.Base(file.Name)), os.O_WRONLY|os.O_CREATE, 0750)
		if err != nil {
			return nil, err
		}

		fdHandle, err := file.Open()
		if err != nil {
			return nil, err
		}

		_, err = io.Copy(ofHandle, fdHandle)
		if err != nil {
			return nil, err
		}

		err = ofHandle.Close()
		if err != nil {
			log.Panic(err)
		}

		err = fdHandle.Close()
		if err != nil {
			log.Panic(err)
		}

	}

	return &newDoc, err
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

func (doc *Document) ReadRelations() error {
	err := doc.readDocuments()
	if err != nil {
		return err
	}

	err = doc.readStyles()
	if err != nil {
		return err
	}

	err = doc.getLinkRelations()
	if err != nil {
		return err
	}

	err = doc.close()
	if err != nil {
		log.Panic(err)
	}

	return nil
}

func (doc *Document) readStyles() error {
	var parsedStyles xmlStyles
	file, err := os.Open(path.Join(doc.tempPath, "styles.xml"))
	if err != nil {
		return err
	}

	readAllContent, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	err = xml.Unmarshal(readAllContent, &parsedStyles)
	if err != nil {
		return err
	}

	for _, style := range parsedStyles.Xstyles {
		doc.styles[style.XstyleId] = style.Xname.Xval
	}

	return nil
}

func (doc *Document) getLinkRelations() error {
	var relationships xmlRelationships
	file, err := os.Open(path.Join(doc.tempPath, "document.xml.rels"))

	if err != nil {
		return err
	}

	byteContent, err := ioutil.ReadAll(file)

	if err != nil {
		return err
	}

	err = xml.Unmarshal(byteContent, &relationships)
	if err != nil {
		return err
	}

	for _, rl := range relationships.Xrelationships {
		doc.linkRelations[rl.Xid] = rl.Xtarget
	}

	return nil
}

func (doc *Document) readDocuments() error {
	file, err := os.Open(path.Join(doc.tempPath, "document.xml"))

	if err != nil {
		return err
	}

	byteContent, err := ioutil.ReadAll(file)

	if err != nil {
		return err
	}

	err = xml.Unmarshal(byteContent, &doc.parsedDocument)
	if err != nil {
		return err
	}

	return nil
}

func (doc *Document) close() error {
	return os.RemoveAll(doc.tempPath)
}

func (doc *Document) HTML() string {
	bufferPara := ""
	html := ""
	for i, p := range doc.parsedDocument.Xbody.Xparagraphs {
		for _, r := range p.paragraph.run {
			if r.urlId != "" {
				bufferPara += fmt.Sprintf("<a href=\"%s\">%s</a>", doc.linkRelations[r.urlId], r.text)
			} else {
				bufferPara += r.text
			}
		}

		if _, exists := htmlElementAliases[doc.styles[p.paragraph.style]]; exists {
			html += fmt.Sprintf(htmlElementAliases[doc.styles[p.paragraph.style]], bufferPara)
		} else {
			html += fmt.Sprintf("<p>%s<p>", bufferPara)
		}

		if i < len(doc.parsedDocument.Xbody.Xparagraphs) {
			html += "\n"
		}

		bufferPara = ""
	}

	return html
}
