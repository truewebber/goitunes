package goitunes

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/groob/plist"
	"github.com/mgutz/logxi/v1"

	"github.com/truewebber/goitunes/helper"
	"github.com/truewebber/goitunes/model"
	"github.com/truewebber/goitunes/store"
)

type (
	GOiTunes struct {
		AppleId     string
		XToken      string
		Kbsync      string
		DSID        string
		Geo         string
		GUID        string
		UserAgent   string
		MachineName string

		httpClient *http.Client
	}
)

const (
	Top200UserAgent   = "AppStore/2.0 iOS/9.0 model/iPhone6,1 hwp/s5l8960x build/13A344 (6; dt:89)"
	TopUserAgent      = "iTunes-iPad/5.1.1 (64GB; dt:28)"
	DownloadUserAgent = "itunesstored/1.0 iOS/9.0 model/iPhone6,1 hwp/s5l8960x build/13A344 (6; dt:89)"

	Top200AppsURL               = "https://itunes.apple.com/WebObjects/MZStore.woa/wa/viewTop"
	TopAppsURL                  = "https://itunes.apple.com/WebObjects/MZStore.woa/wa/topChartFragmentData"
	AppInfoURL                  = "https://uclient-api.itunes.apple.com/WebObjects/MZStorePlatform.woa/wa/lookup"
	NativeAppInfoURL            = "https://itunes.apple.com/app/id%s?mt=8"
	NativeAppRatingInfoURL      = "https://itunes.apple.com/customer-reviews/id%s?dataOnly=true&displayable-kind=11"
	OpenAppOverAllRatingInfoURL = "https://itunes.apple.com/lookup?id=%s&entity=software&country=%s"

	LoginURL        = "https://p%d-buy.itunes.apple.com/WebObjects/MZFinance.woa/wa/authenticate"
	BuyProductURL   = "https://p%d-buy.itunes.apple.com/WebObjects/MZBuy.woa/wa/buyProduct?%s"
	ConfirmDownload = "https://p%d-buy.itunes.apple.com/WebObjects/MZFastFinance.woa/wa/songDownloadDone?%s"
)

func NewGOiTunes(
	appleId string,
	xToken string,
	kbsync string,
	dsid string,
	geo string,
	guid string,
	userAgent string,
	machineName string,
) (*GOiTunes, error) {
	if _, ok := store.GetStoreList()[strings.ToLower(geo)]; !ok {
		return nil, fmt.Errorf("Geo `%s` is't supported yet", geo)
	}

	return &GOiTunes{
		AppleId:     appleId,
		XToken:      xToken,
		Kbsync:      kbsync,
		DSID:        dsid,
		Geo:         strings.ToLower(geo),
		GUID:        guid,
		UserAgent:   userAgent,
		MachineName: machineName,
		httpClient:  &http.Client{},
	}, nil
}

func (g *GOiTunes) getStore() *store.ItunesStore {
	return store.GetStoreList()[g.Geo]
}

func (g *GOiTunes) getLoginUrl() string {
	return fmt.Sprintf(LoginURL, g.getStore().HostPrefix)
}

func (g *GOiTunes) getBuyProductUrl() string {
	params := url.Values{
		"xToken": {g.XToken},
	}

	return fmt.Sprintf(BuyProductURL, g.getStore().HostPrefix, params.Encode())
}

//Unauthorized

