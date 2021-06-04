package main

import (
	"flag"
	"fmt"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
)

type Amount struct {
	Value     string
	Timestamp time.Time
}

var ID, adddress string
var nAmount Amount

const totalBalance = "-Total Balance"

func main() {
	nAmount = Amount{}

	flag.StringVar(&ID, "f", "", "f es el fingerprint")
	flag.StringVar(&adddress, "address", "", "address es la wallet destino")

	flag.Parse()

	if ID == "" || adddress == "" {
		fmt.Println("Fingerprint y Address son requeridos")

		return
	}

	args := []string{
		"wallet",
		"show",
	}

	args = append(args, "-f", ID)

	fmt.Println(color.GreenString(fmt.Sprintf("Fingerprint: %s - Address: %s", ID, adddress)))

	for {
		var cmd *exec.Cmd
		if runtime.GOOS == "windows" {
			cmd = exec.Command(".\\chia.exe", args...)
		} else {
			cmd = exec.Command("./chia.exe", args...)
		}

		stdoutSterr, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Println(color.RedString("Wallet.Cmd.Err: "), err)
			return
		}

		output := string(stdoutSterr)
		lines := strings.Split(output, "\n")

		size := len(lines)
		for i := 0; i < size; i++ {
			line := lines[i]
			if strings.Contains(line, totalBalance) {
				one := strings.Split(line, ":")
				if len(one) >= 2 {
					two := strings.Split(one[1], " ")
					if len(two) >= 2 {
						value := two[1]
						transWallet(value)
					}
				}
			}
		}

		time.Sleep(time.Millisecond * 50)
	}
}

func transWallet(value string) {
	f64, err := strconv.ParseFloat(value, 64)
	if err != nil {
		fmt.Println(color.RedString("Wallet.Cmd.Err.XCH: "), err)
	}

	amount := fmt.Sprintf("%.12f", f64)

	if f64 > 0 && nAmount.Value != amount {
		args := []string{
			"wallet",
			"send",
			"-a",
			amount,
			"-m",
			"0",
			"-i",
			"1",
			"-t",
			adddress,
		}

		fmt.Println(color.YellowString(fmt.Sprint(args)))

		var cmd *exec.Cmd
		if runtime.GOOS == "windows" {
			cmd = exec.Command(".\\chia.exe", args...)
		} else {
			cmd = exec.Command("./chia.exe", args...)
		}

		buffer, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Println(color.RedString("Wallet.Cmd.Send.Err: "), err)
		}

		fmt.Println(color.GreenString("Wallet.Send.Ok: "), color.CyanString(string(buffer)))

		nAmount.Value = amount
	}

	now := time.Now()
	if !nAmount.Timestamp.IsZero() {
		diff := now.Sub(nAmount.Timestamp).Seconds()
		if diff >= 60 {
			nAmount.Timestamp = now
			nAmount.Value = ""
		}

	} else {
		nAmount.Timestamp = now
	}
}
