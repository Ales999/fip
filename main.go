package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/alecthomas/kong"

	"github.com/ales999/cisaccs"
	"github.com/ales999/fip/nxthst"
)

var cli struct {

	// Команда для поиска ARP данных
	Arp struct {
		CheckHosts  []string `arg:"" name:"hosts" help:"Name of cisco hosts for finded ARP"`
		FindIpOrMac string   `name:"find" help:"Поиск по MAC или IP" short:"f"`
	} `cmd:"" help:"Get or find ARP"`

	// Команда для поиска MAC данных
	Mac struct {
		CheckHosts []string `arg:"" name:"hosts" help:"Name of cisco hosts for finded MAC"`
		FindedMac  string   `name:"find" help:"Поиск по MAC" short:"f"`
	} `cmd:"" help:"Get or find MAC"`

	// Команда для поиска сначала ARP по IP, и по найденному ищем порт по MAC
	Combo struct {
		CheckHosts []string `arg:"" name:"hosts" help:"Name of cisco hosts for finded Port"`
		FindByIp   string   `name:"ip" help:"Поиск по IP" short:"f" required:""`
	} `cmd:"" help:"Combo find, first arp, next port"`

	// Больше информации при подключении
	DebugInfo bool `name:"debug" help:"More connect debug info" short:"d"`
	// Номер порта для ssh
	PortSsh int `help:"SSH порт для доступа к cisco" short:"p" default:"22"`
	// TimeOut SSH connect
	PortTimeout int `help:"SSH Timeout in second" short:"t" default:"10"`
	// Путь к файлу конфигурации имя_cisco/группа/ip - env: CISFILE
	CisFileName string `help:"Путь к файлу конфигурации имя_cisco/группа/ip" default:"/etc/cisco/cis.yaml" env:"CISFILE"`
	// Путь к файлу конфигурации имя_группы/имя/пароль - env: CISPWDS
	PwdFileName string `help:"Путь к файлу конфигурации имя_группы/имя/пароль" default:"/etc/cisco/passw.json" env:"CISPWDS"`
}

func main() {

	ctx := kong.Parse(&cli,
		kong.Name("fip"),
		kong.Description("Find IP with ARP and MAC table"),
		kong.UsageOnError(),
	)

	// Если указан MAC, то его нужно очистить от лишних символов и привести к cisco формату
	if len(cli.Arp.FindIpOrMac) > 0 {
		cli.Arp.FindIpOrMac = nxthst.ConvertMACAddress(cli.Arp.FindIpOrMac)
	} else if len(cli.Mac.FindedMac) > 0 {
		cli.Mac.FindedMac = nxthst.ConvertMACAddress(cli.Mac.FindedMac)
	} else if len(cli.Combo.FindByIp) > 0 {
		cli.Combo.FindByIp = nxthst.ConvertMACAddress(cli.Combo.FindByIp)
	}

	// Если в параметрах указана отладка, передадим ее в нашу библиотеку.
	cisaccs.SetMoreOutputConnectInfo(cli.DebugInfo)

	switch ctx.Command() {
	case "arp <hosts>":
		err := FipFindArpCommand()
		ctx.FatalIfErrorf(err)
	case "mac <hosts>":
		err := FipFindMacCommand()
		ctx.FatalIfErrorf(err)
	case "combo <hosts>":
		err := FipFindCombo()
		ctx.FatalIfErrorf(err)
	default:
		panic(ctx.Command())

	}
	fmt.Println("---")
	os.Exit(0)
}

