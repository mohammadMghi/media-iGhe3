package base

type IFileInfo interface {
	GetBaseResult() IFileInfo
}

type FileInfo struct {
	IFileInfo   `json:"-"`
	RelativeURL string `json:"relative_url,omitempty"`
	Width       *int   `json:"width,omitempty"`
	Height      *int   `json:"height,omitempty"`
}

func (r *FileInfo) GetBaseResult() IFileInfo {
	return r
}