func (g *GOiTunes) GetTop200Applications(
	genreId string, popId string, kidPrefix string, from int, chunk int,
) ([]model.TopAppItemResponse, error) {
	if from < 1 {
		from = 1
	}

	var out []model.TopAppItemResponse

	//check popId. Only free for a while
	if !store.PopIdExists(popId) {
		return out, errors.New(fmt.Sprintf("PopId `%s` not exists", popId))
	}

	req, err := http.NewRequest("GET", Top200AppsURL, nil)
	if err != nil {
		return out, fmt.Errorf("Error create request top200 charts iTunes, error: %s", err.Error())
	}

	//query build
	q := req.URL.Query()
	if len(kidPrefix) > 0 {
		q.Add("ageBandId", kidPrefix)
	}
	q.Add("genreId", genreId)
	q.Add("popId", popId)
	q.Add("cc", g.Geo)
	q.Add("l", "en")
	req.URL.RawQuery = q.Encode()

	//set headers
	req.Header.Add("User-Agent", Top200UserAgent)

	//exec request
	resp, err := g.httpClient.Do(req)
	if err != nil {
		return out, fmt.Errorf("Error send request top200 charts iTunes, error: %s", err.Error())
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return out, fmt.Errorf("Error read response body top200 charts request iTunes, error: %s", err.Error())
	}

	//log.Debug("Response top200", "code", resp.StatusCode, "bodyLen", len(data),
	//	"category", genreId, "sub", popId, "geo", g.Geo, "kids", kidPrefix)

	obj := new(model.Top200ApplicationsResponse)
	err = json.Unmarshal(data, obj)
	if err != nil {
		return out, fmt.Errorf("Error unmarshal response top200 charts request iTunes, error: %s", err.Error())
	}

	index := obj.PageData.SegmentedControl.SelectedIndex
	ListAdamIds := obj.PageData.SegmentedControl.Segments[index].PageData.SelectedChart.AdamIds
	if chunk <= 0 {
		chunk = obj.Properties.DI6TopChartsPageNumIdsPerChart
	}

	fromToChunk := from + chunk - 1
	if len(ListAdamIds) < fromToChunk {
		fromToChunk = len(ListAdamIds)
	}

	topResults := obj.StorePlatformData.Lockup.Results
	infoResults := make(map[string]model.AppItem)
	var needGetInfo []string

	for i := from - 1; i < fromToChunk; i++ {
		_, ok := topResults[ListAdamIds[i]]
		if !ok {
			needGetInfo = append(needGetInfo, ListAdamIds[i])
		}
	}

	for len(needGetInfo) > 0 {
		var a []string

		l := 50
		if len(needGetInfo) < l {
			l = len(needGetInfo)
		}

		a, needGetInfo = needGetInfo[0:l], needGetInfo[l:]

		info, err := g.GetApplicationInfoByAdamId(a)
		if err != nil {
			return nil, err
		}

		for adamId, appInfo := range info {
			infoResults[adamId] = appInfo
		}
	}

	for i := from - 1; i < fromToChunk; i++ {
		position := i + 1
		adamId := ListAdamIds[i]
		bundleId := ""
		rating := float64(0)
		price := float64(0)
		currencyLabel := ""
		versionCode := 0
		version := ""

		if appInfo, ok := topResults[adamId]; ok {
			bundleId = appInfo.BundleID
			rating = appInfo.UserRating.Value
			price = appInfo.Offers[0].Price
			currencyLabel = helper.GetCurrency(price, appInfo.Offers[0].PriceFormatted)
			versionCode = appInfo.Offers[0].Version.ExternalID
			version = appInfo.Offers[0].Version.Display
		} else if appInfo, ok := infoResults[adamId]; ok {
			bundleId = appInfo.BundleID
			rating = appInfo.UserRating.Value
			price = appInfo.Offers[0].Price
			currencyLabel = helper.GetCurrency(price, appInfo.Offers[0].PriceFormatted)
			versionCode = appInfo.Offers[0].Version.ExternalID
			version = appInfo.Offers[0].Version.Display
		}

		out = append(out, model.TopAppItemResponse{
			AdamId:        adamId,
			BundleId:      bundleId,
			Position:      position,
			Rating:        rating,
			Price:         price,
			CurrencyLabel: currencyLabel,
			VersionCode:   int64(versionCode),
			Version:       version,
		})
	}

	return out, nil
}

