package service

import (
	"errors"
	"fmt"
	"github.com/sean-tech/webservice/logging"
	"gopkg.in/go-playground/validator.v9"
	"reflect"
	"regexp"
	"sync"
)

type Tag string
type Pattern string

const (
	// 手机号
	validation_tag_phone 	= "phone"
	tag_pattern_phone 		= "^1(3[0-9]|4[5,7]|5[0,1,2,3,5,6,7,8,9]|6[2,5,6,7]|7[0,1,7,8]|8[0-9]|9[1,8,9])\\d{8}$"
	// 邮件或者手机号
	validation_tag_eorp 	= "eorp"
	tag_pattern_eorp 		= "^\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*$|^1(3[0-9]|4[5,7]|5[0,1,2,3,5,6,7,8,9]|6[2,5,6,7]|7[0,1,7,8]|8[0-9]|9[1,8,9])\\d{8}$"
	// 身份证号
	validation_tag_idcard 	= "idcard"
	tag_pattern_idcard 		= "^\\d{15}$)|(^\\d{18}$)|(^\\d{17}(\\d|X|x)$"
	// qq号
	validation_tag_qq 		= "qq"
	tag_pattern_qq 			= "[1-9][0-9]{4,}"
	// 账号
	validation_tag_account 	= "account"
	tag_pattern_account 	= "^[a-zA-Z][a-zA-Z0-9_]{4,15}$"
	// 密码 - 以字母开头，长度在6~18之间，只能包含字母、数字和下划线
	validation_tag_pwd1 	= "pwd1"
	tag_pattern_pwd1	 	= "^[a-zA-Z]\\w{5,17}$"
	// 密码 - 须包含大小写字母和数字的组合，不能使用特殊字符，长度在 6-18 之间
	validation_tag_pwd2 	= "pwd2"
	tag_pattern_pwd2	 	= "^(?=.*\\d)(?=.*[a-z])(?=.*[A-Z])[a-zA-Z0-9]{6,18}$"
	// 密码 - 须包含大小写字母和数字的组合，可以使用特殊字符，长度在 6-18 之间
	validation_tag_pwd3 	= "pwd3"
	tag_pattern_pwd3	 	= "^(?=.*\\d)(?=.*[a-z])(?=.*[A-Z]).{8,18}$"
)

func init() {
	ValidationTagRegexpPatternRegister(map[Tag]Pattern{
		validation_tag_phone 	: tag_pattern_phone,
		validation_tag_eorp 	: tag_pattern_eorp,
		validation_tag_idcard 	: tag_pattern_idcard,
		validation_tag_qq 		: tag_pattern_qq,
		validation_tag_account 	: tag_pattern_account,
		validation_tag_pwd1 	: tag_pattern_pwd1,
		validation_tag_pwd2 	: tag_pattern_pwd2,
		validation_tag_pwd3 	: tag_pattern_pwd3,
	})
}

var _tagPatternMap = sync.Map{}

/**
 * 给验证注册tag和对应的正则
 */
func ValidationTagRegexpPatternRegister(tagPatternMap map[Tag]Pattern)  {
	for tag, pattern := range tagPatternMap {
		_tagPatternMap.Store(string(tag), string(pattern))
	}
}

/**
 * 参数绑定验证
 */
func ValidateParameter(parameter interface{}) error {
	validate := validator.New()

	// register validation func
	_tagPatternMap.Range(func(tag, pattern interface{}) bool {
		err := validate.RegisterValidation(tag.(string), func(fl validator.FieldLevel) bool {
			ok, err := regexp.MatchString(pattern.(string), fl.Field().String())
			if err != nil {
				return false
			}
			return ok
		})
		if err != nil {
			logging.Warning(err)
		}
		return true
	})

	// validate
	err := validate.Struct(parameter)
	if err == nil {
		return nil
	}
	if _, ok := err.(*validator.InvalidValidationError); ok {
		logging.Warning(err)
		return errors.New(STATUS_MSG_INVALID_PARAMS)
	}
	for _, err := range err.(validator.ValidationErrors) {
		info := fmt.Sprintf("validate err.Namespace:%s Field:%s StructNamespace:%s StructField:%s Tag:%s ActualTag:%s Kind:%v Type:%v Value:%v Param:%s",
			err.Namespace(), err.Field(), err.StructNamespace(), err.StructField(), err.Tag(), err.ActualTag(), err.Kind(), err.Type(), err.Value(), err.Param(),
		)
		logging.Warning(info)
		return errors.New(STATUS_MSG_INVALID_PARAMS)
	}
	return nil
}

/**
 * 正则Tag验证
 */
func ValidateRegexpTagParameter(parameter interface{}) error {
	// regexp
	s := reflect.TypeOf(parameter).Elem()
	v := reflect.ValueOf(parameter).Elem()
	for i := 0; i < s.NumField(); i++ {
		pattern := s.Field(i).Tag.Get("regexp")
		if len(pattern) <= 0 {
			continue
		}
		filedName := s.Field(i).Name
		ok, err := regexp.MatchString(pattern, v.FieldByName(filedName).String())
		if err != nil || !ok {
			logging.Warning(err.Error())
			return errors.New(STATUS_MSG_INVALID_PARAMS)
		}
	}
	return nil
}

/**
 * 验证是否是email
 */
func ValidateEmail(str string) bool {
	emailPattern := "^\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*$"
	ok, err := regexp.MatchString(emailPattern, str)
	if err != nil {
		return false
	}
	return ok
}

/**
 * 验证是否是phone
 */
func ValidatePhone(str string) bool {
	phonePattern := "^1(3[0-9]|4[5,7]|5[0,1,2,3,5,6,7,8,9]|6[2,5,6,7]|7[0,1,7,8]|8[0-9]|9[1,8,9])\\d{8}$"
	ok, err := regexp.MatchString(phonePattern, str)
	if err != nil {
		return false
	}
	return ok
}