package xregex

const (
	UUIDGROUP  string = "([a-f0-9]{8}-[a-f0-9]{4}-4[a-f0-9]{3}-[89aAbB][a-f0-9]{3}-[a-f0-9]{12})?"
	URLGROUP   string = `(http|ftp|https):\/\/[\w\-_]+(\.[\w\-_]+)+([\w\-\.,@?^=%&amp;:/~\+#]*[\w\-\@?^=%&amp;/~\+#])?`
	REGEXEnd   string = "$"
	REGEXBegin string = "^"
)

var (
	// uuid required
	UUID string = REGEXBegin + UUIDGROUP + REGEXEnd
	// url required
	URL string = REGEXBegin + URLGROUP + REGEXEnd
)
