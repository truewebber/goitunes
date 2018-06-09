package goitunes

import (
	"fmt"
	"testing"

	"github.com/truewebber/goitunes/model"
	"github.com/truewebber/store"
)

var (
	goitunes *GOiTunes
	AuthResp *model.AuthResponse

	kbsync = `AAQAA5nfBFJ2pQ2dGNcv4hSmXBcv2FaCiqIO6NM5sXeGnZOyD82E1lxsvvhHihO6Xc5r
	CvW/7AMmKJeW3EnkBU6EaxHEYHHsAEu1g3+mTTGxroaKuQfRLsGz+UUImTDlvBr0xhwm
	xsr7+t8o4xIqeRhZZzlBaWjeGRFRPQ5JmJNlqyyH2ByxhyFwqhJhlQsiPK5qt0iymdP2
	KFm0L1b+FvDexd9pDbYGIXDNoqFZrgxtUQkdd0RfgPdZi5QB2DWePkrUnoBF02OEqg81
	BCnh8EWPA+VTCkEbF9mi58rXvXQXZlR66vFmS7o/MdmjrFrRdQEV10ScVTwQbLls26EZ
	FW77Lql5tPSCwvkVRqREOCGiYOdyrTjvH7b69G3NUg3ZyVwZz8I72tioGNHVzsEaEtvC
	WcID3ZD/i//7BoUuwYTwausTLr3YHh6HP2IYKXhZe116ZS8wf6mn1sfFSW/vxEHfqIJN
	xnVp/VMtbUJhCN2k94LiEk89IKEcdcvX/s9S8UaSHdPJH/qCEcKdUBoM5wysDsQsORI3
	C0DtRbU3BSiqZtIt/AGVtiQGnZBoVNiHjgNO5cn3Z1NHva0w/VBO7He0gE69By2PA2AJ
	BLIjVoiZj62w3GuCqbFEO2y6xmYnK34EzK75J3SUv6Tpgxp3rs4R4z4p1/lGuqoapRbe
	Np2SVjoSSxXI0wCoOZzKe0TFgSYp0vHboz+vpxxK7565QWc47C5H3S8ixXJ0xJnTbIjv
	R1nrEJ95DbNxPpnW5jgZRT/s71nkNUoDBlH6S1byKHClrIaKuMLmYTotFywkwt4P9lup`
)

func init() {
	var err error
	goitunes, err = NewGOiTunes(
		"truewebber@apptica.com",
		"",
		kbsync,
		"",
		"RU",
		"9801A7A4ED7B",
		"iTunes/10.6 (Windows; Microsoft Windows 7 x64 Ultimate Edition Service Pack 1 (Build 7601)) AppleWebKit/534.54.16",
		"m83",
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
	AuthResp, err = goitunes.Login("Apptica41")
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
	AuthResp, err = goitunes.Login("Apptica41")
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
