package html2md

import (
	"strings"

	"fmt"

	"regexp"

	"github.com/PuerkitoBio/goquery"
	"github.com/astaxie/beego/logs"
)

var closeTag = map[string]string{
	"del":    "~~",
	"b":      "**",
	"strong": "**",
	"i":      "_",
	"em":     "_",
	"dfn":    "_",
	"var":    "_",
	"cite":   "_",
	"br":     "\n",
	"span":   "",
	"small":  "",
}

var blockTag = []string{
	"div",
	"figure",
	"p",
	"article",
	"nav",
	"footer",
	"header",
	"section",
}
var nextlineTag = []string{
	"pre", "blockquote", "table",
}

//convert html to markdown
//将html转成markdown
func Convert(htmlstr string) (md string) {
	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(htmlstr))
	doc = compressHtml(doc)
	doc = handleNextLine(doc)   //<div>...
	doc = handleBlockTag(doc)   //<div>...
	doc = handleA(doc)          //<a>
	doc = handleImg(doc)        //<img>
	doc = handleHead(doc)       //h1~h6
	doc = handleClosedTag(doc)  //<strong>、<i>、eg..
	doc = handleHr(doc)         //<hr>
	doc = handleLi(doc)         //<li>
	doc = handleTable(doc)      //<table>
	doc = handleBlockquote(doc) //<table>
	md, _ = doc.Find("body").Html()
	return
}

//压缩html
func compressHtml(doc *goquery.Document) *goquery.Document {
	//blockquote、pre、code
	var maps = make(map[string]string)
	if ele := doc.Find("blockquote"); len(ele.Nodes) > 0 {
		ele.Each(func(i int, selection *goquery.Selection) {
			key := fmt.Sprintf("{$blockquote%v}", i)
			cont := "<blockquote>" + getInnerHtml(selection) + "</blockquote>"
			selection.BeforeHtml(key)
			selection.Remove()
			maps[key] = cont
		})
	}
	if ele := doc.Find("pre"); len(ele.Nodes) > 0 {
		ele.Each(func(i int, selection *goquery.Selection) {
			key := fmt.Sprintf("{$pre%v}", i)
			cont := "<pre>" + getInnerHtml(selection) + "</pre>"
			selection.BeforeHtml(key)
			selection.Remove()
			maps[key] = cont
		})
	}
	if ele := doc.Find("code"); len(ele.Nodes) > 0 {
		ele.Each(func(i int, selection *goquery.Selection) {
			key := fmt.Sprintf("{$code%v}", i)
			cont := "<code>" + getInnerHtml(selection) + "</code>"
			selection.BeforeHtml(key)
			selection.Remove()
			maps[key] = cont
		})
	}
	htmlstr, _ := doc.Html()
	htmlstr = strings.Replace(htmlstr, "\n", "", -1)
	htmlstr = strings.Replace(htmlstr, "\r", "", -1)
	htmlstr = strings.Replace(htmlstr, "\t", "", -1)
	//正则匹配，把“>”和“<”直接的空格全部去掉
	//去除标签之间的空格，如果是存在代码预览的页面，不要替换空格，否则预览的代码会错乱
	r, _ := regexp.Compile(">\\s{1,}<")
	htmlstr = r.ReplaceAllString(htmlstr, "><")
	//多个空格替换成一个空格
	r2, _ := regexp.Compile("\\s{1,}")
	htmlstr = r2.ReplaceAllString(htmlstr, " ")
	for key, val := range maps {
		htmlstr = strings.Replace(htmlstr, key, val, -1)
	}
	doc, _ = goquery.NewDocumentFromReader(strings.NewReader(htmlstr))
	return doc
}

func handleBlockTag(doc *goquery.Document) *goquery.Document {
	for _, tag := range blockTag {
		hasTag := true
		for hasTag {
			if tagEle := doc.Find(tag); len(tagEle.Nodes) > 0 {
				tagEle.Each(func(i int, selection *goquery.Selection) {
					selection.BeforeHtml("\n" + getInnerHtml(selection) + "\n")
					selection.Remove()
				})
			} else {
				hasTag = false
			}
		}
	}
	return doc
}

func handleBlockquote(doc *goquery.Document) *goquery.Document {
	if tagEle := doc.Find("blockquote"); len(tagEle.Nodes) > 0 {
		tagEle.Each(func(i int, selection *goquery.Selection) {
			cont := getInnerHtml(selection)
			cont = strings.Replace(cont, "\r", "", -1)
			cont = strings.Replace(cont, "\n", "", -1)
			selection.BeforeHtml("\r\n<blockquote>" + cont + "\n</blockquote>\n")
			selection.Remove()
		})
	}
	return doc
}

