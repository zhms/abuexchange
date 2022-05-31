<template>
	<div class="container">
		<!-- 筛选 -->
		<div class="handle-box">
			<el-button type="primary" icon="el-icon-plus" class="mr10" v-show="auth('增')" @click="handleAdd">添加</el-button>
		</div>
		<!-- 表 -->
		<div>
			<el-table :data="table_data" border max-height="620px" class="table" :cell-style="{ padding: '0' }">
				<el-table-column align="center" prop="GameId" label="游戏id" width="200"></el-table-column>
				<el-table-column align="center" prop="GameName" label="游戏名称" width="150"></el-table-column>
				<el-table-column align="center" label="状态" width="100">
					<template slot-scope="scope">
						<el-link :underline="false" type="danger" v-if="scope.row.IsDisabled == 1">禁用</el-link>
						<el-link :underline="false" type="primary" v-if="scope.row.IsDisabled == 0">启用</el-link>
					</template>
				</el-table-column>
				<el-table-column align="center" prop="GameTax" label="税率" width="100"></el-table-column>
				<el-table-column align="center" prop="GameTag" label="标签" width="100"></el-table-column>
				<el-table-column label="游戏控制" width="100">
					<template slot-scope="scope">
						<el-button type="text" icon="el-icon-edit" v-show="auth('改')" @click="handleControl(scope.$index)">编辑</el-button>
					</template>
				</el-table-column>
				<el-table-column label="机器人配置" width="100">
					<template slot-scope="scope">
						<el-button type="text" icon="el-icon-edit" v-show="auth('改')" @click="handleRobot(scope.$index)">编辑</el-button>
					</template>
				</el-table-column>
				<el-table-column label="游戏配置" width="100">
					<template slot-scope="scope">
						<el-button type="text" icon="el-icon-edit" v-show="auth('改')" @click="handleConfig(scope.$index)">编辑</el-button>
					</template>
				</el-table-column>
				<el-table-column align="center" prop="Sort" label="排序" width="100"></el-table-column>
				<el-table-column label="操作" width="200">
					<template slot-scope="scope">
						<el-button type="text" icon="el-icon-edit" v-show="auth('改')" @click="handleModify(scope.$index)">编辑</el-button>
						<el-button type="text" icon="el-icon-delete" class="red" v-show="auth('删')" @click="handleDel(scope.$index)">删除</el-button>
					</template>
				</el-table-column>
				<!-- <el-table-column align="center" prop="Note" label="备注"></el-table-column> -->
			</el-table>
		</div>
		<!--对话框-->
		<div>
			<el-dialog :title="modify_or_new.title" :visible.sync="modify_or_new.show" width="400px" center>
				<el-form :inline="true" label-width="100px">
					<el-form-item label="游戏id:">
						<el-input v-model="dialog_data.GameId" :disabled="modify_or_new.title == '修改游戏'"></el-input>
					</el-form-item>
					<el-form-item label="游戏名称:">
						<el-input v-model="dialog_data.GameName"></el-input>
					</el-form-item>
					<el-form-item label="排序:">
						<el-input v-model="dialog_data.Sort"></el-input>
					</el-form-item>
					<el-form-item label="选项:">
						<el-checkbox border label="禁用" v-model="dialog_data.IsDisabled" :true-label="1" :false-label="0"></el-checkbox>
						<el-checkbox border label="大图" v-model="modify_or_new.bigicon" :true-label="1" :false-label="0"></el-checkbox>
					</el-form-item>
					<el-form-item label="税率:">
						<el-input v-model="dialog_data.GameTax"></el-input>
					</el-form-item>
					<el-form-item label="标签:">
						<template>
							<el-radio v-model="modify_or_new.radio" label="0">无</el-radio>
							<el-radio v-model="modify_or_new.radio" label="1">最新</el-radio>
							<el-radio v-model="modify_or_new.radio" label="2">火爆</el-radio>
						</template>
					</el-form-item>
				</el-form>
				<span slot="footer" class="dialog-footer">
					<el-button type="primary" @click="handleConfirm">确 定</el-button>
				</span>
			</el-dialog>
		</div>
		<!--游戏配置-->
		<LongHuDaZhanConfig :config_title="config_title" :show="configShow" :confirm="configConfirm"></LongHuDaZhanConfig>
		<HongHeiDaZhanConfig :config_title="config_title" :show="configShow" :confirm="configConfirm"></HongHeiDaZhanConfig>
		<BaiJiaLeConfig :config_title="config_title" :show="configShow" :confirm="configConfirm"></BaiJiaLeConfig>
		<!--机器人配置-->
		<LongHuDaZhanRobot :robot_title="robot_title" :show="robotShow" :confirm="robotConfirm"></LongHuDaZhanRobot>
		<HongHeiDaZhanRobot :robot_title="robot_title" :show="robotShow" :confirm="robotConfirm"></HongHeiDaZhanRobot>
		<BaiJiaLeRobot :robot_title="robot_title" :show="robotShow" :confirm="robotConfirm"></BaiJiaLeRobot>
		<!--游戏控制-->
		<LongHuDaZhanControl :control_title="control_title" :show="controlShow" :confirm="controlConfirm"></LongHuDaZhanControl>
		<HongHeiDaZhanControl :control_title="control_title" :show="controlShow" :confirm="controlConfirm"></HongHeiDaZhanControl>
		<BaiJiaLeControl :control_title="control_title" :show="controlShow" :confirm="controlConfirm"></BaiJiaLeControl>
		<CaiShenDaoControl :control_title="control_title" :show="controlShow" :confirm="controlConfirm"></CaiShenDaoControl>
		<JiuXianLaBaControl :control_title="control_title" :show="controlShow" :confirm="controlConfirm"></JiuXianLaBaControl>
		<QiangZhuangNiuNiuControl :control_title="control_title" :show="controlShow" :confirm="controlConfirm"></QiangZhuangNiuNiuControl>
		<HaiWangBuYuControl :control_title="control_title" :show="controlShow" :confirm="controlConfirm"></HaiWangBuYuControl>
	</div>
