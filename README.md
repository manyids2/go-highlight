# go-highlight

Neovim type highlight groups for tcell.

- colors taken from [zenbones](https://github.com/mcchrish/zenbones.nvim)

Proof of concept

- [x] read a colors file
- [x] index colors
- [x] extract highlight assignments
- [x] assign colors, attributes with `termenv`
- [x] print highlights sheet
- [ ] divide by function ( tree-sitter, neovim, lang, etc. )
- [ ] divide by theme ( light, dark, italics )


Type of output we need
Maybe in the scanner store array of lines? we do have start and end bytes though.

```
package_clause [0, 0] - [0, 14] @include /                                                                                                                        
  package_identifier [0, 8] - [0, 14] @namespace /                                                                                                                
import_declaration [2, 0] - [13, 1] @include /                                                                                                                    
  import_spec_list [2, 7] - [13, 1] @punctuation.bracket /                                                                                                        
    import_spec [3, 1] - [3, 10] @string /                                                                                                                        
      path: interpreted_string_literal [3, 1] - [3, 10] @string /  
```

- tree-sitter queries and highlights
- given file of language, highlights sheet
  - print language with given highlights
- languages ( which go-tree-sitter supports )
  - markdown
  - go
  - python
  - c

## Get vim highlight on file

Change filetype for relevant highlights.

```bash
FILETYPES=("md" "go" "py" "bash")
for FILETYPE in ${FILETYPES[@]}; do
    echo "$FILETYPE"
    nvim "temp.$FILETYPE" -c "pu=execute('highlight')" -c "wq"
    mv "temp.$FILETYPE" "$FILETYPE.hi"
done
```

Check ANSI highlights with fzf.

```bash
./go-highlight -p ./corpus/md.hi |\
    fzf --ansi --reverse --height 30
```

## As blog

1. Create mock highlights file about 50-100 MB in size by duplicating

- Check performance of various scanning strategies
- Create custom parser

2. Go port of antirez-kilo, for line buffers
