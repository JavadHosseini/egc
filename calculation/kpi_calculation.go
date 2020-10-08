package calculation

import (
	"strconv"
	"xml-parser/models"
)

func Calculate(mdc models.MDC) map[string]map[string]int {
	kpis := accessibilityCalculation(mdc)
	return kpis
}

func accessibilityCalculation(mdc models.MDC) map[string]map[string]int {
	md2761 := populateCountersKeyValue(mdc)
	KPIs := calculateAccessibilityKPIs(md2761)

	for key, value := range KPIs {
		println("-------------", key, "---------------------")
		for k, v := range value {
			println(k, ":", v)
		}
	}
	return KPIs
}

func populateCountersKeyValue(mdc models.MDC) map[string]map[string]string {
	md2761 := make(map[string]map[string]string)

	mdList := []uint{27, 61}
	EUtranCellFDD := []string{"LN4277XA0", "LN4277XA10", "LN4277XA2", "LN4277XB0", "LN4277XB1", "LN4277XB2", "LN4277XC0", "LN4277XC1", "LN4277XC2"}
	for moidKey, moid := range EUtranCellFDD {
		md2761[moid] = make(map[string]string)
		for _, md := range mdList {
			mt := mdc.Md[md].Mi.Mt
			mv := mdc.Md[md].Mi.Mv
			for i, el := range mt {
				md2761[moid][el] = mv[moidKey].R[i]
			}
		}
	}

	return md2761
}

func calculateAccessibilityKPIs(md2761 map[string]map[string]string) map[string]map[string]int {
	KPIs := make(map[string]map[string]int)

	for EUtranCellFDDKey, EUtranCellFDDValue := range md2761 {
		counters := convertCounterToInt(EUtranCellFDDValue)

		kpis := calculateKPIs(counters)

		KPIs[EUtranCellFDDKey] = kpis
	}

	return KPIs
}

func convertCounterToInt(EUtranCellFDDValue map[string]string) map[string]int {
	integerCounter := make(map[string]int)
	accessibilityCounters := []string{
		"pmRrcConnEstabSucc", "pmRrcConnEstabAtt", "pmRrcConnEstabAttReatt", "pmRrcConnEstabFailMmeOvlMos",
		"pmRrcConnEstabFailMmeOvlMod", "pmS1SigConnEstabSucc", "pmS1SigConnEstabAtt", "pmS1SigConnEstabFailMmeOvlMos",
		"pmErabEstabSuccInit", "pmErabEstabAttInit", "pmErabEstabSuccAdded", "pmErabEstabAttAdded", "pmErabEstabAttAddedHoOngoing",
		"pmErabEstabFailAddedLic", "pmRrcConnEstabFailLic", "pmRrcConnEstabAttMod", "pmRrcConnEstabAttMos", "pmRrcConnEstabAttEm",
		"pmRrcConnEstabAttMta", "pmRrcConnEstabAttHpa", "pmErabEstabFailInitLic", "pmUeCtxtEstabSucc", "pmUeCtxtEstabAtt",
		"pmPagDiscarded", "pmPagDiscarded", "pmPagReceived", "pmRaSuccCbra", "pmRaAttCbra", "pmRaFailCbraMsg2Disc",
	}

	for _, cnt := range accessibilityCounters {
		integerCounter[cnt], _ = strconv.Atoi(EUtranCellFDDValue[cnt])
	}
	return integerCounter
}

func calculateKPIs(counters map[string]int) map[string]int {
	kpi := make(map[string]int)
	// 2.1 Accessibility (EUtranCellFDD/TDD)
	kpi["Acc_RrcConnSetupSuccRate"] = 100 * counters["pmRrcConnEstabSucc"] / (counters["pmRrcConnEstabAtt"] - counters["pmRrcConnEstabAttReatt"] - counters["pmRrcConnEstabFailMmeOvlMos"] - counters["pmRrcConnEstabFailMmeOvlMod"])
	kpi["Acc_S1SigEstabSuccRate"] = 100 * counters["pmS1SigConnEstabSucc"] / (counters["pmS1SigConnEstabAtt"] - counters["pmS1SigConnEstabFailMmeOvlMos"])
	kpi["Acc_InitialErabSetupSuccRate"] = 100 * counters["pmErabEstabSuccInit"] / counters["pmErabEstabAttInit"]
	//kpi["Acc_InitialERabEstabSuccRate"] = 100 * kpi["Acc_RrcConnSetupSuccRate"] * kpi["Acc_S1SigEstabSuccRate"] * kpi["Acc_InitialErabSetupSuccRate"] / 10000
	//kpi["Acc_AddedERabEstabSuccRate"] = 100 * counters["pmErabEstabSuccAdded"] / (counters["pmErabEstabAttAdded"] - counters["pmErabEstabAttAddedHoOngoing"])
	// kpi["Acc_AddedERabEstabFailRateDueToMultipleLicense"] = 100 * counters["pmErabEstabFailAddedLic"] / counters["pmErabEstabAttAdded"]
	// kpi["Acc_RrcConnSetupFailureRateDueToLackOfConnectedUsersLicense"] = 100 * counters["pmRrcConnEstabFailLic"] / counters["pmRrcConnEstabAtt"]
	// kpi["Acc_RrcConnSetupRatioForMOData"] = 100 * counters["pmRrcConnEstabAttMod"] / counters["pmRrcConnEstabAtt"]
	// kpi["Acc_RrcConnSetupRatioForMOSignalling"] = 100 * counters["pmRrcConnEstabAttMos"] / counters["pmRrcConnEstabAtt"]
	// kpi["Acc_RrcConnSetupRatioForEmergency"] = 100 * counters["pmRrcConnEstabAttEm"] / counters["pmRrcConnEstabAtt"]
	// kpi["Acc_RrcConnSetupRatioForMobileTerminating"] = 100 * counters["pmRrcConnEstabAttMta"] / counters["pmRrcConnEstabAtt"]
	// kpi["Acc_RrcConnSetupRatioForHighPrioAccess"] = 100 * counters["pmRrcConnEstabAttHpa"] / counters["pmRrcConnEstabAtt"]
	// kpi["Acc_InitialERabEstabFailureRateDueToMultipleLicense"] = 100 * counters["pmErabEstabFailInitLic"] / counters["pmErabEstabAttInit"]
	// kpi["Acc_InitialUEContextEstabSuccRate"] = 100 * counters["pmUeCtxtEstabSucc"] / counters["pmUeCtxtEstabAtt"]
	// kpi["Acc_PagingDiscardRate"] = 100 * counters["pmPagDiscarded"] / counters["pmPagReceived"]
	// kpi["Acc_RandomAccessDecodingRate"] = 100 * counters["pmRaSuccCbra"] / counters["pmRaAttCbra"]
	// kpi["Acc_RandomAccessMSG2Congestion"] = 100 * counters["pmRaFailCbraMsg2Disc"] / counters["pmRaAttCbra"]

	return kpi
}
