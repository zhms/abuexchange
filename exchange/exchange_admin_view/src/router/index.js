import Vue from 'vue'
import Router from 'vue-router'

Vue.use(Router)

export default new Router({
	routes: [
		{ path: '/', redirect: '/home' },
		{ path: '/login', component: () => import('../components/login.vue'), meta: { title: '登录' } },
		{ path: '*', redirect: '/404' },
		{
			path: '/',
			component: () => import('../components/home.vue'),
			meta: { title: '自述文件' },
			children: [
				//首页
				{ path: '/home', component: () => import('../../page/home.vue'), meta: { title: '系统首页' } },
				//玩家管理
				{ path: '/user_user', component: () => import('../../page/user/user.vue'), meta: { title: '账号管理' } },
				{ path: '/score_change_list', component: () => import('../../page/user/scorelist.vue'), meta: { title: '金流记录' } },
				//游戏管理
				{ path: '/game_game', component: () => import('../../page/game/game.vue'), meta: { title: '游戏列表' } },
				{ path: '/game_room', component: () => import('../../page/game/room.vue'), meta: { title: '房间列表' } },
				{ path: '/game_record', component: () => import('../../page/game/record.vue'), meta: { title: '游戏记录' } },
				{ path: '/game_detail', component: () => import('../../page/game/detail.vue'), meta: { title: '牌局记录' } },
				//服务管理
				{ path: '/server_list', component: () => import('../../page/server/list.vue'), meta: { title: '服务列表' } },
				//营商管理
				{ path: '/seller_list', component: () => import('../../page/seller/list.vue'), meta: { title: '营商列表' } },
				//个人操盘
				{ path: '/control_wbdef', component: () => import('../../page/control/wbdef.vue'), meta: { title: '黑白定义' } },
				{ path: '/control_list', component: () => import('../../page/control/list.vue'), meta: { title: '个控配置' } },
				{ path: '/control_person', component: () => import('../../page/control/person.vue'), meta: { title: '查看个控' } },
				//兑换管理
				{ path: '/cash_setting', component: () => import('../../page/cash/setting.vue'), meta: { title: '兑换设置' } },
				//统计报表
				{ path: '/statistic_online', component: () => import('../../page/statistic/online.vue'), meta: { title: '在线统计' } },
				{ path: '/statistic_platform', component: () => import('../../page/statistic/platform.vue'), meta: { title: '平台统计' } },
				{ path: '/statistic_game', component: () => import('../../page/statistic/game.vue'), meta: { title: '游戏统计' } },
				{ path: '/statistic_room', component: () => import('../../page/statistic/room.vue'), meta: { title: '房间统计' } },
				//系统管理
				{ path: '/system_login_log', component: () => import('../../page/system/loginlog.vue'), meta: { title: '登录日志' } },
				{ path: '/system_account', component: () => import('../../page/system/account.vue'), meta: { title: '账号管理' } },
				{ path: '/system_role', component: () => import('../../page/system/roles.vue'), meta: { title: '角色管理' } },
				{ path: '/system_log', component: () => import('../../page/system/log.vue'), meta: { title: '操作日志' } },
			],
		},
	],
})
