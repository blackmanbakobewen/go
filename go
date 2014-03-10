#!/bin/bash

GET_CHAR()
{
  SAVEDSTTY=`stty -g`
  stty -echo
  stty raw
  dd if=/dev/tty bs=1 count=1 2> /dev/null
  stty -raw
  stty echo
  stty $SAVEDSTTY
}

GO_HOME=~
if [ -z $1 ];then
  echo -n "请输入要登陆的机器IP/域名： "
  read target
else
  target=$1
  shift
fi

target=".*$target.*"
AUTHFILE=$GO_HOME/.go.conf
count=`grep "$target" $AUTHFILE -c`
targetfullname=`grep "$target" $AUTHFILE | awk '{print $1}' | awk -F ':' '{print $1}'`
port=`grep "$target" $AUTHFILE | awk '{print $1}' | awk -F ':' '{print $2}'`
passwd=`grep "$target" $AUTHFILE | awk '{print $2}' | awk -F ':' '{print $2}'`
user=`grep "$target" $AUTHFILE | awk '{print $2}' | awk -F ':' '{print $1}'`
label=`grep "$target" $AUTHFILE | awk '{print $3}'`
if [ $count -gt 1 ];then
  echo -e '查找到以下主机 (\033[0;31m选择一个\033[0m)'
  arrtarget=($targetfullname)
  arruser=($user)
  arrpasswd=($passwd)
  arrlabel=($label)
  arrport=($port)
  length=${#arrtarget[@]}
  for ((i=0; i<$length; i++))
  do
    echo -e '[\033[4;34m'$(($i+1))'\033[0m]' "${arruser[$i]}@${arrtarget[$i]} ${arrlabel[$i]}"
  done
  echo -n "选择序号："
  choice=`GET_CHAR`
  echo ""

  echo $choice;

  if [[ "$choice" =~ ^[0-9]+$ ]]; then
    echo '';
  else
    exit 1;
  fi

  targetfullname=${arrtarget[$(($choice-1))]}
  passwd=${arrpasswd[$(($choice-1))]}
  user=${arruser[$(($choice-1))]}
  label=${arrencoding[$(($choice-1))]}
  port=${arrport[$(($choice-1))]}
fi

if [ -z $targetfullname ] || [ -z $passwd ] || [ -z $user ];then
  echo "配置文件中没有查找到匹配的信息";
  exit 1;
fi
target=$targetfullname

# Process options

while getopts g ARGS
do
case $ARGS in 
    g)
	extra_options="-D7070"
	;;
    *)
        echo "Unknow option: $ARGS"
        exit 1;
        ;;
esac
done

if [ -z $port ]; then
  port=22
fi

echo "正在登录${user}@${target} ${label}..."

ssh-expect $user $target $passwd $port $extra_options
