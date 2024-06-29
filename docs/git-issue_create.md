## git-issue create

Create a new issue

```
git-issue create [--[no]-stripspace] [-F <file> | -m <message>] [-l <label>] [-a <user>] [-s <status>] [flags]
```

### Options

```
  -a, --assign stringArray    Assign one or more users to the issue. Multiple users can be assigned with multiple -a options.
  -F, --file string           Take the issue message from the given file. Use - to read the issue message from the standard input. Lines starting with # and empty lines other than a single line between paragraphs will be stripped out. If you wish to keep them verbatim, use --no-stripspace.
  -h, --help                  help for create
  -l, --label stringArray     Add one or more labels to the issue. Multiple labels can be added with multiple -l options.
  -m, --message stringArray   Use the given issue message (instead of prompting). If multiple -m options are given, their values are concatenated as separate paragraphs. Lines starting with # and empty lines other than a single line between paragraphs will be stripped out. If you wish to keep them verbatim, use --no-stripspace.
      --no-stripspace         
  -s, --status string         Set the initial issue status. Defaults to open. (default "open")
      --stripspace            Strip leading and trailing whitespace from the issue message. Also strip out empty lines other than a single line between paragraphs. (default true)
```

### SEE ALSO

* [git-issue](git-issue.md)	 - Create, edit, or list issues

###### Auto generated by spf13/cobra on 29-Jun-2024