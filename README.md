# gettrackers

## Purpose 

Fetches, filters, and structures public torrent tracker lists,
grouping them into **priority tiers** optimized for use in qBittorrent and
other compatible clients.

such as [ngosang's trackerslist](https://github.com/ngosang/trackerslist) 

## Usage

```
> gettrackers --help
gettrackers is a CLI tool that downloads tracker URLs from configurable sources, filters them using a blocklist, and outputs them grouped by domain.

Usage:
  gettrackers [flags]
  gettrackers [command]

Available Commands:
  block       Add a pattern to the blocklist
  completion  Generate the autocompletion script for the specified shell
  config      Manage configuration
  fetch       Force download/update the cached sources file
  groups      Output grouped tracker URLs (default command)
  help        Help about any command
  show        Show cached sources or blocklist

Flags:
  -h, --help                 help for gettrackers
  -o, --output string        Write output to file instead of stdout
  -p, --start-priority int   Output N blank lines before tracker groups (default: 0)

Use "gettrackers [command] --help" for more information about a command.
```

note the source download file is cached for 24h. 

## Backstory

found myself reviving a bunch of very old torrents, most of which no
longer had working trackers listed in them.  Some trackers were
rejecting the client contacting too quickly, and realized this was because sometimes the all lists have http, https, and udp  variants of the same tracker

so this program groups those into the same tier

## Tracker Tiering (Based on libtorrent)

in th underlying **libtorrent** library, which powers qBittorrent. Trackers are
organized into **priority tiers** 

### Tracker Tiers Redundancy and Announcement Logic 

libtorrent use of tracker tiers are controlled by [two
settings](https://www.libtorrent.org/reference-Settings.html#announce_to_all_trackers). 

> `announce_to_all_trackers` controls how multi tracker torrents are treated.
If this is set to true, all trackers in the same tier are announced to in
parallel. If all trackers in tier 0 fails, all trackers in tier 1 are
announced as well. If it's set to false, the behavior is as defined by the
multi tracker specification.

> `announce_to_all_tiers` also controls how multi tracker torrents are
treated. When this is set to true, one tracker from each tier is announced
to. This is the uTorrent behavior. To be compliant with the Multi-tracker
specification, set it to false.

this tool groups trackers with the same hostname, into the same priority
tier, with the assumption that qBittorrent's default configuration is
`announce_to_all_trackers` false,  and `announce_to_all_tiers` true. 

so `udp://` `http://` `https://` or different port numbers,  with the same hostname all get grouped into the same tier. 

### Defining Tiers in Tracker Lists

When a multi-line list of trackers is added to a torrent, qBittorrent
interprets the structure as follows:

* **Same Tier (Group):** All consecutive trackers (each on a new line) areassigned to the **same priority tier**.

* **Next Tier (Priority Fallback):** A single **empty line** separates the
  current group, assigning the subsequent block of trackers to the **next
  priority level** (e.g., Tier 0, then Tier 1, etc.).

