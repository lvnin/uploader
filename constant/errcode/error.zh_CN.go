package errcode

var ErrCode_zh_CN = map[string]CodeResult{
	"ERRCODE_UPLOAD_SUCCESS":             {Code: 0, Message: "上传成功"},
	"ERRCODE_DELETE_SUCCESS":             {Code: 0, Message: "删除成功"},
	"ERRCODE_COVERT_SUCCESS":             {Code: 0, Message: "转换成功"},
	"ERRCODE_PREVIEW_SUCCESS":            {Code: 0, Message: "预览成功"},
	"ERRCODE_GENERATE_SUCCESS":           {Code: 0, Message: "生成成功"},
	"ERRCODE_TOKEN_INVALID":              {Code: 10001, Message: "无效令牌"},
	"ERRCODE_TOKEN_MALFORMED":            {Code: 10002, Message: "令牌不完整"},
	"ERRCODE_TOKEN_NOTVALIDYET":          {Code: 10003, Message: "令牌未生效"},
	"ERRCODE_TOKEN_EXPIRED":              {Code: 10004, Message: "令牌过期"},
	"ERRCODE_PARAMETER_FAILURE":          {Code: 10010, Message: "参数错误"},
	"ERRCODE_FILE_DELETE_FAILURE":        {Code: 20001, Message: "文件删除失败"},
	"ERRCODE_FILE_RECEIVE_FAILURE":       {Code: 20002, Message: "接收文件失败"},
	"ERRCODE_FILE_UPLOAD_FAILURE":        {Code: 20003, Message: "上传文件失败"},
	"ERRCODE_FILE_UPLOAD_NOT_COMPLIANCE": {Code: 20004, Message: "上传文件有存在违规"},
	"ERRCODE_FILE_COVERT_FAILURE":        {Code: 20005, Message: "文件转换失败"},
	"ERRCODE_FILE_PREVIEW_FAILURE":       {Code: 20006, Message: "文件预览失败"},
	"ERRCODE_FILE_GENERATE_FAILURE":      {Code: 20007, Message: "文件生成失败"},
	"ERRCODE_DBLINK_UPDATE_FAILURE":      {Code: 20101, Message: "修改数据库链接失败"},
	"ERRCODE_GENERATE_CONTRACT_FAILURE":  {Code: 20102, Message: "生成合同失败"},
}
