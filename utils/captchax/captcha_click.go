package captchax

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/go-xuan/quanx/os/errorx"
	"github.com/wenlng/go-captcha/captcha"
)

// NewClickCaptcha 初始化点击行为验证码
func NewClickCaptcha() *ClickCaptcha {
	return &ClickCaptcha{
		capt:  getClickCaptcha(),
		store: DefaultStore(),
	}
}

type ClickCaptcha struct {
	capt  *captcha.Captcha
	store *CaptchaStore
}

func (impl *ClickCaptcha) New(ctx context.Context) (string, string, string, error) {
	if dots, image, thumb, key, err := impl.capt.Generate(); err != nil {
		return "", "", "", errorx.Wrap(err, "generate captcha failed")
	} else {
		expiration := time.Duration(impl.store.expired) * time.Second
		if err = impl.store.client.Set(ctx, key, dots, expiration); err != nil {
			return "", "", "", errorx.Wrap(err, "store captcha failed")
		}
		return key, image, thumb, nil
	}
}

func (impl *ClickCaptcha) Verify(ctx context.Context, key, answer string) bool {
	var dots = make(map[int]captcha.CharDot)
	dotsCache := impl.store.client.GetString(ctx, key)
	if err := json.Unmarshal([]byte(dotsCache), &dots); err != nil {
		return false
	}
	var ok, points = false, strings.Split(answer, ",")
	if (len(dots) * 2) == len(points) {
		for i, dot := range dots {
			x, _ := strconv.ParseFloat(points[i*2], 64)
			y, _ := strconv.ParseFloat(points[i*2+1], 64)
			// 校验点的位置,在原有的区域上添加额外边距进行扩张计算区域,不推荐设置过大的padding
			// 例如：文本的宽和高为30，校验范围x为10-40，y为15-45，此时扩充5像素后校验范围宽和高为40，则校验范围x为5-45，位置y为10-50
			if ok = captcha.CheckPointDistWithPadding(int64(x), int64(y), int64(dot.Dx), int64(dot.Dy), int64(dot.Width), int64(dot.Height), 5); !ok {
				break
			}
		}
	}
	if ok {
		impl.store.client.Delete(ctx, key)
	}
	return ok
}

func getClickCaptcha() *captcha.Captcha {
	capt := captcha.GetCaptcha()
	// ========================主图配置============================
	// 设置验证码主图的尺寸
	//capt.SetImageSize(captcha2.Size{Width: 300, Height: 300})
	// 设置验证码主图清晰度，压缩级别范围 QualityCompressLevel1 - 5，QualityCompressNone无压缩，默认为最低压缩级别
	//capt.SetImageQuality(captcha.QualityCompressNone)
	// 设置字体Hinting值 (HintingNone,HintingVertical,HintingFull)
	//capt.SetFontHinting(font.HintingFull)
	// 设置验证码文本显示的总数随机范围
	//capt.SetTextRangLen(captcha.RangeVal{Min: 6, Max: 7})
	// 设置验证码文本的随机字体大小
	//capt.SetRangFontSize(captcha.RangeVal{Min: 32, Max: 42})
	// 设置验证码文本的随机十六进制颜色
	//capt.SetTextRangFontColors([]string{"#1d3f84", "#3a6a1e"})
	// 设置验证码字体的透明度
	//capt.SetImageFontAlpha(0.5)
	// 设置字体阴影
	//capt.SetTextShadow(true)
	// 设置字体阴影颜色
	//capt.SetTextShadowColor("#101010")
	// 设置字体阴影偏移位置
	//capt.SetTextShadowPoint(captcha.Point{X: 1, Y: 1})
	// 设置验证码字体的扭曲程度
	//capt.SetImageFontDistort(captcha.DistortLevel2)

	// ========================缩略图配置============================
	// 设置缩略图的尺寸
	//capt.SetThumbSize(captcha.Size{Width: 150, Height: 40})
	// 设置缩略图校验文本的随机长度范围
	//capt.SetRangCheckTextLen(captcha.RangeVal{Min: 2, Max: 4})
	// 设置缩略图校验文本的随机大小
	//capt.SetRangCheckFontSize(captcha.RangeVal{Min: 24, Max: 30})
	// 设置缩略图文本的随机十六进制颜色
	//capt.SetThumbTextRangFontColors([]string{"#1d3f84", "#3a6a1e"})
	//  设置缩略图的背景随机十六进制颜色
	//capt.SetThumbBgColors([]string{"#1d3f84", "#3a6a1e"})
	// 设置缩略图背景的扭曲程度
	//capt.SetThumbBgDistort(captcha.DistortLevel2)
	// 设置缩略图字体的扭曲程度
	//capt.SetThumbFontDistort(captcha.DistortLevel2)
	// 设置缩略图背景的圈点数
	//capt.SetThumbBgCirclesNum(20)
	// 设置缩略图背景的线条数
	//capt.SetThumbBgSlimLineNum(3)
	return capt
}
