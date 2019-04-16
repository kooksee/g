package bookmarks

//var (
//	ErrBookmarkEmpty    = fmt.Errorf("bookmark empty")
//	DefaultParseOptions = ParseOptions{
//		FoldersAsTags: false,
//	}
//)
//
//type Bookmark struct {
//	Title   string
//	Url     string
//	Created time.Time
//	Icon    string
//	Tags    []string
//}
//
//type ParseOptions struct {
//	// Converts the folder hierarchy into tags
//	// e.g.
//	// - Folder1
//	// -- Folder2
//	// --- Bookmark Tags[Tag1, Tag2]
//	//
//	// Will return a bookmark with tags Tag1, Tag2, Folder1, Folder2
//	FoldersAsTags bool
//}
//
//func parseBookmark(r string) (Bookmark, error) {
//	var bm Bookmark
//
//	tr := regexp.MustCompile(`(?i)<a.*>(.*?)<\/a>`)
//	ur := regexp.MustCompile(`(?i)href="(.*?)"`)
//	tsr := regexp.MustCompile(`(?i)add_date="(.*?)"`)
//	ir := regexp.MustCompile(`(?i)icon="(.*?)"`)
//	tagr := regexp.MustCompile(`(?i)tags="(.*?)"`)
//
//	titlematch := tr.FindStringSubmatch(r)
//	if len(titlematch) > 1 {
//		bm.Title = titlematch[1]
//	}
//
//	urlmatch := ur.FindStringSubmatch(r)
//	if len(urlmatch) > 1 {
//		bm.Url = urlmatch[1]
//	}
//
//	ts := tsr.FindStringSubmatch(r)
//	if len(ts) > 1 {
//		tsi, err := strconv.ParseInt(ts[1], 10, 64)
//		if err == nil {
//			bm.Created = time.Unix(tsi, 0)
//		}
//	}
//
//	iconmatch := ir.FindStringSubmatch(r)
//	if len(iconmatch) > 1 {
//		bm.Icon = iconmatch[1]
//	}
//
//	tagsmatch := tagr.FindStringSubmatch(r)
//	if len(tagsmatch) > 1 {
//		tags := strings.Split(tagsmatch[1], ",")
//		if len(tags) >= 1 && tagsmatch[1] != "" {
//			for i, tag := range tags {
//				tags[i] = strings.TrimSpace(tag)
//			}
//			bm.Tags = tags
//		}
//	}
//
//	if reflect.DeepEqual(Bookmark{}, bm) || bm.Url == "" {
//		return bm, ErrBookmarkEmpty
//	}
//
//	return bm, nil
//}
//
//func ParseWithOptions(r io.Reader, opts ParseOptions) ([]Bookmark, error) {
//	b, err := ioutil.ReadAll(r)
//	if err != nil {
//		return []Bookmark{}, err
//	}
//
//	return parseLines(string(b), opts)
//}
//
//func Parse(r io.Reader) ([]Bookmark, error) {
//	b, err := ioutil.ReadAll(r)
//	if err != nil {
//		return []Bookmark{}, err
//	}
//
//	return parseLines(string(b), DefaultParseOptions)
//}
//
//func parseLines(str string, opts ParseOptions) ([]Bookmark, error) {
//	lines := strings.Split(sanatize(str), "\n")
//	var bms []Bookmark
//	var folders []string
//
//	isFolder := regexp.MustCompile(`(?i)<h\d.*>(.*)<\/h\d>`)
//	isFolderClose := regexp.MustCompile(`(?i)<\/dl>`)
//	isLink := regexp.MustCompile(`(?i)<a`)
//
//	for _, line := range lines {
//		// Skip empty
//		if line == "" {
//			continue
//		}
//
//		if isFolder.MatchString(line) {
//			match := isFolder.FindStringSubmatch(line)
//			if len(match) >= 1 {
//				folders = append(folders, match[1])
//			}
//		}
//
//		if isFolderClose.MatchString(line) {
//			folders = folders[0 : len(folders)-1]
//		}
//
//		// Parse bookmark
//		if isLink.MatchString(line) {
//			bm, err := parseBookmark(line)
//			if err != nil {
//				continue
//			}
//			if opts.FoldersAsTags {
//				bm.Tags = append(bm.Tags, folders...)
//			}
//			bms = append(bms, bm)
//		}
//	}
//
//	return bms, nil
//}
//
//// Normalizes the bookmark file contents
//func sanatize(str string) string {
//	// Trim spaces and and newlines from beginning and end
//	s := strings.Trim(str, " \n\t")
//
//	// Remove carriage returns
//	s = strings.Replace(s, "\r", "", -1)
//
//	// Replace tabs with a space
//	s = strings.Replace(s, "\t", " ", -1)
//	return s
//}