func (g *GOiTunes) GetTop1500Applications(
	genreId string, popId string, page int, chunk int,
) ([]model.TopAppItemResponse, error) {
	// page - starts from 0
	// chunk - num of apps in every page
	// app index starts from 0

	var out []model.TopAppItemResponse

	//check popId. Only free for a while
	if !store.PopIdExists(popId) {
		return out, errors.New(fmt.Sprintf("PopId `%s` not exists", popId))
	}

	//genreId=36&popId=27&pageNumbers=0&pageSize=1&cc=us
	q := url.Values{}
	q.Add("genreId", genreId)
	q.Add("popId", popId)
	q.Add("pageNumbers", fmt.Sprintf("%d", page))
	q.Add("pageSize", fmt.Sprintf("%d", chunk))
	q.Add("cc", strings.ToLower(g.Geo))

	u := fmt.Sprintf("%s?%s", TopAppsURL, q.Encode())

	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return out, fmt.Errorf("Error create request top1500 charts iTunes, error: %s", err.Error())
	}

	//set headers
	req.Header.Add("User-Agent", TopUserAgent)

	//exec request
	resp, err := g.httpClient.Do(req)
	if err != nil {
		return out, fmt.Errorf("Error send request top1500 charts iTunes, error: %s", err.Error())
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return out, fmt.Errorf("Error read response body top1500 charts request iTunes, error: %s", err.Error())
	}

	obj := make([]*model.Top1500ApplicationsResponse, 0)
	err = json.Unmarshal(data, &obj)
	if err != nil {
		return out, fmt.Errorf("Error unmarshal response top1500 charts request iTunes, error: %s", err.Error())
	}

	if len(obj) != 1 {
		return out, fmt.Errorf("Oh FUCK, that's strange, error: %s", err.Error())
	}

	for i, item := range obj[0].ContentData {
		position := page*chunk + i + 1
		rating, _ := strconv.ParseFloat(item.UserRating, 64)

		buyParams, _ := url.ParseQuery(item.BuyData.ActionParams)
		paramPrice, _ := strconv.ParseInt(buyParams.Get("price"), 10, 64)

		currencyLabel := helper.GetCurrency(float64(paramPrice), item.ButtonText)
		price := float64(paramPrice) / 1000

		versionCode, err := strconv.ParseInt(item.BuyData.VersionID, 10, 64)
		if err != nil {
			log.Error("goitunes lib: Error parse version code, error: " + err.Error())
		}

		out = append(out, model.TopAppItemResponse{
			AdamId:        item.ID,
			BundleId:      item.BuyData.BundleID,
			Position:      position,
			Rating:        rating,
			Price:         price,
			CurrencyLabel: currencyLabel,
			VersionCode:   versionCode,
		})
	}

	return out, nil
}

func (g *GOiTunes) GetApplicationInfoByAdamId(adamIdList []string) (map[string]model.AppItem, error) {
	var appItemList map[string]model.AppItem

	query := url.Values{
		"version":  []string{"2"},
		"id":       []string{strings.Join(adamIdList, ",")},
		"p":        []string{"mdm-lockup"},
		"caller":   []string{"MDM"},
		"platform": []string{"itunes"},
		"cc":       []string{g.Geo},
		"l":        []string{"en_us"},
	}
	req, err := http.NewRequest("GET", AppInfoURL+"?"+query.Encode(), nil)
	if err != nil {
		return appItemList, fmt.Errorf("Error create request application info iTunes, error: %s", err.Error())
	}

	//exec request
	resp, err := g.httpClient.Do(req)
	if err != nil {
		return appItemList, fmt.Errorf("Error send request get application info iTunes, error: %s", err.Error())
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)

	obj := &model.ApplicationInfoByBundleIdResponse{}
	err = json.Unmarshal(data, &obj)
	if err != nil {
		return appItemList, fmt.Errorf("Error unmarshal response application info request iTunes, error: %s", err.Error())
	}

	if len(obj.Results) == 0 {
		return appItemList, fmt.Errorf(
			"This adamIds: %s was not found in iTunes",
			strings.Join(adamIdList, ","),
		)
	}

	return obj.Results, nil
}

