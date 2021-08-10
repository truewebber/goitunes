package goitunes

import (
	"fmt"
	"os"
	"testing"

	"github.com/truewebber/goitunes/model"
	"github.com/truewebber/goitunes/store"
)

var (
	goitunes *GOiTunes
	AuthResp *model.AuthResponse
)

var (
	testKbsync      = os.Getenv("kbsync")
	testPassword    = os.Getenv("password")
	testAppleID     = os.Getenv("appleID")
	testMachineName = os.Getenv("machineName")
	testGuid        = os.Getenv("guid")
	testGeo         = os.Getenv("geo")
)

func init() {
	var err error
	goitunes, err = NewGOiTunes(
		testAppleID,
		"",
		testKbsync,
		"",
		testGeo,
		testGuid,
		"iTunes/10.6 (Windows; Microsoft Windows 7 x64 Ultimate Edition Service Pack 1 (Build 7601)) AppleWebKit/534.54.16",
		testMachineName,
	)

	if err != nil {
		panic(err)
	}
}

func TestGetTopApplications(t *testing.T) {
	genreId := "6005"

	result1500, err := goitunes.GetTop1500Applications(genreId, store.TopFree, 0, 600)
	if err != nil {
		t.Errorf("Error get top1500 apps, %s", err.Error())

		return
	}

	result200, err := goitunes.GetTop200Applications(genreId, store.TopFree, "", 1, 200)
	if err != nil {
		t.Errorf("Error get top200 apps, %s", err.Error())

		return
	}

	indent := 0
	for i := 199; i > 190; i-- {
		if result200[199].AdamId == result1500[i].AdamId {
			indent = 199 - i

			break
		}
	}

	t.Log("----")
	t.Log("indent", indent)
	t.Log("----", "RESULT")

	for i := 195; i < 200; i++ {
		t.Log("0200", result200[i].Position, ":", result200[i].BundleId)
	}

	t.Log("----")

	for i := 200 - indent; i < 205; i++ {
		t.Log("1500", result1500[i].Position+indent, ":", result1500[i].BundleId)
	}
}

func TestGetApplicationInfoByAdamId(t *testing.T) {
	result, err := goitunes.GetApplicationInfoByAdamId([]string{"564177498"})
	if err != nil {
		t.Errorf("Error get info, %s", err.Error())

		return
	}

	t.Log(result)
}

func TestGetApplicationInfoByBundleId(t *testing.T) {
	data, err := goitunes.GetApplicationInfoByBundleId([]string{"com.vk.vkclient", "DisneyDigitalBooks.StoryHub"})
	if err != nil {
		t.Errorf("Can't check applications in iTunes, error: %s", err.Error())

		return
	}

	t.Log(data)
}

func TestLogin(t *testing.T) {
	var err error
	AuthResp, err = goitunes.Login(testPassword)
	if err != nil {
		t.Errorf("Can't login iTunes, error: %s", err.Error())

		return
	}

	t.Log(AuthResp)
}

func TestBuyApplication(t *testing.T) {
	//BUY
	ipa, err := goitunes.BuyApplication("1118882627", 822467210)
	if err != nil {
		t.Errorf("Can't buy application iTunes, error: %s", err.Error())

		return
	}

	t.Log(ipa)
}

func TestDownload(t *testing.T) {
	bundle := "com.zynga.crosswordswithfriends"

	var err error
	AuthResp, err = goitunes.Login(testPassword)
	if err != nil {
		t.Errorf("Can't login iTunes, error: %s", err.Error())

		return
	}

	t.Log(AuthResp)

	info, err := goitunes.GetApplicationInfoByBundleId([]string{bundle})
	if err != nil {
		t.Errorf(err.Error())

		return
	}

	ipa, err := goitunes.BuyApplication(info[bundle].ID, info[bundle].Offers[0].Version.ExternalID)
	if err != nil {
		t.Errorf("Can't buy application iTunes, error: %s", err.Error())

		return
	}

	t.Log("BundleId:", ipa.BundleId)
	t.Log("URL:", ipa.Url)
	t.Log("Headers:", fmt.Sprintf("%v", ipa.Headers))
	t.Log("Sinf:", ipa.Sinf)
}
