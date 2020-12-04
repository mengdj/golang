#!/usr/bin/env bash
#Author:mengdj@outlook.com
#Created Time:2020.12.04 11:56
#Description:execute go mod tidy in current directory

function search
{
    for file in `ls $1`  
    do  
        if test -d $1"/"$file;then
	     if test -f $1"/"$file"/go.mod";then
		echo "process "$1"/"$file
		cd $1"/"$file  
		go mod tidy
	     fi
             search $1"/"$file  
        fi  
    done  
}

CURRENT_DIR=`pwd`
SEARCH_DIR=$CURRENT_DIR
if test $# -ne 0;then
	if test -d $1;then
		cd $1
		SEARCH_DIR=`pwd`
	else
		echo $1" is not exist directory."
		exit
	fi
fi
echo "start in "$SEARCH_DIR
search $SEARCH_DIR
cd $CURRENT_DIR