func (g *GOiTunes) GetApplicationInfoByBundleId(bundleIdList []string) (map[string]model.AppItem, error) {
	out := map[string]model.AppItem{}

	query := url.Values{
		"version":  []string{"2"},
		"bundleId": []string{strings.Join(bundleIdList, ",")},
		"p":        []string{"mdm-lockup"},
		"caller":   []string{"MDM"},
		"platform": []string{"itunes"},
		"cc":       []string{g.Geo},
		"l":        []string{"en_us"},
	}
	req, err := http.NewRequest("GET", AppInfoURL+"?"+query.Encode(), nil)
	if err != nil {
		return out, fmt.Errorf("Error create request application info iTunes, error: %s", err.Error())
	}

	resp, err := g.httpClient.Do(req)
	if err != nil {
		return out, fmt.Errorf("Error request application info iTunes, error: %s", err.Error())
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)

	obj := &model.ApplicationInfoByBundleIdResponse{}
	err = json.Unmarshal(data, &obj)
	if err != nil {
		return out, fmt.Errorf("Error unmarshal response application info iTunes, error: %s", err.Error())
	}

	if len(obj.Results) == 0 {
		return out, fmt.Errorf(
			"This bundleIds: %s was not found in iTunes",
			strings.Join(bundleIdList, ","),
		)
	}

	return obj.Results, nil
}

func (g *GOiTunes) GetFullInfoApplicationByAdamId(adamId string) (*model.AppItemFullResponse, error) {
	xAppleStore := g.getStore()

	req, err := http.NewRequest("GET", fmt.Sprintf(NativeAppInfoURL, adamId), nil)
	if err != nil {
		return nil, fmt.Errorf("Error create request full application info, error: %s", err.Error())
	}
	headers := http.Header{
		"X-Apple-Store-Front": []string{
			fmt.Sprintf(store.XAppleStoreFrontTemplate, xAppleStore.StoreFront, store.IPhoneDeviceCode),
		},
	}
	req.Header = headers

	resp, err := g.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Error request full application info iTunes, error: %s", err.Error())
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)

	obj := &model.AppItemFullResponse{}
	err = json.Unmarshal(data, &obj)
	if err != nil {
		return nil, fmt.Errorf("Error unmarshal response full application info iTunes, error: %s", err.Error())
	}

	_, ok := obj.StorePlatformData.ProductDv.Results[adamId]
	if !ok {
		return nil, fmt.Errorf("Response doesn't contains requested adamId: %s", adamId)
	}

	return obj, nil
}

func (g *GOiTunes) GetApplicationRatingByAdamId(adamId string) (*model.ApplicationRatingResponse, error) {
	xAppleStore := g.getStore()

	req, err := http.NewRequest("GET", fmt.Sprintf(NativeAppRatingInfoURL, adamId), nil)
	if err != nil {
		return nil, fmt.Errorf("Error create request application rating info, error: %s", err.Error())
	}
	headers := http.Header{
		"X-Apple-Store-Front": []string{
			fmt.Sprintf(store.XAppleStoreFrontTemplate, xAppleStore.StoreFront, store.IPhoneDeviceCode),
		},
	}
	req.Header = headers

	resp, err := g.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Error request application rating info iTunes, error: %s", err.Error())
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)

	obj := &model.ApplicationRatingResponse{}
	err = json.Unmarshal(data, &obj)
	if err != nil {
		return nil, fmt.Errorf("Error unmarshal response application rating info iTunes, error: %s", err.Error())
	}

	if fmt.Sprintf("%d", obj.AdamID) != adamId {
		return nil, fmt.Errorf("Response doesn't contains requested adamId: %s", adamId)
	}

	return obj, nil
}

