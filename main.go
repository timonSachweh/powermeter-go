package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type PowerMetrics struct {
	ConsumedEnergy float64 `json:"consumed_energy,omitempty"`
	ProvidedEnergy float64 `json:"provided_energy,omitempty"`
	CurrentPower   float64 `json:"current_power_usage,omitempty"`
	CurrentPowerL1 float64 `json:"current_power_l1,omitempty"`
	CurrentPowerL2 float64 `json:"current_power_l2,omitempty"`
	CurrentPowerL3 float64 `json:"current_power_l3,omitempty"`
	VoltageL1      float64 `json:"voltage_l1,omitempty"`
	VoltageL2      float64 `json:"voltage_l2,omitempty"`
	VoltageL3      float64 `json:"voltage_l3,omitempty"`
}

type KeyNotFound struct{}

func (m *KeyNotFound) Error() string {
	return "key not found"
}

func main() {
	ip := flag.String("ip", "", "IP address of the server")
	currentPower := flag.String("current-power", "", "Current Power Serial")
	currentPowerL1 := flag.String("current-powerL1", "", "Current Power L1")
	currentPowerL2 := flag.String("current-powerL2", "", "Current Power L2")
	currentPowerL3 := flag.String("current-powerL3", "", "Current Power L3")
	voltageL1 := flag.String("voltageL1", "", "Voltage L1")
	voltageL2 := flag.String("voltageL2", "", "Voltage L2")
	voltageL3 := flag.String("voltageL3", "", "Voltage L3")
	consumedEnergy := flag.String("consumed-energy", "", "Consumed Energy")
	providedEnergy := flag.String("provided-energy", "", "Provided Energy")

	flag.Parse()

	if "" == *ip {
		os.Exit(1)
	}

	resp, err := http.Get(fmt.Sprintf("http://%s/metrics", *ip))

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	bytes, _ := io.ReadAll(resp.Body)
	stringResp := string(bytes)

	metrics := PowerMetrics{}

	if value, err := getValueFor(stringResp, *currentPower); err == nil {
		metrics.CurrentPower = value
	}
	if value, err := getValueFor(stringResp, *currentPowerL1); err == nil {
		metrics.CurrentPowerL1 = value
	}
	if value, err := getValueFor(stringResp, *currentPowerL2); err == nil {
		metrics.CurrentPowerL2 = value
	}
	if value, err := getValueFor(stringResp, *currentPowerL3); err == nil {
		metrics.CurrentPowerL3 = value
	}
	if value, err := getValueFor(stringResp, *voltageL1); err == nil {
		metrics.VoltageL1 = value
	}
	if value, err := getValueFor(stringResp, *voltageL2); err == nil {
		metrics.VoltageL2 = value
	}
	if value, err := getValueFor(stringResp, *voltageL3); err == nil {
		metrics.VoltageL3 = value
	}
	if value, err := getValueFor(stringResp, *consumedEnergy); err == nil {
		metrics.ConsumedEnergy = value
	}
	if value, err := getValueFor(stringResp, *providedEnergy); err == nil {
		metrics.ProvidedEnergy = value
	}

	jsonMetrics, err := json.Marshal(metrics)
	if err != nil {
		os.Exit(1)
	}
	fmt.Println(string(jsonMetrics))
}

func getValueFor(data string, metric string) (float64, error) {
	if metric == "" {
		return 0, &KeyNotFound{}
	}
	lines := strings.Split(data, "\n")
	for _, line := range lines {
		if strings.Contains(line, metric) && strings.Contains(line, "serial") {
			number := strings.Split(line, " ")[1]
			if s, err := strconv.ParseFloat(number, 64); err == nil {
				return s, nil
			}
			return 0, &KeyNotFound{} // Placeholder return value
		}
	}
	return 0, &KeyNotFound{}
}
