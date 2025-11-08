package entity

// DownloadInfo contains all necessary information to download an application.
type DownloadInfo struct {
	bundleID    string
	url         string
	downloadKey string
	sinf        string
	metadata    string
	headers     map[string]string
	downloadID  string
	versionID   int64
	fileSize    int64
}

// NewDownloadInfo creates a new DownloadInfo entity.
func NewDownloadInfo(bundleID, url, downloadKey string) *DownloadInfo {
	return &DownloadInfo{
		bundleID:    bundleID,
		url:         url,
		downloadKey: downloadKey,
		headers:     make(map[string]string),
	}
}

// Getters.
func (d *DownloadInfo) BundleID() string           { return d.bundleID }
func (d *DownloadInfo) URL() string                { return d.url }
func (d *DownloadInfo) DownloadKey() string        { return d.downloadKey }
func (d *DownloadInfo) Sinf() string               { return d.sinf }
func (d *DownloadInfo) Metadata() string           { return d.metadata }
func (d *DownloadInfo) Headers() map[string]string { return d.headers }
func (d *DownloadInfo) DownloadID() string         { return d.downloadID }
func (d *DownloadInfo) VersionID() int64           { return d.versionID }
func (d *DownloadInfo) FileSize() int64            { return d.fileSize }

// Setters.
func (d *DownloadInfo) SetSinf(sinf string) *DownloadInfo {
	d.sinf = sinf

	return d
}

func (d *DownloadInfo) SetMetadata(metadata string) *DownloadInfo {
	d.metadata = metadata

	return d
}

func (d *DownloadInfo) SetHeaders(headers map[string]string) *DownloadInfo {
	d.headers = headers

	return d
}

func (d *DownloadInfo) AddHeader(key, value string) *DownloadInfo {
	d.headers[key] = value

	return d
}

func (d *DownloadInfo) SetDownloadID(id string) *DownloadInfo {
	d.downloadID = id

	return d
}

func (d *DownloadInfo) SetVersionID(id int64) *DownloadInfo {
	d.versionID = id

	return d
}

func (d *DownloadInfo) SetFileSize(size int64) *DownloadInfo {
	d.fileSize = size

	return d
}
