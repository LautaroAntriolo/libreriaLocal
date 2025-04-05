package modelos

type BookInfo struct {
	URL           string       `json:"url"`
	Key           string       `json:"key"`
	Title         string       `json:"title"`
	Authors       []Author     `json:"authors"`
	NumberOfPages int          `json:"number_of_pages"`
	Publishers    []Publisher  `json:"publishers"`
	Cover         Cover        `json:"cover"`
}

type BookData struct {
	URL            string       `json:"url"`
	Title          string       `json:"title"`
	Authors        []Author     `json:"authors"`
	NumberOfPages  int          `json:"number_of_pages"`
	Publishers     []Publisher  `json:"publishers"`
	Cover          Cover        `json:"cover"`
	Isbn           string       `json:"isbn"`
}
type Author struct {
	Name string `json:"name"`
}

type Publisher struct {
	Name string `json:"name"`
}

type Cover struct {
	Small  string `json:"small"`
	Medium string `json:"medium"`
	Large  string `json:"large"`
}