@echo off

set PROJECT_ROOT=F:\work_home\server_projects\trcell

set PROTO_ROOT_PATH=%PROJECT_ROOT%\protos\

set SRC_CLIENT=%PROJECT_ROOT%\protos\proto_client
@REM set TAR_CLIENT=%PROJECT_ROOT%\pkg\pb\
set TAR_CLIENT_PKG_ROOT=F:\work_home\server_projects\\

set SRC_SERVER=%PROJECT_ROOT%\protos\proto_server
set TAR_SERVER=%PROJECT_ROOT%\pkg\pb\

:: 客户端
protoc --proto_path=%SRC_CLIENT% --go_out=%TAR_CLIENT_PKG_ROOT%  %SRC_CLIENT%\c_base.proto
protoc --proto_path=%SRC_CLIENT% --go_out=%TAR_CLIENT_PKG_ROOT%  %SRC_CLIENT%\c_base_ext.proto
protoc --proto_path=%SRC_CLIENT% --go_out=%TAR_CLIENT_PKG_ROOT%  %SRC_CLIENT%\c_common.proto
protoc --proto_path=%SRC_CLIENT% --go_out=%TAR_CLIENT_PKG_ROOT%  %SRC_CLIENT%\c_player.proto

:: 服务器
protoc --proto_path=%SRC_SERVER% --go_out=%TAR_SERVER% %SRC_SERVER%\s_frame.proto
protoc --proto_path=%SRC_SERVER% --proto_path=%PROJECT_ROOT% --proto_path=%PROTO_ROOT_PATH% --go_out=%TAR_SERVER% %SRC_SERVER%\s_common.proto
protoc --proto_path=%SRC_SERVER% --proto_path=%PROJECT_ROOT% --proto_path=%PROTO_ROOT_PATH% --go_out=%TAR_SERVER% %SRC_SERVER%\s_player.proto
protoc --proto_path=%SRC_SERVER% --proto_path=%PROJECT_ROOT% --proto_path=%PROTO_ROOT_PATH% --go_out=%TAR_SERVER% %SRC_SERVER%\s_db.proto
protoc --proto_path=%SRC_SERVER% --proto_path=%PROJECT_ROOT% --proto_path=%PROTO_ROOT_PATH% --go_out=%TAR_SERVER% %SRC_SERVER%\s_fieldobj_save.proto
pause