</template>

<script>
import { app } from '@/api/app.js'
import '@/assets/css/k.css'
import LongHuDaZhanConfig from './game_config_long_hu_da_zhan.vue'
import HongHeiDaZhanConfig from './game_config_hong_hei_da_zhan.vue'
import BaiJiaLeConfig from './game_config_bai_jia_le.vue'
import LongHuDaZhanRobot from './game_robot_long_hu_da_zhan.vue'
import HongHeiDaZhanRobot from './game_robot_hong_hei_da_zhan.vue'
import BaiJiaLeRobot from './game_robot_bai_jia_le.vue'
import LongHuDaZhanControl from './game_control_long_hu_da_zhan.vue'
import HongHeiDaZhanControl from './game_control_hong_hei_da_zhan.vue'
import CaiShenDaoControl from './game_control_cai_shen_dao.vue'
import BaiJiaLeControl from './game_control_bai_jia_le.vue'
import JiuXianLaBaControl from './game_control_jiu_xian_la_ba.vue'
import QiangZhuangNiuNiuControl from './game_control_qiang_zhuang_niu_niu.vue'
import HaiWangBuYuControl from './game_control_hai_wang_bu_yu.vue'
export default {
	data() {
		return {
			table_data: [],
			dialog_data: {},
			current_row: null,
			config_title: '',
			control_title: '',
			robot_title: '',
			config_show: {},
			robot_show: {},
			control_show: {},
			modify_or_new: {
				show: false,
				title: '',
				radio: '0',
				bigicon: 0,
			},
		}
	},
	components: {
		LongHuDaZhanConfig,
		HongHeiDaZhanConfig,
		BaiJiaLeConfig,
		LongHuDaZhanRobot,
		HongHeiDaZhanRobot,
		BaiJiaLeRobot,
		LongHuDaZhanControl,
		HongHeiDaZhanControl,
		BaiJiaLeControl,
		CaiShenDaoControl,
		JiuXianLaBaControl,
		QiangZhuangNiuNiuControl,
		HaiWangBuYuControl,
	},
	computed: {},
	created() {
		this.handleQuery()
	},
	methods: {
		auth(o) {
			return app.auth2('游戏管理', '游戏列表', o)
		},
		handleQuery() {
			app.post('/game/game/query', {}, (result) => {
				this.table_data = result
			})
		},
		handleAdd() {
			this.modify_or_new.title = '新增游戏'
			this.dialog_data = { IsDisabled: 0 }
			this.modify_or_new.show = true
		},
		handleModify(index) {
			this.modify_or_new.title = '修改游戏'
			this.dialog_data = app.clone(this.table_data[index])
			this.current_row = index
			this.modify_or_new.radio = '0'
			this.modify_or_new.bigicon = 0
			if (this.dialog_data.GameTag == null) {
				this.dialog_data.GameTag = ''
			}
			if (this.dialog_data.GameTag.indexOf('big') >= 0) {
				this.modify_or_new.bigicon = 1
			}
			if (this.dialog_data.GameTag.indexOf('new') >= 0) {
				this.modify_or_new.radio = '1'
			}
			if (this.dialog_data.GameTag.indexOf('hot') >= 0) {
				this.modify_or_new.radio = '2'
			}
			this.modify_or_new.show = true
		},
		handleDel(index) {
			if (confirm('确定删除该游戏?')) {
				app.post('/game/game/delete', { GameId: this.table_data[index].GameId }, () => {
					this.table_data.splice(index, 1)
				})
			}
		},
		handleRobot(index) {
			this.current_row = index
			this.dialog_data = app.clone(this.table_data[index])
			this.robot_title = this.dialog_data.GameName + '机器人配置'
			if (this.robot_show[this.dialog_data.GameId]) {
				this.robot_show[this.dialog_data.GameId](this.dialog_data)
			} else {
				this.$message.error('该游戏无机器人')
			}
		},
		handleConfig(index) {
			this.current_row = index
			this.dialog_data = app.clone(this.table_data[index])
			this.config_title = this.dialog_data.GameName + '配置'
			if (this.config_show[this.dialog_data.GameId]) {
				this.config_show[this.dialog_data.GameId](this.dialog_data)
			} else {
				this.$message.error('该游戏无需配置')
			}
		},
		handleControl(index) {
			this.current_row = index
			this.dialog_data = app.clone(this.table_data[index])
			this.control_title = this.dialog_data.GameName + '控制'
			if (this.control_show[this.dialog_data.GameId]) {
				this.control_show[this.dialog_data.GameId](this.dialog_data)
			}
		},
		configShow(name, callback) {
			this.config_show[name] = callback
		},
		robotShow(name, callback) {
			this.robot_show[name] = callback
		},
		controlShow(name, callback) {
			this.control_show[name] = callback
		},
		configConfirm(config_str) {
			this.table_data[this.current_row].GameConfig = config_str
		},
		robotConfirm(robot_str) {
			this.table_data[this.current_row].GameRobot = robot_str
		},
		controlConfirm(control_str) {
			this.table_data[this.current_row].GameControl = control_str
		},
		handleConfirm() {
			if (this.modify_or_new.title == '新增游戏') {
				this.dialog_data.GameTag = ''
				if (this.modify_or_new.bigicon == 1) {
					if (this.dialog_data.GameTag.length > 0) this.dialog_data.GameTag += ','
					this.dialog_data.GameTag += 'big'
				}
				if (this.modify_or_new.radio == '1') {
					if (this.dialog_data.GameTag.length > 0) this.dialog_data.GameTag += ','
					this.dialog_data.GameTag += 'new'
				}
				if (this.modify_or_new.radio == '2') {
					if (this.dialog_data.GameTag.length > 0) this.dialog_data.GameTag += ','
					this.dialog_data.GameTag += 'hot'
				}
				app.post('/game/game/add', this.dialog_data, (result) => {
					this.table_data.push(app.clone(this.dialog_data))
					this.modify_or_new.show = false
					this.$message.success('操作成功')
				})
			}
			if (this.modify_or_new.title == '修改游戏') {
				this.dialog_data.GameTag = ''
				if (this.modify_or_new.bigicon == 1) {
					if (this.dialog_data.GameTag.length > 0) this.dialog_data.GameTag += ','
					this.dialog_data.GameTag += 'big'
				}
				if (this.modify_or_new.radio == '1') {
					if (this.dialog_data.GameTag.length > 0) this.dialog_data.GameTag += ','
					this.dialog_data.GameTag += 'new'
				}
				if (this.modify_or_new.radio == '2') {
					if (this.dialog_data.GameTag.length > 0) this.dialog_data.GameTag += ','
					this.dialog_data.GameTag += 'hot'
				}
				app.post('/game/game/modify', this.dialog_data, (result) => {
					this.table_data[this.current_row] = app.clone(this.dialog_data)
					this.modify_or_new.show = false
					this.$message.success('操作成功')
				})
			}
		},
	},
}
</script>
