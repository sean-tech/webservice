package services

type StatusCode int

const (
	_ StatusCode = 0
	// base
	STATUS_CODE_SUCCESS        = 200
	STATUS_CODE_INVALID_PARAMS = 400
	STATUS_CODE_ERROR          = 500
	STATUS_CODE_FAILED         = 800

	// jwt token
	STATUS_CODE_AUTH_CHECK_TOKEN_FAILED    = 801
	STATUS_CODE_AUTH_CHECK_TOKEN_TIMEOUT   = 802
	STATUS_CODE_AUTH_TOKEN_GENERATE_FAILED = 803
	STATUS_CODE_AUTH_FAILED                = 804

	// upload
	STATUS_CODE_UPLOAD_FILE_SAVE_FAILED        = 811
	STATUS_CODE_UPLOAD_FILE_CHECK_FAILED       = 812
	STATUS_CODE_UPLOAD_FILE_CHECK_FORMAT_WRONG = 813
)

var StatusCodeMsgMap = map[StatusCode]string {
	// base
	STATUS_CODE_SUCCESS:        "ok",
	STATUS_CODE_INVALID_PARAMS: "请求参数错误",
	STATUS_CODE_ERROR:          "system error",
	STATUS_CODE_FAILED:         "操作失败",

	// jwt token
	STATUS_CODE_AUTH_CHECK_TOKEN_FAILED:    "Token鉴权失败",
	STATUS_CODE_AUTH_CHECK_TOKEN_TIMEOUT:   "Token已超时",
	STATUS_CODE_AUTH_TOKEN_GENERATE_FAILED: "Token生成失败",
	STATUS_CODE_AUTH_FAILED:                "Token错误",

	// upload
	STATUS_CODE_UPLOAD_FILE_SAVE_FAILED:        "保存文件失败",
	STATUS_CODE_UPLOAD_FILE_CHECK_FAILED:       "检查文件失败",
	STATUS_CODE_UPLOAD_FILE_CHECK_FORMAT_WRONG: "校验文件错误，文件格式或大小不正确",
}

func (code StatusCode) Msg() string {
	msg, ok := StatusCodeMsgMap[code]
	if ok {
		return msg
	}
	return StatusCodeMsgMap[STATUS_CODE_ERROR]
}