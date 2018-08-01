package helper

import (
	"fmt"
	"net/url"
	"strings"
	"time"
)

var (
	BuyTemplate = `<?xml version="1.0" encoding="UTF-8"?>
<plist version="1.0">
<dict>
	<key>appExtVrsId</key><string>%d</string>
	<key>guid</key><string>%s</string>
	<key>kbsync</key><data>%s</data>
	<key>machineName</key><string>%s</string>
	<key>mtApp</key><string>com.apple.iTunes</string>
	<key>mtClientId</key><string>3z30dhYIz29Wz4gvz9AEz1NIUDKelm</string>
	<key>mtEventTime</key><string>%d</string>
	<key>mtPageContext</key><string>App Store</string>
	<key>mtPageId</key><string>1140828062</string>
	<key>mtPageType</key><string>Software</string>
	<key>mtPrevPage</key><string>Genre_134583</string>
	<key>mtRequestId</key><string>3z30dhYIz29Wz4gvz9AEz1NIUDKelmzJ4H6DIUSz1HZC</string>
	<key>mtTopic</key><string>xp_its_main</string>
	<key>needDiv</key><string>0</string>
	<key>pg</key><string>default</string>
	<key>price</key><string>0</string>
	<key>pricingParameters</key><string>%s</string>
	<key>rebuy</key><string>%s</string> 
	<key>productType</key><string>C</string>
	<key>salableAdamId</key><string>%s</string>
	<key>uuid</key><string>353F3F00-9D87-5BB1-9055-B7761CCD57AA</string>
</dict>
</plist>`
	//<key>pricingParameters</key>
	//<string>STDRDL</string>//STDQ
	//<key>ownerDsid</key>
	//<string></string>
)

const (
	PricingParametersBuy        = "STDQ"
	PricingParametersReDownload = "STDRDL"
)

func GenerateBuyProductBody(
	adamId string, extVersionId int, pricingParameters string, guid string, kbsync string, machineName string,
) *strings.Reader {
	unixAppleTime := time.Now().UnixNano() / 1000000

	rebuy := "false"
	if pricingParameters == PricingParametersReDownload {
		rebuy = "true"
	}

	body := fmt.Sprintf(
		BuyTemplate,
		extVersionId,
		guid,
		kbsync,
		machineName,
		unixAppleTime,
		pricingParameters,
		rebuy,
		adamId,
	)

	return strings.NewReader(body)
}

func GenerateLoginBody(password string, machineName string, guid string, appleId string) *strings.Reader {
	params := url.Values{
		"matchineName":  {machineName},
		"why":           {"signin"},
		"attempt":       {"1"},
		"createSession": {"true"},
		"guid":          {guid},
		"appleId":       {appleId},
		"password":      {password},
	}

	return strings.NewReader(params.Encode())
}