func FipFindCombo() error {

	cli.Arp.CheckHosts = cli.Combo.CheckHosts
	cli.Arp.FindIpOrMac = cli.Combo.FindByIp

	err := FipFindArpCommand()
	if err != nil {
		return err
	} else {
		// Если успешно нашли MAC через IP
		if len(cli.Mac.FindedMac) > 0 {
			cli.Mac.CheckHosts = cli.Combo.CheckHosts
			err := FipFindMacCommand()
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// Поиск ARP-данных
func FipFindArpCommand() error {

	// Что будем выполнять на cisco
	cmds := []string{`sh arp`}

	// Подготовка к подключению.
	acc := cisaccs.NewCisAccount(cli.CisFileName, cli.PwdFileName)

	// Если мы хотим найти конкретный IP или MAC в ARP то нужно парсить вывод
	if len(strings.TrimSpace(cli.Arp.FindIpOrMac)) > 0 {

		for _, hst := range cli.Arp.CheckHosts {

			// Получаем данные с каждого хоста
			out, err := acc.OneCisExecuteSsh(hst, cli.PortSsh, cmds, cli.PortTimeout)
			if err != nil {
				fmt.Println(err.Error())
				continue
			}

			found, arpsm := cisaccs.CisFindArp(out, cli.Arp.FindIpOrMac)
			if found {
				// Сохраним во временную переменную найденный MAC
				foundedMac := arpsm.GetMac()
				// Печать результата поиска
				fmt.Printf("ARP %s found, Host: %s Port %s IP: %s\n", foundedMac, hst, arpsm.GetIface(), arpsm.GetIp())
				// Запоминаем для Combo-режима
				cli.Mac.FindedMac = foundedMac
				// ARP ищем до первого совпадения, на оствшиеся хосты можно не ходить.
				break
			}

		}

	} else {
		out, err := acc.MultiCisExecuteSsh(cli.Arp.CheckHosts, cli.PortSsh, cmds)
		if err != nil {
			return err
		}

		for _, line := range out {
			if !strings.Contains(line, "Incomplete") {
				// Выводим всю arp-таблицу на экран
				fmt.Println(line)
			}
		}
	}
	return nil

}

/*
func FipFindNextHost() {

}
*/

// Поиск MAC-данных
func FipFindMacCommand() error {

	cmds := []string{"sh mac address-table dynamic"}
	// Prepare cisco account
	acc := cisaccs.NewCisAccount(cli.CisFileName, cli.PwdFileName)

	if len(strings.TrimSpace(cli.Mac.FindedMac)) > 0 {
		// Бежим по указанным хостам
		cmds = append(cmds, `sh etherchannel detail | i Group:|Port:`)
		cmds = append(cmds, `sh cdp entry * | i  Device|Interface`)
		for _, hst := range cli.Mac.CheckHosts {

			cisout, err := acc.OneCisExecuteSsh(strings.ToLower(hst), cli.PortSsh, cmds)
			if err != nil {
				fmt.Println("Host", hst, ":", err.Error())
				continue
			}
			macfound, macstrs := cisaccs.CisFindMac(cisout, cli.Mac.FindedMac)
			// Если что-то нашли то перебираем
			if macfound {
				for _, macstr := range macstrs {
					// Если найденный порт - это линк к другому коммутатору
					if strings.Contains(macstr.GetIface(), "Port-channel") {
						//debug - don't delete this!
						fmt.Print("External host: ")
					}

					// Печать результата поиска
					fmt.Printf("Mac %s found, Host: %s Port: %s Vlan: %s\n", macstr.GetMac(), hst, macstr.GetIface(), macstr.GetVlan())

					bnxtfns, nexthost := nxthst.FindNextPortbyMac(cisout, hst, *nxthst.NewLocMacLineData(macstr.GetVlan(), cli.Mac.FindedMac, macstr.GetIface()))
					if bnxtfns { // Если что-то нашли
						//var tstd *cisaccs.CisAccount
						//fmt.Println("Next host:", nexthost)

						// Рекурсивно типа перебираем
						for {
							//fmt.Println("Connected to", nexthost)
							if len(nexthost) == 0 || !bnxtfns {
								break
							}
							fmt.Println("Next host:", nexthost)
							cisout, err := acc.OneCisExecuteSsh(nexthost, cli.PortSsh, cmds)
							if err != nil {
								fmt.Println("Host", nexthost, ":", err.Error())
								break
							}
							macfound, macstrs := cisaccs.CisFindMac(cisout, cli.Mac.FindedMac)
							if !macfound {
								break
							}
							for _, macstr := range macstrs {

								if strings.Contains(macstr.GetIface(), "Port-channel") {
									//debug
									fmt.Print("External host: ")
								}

								fmt.Printf("Mac %s found, Host: %s Port: %s Vlan: %s\n", macstr.GetMac(), nexthost, macstr.GetIface(), macstr.GetVlan())
								bnxtfns, nexthost = nxthst.FindNextPortbyMac(cisout, hst, *nxthst.NewLocMacLineData(macstr.GetVlan(), cli.Mac.FindedMac, macstr.GetIface()))
								if !bnxtfns {
									break
								}

							}
							//if !bnxtfns {
							//	break
							//}
						} // End рекурсия!

					}

				}
			}
		}
	} else {
		out, err := acc.MultiCisExecuteSsh(cli.Mac.CheckHosts, cli.PortSsh, cmds)
		if err != nil {
			fmt.Println(err.Error())
		}

		for _, line := range out {
			// Выводим всю mac-таблицу на экран
			fmt.Println(line)
		}

	}

	return nil
}
