package service

import (
	"go-final-dpo/src/enum"
	st "go-final-dpo/src/structure"
	"sort"
)

func GetResultData(data st.Data) (res st.ResultT) {

	// SMS
	getSortSmsProvider(data, &res)
	getSortSmsCountry(data, &res)
	// MMS
	getSortMmsProvider(data, &res)
	getSortMmsCountry(data, &res)
	// VOICE CALL
	res.Data.VoiceCall = data.VoiceCall
	// E-MAIL
	getSortEmail(data, &res)
	// BILLING
	//res.Data.Billing = data.Billing
	// SUPPORT
	getSupport(data, &res)
	// INCIDENT
	getSortIncident(data, &res)

	status, errMessage := getStatusResult(&res)
	res.Status = status
	res.Error = errMessage

	return
}

// getSortSmsProvider Возвращает отсортированные данные SMS по провайдерам
func getSortSmsProvider(data st.Data, res *st.ResultT) {
	var smsData []st.SMSData
	for _, sms := range data.Sms {
		sms.Country = enum.CountryCode[sms.Country]
		smsData = append(smsData, sms)
	}
	// сортировка
	sort.SliceStable(smsData, func(i, j int) bool {
		return smsData[i].Provider < smsData[j].Provider
	})
	res.Data.SMS = append(res.Data.SMS, smsData)
}

// getSortSmsCountry Возвращает отсортированные данные SMS по странам
func getSortSmsCountry(data st.Data, res *st.ResultT) {
	var smsData []st.SMSData
	for _, sms := range data.Sms {
		sms.Country = enum.CountryCode[sms.Country]
		smsData = append(smsData, sms)
	}
	// сортировка
	sort.SliceStable(smsData, func(i, j int) bool {
		return smsData[i].Country < smsData[j].Country
	})
	res.Data.SMS = append(res.Data.SMS, smsData)
}

// getSortMmsProvider Возвращает отсортированные данные MMS по провайдерам
func getSortMmsProvider(data st.Data, res *st.ResultT) {
	var mmsData []st.MMSData
	for _, mms := range data.Mms {
		mms.Country = enum.CountryCode[mms.Country]
		mmsData = append(mmsData, mms)
	}
	// сортировка Provider
	sort.SliceStable(mmsData, func(i, j int) bool {
		return mmsData[i].Provider < mmsData[j].Provider
	})
	res.Data.MMS = append(res.Data.MMS, mmsData)
}

// getSortMmsCountry Возвращает отсортированные данные MMS по странам
func getSortMmsCountry(data st.Data, res *st.ResultT) {
	var mmsData []st.MMSData
	for _, mms := range data.Mms {
		mms.Country = enum.CountryCode[mms.Country]
		mmsData = append(mmsData, mms)
	}
	// сортировка Country
	sort.SliceStable(mmsData, func(i, j int) bool {
		return mmsData[i].Country < mmsData[j].Country
	})
	res.Data.MMS = append(res.Data.MMS, mmsData)
}

// getSortEmail Возвращает отсортированные данные по Email
func getSortEmail(data st.Data, res *st.ResultT) {
	emailMap := getMapEmail(data)
	emailMapSort := getSortMapEmail(emailMap)
	dataEmailResult := getDataEmailResult(emailMapSort)

	res.Data.Email = dataEmailResult
}

// getMapEmail Собирает Map по Email
func getMapEmail(data st.Data) (emailMap map[string][]st.EmailData) {
	emailMap = make(map[string][]st.EmailData)
	for _, email := range data.Email {
		emailMap[email.Country] = append(emailMap[email.Country], email)
	}
	return
}

// getSortMapEmail Сортирует Map по Provider
func getSortMapEmail(emailMap map[string][]st.EmailData) map[string][]st.EmailData {
	for key := range emailMap {
		emailMap[key] = getSortEmailProvider(emailMap[key])
	}
	return emailMap
}

// getSortEmailProvider Сортировка по Provider
func getSortEmailProvider(email []st.EmailData) []st.EmailData {
	sort.SliceStable(email, func(i, j int) bool {
		return email[i].DeliveryTime > email[j].DeliveryTime
	})
	return email
}

// getDataEmailResult Возвращает данные по EmailData
func getDataEmailResult(emailMap map[string][]st.EmailData) (dataEmailResult map[string][][]st.EmailData) {
	dataEmailResult = make(map[string][][]st.EmailData)
	for key := range emailMap {
		var providers [][]st.EmailData
		providers = append(providers, emailMap[key][:3])
		providers = append(providers, emailMap[key][len(emailMap[key])-4:len(emailMap[key])-1])

		dataEmailResult[key] = providers
	}
	return dataEmailResult
}

// getSupport Возвращает данное по Support
func getSupport(data st.Data, res *st.ResultT) {
	support := data.Support
	res.Data.Support = make([]int, 2)
	ticketsCount := 0
	averageTime := 60 / 18
	for i := range support {
		ticketsCount = ticketsCount + support[i].ActiveTickets
	}
	if ticketsCount < 9 {
		res.Data.Support[0] = 1
	} else if ticketsCount >= 8 && ticketsCount < 16 {
		res.Data.Support[0] = 2
	} else {
		res.Data.Support[0] = 3
	}
	res.Data.Support[1] = averageTime * ticketsCount
}

func getSortIncident(data st.Data, res *st.ResultT) {
	incidentData := data.Incident
	// сортировка
	sort.SliceStable(incidentData, func(i, j int) bool {
		return incidentData[i].Status < incidentData[j].Status
	})
	res.Data.Incidents = incidentData
}

// getStatusResult Проверяет данные и возвращает статус и ошибки
func getStatusResult(res *st.ResultT) (status bool, errMessage string) {
	data := res.Data
	status = true

	if data.SMS[0] == nil || data.SMS[1] == nil {
		errMessage += "Данны по SMS пустые "
		status = false
	}

	if data.MMS[0] == nil || data.MMS[1] == nil {
		errMessage += "Данны по MMS пустые "
		status = false
	}

	if data.VoiceCall == nil {
		errMessage += "Данны по Voice Call пустые "
		status = false
	}

	if data.Email == nil {
		errMessage += "Данны по Email пустые "
		status = false
	}

	if data.Support[0] == 0 && data.Support[1] == 0 {
		errMessage += "Данны по Support пустые "
		status = false
	}

	if data.Incidents == nil {
		errMessage += "Данны по Incidents пустые "
		status = false
	}

	return
}
