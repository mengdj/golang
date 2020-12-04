#!/usr/bin/env bash
###########################################################
#Author:mengdj@outlook.com
#Created Time:2020.12.04 11:56
#Description:execute go mod tidy in current directory
#Version:0.0.1
###########################################################

CURRENT_DIR=$(pwd)
SEARCH_DIR=$CURRENT_DIR
SEARCH_TOTAL=0

function GoTidy() {
	for file in $(ls $1); do
		local target=$1"/"$file
		if test -d $target; then
			if test -f $target"/go.mod"; then
				echo "process "$target
				let "SEARCH_TOTAL+=1"
				cd $target
				go mod tidy
			fi
			GoTidy $target
		fi
	done
}

if test $# -ne 0; then
	if test -d $1; then
		cd $1
		SEARCH_DIR=$(pwd)
	else
		echo $1" is not exist directory."
		exit
	fi
fi
GoTidy $SEARCH_DIR
echo "processed("$SEARCH_TOTAL")."
cd $CURRENT_DIR
