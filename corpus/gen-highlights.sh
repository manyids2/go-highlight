#!/bin/bash

# Add necessary filetypes here
FILETYPES=("md" "go")

# Open in nvim and write highlights to file
for FILETYPE in "${FILETYPES[@]}"; do
	echo "$FILETYPE"
	nvim "temp.$FILETYPE" -c "pu=execute('highlight')" -c "wq"
	mv "temp.$FILETYPE" "$FILETYPE.hi"
done
