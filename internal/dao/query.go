package dao

import (
	"Crocodile6/internal/consts"
	"Crocodile6/internal/utils"
	"github.com/pkg/errors"
	"log"
	"sync"
)

type QueryDao struct{}

func NewQueryDao() *QueryDao {
	return &QueryDao{}
}

// QueryCenter 查询调度中心
func (s *QueryDao) QueryCenter(typeString string, queryList []string, ref *[]consts.QueryDataDto) error {
	// 判断type和qqList是否正确
	if len(queryList) <= 0 {
		return errors.New("查询数据为空")
	}
	validTypes := map[string]bool{
		"GetQQUidForMobile": true,
		"GetMobileUidForQQ": true,
		"GetQQMobileForUid": true,
	}
	if !validTypes[typeString] {
		return errors.New("查询类型错误")
	}
	log.Println("测试")
	// 查询主键数据
	var (
		// SQL查询结果
		queryResult []consts.QueryDto
		err         error
	)
	// 根据手机号获取QQ和UId
	if typeString == "GetQQUidForMobile" {
		err = consts.DB.Table("weibo.weibo").
			Select("mobile", "uid").
			Where("mobile IN (?)", queryList).
			Scan(&queryResult).Error
		if err != nil {
			return errors.Wrap(err, "查询中心调用SQL1失败")
		}
		log.Println("查询1,Data->", queryResult)
		mobileList := make([]string, 0, len(queryResult))
		for _, p := range queryResult {
			mobileList = append(mobileList, p.Mobile)
		}
		err = consts.DB.Table("8eqq.8eqq").
			Select("username", "mobile").
			Where("username IN (?)", mobileList).
			Scan(&queryResult).Error
		if err != nil {
			return errors.Wrap(err, "查询中心调用SQL2失败")
		}
		log.Println("查询2,Data->", queryResult)
	}
	// 根据QQ获取手机号和UId
	if typeString == "GetMobileUidForQQ" {
		err = consts.DB.Table("8eqq.8eqq").
			Select("username", "mobile").
			Where("username IN (?)", queryList).
			Scan(&queryResult).Error
		if err != nil {
			return errors.Wrap(err, "查询中心调用SQL1失败")
		}
		log.Println("查询1,Data->", queryResult)
		mobileList := make([]string, 0, len(queryResult))
		for _, p := range queryResult {
			mobileList = append(mobileList, p.Mobile)
		}
		err = consts.DB.Table("weibo.weibo").
			Select("mobile", "uid").
			Where("mobile IN (?)", mobileList).
			Scan(&queryResult).Error
		if err != nil {
			return errors.Wrap(err, "查询中心调用SQL2失败")
		}
		log.Println("查询2,Data->", queryResult)

	}
	// 根据uid获取手机号和QQ
	if typeString == "GetQQMobileForUid" {
		err = consts.DB.Table("weibo.weibo").
			Select("mobile", "uid").
			Where("uid IN (?)", queryList).
			Scan(&queryResult).Error
		if err != nil {
			return errors.Wrap(err, "查询中心调用SQL1失败")
		}
		log.Println("查询1,Data->", queryResult)
		mobileList := make([]string, 0, len(queryResult))
		for _, p := range queryResult {
			mobileList = append(mobileList, p.Mobile)
		}
		err = consts.DB.Table("8eqq.8eqq").
			Select("username", "mobile").
			Where("username IN (?)", mobileList).
			Scan(&queryResult).Error
		if err != nil {
			return errors.Wrap(err, "查询中心调用SQL2失败")
		}
		log.Println("查询2,Data->", queryResult)
	}

	indexKB := &sync.Map{}
	wg1 := &sync.WaitGroup{}
	// 获取归属地
	for _, r := range queryResult {
		log.Println("QQ->", r.UserName)
		wg1.Add(1)
		go func(p consts.QueryDto) {
			defer wg1.Done()
			// 使用独立的 goroutine 查询手机号归属地
			address, err := utils.GetPrForMobile(p.Mobile)
			if err != nil {
				log.Println(err)
				return
			}
			p.UserName = utils.SetCustomValue(p.UserName, "空")
			p.Mobile = utils.SetCustomValue(p.Mobile, "空")
			p.Uid = utils.SetCustomValue(p.Uid, "空")
			indexKB.Store(p.Mobile, consts.QueryDataDto{
				Mobile:  p.Mobile,
				Uid:     p.Uid,
				QQ:      p.UserName,
				Address: address,
			})
		}(r)
	}
	wg1.Wait()
	indexKB.Range(func(key, value interface{}) bool {
		*ref = append(*ref, value.(consts.QueryDataDto))
		return true
	})

	return nil
}
