package calculation

import (
	"strconv"
	"strings"
	"xml-parser/models"
)

func Calculate(mdc models.MDC) []models.KPIJsonModel {
	siteName := strings.Split(strings.Split(mdc.Mfh.Sn, ",")[2], "=")[1]

	KPIsModelList := accessibilityCalculation(mdc, siteName)
	return KPIsModelList
}

func accessibilityCalculation(mdc models.MDC, siteName string) []models.KPIJsonModel {
	countersModelList := PopulateCountersKeyValue(mdc, siteName)
	KPIsModelList := calculateAccessibilityKPIs(countersModelList)

	return KPIsModelList
}

func PopulateCountersKeyValue(mdc models.MDC, siteName string) []models.CountersJsonModel {
	var outputCountersList []models.CountersJsonModel
	mdList := []uint{27, 61}
	cells := []string{"A0", "A10", "A2", "B0", "B1", "B2", "C0", "C1", "C2"}
	for cellKey, cell := range cells {
		output := new(models.CountersJsonModel)
		output.SiteName = siteName
		output.SectorName = cell[0:1]
		output.CellName = cell
		output.Counters = make(map[string]string)
		for _, md := range mdList {
			mt := mdc.Md[md].Mi.Mt
			mv := mdc.Md[md].Mi.Mv
			for i, el := range mt {
				output.Counters[el] = mv[cellKey].R[i]
			}
		}
		outputCountersList = append(outputCountersList, *output)
	}

	return outputCountersList
}

func calculateAccessibilityKPIs(countersModelList []models.CountersJsonModel) []models.KPIJsonModel {
	var outputKPIList []models.KPIJsonModel

	for _, counterModel := range countersModelList {
		kpiOutput := new(models.KPIJsonModel)
		kpiOutput.SiteName = counterModel.SiteName
		kpiOutput.SectorName = counterModel.SectorName
		kpiOutput.CellName = counterModel.CellName
		kpiOutput.KPIs = make(map[string]int)

		counters := ConvertCounterToInt(counterModel.Counters)
		kpis := calculateKPIs(counters)

		kpiOutput.KPIs = kpis
		outputKPIList = append(outputKPIList, *kpiOutput)
	}

	return outputKPIList
}

