package core

type ImageInfo struct {
	Path   string `json:"path"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
	Size   int    `json:"size"`
}

func (i *ImageInfo) GetPath() {

}