//[ok]handle tag <a>
func handleA(doc *goquery.Document) *goquery.Document {
	doc.Find("a").Each(func(i int, selection *goquery.Selection) {
		if href, ok := selection.Attr("href"); ok {
			if cont, err := selection.Html(); err == nil {
				md := fmt.Sprintf(`[%v](%v)`, cont, href)
				selection.BeforeHtml(md)
				selection.Remove()
			}
		}
	})
	return doc
}

//[ok]handle tag ul、ol、li
//处理步骤：
//1、先给每个li标签里面的内容加上"- "或者"\t- "
//2、提取li内容
func handleLi(doc *goquery.Document) *goquery.Document {
	var tags = []string{"ul", "li"}
	doc.Find("li").Each(func(i int, selection *goquery.Selection) {
		l := len(selection.ParentsFiltered("li").Nodes)
		tab := strings.Join(make([]string, l+2), "{$space}")
		selection.PrependHtml("\r$" + tab)
	})
	for _, tag := range tags {
		doc.Find(tag).Each(func(i int, selection *goquery.Selection) {
			selection.BeforeHtml(selection.Text())
			selection.Remove()
		})
	}
	htmlstr, _ := doc.Find("body").Html()
	for i := 10; i > 0; i-- {
		oldTab := "$" + strings.Join(make([]string, i), "{$space}")
		newTab := strings.Join(make([]string, i-1), "  ") + "- "
		htmlstr = strings.Replace(htmlstr, oldTab, newTab, -1)
	}
	doc, _ = goquery.NewDocumentFromReader(strings.NewReader(htmlstr))
	return doc
}

//[ok]handle tag <hr/>
func handleHr(doc *goquery.Document) *goquery.Document {
	doc.Find("hr").Each(func(i int, selection *goquery.Selection) {
		selection.BeforeHtml("\n- - -\n")
		selection.Remove()
	})
	return doc
}

//[ok]handle tag <img/>
func handleImg(doc *goquery.Document) *goquery.Document {
	doc.Find("img").Each(func(i int, selection *goquery.Selection) {
		if src, ok := selection.Attr("src"); ok {
			alt := ""
			if val, ok := selection.Attr("alt"); ok {
				alt = val
			}
			md := fmt.Sprintf(`![%v](%v)`, alt, src)
			selection.BeforeHtml(md)
			selection.Remove()
		}
	})
	return doc
}

//[ok]handle tag h1~h6
func handleHead(doc *goquery.Document) *goquery.Document {
	heads := map[string]string{
		"title": "# ",
		"h1":    "# ",
		"h2":    "## ",
		"h3":    "### ",
		"h4":    "#### ",
		"h5":    "##### ",
		"h6":    "###### ",
	}
	for tag, replace := range heads {
		doc.Find(tag).Each(func(i int, selection *goquery.Selection) {
			text, _ := selection.Html()
			selection.BeforeHtml("\n" + replace + text + "\n")
			selection.Remove()
		})
	}
	return doc
}

func handleClosedTag(doc *goquery.Document) *goquery.Document {
	for tag, close := range closeTag {
		doc.Find(tag).Each(func(i int, selection *goquery.Selection) {
			if text, _ := selection.Html(); strings.TrimSpace(text) != "" {
				selection.BeforeHtml(close + text + close)
			}
			selection.Remove()
		})
	}
	return doc
}

func handleNextLine(doc *goquery.Document) *goquery.Document {
	for _, tag := range nextlineTag {
		doc.Find(tag).Each(func(i int, selection *goquery.Selection) {
			selection.BeforeHtml("\n")
			selection.AfterHtml("\n")
		})
	}
	return doc
}

func handleTable(doc *goquery.Document) *goquery.Document {
	doc.Find("table").Each(func(i int, table *goquery.Selection) {
		rows := []string{}
		table.Find("tr").Each(func(i int, tr *goquery.Selection) {
			ths := []string{}
			tr.Find("th").Each(func(i int, trth *goquery.Selection) {
				ths = append(ths, getInnerHtml(trth))
			})
			if len(ths) > 0 {
				rows = append(rows, "|"+strings.Join(ths, "|")+"\n|-----\n")
			}
			tds := []string{}
			tr.Find("td").Each(func(i int, trtd *goquery.Selection) {
				tds = append(tds, getInnerHtml(trtd))
			})
			if len(tds) > 0 {
				rows = append(rows, "|"+strings.Join(tds, "|")+"\n")
			}
		})
		table.BeforeHtml(strings.Join(rows, ""))
		table.Remove()
	})
	return doc
}

func getInnerHtml(selection *goquery.Selection) (html string) {
	var err error
	html, _ = selection.Html()
	if err != nil {
		logs.Error(err)
	}
	return
}
