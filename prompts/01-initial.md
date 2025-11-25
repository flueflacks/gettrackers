# Prompt for AI Agent: Build gettrackers CLI Tool

Write a Go program called `gettrackers` that downloads, filters, and outputs tracker URLs grouped by domain. The program should be idiomatic Go, simple, and maintainable.

## Program Overview

`gettrackers` is a CLI tool that:
1. Downloads a list of tracker URLs from configurable mirror sources
2. Caches the downloaded list locally
3. Filters URLs using a blocklist
4. Groups URLs by domain (treating subdomains as separate groups)
5. Outputs the grouped URLs in random group order

## Command Structure

Use the Cobra library for command-line parsing. Implement these commands:

### Main Commands
- `gettrackers` or `gettrackers groups` - Default command that outputs the final grouped tracker list
- `gettrackers fetch` - Force download/update the cached sources file
- `gettrackers block <pattern>` - Add a pattern to the blocklist (check for duplicates, don't add if exists)
- `gettrackers show sources` - Display the cached source URLs
- `gettrackers show blocklist` - Display blocklist entries
- `gettrackers config get [key]` - Show config values (all config if no key specified, specific value if key given)
- `gettrackers config set <key> <value>` - Set a config value

### Flags
- `-o, --output <file>` - Write output to file instead of stdout (works with `groups` command)

## File Locations (XDG Standard)

Use XDG environment variables with fallback defaults:
- **Config file**: `$XDG_CONFIG_HOME/gettrackers/config.yaml` (default: `~/.config/gettrackers/config.yaml`)
- **Blocklist**: `$XDG_CONFIG_HOME/gettrackers/blocklist.txt` (default: `~/.config/gettrackers/blocklist.txt`)
- **Cache**: `$XDG_CACHE_HOME/gettrackers/sources.txt` (default: `~/.cache/gettrackers/sources.txt`)

Create directories as needed if they don't exist.

## Configuration (YAML)

The config file should be in YAML format with this structure:

```yaml
source_urls:
  - https://raw.githubusercontent.com/ngosang/trackerslist/master/trackers_all.txt
  - https://ngosang.github.io/trackerslist/trackers_all.txt
  - https://cdn.jsdelivr.net/gh/ngosang/trackerslist@master/trackers_all.txt
```

**Default values**: If no config file exists, use the three URLs above as defaults.

## Core Behavior

### Caching Logic
- Download source URLs only when:
  - Cache file doesn't exist, OR
  - Cache file is older than 24 hours (based on modification time), OR
  - `fetch` command is explicitly run
- Try mirror URLs in **random order** until one succeeds
- Update cache file modification time on successful download

### Filtering
- Read blocklist file (if it exists; if not, proceed without filtering)
- Remove any source URL line that contains any blocklist pattern (simple substring matching)
- Remove blank lines from source list

### Grouping
- Parse each URL and extract the domain
- Group URLs by their full domain (e.g., `tracker.example.com` and `api.example.com` are separate groups)
- URLs that fail to parse should be grouped together in a "no domain" group

### Output Format
- Output groups in **random order**
- Within each group, list URLs one per line
- Separate groups with a single blank line
- Example:
```
http://tracker1.example.com:8080/announce
http://tracker1.example.com:9090/announce

http://api.other.com/tracker
http://api.other.com/announce

http://subdomain.site.org/path
```

### Error Handling
- If download fails for all mirrors: output clear error message to stderr, exit with non-zero status
- If cache is stale but download fails: use stale cache and warn to stderr
- If all URLs are blocked: output nothing to stdout, write "all urls blocked" warning to stderr
- For any fatal errors: output clear error message to stderr, no output to stdout, exit non-zero

## Requirements

- Use idiomatic Go code
- Keep it simple and maintainable
- Use Cobra for CLI commands
- Use YAML library for config (e.g., `gopkg.in/yaml.v3`)
- Use standard library where possible
- Include proper error handling
- Create necessary directories if they don't exist
- Add helpful comments for clarity

## Implementation Notes

- When adding to blocklist with `block` command, check for duplicates and skip if already exists (but return successfully)
- The cache file modification time should reflect the last successful download
- Random group ordering should use Go's `math/rand` or `crypto/rand`
- Subdomains are treated as distinct domains (www.example.com ≠ example.com)
- Maintain clean separation between config loading, fetching, filtering, grouping, and output logic