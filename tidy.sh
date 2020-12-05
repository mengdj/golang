#!/usr/bin/env bash
###########################################################
#Author:mengdj@outlook.com
#Created Time:2020.12.04 11:56
#Description:execute go mod tidy in current directory
#Version:0.0.3
#File:tidy.sh
###########################################################

CURRENT_DIR=$(pwd)
SEARCH_DIR=$CURRENT_DIR
SEARCH_TOTAL=0
EXECUTE_CMD="go mod tidy"
EXECUTE_TIMESTAMP=`date +%s`

function GoTidy() {
	for file in $(ls $1); do
		local target="$1/$file"
		if [ -d $target ]; then
			cd $target
			#case
			if [ -f "go.mod" ];then
				`$EXECUTE_CMD`
				if [ $? -ne 0 ]; then
					break
				fi
				echo "process $target"
				let "SEARCH_TOTAL+=1"
			fi
			GoTidy $target
		fi
	done
}

#test
if [ $# -ne 0 ]; then
	if [ -d $1 ]; then
		cd $1
		SEARCH_DIR=$(pwd)
	else
		echo "$1 is not exist directory."
		exit
	fi
fi
#start
GoTidy $SEARCH_DIR
let "EXECUTE_TIMESTAMP=`date +%s`-EXECUTE_TIMESTAMP"
echo "processed($SEARCH_TOTAL),loss $EXECUTE_TIMESTAMP seconds."
#back directory
cd $CURRENT_DIR
