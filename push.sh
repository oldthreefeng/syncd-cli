#!/usr/bin/env bash

git add .

echo -n "enter git commit message:"
read name
git commit -m "$name"

git push origin master
git remote remove origin
git remote add https://github.com/oldthreefeng/syncd-cli.git
git push origin master && \
git remote remove origin
git remote add git@gogs.wangke.co:go/syncd-cli.git