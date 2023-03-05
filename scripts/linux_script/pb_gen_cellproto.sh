#/bin/bash

#@Author: mknight
#@Mail: 824338670@qq.com
#@Date: 2023-01-17 14:24:50
#@LastEditTime: 2022-09-29 18:22:26
#@Describe: protoc执行文件需要在系统路径能找到或者在当前文件夹下面


#使用说明
#./pb_gen_cellproto.sh [文件名]
# 如果不带文件名参数,生成所有
# 如果带文件参数名,生成单个文件 eg: ./pb_gen_cellproto.sh c_misc.proto


CMD_MODE="all"

if [ $# -lt 1 ]; then
	echo -e "\e[;31mmode is all\e[0m"
else
    CMD_MODE="single"
    echo -e "\e[;31mmode is single\e[0m"
fi

SCRIPT_PATH=`(pwd)`
echo -e "SCRIPT_PATH=$SCRIPT_PATH"
cd $SCRIPT_PATH/../../
PROJECT_ROOT=`(pwd)`
cd $PROJECT_ROOT
echo -e "PROJECT_ROOT"=$PROJECT_ROOT
PROTO_ROOT_PATH=$PROJECT_ROOT/protos
echo -e "PROTO_ROOT_PATH=$PROTO_ROOT_PATH"
PKG_PB_ROOT_PATH=$PROJECT_ROOT/pkg/pb/

SRC_CLIENT=$PROTO_ROOT_PATH/proto_client
TAR_CLIENT_PKG=$PKG_PB_ROOT_PATH/pbclient
SRC_SERVER=$PROTO_ROOT_PATH/proto_server

# 客户端
CLIENT_PROTOS=`(ls -a $SRC_CLIENT)`
for proto_file in $CLIENT_PROTOS; do
    if [ $proto_file == "." ] || [ $proto_file == ".." ]; then
        continue
    fi
    if [ $CMD_MODE != "all" ] && [ $proto_file != "$1" ] ; then
        continue
    fi
    echo -e "$proto_file"
    protoc --proto_path=$SRC_CLIENT --go_out=$PKG_PB_ROOT_PATH  $proto_file
    if [ $CMD_MODE != "all" ]; then
        break
    fi
done

# 服务器
SERVER_PROTOS=`(ls -a $SRC_SERVER)`
for proto_file in $SERVER_PROTOS; do
    if [ $proto_file == "." ]; then
        continue
    fi
    if [ $proto_file == ".." ]; then
        continue
    fi
    if [ $CMD_MODE != "all" ] && [ $proto_file != "$1" ] ; then
        continue
    fi
    echo -e "$proto_file"
    protoc --proto_path=$SRC_SERVER --go_out=$PKG_PB_ROOT_PATH  $proto_file
    if [ $CMD_MODE != "all" ]; then
        break
    fi
done
# protoc --proto_path=%SRC_SERVER% --go_out=%TAR_SERVER% %SRC_SERVER%\s_frame.proto
# protoc --proto_path=%SRC_SERVER% --proto_path=%PROJECT_ROOT% --proto_path=%PROTO_ROOT_PATH% --go_out=%TAR_SERVER% %SRC_SERVER%\s_common.proto
# protoc --proto_path=%SRC_SERVER% --proto_path=%PROJECT_ROOT% --proto_path=%PROTO_ROOT_PATH% --go_out=%TAR_SERVER% %SRC_SERVER%\s_player.proto
# protoc --proto_path=%SRC_SERVER% --proto_path=%PROJECT_ROOT% --proto_path=%PROTO_ROOT_PATH% --go_out=%TAR_SERVER% %SRC_SERVER%\s_db.proto
# protoc --proto_path=%SRC_SERVER% --proto_path=%PROJECT_ROOT% --proto_path=%PROTO_ROOT_PATH% --go_out=%TAR_SERVER% %SRC_SERVER%\s_fieldobj_save.proto
# protoc --proto_path=%SRC_SERVER% --proto_path=%PROJECT_ROOT% --proto_path=%PROTO_ROOT_PATH% --go_out=%TAR_SERVER% %SRC_SERVER%\s_servdata.proto
exit
