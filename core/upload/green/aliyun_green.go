package green

import (
	"encoding/json"
	"path"
	"strconv"
	"strings"
	"uploader/global"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/green"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
)

// 审核场景
const (
	AliyunGreenScenePorn      = "porn"      // 智能鉴黄
	AliyunGreenSceneTerrorism = "terrorism" //暴恐涉政
	AliyunGreenSceneAd        = "ad"        // 图文
	AliyunGreenSceneQrcode    = "qrcode"    // 二维码图片
	AliyunGreenSceneLive      = "live"      // 不良场景
	AliyunGreenSceneLogo      = "logo"      // 商标
)

type AliyunGreen struct{}

func (*AliyunGreen) GreenErrorImageDefault() string {
	return global.ServerConfig.AliyunOSS.Bucket.BucketUrl + "/" + global.ServerConfig.AliyunOSS.Green.ErrorImagePath
}

func (*AliyunGreen) GreenImage(url string) bool {
	if !checkGreenEnabled() || !isImage(url) {
		return true
	}

	client, err := NewGreenClient()
	if err != nil {
		global.Logger.Error("aliyun green new client is failed", zap.Any("err", err.Error()))
		return false
	}

	var tasks []interface{}

	tasks = append(tasks, map[string]interface{}{"dataId": uuid.NewV4(), "url": url})

	content, _ := json.Marshal(
		map[string]interface{}{
			"tasks":  tasks,
			"scenes": global.ServerConfig.AliyunOSS.Green.Scenes,
		},
	)

	request := green.CreateImageSyncScanRequest()
	request.SetContent(content)
	response, _err := client.ImageSyncScan(request)
	if _err != nil {
		global.Logger.Error(_err.Error())
		return false
	}
	if response.GetHttpStatus() != 200 {
		global.Logger.Error("image green response failed status:" + strconv.Itoa(response.GetHttpStatus()))
		return false
	}

	return scoreGreenScan(response.GetHttpContentBytes())
}

func scoreGreenScan(data []byte) bool {
	var result struct {
		Code int `json:"code"`
		Data []struct {
			Code    int         `json:"code"`
			DataID  string      `json:"dataId"`
			Extras  interface{} `json:"extras"`
			Msg     string      `json:"msg"`
			Results []struct {
				Label      string  `json:"label"`
				Rate       float64 `json:"rate"`
				Scene      string  `json:"scene"`
				Suggestion string  `json:"suggestion"`
			} `json:"results"`
			TaskID string `json:"taskId"`
			Url    string `json:"url"`
		} `json:"data"`
		Msg       string `json:"msg"`
		RequestID string `json:"requestId"`
	}
	err := json.Unmarshal(data, &result)
	if err != nil {

		return false
	}

	for _, v := range result.Data {
		for _, res := range v.Results {
			if res.Suggestion == "block" {
				return false
			}
		}
	}

	return true
}

// isImage - 是否为图片
// @param {string} url
// @returns bool
func isImage(url string) bool {
	ext := strings.ToLower(path.Ext(url))
	return strings.Contains(ext, "png") ||
		strings.Contains(ext, "jpg") ||
		strings.Contains(ext, "jpeg") ||
		strings.Contains(ext, "tiff") ||
		strings.Contains(ext, "bmp") ||
		strings.Contains(ext, "gif") ||
		strings.Contains(ext, "webp") ||
		strings.Contains(ext, "svg")
}

func checkGreenEnabled() bool {
	return global.ServerConfig.AliyunOSS.Green.Score > 0.0
}

func NewGreenClient() (client *green.Client, err error) {
	return green.NewClientWithAccessKey(global.ServerConfig.AliyunOSS.Green.Region,
		global.ServerConfig.AliyunOSS.Green.AccessKey,
		global.ServerConfig.AliyunOSS.Green.AccessKeySecret)
}

