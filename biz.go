func (c *BizDetailController) getOneBiz(v map[string]string) map[string]interface{} {
	biz := make(map[string]interface{})

	// 运营态：-1-已下线（未接入）；0-N/A（初始）；1-预约中；2-已上线（未商业化）；3-商业化（道具售卖）
	status, _ := strconv.Atoi(v["iStatus"])
	djlog.LogDebug("biz", v["sBizCode"], "iStatus:", v["iStatus"], "status:", status)
	biz["iStatus"] = status

	biz["sBizCode"] = v["sBizCode"]
	biz["sBizName"] = v["sBizName"]
	biz["iGameType"] = v["iGameType"]
	biz["iDemand"] = v["iDemand"]
	biz["iGameFriend"] = v["iGameFriend"]
	biz["video"] = 0                // 视频，功能已废弃
	biz["tuan"] = 0                 // 团购，功能已废弃
	biz["iBian"] = v["iBian"]       // 周边
	biz["Abtest"] = 1               // 新旧版活动标识，功能已废弃
	biz["iPayType"] = v["iPayType"] // QQ钱包

	if c.source == "ios" {
		biz["activityId"] = v["iActivityIos"]
	} else {
		biz["activityId"] = v["iActivityAnd"]
	}

	//贵族积分，[]类型
	var vipList []map[string]interface{}
	if c.version < 126 {
		if len(v["sVipPoint"]) > 0 {
			_ = json.Unmarshal([]byte(v["sVipPoint"]), &vipList)
		}
	} else if status == 3 { // 仅商业化的业务读取vip点配置
		var err error
		if vipList, err = c.getVipPointConfig(v["sBizCode"]); err != nil {
			djlog.LogError("getVipPointConfig failed:", err.Error())
		}
	}
	if len(vipList) > 0 {
		biz["vipPoint"] = vipList
	} else {
		biz["vipPoint"] = []interface{}{} // 保持数组格式
	}

	// biz["saleText"] = v["sSaleText"]
	var text map[string]string
	if len(v["sSaleText"]) > 0 {
		json.Unmarshal([]byte(v["sSaleText"]), &text)
	}
	biz["saleText"] = text

	biz["sGameCode"] = v["sGameCode"]
	biz["wxAppid"] = ""
	// 如果是手游(iGameType：0-端游；1-手游)且已上线，则查询appid
	if v["iGameType"] == "1" && status > 1 {
		qqAppID, wxAppID, gopenid, _ := c.getBizAppID(v["sBizCode"])
		biz["qqAppid"] = qqAppID
		biz["wxAppid"] = wxAppID
		if gopenid == 1 && c.version <= 120 {
			gopenid = 2
		}
		biz["gopenid"] = gopenid
	}
	biz["sGameIcon"] = v["sGameIcon"]
	biz["iRecentRole"] = v["iRecentRole"] // 开启最近登录角色

	// 客户端下载地址和App拉起scheme
	if c.source == "ios" {
		biz["sIosDownloadUrl"] = v["sIosDownloadUrl"]
		biz["sIosScheme"] = v["sIosScheme"]
	} else {
		biz["sAndroidDownloadUrl"] = v["sAndroidDownloadUrl"]
		biz["sAndroidPkgName"] = v["sAndroidPkgName"]
	}

	biz["iosDqPay"] = 0                     // ios点券，已废弃
	biz["payment"] = c.getPaymentSetting(v) // 支付配置
	biz["iCategory"] = v["iCategory"]       // 游戏品类id
	biz["sCateName"] = v["sCateName"]       // 游戏品类
	biz["iHot"] = v["iHot"]                 // 热门
	biz["iTop"] = v["iTop"]                 // 置顶

	biz["sDesc"] = "" // 游戏介绍
	if status == 1 && c.version < 129 {
		biz["sDesc"] = v["sDesc"] // 预约状态游戏返回游戏介绍
	}

	biz["sMktText"] = v["sMktText"]     // 营销文案
	biz["sCloudGame"] = v["sCloudGame"] // 云游戏

	return biz
}
