package bookmarks

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/kooksee/go-assert"
	"github.com/spf13/cobra"
	"html/template"
	"os"
	"strconv"
	"strings"
	"time"
)

// Database is interface for manipulating data in database.
type Database interface {
	// SaveBookmark saves new bookmark to database.
	CreateBookmark(bookmark Bookmark) (int64, error)

	// GetBookmarks fetch list of bookmarks based on submitted indices.
	GetBookmarks(withContent bool, indices ...string) ([]Bookmark, error)

	//GetTags fetch list of tags and their frequency
	GetTags() ([]Tag, error)

	// DeleteBookmarks removes all record with matching indices from database.
	DeleteBookmarks(indices ...string) error

	// SearchBookmarks search bookmarks by the keyword or tags.
	SearchBookmarks(orderLatest bool, keyword string, tags ...string) ([]Bookmark, error)

	// UpdateBookmarks updates the saved bookmark in database.
	UpdateBookmarks(bookmarks ...Bookmark) ([]Bookmark, error)

	// CreateAccount creates new account in database
	CreateAccount(username, password string) error

	// GetAccount fetch account with matching username
	GetAccount(username string) (Account, error)

	// GetAccounts fetch list of accounts with matching keyword
	GetAccounts(keyword string) ([]Account, error)

	// DeleteAccounts removes all record with matching usernames
	DeleteAccounts(usernames ...string) error
}

// Tag is tag for the bookmark
type Tag struct {
	ID         int64  `db:"id"          json:"id"`
	Name       string `db:"name"        json:"name"`
	NBookmarks int64  `db:"n_bookmarks" json:"nBookmarks"`
	Deleted    bool   `json:"-"`
}

// Bookmark is record of a specified URL
type Bookmark struct {
	ID          int64  `db:"id"            json:"id"`
	URL         string `db:"url"           json:"url"`
	Title       string `db:"title"         json:"title"`
	ImageURL    string `db:"image_url"     json:"imageURL"`
	Excerpt     string `db:"excerpt"       json:"excerpt"`
	Author      string `db:"author"        json:"author"`
	MinReadTime int    `db:"min_read_time" json:"minReadTime"`
	MaxReadTime int    `db:"max_read_time" json:"maxReadTime"`
	Modified    string `db:"modified"      json:"modified"`
	Content     string `db:"content"       json:"-"`
	HTML        string `db:"html"          json:"-"`
	Tags        []Tag  `json:"tags"`
}

// Account is account for accessing bookmarks from web interface
type Account struct {
	ID       int64  `db:"id"       json:"id"`
	Username string `db:"username" json:"username"`
	Password string `db:"password" json:"password"`
}

// LoginRequest is login request
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Remember bool   `json:"remember"`
}

// exportBookmarks is handler for exporting bookmarks.
// Accept exactly one argument, the file to be exported.
func exportBookmarks(cmd *cobra.Command, args []string) {
	// Fetch bookmarks from database
	// Make sure destination directory exist
	dstDir := fp.Dir(args[0])
	os.MkdirAll(dstDir, os.ModePerm)

	// Open destination file
	dstFile, err := os.Create(args[0])
	if err != nil {
		cError.Println(err)
		return
	}
	defer dstFile.Close()

	// Create template
	funcMap := template.FuncMap{
		"unix": func(str string) int64 {
			t, err := time.Parse("2006-01-02 15:04:05", str)
			if err != nil {
				return time.Now().Unix()
			}

			return t.Unix()
		},
		"combine": func(tags []Tag) string {
			strTags := make([]string, len(tags))
			for i, tag := range tags {
				strTags[i] = tag.Name
			}

			return strings.Join(strTags, ",")
		},
	}

	tplContent := `<!DOCTYPE NETSCAPE-Bookmark-file-1>` +
		`<META HTTP-EQUIV="Content-Type" CONTENT="text/html; charset=UTF-8">` +
		`<TITLE>Bookmarks</TITLE>` +
		`<H1>Bookmarks</H1>` +
		`<DL><p>` +
		`{{range $book := .}}` +
		`<DT><A HREF="{{$book.URL}}" ADD_DATE="{{unix $book.Modified}}" TAGS="{{combine $book.Tags}}">{{$book.Title}}</A>` +
		`{{if gt (len $book.Excerpt) 0}}<DD>{{$book.Excerpt}}{{end}}{{end}}` +
		`</DL><p>`

	tpl, err := template.New("export").Funcs(funcMap).Parse(tplContent)
	if err != nil {
		assert.MustNotError(err)
		return
	}

	// Execute template
	err = tpl.Execute(dstFile, &bookmarks)
	if err != nil {
		cError.Println(err)
		return
	}

	fmt.Println("Export finished")
}

// importBookmarks is handler for importing bookmarks.
// Accept exactly one argument, the file to be imported.
func importBookmarks(cmd *cobra.Command, args []string) {
	// Parse flags
	generateTag := cmd.Flags().Changed("generate-tag")

	// If user doesn't specify, ask if tag need to be generated
	if !generateTag {
		var submit string
		fmt.Print("Add parents folder as tag? (y/n): ")
		fmt.Scanln(&submit)

		generateTag = submit == "y"
	}

	// Open bookmark's file
	srcFile, err := os.Open(args[0])
	if err != nil {
		cError.Println(err)
		return
	}
	defer srcFile.Close()

	// Parse bookmark's file
	doc, err := goquery.NewDocumentFromReader(srcFile)
	if err != nil {
		cError.Println(err)
		return
	}

	bookmarks := []Bookmark{}
	doc.Find("dt>a").Each(func(_ int, a *goquery.Selection) {
		// Get related elements
		dt := a.Parent()
		dl := dt.Parent()

		// Get metadata
		title := a.Text()
		url, _ := a.Attr("href")
		strTags, _ := a.Attr("tags")
		strModified, _ := a.Attr("last_modified")
		intModified, _ := strconv.ParseInt(strModified, 10, 64)
		modified := time.Unix(intModified, 0)

		// Get bookmark tags
		tags := []Tag{}
		for _, strTag := range strings.Split(strTags, ",") {
			if strTag != "" {
				tags = append(tags, Tag{Name: strTag})
			}
		}

		// Get bookmark excerpt
		excerpt := ""
		if dd := dt.Next(); dd.Is("dd") {
			excerpt = dd.Text()
		}

		// Get category name for this bookmark
		// and add it as tags (if necessary)
		category := ""
		if dtCategory := dl.Prev(); dtCategory.Is("h3") {
			category = dtCategory.Text()
			category = normalizeSpace(category)
			category = strings.ToLower(category)
			category = strings.Replace(category, " ", "-", -1)
		}

		if category != "" && generateTag {
			tags = append(tags, Tag{Name: category})
		}

		// Add item to list
		bookmark := Bookmark{
			URL:      url,
			Title:    normalizeSpace(title),
			Excerpt:  normalizeSpace(excerpt),
			Modified: modified.Format("2006-01-02 15:04:05"),
			Tags:     tags,
		}

		bookmarks = append(bookmarks, bookmark)
	})

	// Save bookmarks to database
	for _, book := range bookmarks {
		// Make sure URL valid
		parsedURL, err := nurl.ParseRequestURI(book.URL)
		if err != nil || parsedURL.Host == "" {
			cError.Println("URL is not valid")
			continue
		}

		// Clear UTM parameters from URL
		book.URL = clearUTMParams(parsedURL)

		// Save book to database
		book.ID, err = h.db.CreateBookmark(book)
		if err != nil {
			cError.Println(err)
			continue
		}

		printBookmarks(book)
	}
}