/*
// 是否启用图片审核
func (service *greenPlusService) greenEnabled() bool {
	service.init()
	return global.ServerConfig.AliyunOSS.Green.Score > 0.0 && service.client != nil
}

// 初期化图片审核
func (service *greenPlusService) init() bool {
	if global.ServerConfig.AliyunOSS.Green.Score > 0.0 || service.client != nil {
		return true
	}

	service.config = &openapi.Config{
		// 您的AccessKey ID。
		AccessKeyId: tea.String(global.ServerConfig.AliyunOSS.Green.AccessKey),
		// 您的AccessKey Secret。
		AccessKeySecret: tea.String(global.ServerConfig.AliyunOSS.Green.AccessKeySecret),
		// 设置HTTP代理。
		// HttpProxy: tea.String("http://xx.xx.xx.xx:xxxx"),
		// 设置HTTPS代理。
		// HttpsProxy: tea.String("https://username:password@xxx.xxx.xxx.xxx:9999"),
		RegionId: tea.String(global.ServerConfig.AliyunOSS.Green.Region),
		Endpoint: tea.String(global.ServerConfig.AliyunOSS.Green.Endpoint),
		//请设置超时时间。服务端全链路处理超时时间为10秒，请做相应设置。
		//如果您设置的ReadTimeout小于服务端处理的时间，程序中会获得一个ReadTimeout异常。
     ConnectTimeout: tea.Int(3000),
     ReadTimeout:    tea.Int(6000),
   }

   var err error
   service.client, err = green20220302.NewClient(service.config)
   if err != nil {
     global.Logger.Error(err.Error())
     return false
   }
   return true
 }

 // 获取图片审核运行参数
 func (service *greenPlusService) getRuntimeOptions() *util.RuntimeOptions {
   // 创建RuntimeObject实例并设置运行参数。
   runtime := &util.RuntimeOptions{}
   runtime.ReadTimeout = tea.Int(global.ServerConfig.AliyunOSS.Green.ReadTimeout)
   runtime.ConnectTimeout = tea.Int(global.ServerConfig.AliyunOSS.Green.ConnectTimeout)
   return runtime
 }

 // 图片审核接口
 func (service *greenPlusService) ImageGreenPlus(sources []model.UploadModel, noGreen string) bool {
   if noGreen == "1" {
     return true
   }

   if !service.greenEnabled() {
     return false
   }

   for _, source := range sources {
     if ok := service.checkImageGreenPlus(service.getImageModerationRequest(source)); !ok {
       return false
     }
   }

   return true
 }

 // 获取图片审核请求
 func (service *greenPlusService) getImageModerationRequest(source model.UploadModel) *green20220302.ImageModerationRequest {
   serviceParameters, _ := json.Marshal(
     map[string]interface{}{
       "imageUrl": source.Url,
       "dataId":   uuid.NewV4(),
     },
   )
   imageModerationRequest := &green20220302.ImageModerationRequest{
     //图片检测service:通用基线检测
     Service:           tea.String(global.ServerConfig.AliyunOSS.Green.Service),
     ServiceParameters: tea.String(string(serviceParameters)),
   }

   return imageModerationRequest
 }

 // 图片审核验证
 func (service *greenPlusService) checkImageGreenPlus(request *green20220302.ImageModerationRequest) bool {
   runtime := service.getRuntimeOptions()
   response, _err := service.client.ImageModerationWithOptions(request, runtime)
   //自动路由，服务端错误，区域切换至cn-beijing
   flag := false
   if _err != nil {
     var err = &tea.SDKError{}
     if _t, ok := _err.(*tea.SDKError); ok {
       err = _t
       // 系统异常，切换到下个地域调用。
       if *err.StatusCode == 500 {
         flag = true
       }
     }
   }

   if response == nil || *response.StatusCode == 500 || *response.Body.Code == 500 {
     flag = true
   }
   if flag {
     service.config.SetRegionId(global.ServerConfig.AliyunOSS.Green.SpareRegion)
     service.config.SetEndpoint(global.ServerConfig.AliyunOSS.Green.SpareEndpoint)
     service.client, _err = green20220302.NewClient(service.config)
     if _err != nil {
       return false
     }
     response, _err = service.client.ImageModerationWithOptions(request, runtime)
     if _err != nil {
       global.Logger.Error(_err.Error())
       return false
     }
   }

   if response != nil {
     statusCode := tea.IntValue(tea.ToInt(response.StatusCode))
     body := response.Body
     imageModerationResponseData := body.Data
     fmt.Println("requestId:" + tea.StringValue(body.RequestId))
     if statusCode == http.StatusOK {
       fmt.Println("response success. response:" + body.String())
       if tea.IntValue(tea.ToInt(body.Code)) == 200 {
         result := imageModerationResponseData.Result
         global.Logger.Info("response dataId:" + tea.StringValue(imageModerationResponseData.DataId))
         for i := 0; i < len(result); i++ {
           var score = tea.Float32Value(result[i].Confidence)
           if ok := service.imageConfidenceRuler(tea.StringValue(result[i].Label), score); !ok {
             return false
           }
         }
       } else {
         global.Logger.Info("image moderation not success. status" + tea.ToString(body.Code))
         return false
       }
     } else {
       global.Logger.Info("response not success. status:" + tea.ToString(statusCode))
       return false
     }
   }

   return true
 }

 // 图片审核规则
 func (service *greenPlusService) imageConfidenceRuler(label string, score float32) bool {
   fmt.Println("response label:" + label)
   fmt.Println("response confidence:" + tea.ToString(score))
   return score >= 60.0
 }

*/
