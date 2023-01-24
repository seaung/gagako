package snmp

import (
	"fmt"

	"github.com/alouca/gosnmp"
	"github.com/seaung/gagako/pkg/utils"
)

func SNMPOK(host string) {
	snmp, err := gosnmp.NewGoSNMP(host, "public", gosnmp.Version2c, 5)
	if err != nil {
		utils.New().Info(fmt.Sprintf("Get SNMP error %v", err))
	}

	response, err := snmp.Get(".1.3.6.1.2.1.1.1.0")
	if err != nil {
		utils.New().Info(fmt.Sprintf("Not Found SNMP host : %s", host))
	}

	for _, val := range response.Variables {
		switch val.Type {
		case gosnmp.OctetString:
			utils.New().Success(fmt.Sprintf("Found SNMP Host : %s", host))
		}
	}
}

func GetHostInfoFromSNMP(host string) {
	snmp, err := gosnmp.NewGoSNMP(host, "public", gosnmp.Version2c, 5)
	if err != nil {
		utils.New().Warnning(fmt.Sprintf("GET SNMP error %v", err))
	}

	response, err := snmp.Get(".1.3.6.1.2.1.1.1.0")
	if err != nil {
		utils.New().Warnning(fmt.Sprintf("GET SNMP error from response : %v", err))
	}

	for _, val := range response.Variables {
		switch val.Type {
		case gosnmp.OctetString:
			utils.New().Success(fmt.Sprintf("SNMP : %s \t HOST NAME : %s  SNMP VALUE : %s\t", host, GetHostname(host), val.Value.(string)))
		}
	}
}

func GetHostname(host string) string {
	snmp, err := gosnmp.NewGoSNMP(host, "public", gosnmp.Version2c, 5)
	if err != nil {
		return "SNMP NOT FOUND"
	}

	response, err := snmp.Get(".1.3.6.1.2.1.1.1.0")
	if err != nil {
		return "SNMP NOT FOUND"
	}

	for _, val := range response.Variables {
		switch val.Type {
		case gosnmp.OctetString:
			return val.Value.(string)
		}
	}
	return "SNMP NOT FOUND"
}
