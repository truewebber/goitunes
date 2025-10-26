package mapper

import (
	"github.com/truewebber/goitunes/internal/application/dto"
	"github.com/truewebber/goitunes/internal/domain/entity"
)

// ApplicationMapper handles mapping between domain entities and DTOs
type ApplicationMapper struct{}

// NewApplicationMapper creates a new application mapper
func NewApplicationMapper() *ApplicationMapper {
	return &ApplicationMapper{}
}

// ToDTO maps an Application entity to ApplicationDTO
func (m *ApplicationMapper) ToDTO(app *entity.Application) dto.ApplicationDTO {
	return dto.ApplicationDTO{
		AdamID:           app.AdamID(),
		BundleID:         app.BundleID(),
		Name:             app.Name(),
		ArtistName:       app.ArtistName(),
		ArtistID:         app.ArtistID(),
		Version:          app.Version(),
		VersionID:        app.VersionID(),
		Price:            app.Price(),
		Currency:         app.Currency(),
		Rating:           app.Rating(),
		RatingCount:      app.RatingCount(),
		ReleaseDate:      app.ReleaseDate(),
		GenreID:          app.GenreID(),
		GenreName:        app.GenreName(),
		DeviceFamilies:   app.DeviceFamilies(),
		FileSize:         app.FileSize(),
		MinimumOSVersion: app.MinimumOSVersion(),
		Description:      app.Description(),
		IconURL:          app.IconURL(),
		ScreenshotURLs:   app.ScreenshotURLs(),
		IsFree:           app.IsFree(),
		IsUniversal:      app.IsUniversal(),
	}
}

// ToDTOList maps a list of Application entities to DTOs
func (m *ApplicationMapper) ToDTOList(apps []*entity.Application) []dto.ApplicationDTO {
	dtos := make([]dto.ApplicationDTO, 0, len(apps))
	for _, app := range apps {
		dtos = append(dtos, m.ToDTO(app))
	}
	return dtos
}

// ChartItemToDTO maps a ChartItem entity to ChartItemDTO
func (m *ApplicationMapper) ChartItemToDTO(item *entity.ChartItem) dto.ChartItemDTO {
	return dto.ChartItemDTO{
		Position: item.Position(),
		App:      m.ToDTO(item.Application()),
	}
}

// ChartItemsToDTOList maps a list of ChartItem entities to DTOs
func (m *ApplicationMapper) ChartItemsToDTOList(items []*entity.ChartItem) []dto.ChartItemDTO {
	dtos := make([]dto.ChartItemDTO, 0, len(items))
	for _, item := range items {
		dtos = append(dtos, m.ChartItemToDTO(item))
	}
	return dtos
}

// DownloadInfoToDTO maps a DownloadInfo entity to DownloadInfoDTO
func (m *ApplicationMapper) DownloadInfoToDTO(info *entity.DownloadInfo) dto.DownloadInfoDTO {
	return dto.DownloadInfoDTO{
		BundleID:    info.BundleID(),
		URL:         info.URL(),
		DownloadKey: info.DownloadKey(),
		Sinf:        info.Sinf(),
		Metadata:    info.Metadata(),
		Headers:     info.Headers(),
		DownloadID:  info.DownloadID(),
		VersionID:   info.VersionID(),
		FileSize:    info.FileSize(),
	}
}

