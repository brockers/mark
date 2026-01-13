# Mark

Mark is an opinionated go commandline unix utility to create bookmarks for easy access to other folders.

## Examples

### Create Default

Create a new bookmark in of the current folders

```bash
mark 
```

This will create a symbolic link in your ~/.marks/ folder, named the same as the folder you are in, that points to the folder you are in.

### Create Named

Create a new bookmark but with a specified name

```bash
mark downloads
```

Creates a symbolic link in your ~/.marks/ folder, named downloads that points to the folder you are currently in.

### Create with Custom Path

Create a bookmark pointing to a custom directory (not current directory)

```bash
mark work ~/projects/work
mark tmp /tmp
```

Creates a symbolic link in your ~/.marks/ folder with the specified name that points to the given path. The path can be absolute or relative, and tilde expansion is supported.

### Show Bookmarks

List all of your bookmarks and where they point to:

```bash
mark -l

  downloads -> /home/jsmith/Downloads
  mark      -> /home/jsmith/Project/mark
```

Cleanly displays all of the symbolic links in your ~/.marks/ folder.

### Delete Bookmark

```bash 
mark -d downloads 
```

Removes the sybmolic link in your ~/.marks/ folder named downloads.

### Go to Bookmark 

Jump to your bookmarked folder

```bash
mark -j downloads
```

Does a `cd ~/.marks/downloads` to send you to the named bookmark.

### Alias

mark also has a couple built in aliases including

```bash
marks  #same as mark -l 
unmark #same as mark -d 
jump   #same as mark -j
```

### Autocomplete

Finally mark has built in autocomplete so you can alway double tab to see which mark you will jump to or delete.

## Related Projects

- [note](https://github.com/brockers/note) - A minimalist, opinionated command line note creation/management application.

## License

This program is free software licensed under GPL-3.0.
See https://www.gnu.org/licenses/ for details.

