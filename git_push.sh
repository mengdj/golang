#!/usr/bin/env bash
#
#Author:1824854886@qq.com
#Created Time:2020.12.04 15:46
#Description:batch push file to github
#

git checkout dev
git pull
git add .
git rm -rf --cached .idea
git rm -rf --cached bw/boss/boss
git rm -rf --cached bw/worker/worker
#输入提交注释语句，15秒没有输入就自动终止 
read -p "please enter comment:" -t 15 -a comment

if test ${#comment} -gt 0;then
	git commit -m "$comment"
	git push
else
	echo "comment can't null"
	exit
fi

