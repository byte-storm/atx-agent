package main

import (
	"encoding/xml"
	"fmt"
	"strconv"
	"strings"
)

var sendButtonPackage = []string{
	"com.tencent.mm",
	"com.tencent.wework",
	"com.alibaba.android.rimet",
}

type Node struct {
	Index      string `xml:"index,attr"`
	Text       string `xml:"text,attr"`
	ResourceID string `xml:"resource-id,attr"`
	Class      string `xml:"class,attr"`
	Package    string `xml:"package,attr"`
	Bounds     Bounds `xml:"bounds,attr"`
	Children   []Node `xml:",any"`
}

type Bounds struct {
	Left   int `xml:"-"`
	Top    int `xml:"-"`
	Right  int `xml:"-"`
	Bottom int `xml:"-"`
}

type Position struct {
	X int `json:"x"`
	Y int `json:"y"`
}

func (b *Bounds) UnmarshalXMLAttr(attr xml.Attr) error {
	//values := strings.Fields(attr.Value)
	values := strings.Split(attr.Value, "][")

	if len(values) != 2 {
		return fmt.Errorf("invalid bounds format: %s", attr.Value)
	}

	rect1 := strings.Trim(values[0], "[]")
	rect2 := strings.Trim(values[1], "[]")

	leftTop := strings.Split(rect1, ",")
	rightBottom := strings.Split(rect2, ",")

	if len(leftTop) != 2 || len(rightBottom) != 2 {
		return fmt.Errorf("invalid bounds format: %s", attr.Value)
	}

	left, err := strconv.Atoi(leftTop[0])
	if err != nil {
		return err
	}
	top, err := strconv.Atoi(leftTop[1])
	if err != nil {
		return err
	}
	right, err := strconv.Atoi(rightBottom[0])
	if err != nil {
		return err
	}
	bottom, err := strconv.Atoi(rightBottom[1])
	if err != nil {
		return err
	}

	b.Left = left
	b.Top = top
	b.Right = right
	b.Bottom = bottom

	return nil
}

func reverseFindNodeByText(node Node, searchText string, resultChan chan<- Node) {
	for i := len(node.Children) - 1; i >= 0; i-- {
		reverseFindNodeByText(node.Children[i], searchText, resultChan)
	}

	if node.Text == searchText {
		resultChan <- node
	}
}

func GetSendButtonBounds(xmldata []byte) (pos Position, err error) {
	var root Node
	err = xml.Unmarshal(xmldata, &root)
	if err != nil {
		err = fmt.Errorf("Error decoding XML:%w", err)
		return
	}

	searchText := "发送"
	resultChan := make(chan Node)

	go func() {
		reverseFindNodeByText(root, searchText, resultChan)
		close(resultChan)
	}()

	for node := range resultChan {
		pos.X = node.Bounds.Left
		pos.Y = node.Bounds.Top
		return
	}
	return
}
