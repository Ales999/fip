# fip

Show all OR find ARP by IP and Find port by MAC, and usind CDP for find next host

If the location of the configuration files and their names differ from the default ones, specify them through the necessary switches or define the corresponding variables CISFILE and CISPWDS.
Example set environment with .profile:

```bash
export CISFILE="$HOME/cisco/myciscos.yaml"
export CISPWDS="$HOME/cisco/mypwds.json"
```

---

Вывод ARP и MAC таблиц, а так-же поиск по ним (-f or --find)
При этом программа использует CDP для поиска следующего коммутатора при поиске кокретного порта при поиске по MAC.

Примечание: Если местоположение файлов конфигурации и их имена отличаются от указанных дефолтных, укажите их через необходимые ключили, или определите соответствующие переменные CISFILE и CISPWDS. Пример как можно указать есть выше.

```bash
$ fip
Usage: fip <command> [flags]

Find IP with ARP and MAC table

Flags:
  -h, --help                                     Show context-sensitive help.
  -d, --debug                                    More connect debug info
  -p, --port-ssh=22                              SSH порт для доступа к cisco
      --cis-file-name="/etc/cisco/cis.yaml"      Путь к файлу конфигурации имя_cisco/группа/ip ($CISFILE)
      --pwd-file-name="/etc/cisco/passw.json"    Путь к файлу конфигурации имя_группы/имя/пароль ($CISPWDS)

Commands:
  arp <hosts> ... [flags]
    Get or find ARP

  mac <hosts> ... [flags]
    Get or find MAC

Run "fip <command> --help" for more information on a command.

fip: error: expected one of "arp",  "mac"
```

Command ARP:

```bash
$ fip arp
Usage: fip arp <hosts> ... [flags]

Get or find ARP

Arguments:
  <hosts> ...    Name of cisco hosts for finded ARP

Flags:
  -h, --help                                     Show context-sensitive help.
  -d, --debug                                    More connect debug info
  -p, --port-ssh=22                              SSH порт для доступа к cisco
      --cis-file-name="/etc/cisco/cis.yaml"      Путь к файлу конфигурации имя_cisco/группа/ip ($CISFILE)
      --pwd-file-name="/etc/cisco/passw.json"    Путь к файлу конфигурации имя_группы/имя/пароль ($CISPWDS)

  -f, --find=STRING                              Поиск по MAC или IP

fip: error: expected "<hosts> ..."

```

Command MAC:

```bash
$ fip mac
Usage: fip mac <hosts> ... [flags]

Get or find MAC

Arguments:
  <hosts> ...    Name of cisco hosts for finded MAC

Flags:
  -h, --help                                     Show context-sensitive help.
  -d, --debug                                    More connect debug info
  -p, --port-ssh=22                              SSH порт для доступа к cisco
      --cis-file-name="/etc/cisco/cis.yaml"      Путь к файлу конфигурации имя_cisco/группа/ip ($CISFILE)
      --pwd-file-name="/etc/cisco/passw.json"    Путь к файлу конфигурации имя_группы/имя/пароль ($CISPWDS)

  -f, --find=STRING                              Поиск по MAC

fip: error: expected "<hosts> ..."

```
