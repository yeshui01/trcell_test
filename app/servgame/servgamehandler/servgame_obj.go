/*
 * @Author: mknight(tianyh)
 * @Mail: 824338670@qq.com
 * @Date: 2022-09-19 15:20:02
 * @LastEditTime: 2022-09-19 15:21:11
 * @FilePath: \trcell\app\servgame\servgamehandler\servgame_obj.go
 */
package servgamehandler

import "trcell/app/servgame/iservgame"

var (
	servGame iservgame.IServGame
)

func InitServGameObj(iserv iservgame.IServGame) {
	servGame = iserv
}