func (g *GOiTunes) GetApplicationOverAllRatingByAdamId(adamId string) (*model.ApplicationOverAllRatingItem, error) {
	req, err := http.NewRequest("GET",
		fmt.Sprintf(OpenAppOverAllRatingInfoURL, adamId, strings.ToLower(g.Geo)), nil)
	if err != nil {
		return nil, fmt.Errorf("Error create request application overall rating info, error: %s", err.Error())
	}

	resp, err := g.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Error request application overall rating info iTunes, error: %s", err.Error())
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)

	obj := &model.ApplicationOverAllRatingResponse{}
	err = json.Unmarshal(data, &obj)
	if err != nil {
		return nil, fmt.Errorf("Error unmarshal response application overall rating info iTunes, error: %s", err.Error())
	}

	if len(obj.Results) != 1 {
		return nil, fmt.Errorf("Response doesn't contains requested adamId: %s", adamId)
	}

	return obj.Results[0], nil
}

//Login and Request with authorization

func (g *GOiTunes) Login(password string) (*model.AuthResponse, error) {
	if len(password) == 0 {
		return nil, errors.New("Password can not be empty")
	}

	req, err := http.NewRequest("POST", g.getLoginUrl(),
		helper.GenerateLoginBody(password, g.MachineName, g.GUID, g.AppleId))
	if err != nil {
		return nil, fmt.Errorf("Error create request authorization iTunes, error: %s", err.Error())
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("User-Agent", g.UserAgent)
	req.Header.Add("X-Apple-Store-Front", g.getStore().XAppleStoreFront)

	resp, err := g.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Error send request authorization iTunes, error: %s", err.Error())
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)

	obj := &model.AuthResponse{}
	err = plist.Unmarshal(data, &obj)
	if err != nil {
		return nil, fmt.Errorf("Error unmarshal authorization iTunes response, error: %s", err.Error())
	}

	var errorText string
	if len(obj.PasswordToken) == 0 {
		errorText = "passwordToken not found in response"
	} else if len(obj.DSID) == 0 {
		errorText = "DSID not found in response"
	}

	if len(errorText) > 0 {
		return nil, errors.New("Auth error, " + errorText)
	}

	g.XToken = obj.PasswordToken
	g.DSID = obj.DSID

	return obj, nil
}

func (g *GOiTunes) BuyApplication(adamId string, extVersionId int) (*model.IPA, error) {
	obj, err := g.buy(adamId, extVersionId, helper.PricingParametersBuy)
	if err != nil {
		return nil, err
	}

	if obj.Metrics.DialogId == "MZCommerceSoftware.OwnsSupersededMinorSoftwareApplicationForUpdate" {
		// пока сила Господа не снизайдет на раба его Алексея, и он не найдет способ генерить kbsync
		// хер нам а не скачка при этом кейсе
		// на метод STDRDL нужен другой kbsync, а мы в базе храним один и только для STDQ

		//obj, err = g.buy(adamId, extVersionId, pricingParametersReDownload)
		//if err != nil {
		//	return nil, err
		//}

		return nil, fmt.Errorf("iTunes says, use STDRDL, but we stil can not =(")
	}

	if len(obj.SongPlist) == 0 {
		return nil, fmt.Errorf("Download url not found in buy application reaponse, adamId: %s, account: %s",
			adamId, g.AppleId)
	} else if len(obj.SongPlist) > 1 {
		return nil, fmt.Errorf("iTunes application reaponse contains more than 1 download url, adamId: %s, account: %s",
			adamId, g.AppleId)
	}

	err = g.confirmDownload(obj.SongPlist[0].DownloadId)
	if err != nil {
		return nil, fmt.Errorf("Error confirm download app form iTunes, error: %s", err.Error())
	}

	if len(obj.SongPlist[0].Sinfs) == 0 {
		return nil, fmt.Errorf("No sinf found for application, adamId: %s, account: %s",
			adamId, g.AppleId)
	} else if len(obj.SongPlist[0].Sinfs) > 1 {
		return nil, fmt.Errorf("iTunes application reaponse contains more than 1 sinf, adamId: %s, account: %s",
			adamId, g.AppleId)
	}

	sinf := bytes.TrimSpace(obj.SongPlist[0].Sinfs[0].Sinf)
	if len(sinf) == 0 {
		return nil, fmt.Errorf("Sinf is empty, adamId: %s, account: %s",
			adamId, g.AppleId)
	}

	//FUCKING APPLE!!!
	//BundleId may be in field - BundleId(`softwareVersionBundleId`) or Q(`q`)
	var bundleId string
	if len(obj.SongPlist[0].MetaData.BundleId) != 0 {
		bundleId = obj.SongPlist[0].MetaData.BundleId
	} else {
		bundleId = obj.SongPlist[0].MetaData.Q
	}

	//iTunesMetadata.plist
	iTunesMetadata := helper.GenerateMetadataPlist(bundleId, obj.SongPlist[0],
		g.getStore().XAppleStoreFront, g.getStore().StoreFront.Int())

	return &model.IPA{
		Url: obj.SongPlist[0].Url,
		Headers: map[string]string{
			"User-Agent":          DownloadUserAgent,
			"Cookie":              "downloadKey=" + obj.SongPlist[0].DownloadKey,
			"X-Apple-Store-Front": g.getStore().XAppleStoreFront,
			"X-Dsid":              g.DSID,
		},
		Sinf:           base64.StdEncoding.EncodeToString(sinf),
		ITunesMetadata: base64.StdEncoding.EncodeToString(iTunesMetadata),
		BundleId:       bundleId,
	}, nil
}

