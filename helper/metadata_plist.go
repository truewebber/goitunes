package helper

import (
	"fmt"

	"github.com/truewebber/goitunes/model"
)

var (
	plistManifest = `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>UIRequiredDeviceCapabilities</key>
	<dict>
		<key>armv7</key><true/>
	</dict>
	<key>appleId</key><string>VAR_APPLE_ID</string>
	<key>artistId</key><integer>%d</integer>
	<key>artistName</key><string>%s</string>
	<key>bundleDisplayName</key><string>%s</string>
	<key>bundleShortVersionString</key><string>%s</string>
	<key>bundleVersion</key><string>%s</string>
	<key>com.apple.iTunesStore.downloadInfo</key>
	<dict>
		<key>accountInfo</key>
		<dict>
			<key>AccountStoreFront</key><string>%s</string>
			<key>AppleID</key><string>VAR_APPLE_ID</string>
			<key>DSPersonID</key><integer>VAR_DSID</integer>
			<key>DownloaderID</key><integer>0</integer>
			<key>FamilyID</key><integer>0</integer>
			<key>FirstName</key><string>VAR_USER_FIRSTNAME</string>
			<key>LastName</key><string>VAR_USER_LASTNAME</string>
			<key>PurchaserID</key><integer>VAR_DSID</integer>
			<key>UserName</key><string>VAR_USER_FULL_NAME</string>
		</dict>
		<key>purchaseDate</key>
		<string>%s</string>
	</dict>
	<key>copyright</key><string>%s</string>
	<key>drmVersionNumber</key><integer>%d</integer>
	<key>fileExtension</key><string>.app</string>
	<key>gameCenterEnabled</key><false/>
	<key>gameCenterEverEnabled</key><false/>
	<key>genre</key><string>%s</string>
	<key>genreId</key><integer>%d</integer>
	<key>itemId</key><integer>%d</integer>
	<key>itemName</key><string>%s</string>
	<key>kind</key><string>software</string>
	<key>playlistName</key><string>%s</string>
	<key>product-type</key><string>ios-app</string>
	<key>purchaseDate</key><date>%s</date>
	<key>rating</key>
	<dict>
		<key>content</key><string>%s</string>
		<key>label</key><string>%s</string>
		<key>rank</key><integer>%d</integer>
		<key>system</key><string>%s</string>
	</dict>
	<key>releaseDate</key><string>%s</string>
	<key>s</key><integer>%d</integer>
	<key>softwareIcon57x57URL</key><string>%s</string>
	<key>softwareIconNeedsShine</key><true/>
	<key>softwareSupportedDeviceIds</key>
	<array>%s</array>
	<key>softwareVersionBundleId</key><string>%s</string>
	<key>softwareVersionExternalIdentifier</key><integer>%d</integer>
	<key>softwareVersionExternalIdentifiers</key>
	<array>%s</array>
	<key>userName</key><string>VAR_USER_FULL_NAME</string>
	<key>vendorId</key><integer>%d</integer>
	<key>versionRestrictions</key><integer>%d</integer>
</dict>
</plist>`

	//<key>subgenres</key>
	//<array>
	//<dict>
	//<key>genre</key><string>Симуляторы</string>
	//<key>genreId</key><integer>7015</integer>
	//</dict>
	//<dict>
	//<key>genre</key><string>Гонки</string>
	//<key>genreId</key><integer>7013</integer>
	//</dict>
	//</array>
)

func GenerateMetadataPlist(
	bundleId string, songPlist model.SongPlistSlice, xAppleStoreFront string, storeFront int,
) []byte {
	softwareSupportedDeviceIds := ""
	for _, val := range songPlist.MetaData.SoftwareSupportedDeviceIds {
		softwareSupportedDeviceIds += fmt.Sprintf("<integer>%d</integer>", val)
	}

	softwareVersionExternalIdentifiers := ""
	for _, val := range songPlist.MetaData.ExternalVersionIdList {
		softwareVersionExternalIdentifiers += fmt.Sprintf("<integer>%d</integer>", val)
	}

	return []byte(fmt.Sprintf(plistManifest,
		songPlist.MetaData.ArtistId,
		songPlist.MetaData.ArtistName,
		songPlist.MetaData.DisplayBundleId,
		songPlist.MetaData.BundleShortVersionString,
		songPlist.MetaData.BundleVersion,
		xAppleStoreFront,
		songPlist.PurchaseDate,
		songPlist.MetaData.Copyright,
		songPlist.MetaData.DrmVersionNumber,
		songPlist.MetaData.Genre,
		songPlist.MetaData.GenreId,
		songPlist.MetaData.ItemId,
		songPlist.MetaData.ItemName,
		songPlist.MetaData.PlaylistName,
		songPlist.PurchaseDate,
		songPlist.MetaData.Rating.Content,
		songPlist.MetaData.Rating.Label,
		songPlist.MetaData.Rating.Rank,
		songPlist.MetaData.Rating.System,
		songPlist.MetaData.ReleaseDate,
		storeFront,
		songPlist.MetaData.SoftwareIcon57x57URL,
		softwareSupportedDeviceIds,
		bundleId,
		songPlist.MetaData.ExternalVersionId,
		softwareVersionExternalIdentifiers,
		songPlist.MetaData.VendorId,
		songPlist.MetaData.VersionRestrictions,
	))
}
