#!/bin/bash

# @Author: mknight(tianyh)
# @Mail: 824338670@qq.com
# @Date: 2022-08-19 14:14:17
# @LastEditTime: 2022-06-15 14:14:17
# @Brief: 服务进程管理

# 启动服务器进程

WORK_DIR=$(pwd)
CONFIG_PATH=$WORK_DIR/app_config
checkRet=0

echo "work_dir="$WORK_DIR
echo "config_path="$CONFIG_PATH


function startServer() {
	#echo "first param is $1"
	#echo "second paramis $2"
 	nohup "$WORK_DIR"/bin/"$1" --configPath="$WORK_DIR/app_config" --index="$2" >> "$WORK_DIR"/logs/"$1".log 2>&1 &
}

function showServers() {
	queryResult=$(ps -ef | grep "$WORK_DIR/app_config" | grep "roomcell" | awk '{print $2,$8,$10}')	
	echo "showServers queryResult"
	echo $queryResult
}

function checkServer() {
	server_name=$1
	server_index=$2
	#echo "function checkServer,param1:$1, param2:$2"
	pid=$(pgrep -f "$WORK_DIR/bin/$server_name --configPath=$WORK_DIR/app_config --index=$server_index")
	if [ -z "$pid" ]; then
		#progress not running
		return 0 
	else
		#progress is running 
		#echo "find server $1 $2"
		return 1 
	fi
	return 0 
}

function checkServerPid() {
	server_name=$1
	server_index=$2
	#echo "function checkServer,param1:$1, param2:$2"
	pid=$(pgrep -f "$WORK_DIR/bin/$server_name --configPath=$WORK_DIR/app_config --index=$server_index")
	if [ -z "$pid" ]; then
		#progress not running
		return "0" 
	else
		#progress is running 
		#echo "find server $1 $2"
		return $pid 
	fi
	return 0 
}

function stopServerPid() {
	server_name=$1
	server_index=$2
	#echo "function checkServer,param1:$1, param2:$2"
	for(( j=0;j<100;j++)) do
		pid=$(pgrep -f "$WORK_DIR/bin/$server_name --configPath=$WORK_DIR/app_config --index=$server_index")
		if [ -z "$pid" ]; then
			#progress not running
			echo -e "\e[;35m$1 $2 stopped\e[0m"
			j=100
		else
			echo -e "stoping \e[;35mstoping $server_name $server_index pid=$pid\e[0m"
			kill -9 $pid
		fi
		sleep 1
	done
	
	return 0
}


# 服务器列表定义
#declare -A serverList
serverList[0]="cellserv_account 0" 
serverList[1]="cellserv_root 0" 
serverList[2]="cellserv_data 0"
serverList[3]="cellserv_log 0"
serverList[4]="cellserv_view 0" 
serverList[5]="cellserv_center 0" 
serverList[6]="cellserv_game 0" 
serverList[7]="cellserv_logic 0" 
serverList[8]="cellserv_gate 0" 
# 服务器关闭列表
stopList[0]="cellserv_account 0" 
stopList[1]="cellserv_gate 0"
stopList[2]="cellserv_view 0"
stopList[3]="cellserv_logic 0" 
stopList[4]="cellserv_game 0"
stopList[5]="cellserv_center 0"
stopList[6]="cellserv_log 0"
stopList[7]="cellserv_data 0"
stopList[8]="cellserv_root 0" 


# 开启服务器
#startServer cellserv_account 0
# 查看服务器
#showServers

if [ $# -lt 1 ]; then
	echo -e "\e[;31mparam is less\e[0m"
	exit
fi


case "$1" in 
	start)
		echo "------ [start servers] -----"
		for(( i=0;i<${#serverList[@]};i++)) do
		#for one_server in ${serverList[*]}; do
			one_server=${serverList[i]}
			checkServerPid ${serverList[i]} 
			checkRet=$?	
			#echo "function return:$?"
			if [ $checkRet -eq 0 ]; then
				echo -e "\e[;33mstart $one_server ...\e[0m"
				startServer $one_server
				sleep 1
			else			
				echo -e "$checkRet $one_server\e[;32m is running\e[0m"
			fi
			sleep 1
		done
		echo "----- [servers status] ------"
		for(( i=0;i<${#serverList[@]};i++)) do
		#for one_server in ${serverList[*]}; do
			one_server=${serverList[i]}
			checkServerPid ${serverList[i]} 
			checkRet=$?	
			#echo "function return:$?"
			if [ $checkRet -eq 0 ]; then
				echo -e "pid=$checkRet   $one_server\e[;31m is not running\e[0m"
			else			
				echo -e "pid=$checkRet $one_server\e[;32m is running\e[0m"
			fi
		done
		;;
	status)
		echo "----- [servers status] ------"
		for(( i=0;i<${#serverList[@]};i++)) do
		#for one_server in ${serverList[*]}; do
			one_server=${serverList[i]}
			checkServerPid ${serverList[i]} 
			checkRet=$?	
			#echo "function return:$?"
			if [ $checkRet -eq 0 ]; then
				echo -e "pid=$checkRet   $one_server\e[;31m is not running\e[0m"
			else			
				echo -e "pid=$checkRet $one_server\e[;32m is running\e[0m"
			fi
		done
		;;
	check)
		for(( i=0;i<${#serverList[@]};i++)) do
		#for one_server in ${serverList[*]}; do
			one_server=${serverList[i]}
			checkServer ${serverList[i]} 
			checkRet=$?	
			echo "function return:$?"
			if [ $checkRet -eq 1 ]; then
				echo -e "$one_server\e[;32m is running\e[0m"
			else			
				echo -e "$one_server\e[;32m is not running\e[0m"
			fi
		done
		;;
	stop)
		for(( i=0;i<${#stopList[@]};i++)) do
			one_server=${stopList[i]}
			#echo -e "\e[;36mstop $one_server\e[0m"
			stopServerPid ${stopList[i]} 
			checkRet=$?	
			#echo "stop $one_server $checkRet"
			sleep 1
		done
		;;
	*)
		echo -e "\e[;31merror option\e[0m"
		;;
esac