//private methods
func (g *GOiTunes) buy(adamId string, extVersionId int, pricing string) (*model.BuyProductResponse, error) {
	req, err := http.NewRequest("POST", g.getBuyProductUrl(),
		helper.GenerateBuyProductBody(adamId, extVersionId, pricing, g.GUID, g.Kbsync, g.MachineName))
	if err != nil {
		return nil, fmt.Errorf("Error create request buyProduct iTunes, error: %s", err.Error())
	}

	req.Header.Add("Content-Type", "application/x-apple-plist")
	req.Header.Add("Referer", fmt.Sprintf("http://itunes.apple.com/app/id%s", adamId))
	req.Header.Add("User-Agent", g.UserAgent)
	req.Header.Add("X-Apple-Store-Front", g.getStore().XAppleStoreFront)
	req.Header.Add("X-Apple-Tz", "10800")
	req.Header.Add("X-Dsid", g.DSID)
	req.Header.Add("X-Token", g.XToken)

	resp, err := g.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Error request buyProduct iTunes, error: %s", err.Error())
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)

	obj := &model.BuyProductResponse{}
	err = plist.Unmarshal(data, &obj)
	if err != nil {
		return nil, fmt.Errorf("Error unmarshal buyProduct iTunes response, error: %s", err.Error())
	}

	return obj, nil
}

func (g *GOiTunes) confirmDownload(downloadId string) error {
	query := url.Values{
		"download-id": []string{downloadId},
		"guid":        []string{g.GUID},
	}

	requestUrl := fmt.Sprintf(ConfirmDownload, g.getStore().HostPrefix, query.Encode())
	req, err := http.NewRequest("GET", requestUrl, nil)
	if err != nil {
		return fmt.Errorf("Error create request buyProduct iTunes, error: %s", err.Error())
	}

	req.Header.Add("User-Agent", DownloadUserAgent)
	req.Header.Add("X-Apple-Store-Front", g.getStore().XAppleStoreFront)
	req.Header.Add("X-Dsid", g.DSID)
	req.Header.Add("X-Token", g.XToken)

	resp, err := g.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("Error request confirmDownload iTunes, error: %s", err.Error())
	}
	defer resp.Body.Close()

	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("Error read response from confrim download request, error: %s", err.Error())
	}

	return nil
}
