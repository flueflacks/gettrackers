## Feature Addition: Start Priority Flag

Add a new flag to the `groups` command that outputs N blank lines before the grouped tracker URLs.

### Flag Specification
- `--start-priority <number>` or `-p <number>`
- Takes an integer argument (number of blank lines to output)
- Only applies to the `groups` command (default command)
- Works with both stdout and file output (`-o` flag)

### Behavior
- Output N blank lines at the very beginning, before any tracker groups
- Then output the normal grouped tracker URLs as usual
- Example: `gettrackers groups --start-priority 10` outputs 10 blank lines, then the tracker groups

### Use Case
This creates N "empty groups" that offset the priority of added trackers in clients like qBittorrent, effectively lowering the priority of the new trackers.

### Implementation Notes
- Validate that the number is non-negative (0 or positive integer)
- If invalid number provided, show clear error message
- Default value should be 0 (no blank lines) if flag not specified
- The blank lines should come before any other output, including the first tracker group

### Example Output
```
gettrackers groups -p 3
```

Would output:
```
[blank line]
[blank line]
[blank line]
http://tracker1.example.com:8080/announce
http://tracker1.example.com:9090/announce

http://api.other.com/tracker
...
```