func ConvertCounterToInt(EUtranCellFDDValue map[string]string) map[string]int {
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

/// Cell level calculation
//TODO: Sector level and site level calculation
func calculateKPIs(counters map[string]int) map[string]int {
	kpi := make(map[string]int)

	// 2.1 Accessibility (EUtranCellFDD/TDD)
	if counters["pmRrcConnEstabAtt"]-counters["pmRrcConnEstabAttReatt"]-counters["pmRrcConnEstabFailMmeOvlMos"]-counters["pmRrcConnEstabFailMmeOvlMod"] == 0 {
		kpi["Acc_RrcConnSetupSuccRate"] = 0
	} else {
		kpi["Acc_RrcConnSetupSuccRate"] = 100 * counters["pmRrcConnEstabSucc"] / (counters["pmRrcConnEstabAtt"] - counters["pmRrcConnEstabAttReatt"] - counters["pmRrcConnEstabFailMmeOvlMos"] - counters["pmRrcConnEstabFailMmeOvlMod"])
	}

	if counters["pmS1SigConnEstabAtt"]-counters["pmS1SigConnEstabFailMmeOvlMos"] == 0 {
		kpi["Acc_S1SigEstabSuccRate"] = 0
	} else {
		kpi["Acc_S1SigEstabSuccRate"] = 100 * counters["pmS1SigConnEstabSucc"] / (counters["pmS1SigConnEstabAtt"] - counters["pmS1SigConnEstabFailMmeOvlMos"])
	}

	if counters["pmErabEstabAttInit"] == 0 {
		kpi["Acc_InitialErabSetupSuccRate"] = 0
	} else {
		kpi["Acc_InitialErabSetupSuccRate"] = 100 * counters["pmErabEstabSuccInit"] / counters["pmErabEstabAttInit"]
	}

	kpi["Acc_InitialERabEstabSuccRate"] = kpi["Acc_RrcConnSetupSuccRate"] * kpi["Acc_S1SigEstabSuccRate"] * kpi["Acc_InitialErabSetupSuccRate"] / 10000

	if counters["pmErabEstabAttAdded"]-counters["pmErabEstabAttAddedHoOngoing"] == 0 {
		kpi["Acc_AddedERabEstabSuccRate"] = 0
	} else {
		kpi["Acc_AddedERabEstabSuccRate"] = 100 * counters["pmErabEstabSuccAdded"] / (counters["pmErabEstabAttAdded"] - counters["pmErabEstabAttAddedHoOngoing"])
	}

	if counters["pmErabEstabAttAdded"] == 0 {
		kpi["Acc_AddedERabEstabFailRateDueToMultipleLicense"] = 0
	} else {
		kpi["Acc_AddedERabEstabFailRateDueToMultipleLicense"] = 100 * counters["pmErabEstabFailAddedLic"] / counters["pmErabEstabAttAdded"]
	}

	if counters["pmRrcConnEstabAtt"] == 0 {
		kpi["Acc_RrcConnSetupFailureRateDueToLackOfConnectedUsersLicense"] = 0
		kpi["Acc_RrcConnSetupRatioForMOData"] = 0
		kpi["Acc_RrcConnSetupRatioForMOSignalling"] = 0
		kpi["Acc_RrcConnSetupRatioForEmergency"] = 0
		kpi["Acc_RrcConnSetupRatioForMobileTerminating"] = 0
		kpi["Acc_RrcConnSetupRatioForHighPrioAccess"] = 0
	} else {
		kpi["Acc_RrcConnSetupFailureRateDueToLackOfConnectedUsersLicense"] = 100 * counters["pmRrcConnEstabFailLic"] / counters["pmRrcConnEstabAtt"]
		kpi["Acc_RrcConnSetupRatioForMOData"] = 100 * counters["pmRrcConnEstabAttMod"] / counters["pmRrcConnEstabAtt"]
		kpi["Acc_RrcConnSetupRatioForMOSignalling"] = 100 * counters["pmRrcConnEstabAttMos"] / counters["pmRrcConnEstabAtt"]
		kpi["Acc_RrcConnSetupRatioForEmergency"] = 100 * counters["pmRrcConnEstabAttEm"] / counters["pmRrcConnEstabAtt"]
		kpi["Acc_RrcConnSetupRatioForMobileTerminating"] = 100 * counters["pmRrcConnEstabAttMta"] / counters["pmRrcConnEstabAtt"]
		kpi["Acc_RrcConnSetupRatioForHighPrioAccess"] = 100 * counters["pmRrcConnEstabAttHpa"] / counters["pmRrcConnEstabAtt"]
	}

	if counters["pmErabEstabAttInit"] == 0 {
		kpi["Acc_InitialERabEstabFailureRateDueToMultipleLicense"] = 0
	} else {
		kpi["Acc_InitialERabEstabFailureRateDueToMultipleLicense"] = 100 * counters["pmErabEstabFailInitLic"] / counters["pmErabEstabAttInit"]
	}

	if counters["pmUeCtxtEstabAtt"] == 0 {
		kpi["Acc_InitialUEContextEstabSuccRate"] = 0
	} else {
		kpi["Acc_InitialUEContextEstabSuccRate"] = 100 * counters["pmUeCtxtEstabSucc"] / counters["pmUeCtxtEstabAtt"]
	}

	if counters["pmPagReceived"] == 0 {
		kpi["Acc_PagingDiscardRate"] = 0
	} else {
		kpi["Acc_PagingDiscardRate"] = 100 * counters["pmPagDiscarded"] / counters["pmPagReceived"]
	}

	if counters["pmRaAttCbra"] == 0 {
		kpi["Acc_RandomAccessDecodingRate"] = 0
		kpi["Acc_RandomAccessMSG2Congestion"] = 0
	} else {
		kpi["Acc_RandomAccessDecodingRate"] = 100 * counters["pmRaSuccCbra"] / counters["pmRaAttCbra"]
		kpi["Acc_RandomAccessMSG2Congestion"] = 100 * counters["pmRaFailCbraMsg2Disc"] / counters["pmRaAttCbra"]
	}

	return kpi
}
