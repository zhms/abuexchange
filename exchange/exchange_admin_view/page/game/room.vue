<template>
	<div class="container">
		<!-- 筛选 -->
		<div>
			<el-form :inline="true" :model="filters">
				<el-form-item label="游戏:">
					<el-select v-model="filters.GameId" placeholder="请选择" style="width:150px;margin-right: 10px;">
						<el-option v-for="item in games" :key="item.GameName" :label="item.GameName" :value="item.GameId"> </el-option>
					</el-select>
					<el-form-item>
						<el-button type="primary" v-on:click="handleQuery">查询</el-button>
					</el-form-item>
					<el-form-item>
						<el-button type="primary" icon="el-icon-plus" v-show="auth('增')" @click="handleAdd">添加</el-button>
					</el-form-item>
				</el-form-item>
			</el-form>
		</div>
		<!-- 表 -->
		<div>
			<el-table :data="table_data" border max-height="620px" class="table" :cell-style="{ padding: '0' }">
				<el-table-column align="center" prop="GameId" label="游戏" width="200"></el-table-column>
				<el-table-column align="center" prop="RoomLevel" label="等级" width="80"></el-table-column>
				<el-table-column align="center" prop="RoomName" label="名称" width="150"></el-table-column>
				<el-table-column align="center" label="状态" width="100">
					<template slot-scope="scope">
						<el-link :underline="false" type="danger" v-if="scope.row.IsDisabled == 1">禁用</el-link>
						<el-link :underline="false" type="primary" v-if="scope.row.IsDisabled == 0">启用</el-link>
					</template>
				</el-table-column>
				<el-table-column align="center" prop="BasePoint" label="底分" width="100"></el-table-column>
				<el-table-column align="center" prop="EnterMinScore" label="进入最少金币" width="100"></el-table-column>
				<el-table-column align="center" prop="EnterMaxScore" label="最多最多金币" width="100"></el-table-column>
				<el-table-column align="center" label="机器人投放" width="100">
					<template slot-scope="scope">
						<el-button type="text" icon="el-icon-edit" @click="handleModifyRobot(scope.$index)">编辑</el-button>
					</template>
				</el-table-column>
				<el-table-column label="操作">
					<template slot-scope="scope">
						<el-button type="text" icon="el-icon-edit" v-show="auth('改')" @click="handleModify(scope.$index)">编辑</el-button>
						<el-button type="text" icon="el-icon-delete" class="red" v-show="auth('删')" @click="handleDel(scope.$index)">删除</el-button>
					</template>
				</el-table-column>
			</el-table>
		</div>
		<!--对话框-->
		<div>
			<el-dialog :title="dialog_title" :visible.sync="dialog" width="400px" center>
				<el-form :inline="true" label-width="100px">
					<el-form-item label="游戏:">
						<el-select v-model="dialog_data.GameId" placeholder="请选择" style="width:200px;margin-right: 10px;" :disabled="dialog_title != '添加房间'">
							<el-option v-for="item in games_add" :key="item.GameName" :label="item.GameName" :value="item.GameId"> </el-option>
						</el-select>
					</el-form-item>
					<el-form-item label="房间等级:">
						<el-input v-model="dialog_data.RoomLevel" :disabled="dialog_title != '添加房间'"></el-input>
					</el-form-item>
					<el-form-item label="房间名称:">
						<el-input v-model="dialog_data.RoomName"></el-input>
					</el-form-item>
					<el-form-item label="底分:">
						<el-input v-model="dialog_data.BasePoint"></el-input>
					</el-form-item>
					<el-form-item label="进入最低金币:">
						<el-input v-model="dialog_data.EnterMinScore" placeholder="0表示不限"></el-input>
					</el-form-item>
					<el-form-item label="进入最高金币:">
						<el-input v-model="dialog_data.EnterMaxScore" placeholder="0表示不限"></el-input>
					</el-form-item>
					<el-form-item label="标签:">
						<template>
							<el-radio v-model="dialog_radio" label="0">无</el-radio>
							<el-radio v-model="dialog_radio" label="1">最新</el-radio>
							<el-radio v-model="dialog_radio" label="2">火爆</el-radio>
						</template>
					</el-form-item>
					<el-form-item label="选项:">
						<el-checkbox border label="禁用" v-model="dialog_data.IsDisabled" :true-label="1" :false-label="0"></el-checkbox>
					</el-form-item>
					<el-form-item label="备注:">
						<el-input v-model="dialog_data.Note"></el-input>
					</el-form-item>
				</el-form>
				<span slot="footer" class="dialog-footer">
					<el-button type="primary" @click="handleConfirm">确 定</el-button>
				</span>
			</el-dialog>
		</div>
		<div>
			<el-dialog :title="dialog_title" :visible.sync="dialog_robot" width="400px" center>
				<el-form :inline="true" label-width="100px">
					<el-form-item label="机器人数量:">
						<el-input v-model="dialog_robot_data.RobotCount" placeholder="0表示不投放机器人"></el-input>
					</el-form-item>
					<el-form-item label="最少携带金币:">
						<el-input v-model="dialog_robot_data.EnterMinScore" placeholder="不能低于房间要求"></el-input>
					</el-form-item>
					<el-form-item label="最多携带金币:">
						<el-input v-model="dialog_robot_data.EnterMaxScore" placeholder="不能高于房间要求"></el-input>
					</el-form-item>
					<el-form-item label="最低局数:">
						<el-input v-model="dialog_robot_data.PlayMinRound" placeholder=""></el-input>
					</el-form-item>
					<el-form-item label="最多局数:">
						<el-input v-model="dialog_robot_data.PlayMaxRound" placeholder=""></el-input>
					</el-form-item>
				</el-form>
				<span slot="footer" class="dialog-footer">
					<el-button type="primary" @click="handleConfirm">确 定</el-button>
				</span>
			</el-dialog>
		</div>
		<LongHuDaZhanRoomConfig :config_title="config_title" :show="configShow" :confirm="configConfirm"> </LongHuDaZhanRoomConfig>
	</div>
</template>

<script>
import { app } from '@/api/app.js'
import '@/assets/css/k.css'
import LongHuDaZhanRoomConfig from './room_long_hu_da_zhan.vue'
export default {
	data() {
		return {
			filters: {
				GameId: null,
			},
			games_add: [],
			games: [],
			table_data: [],
			dialog: false,
			dialog_title: '',
			dialog_data: {},
			current_row: null,
			dialog_robot: false,
			dialog_robot_data: {},
			dialog_config: false,
			dialog_config_data: null,
			dialog_radio: '0',
			config_title: '',
			config_show: {},
			robot: null,
		}
	},
	components: { LongHuDaZhanRoomConfig },
	computed: {},
	created() {
		app.getInstance().post('/game/game/query', {}, (games) => {
			this.games = [{ GameName: '全部', GameId: null }]
			for (var i = 0; i < games.length; i++) {
				var v = { GameName: games[i].GameName, GameId: games[i].GameId }
				this.games.push(v)
				this.games_add.push(v)
			}
		})
	},
	methods: {
		auth(o) {
			return app.getInstance().auth('游戏管理', '房间列表', o)
		},
		configShow(name, callback) {
			this.config_show[name] = callback
		},
		configConfirm(config_str) {
			this.table_data[this.current_row].config = config_str
		},
		handleConfig(index) {
			this.current_row = index
			this.dialog_data = app.getInstance().clone(this.table_data[index])
			this.config_title = this.dialog_data.name + '配置'
			if (this.config_show[this.dialog_data.gameid]) this.config_show[this.dialog_data.gameid](this.dialog_data)
		},
		handleQuery() {
			app.getInstance().post('/game/room/query', this.filters, (result) => {
				this.table_data = result.rdata
				this.robot = result.rcdata
			})
		},
		handleAdd() {
			this.dialog_data = { disabled: 0 }
			this.dialog_title = '添加房间'
			this.dialog_radio = '0'
			this.dialog = true
		},
		handleModify(index) {
			this.current_row = index
			this.dialog_title = '修改房间'
			this.dialog_data = app.getInstance().clone(this.table_data[index])
			this.dialog_radio = '0'
			if (this.dialog_data.tag == null) {
				this.dialog_data.tag = ''
			}
			if (this.dialog_data.tag.indexOf('new') >= 0) {
				this.dialog_radio = '1'
			}
			if (this.dialog_data.tag.indexOf('hot') >= 0) {
				this.dialog_radio = '2'
			}
			this.dialog = true
		},
		handleDel(index) {
			if (confirm('确定删除该房间?')) {
				app.getInstance().post('/game/room/delete', { GameId: this.table_data[index].GameId, RoomLevel: this.table_data[index].RoomLevel }, () => {
					this.table_data.splice(index, 1)
				})
			}
		},
		handleConfirm() {
			if (this.dialog_title == '添加房间') {
				if (!this.dialog_data.IsDisabled) this.dialog_data.IsDisabled = 0
				app.getInstance().post('/game/room/add', this.dialog_data, (result) => {
					this.table_data.push(app.getInstance().clone(this.dialog_data))
					this.dialog = false
					this.$message.success('操作成功')
				})
			}
			if (this.dialog_title == '修改房间') {
				this.dialog_data.RoomTag = ''
				if (this.dialog_radio == '1') {
					if (this.dialog_data.tag.length > 0) this.dialog_data.RoomTag += ','
					this.dialog_data.RoomTag += 'new'
				}
				if (this.dialog_radio == '2') {
					if (this.dialog_data.RoomTag.length > 0) this.dialog_data.RoomTag += ','
					this.dialog_data.RoomTag += 'hot'
				}
				app.getInstance().post('/game/room/modify', this.dialog_data, (result) => {
					this.table_data[this.current_row] = app.getInstance().clone(this.dialog_data)
					this.dialog = false
					this.$message.success('操作成功')
				})
			}
			if (this.dialog_title.indexOf('编辑机器人') >= 0) {
				app.getInstance().post('/game/room/robot', this.dialog_robot_data, (result) => {
					for (var i = 0; i < this.robot.length; i++) {
						if (this.robot[i].GameId == this.table_data[this.current_row].GameId && this.robot[i].RoomLevel == this.table_data[this.current_row].RoomLevel) {
							this.robot[i] = app.getInstance().clone(this.dialog_robot_data)
						}
					}
					this.dialog_robot = false
					this.$message.success('操作成功')
				})
			}
			if (this.dialog_title == '房间配置') {
				var data = {
					gameid: this.table_data[this.current_row].gameid,
					level: this.table_data[this.current_row].level,
					config: this.dialog_config_data,
				}
				app.getInstance().post('/game/room/config', data, (result) => {
					this.table_data[this.current_row].config = this.dialog_config_data
					this.dialog_config = false
				})
			}
		},
		handleModifyConfig(index) {
			this.current_row = index
			this.dialog_title = '房间配置'
			this.dialog_config_data = app.getInstance().clone(this.table_data[index]).config
			this.dialog_config = true
		},
		handleModifyRobot(index) {
			this.current_row = index
			this.dialog_title = '编辑机器人 - ' + this.table_data[index].RoomName
			this.dialog_robot_data = { GameId: this.table_data[index].GameId, RoomLevel: this.table_data[index].RoomLevel }
			for (var i = 0; i < this.robot.length; i++) {
				if (this.robot[i].GameId == this.table_data[index].GameId && this.robot[i].RoomLevel == this.table_data[index].RoomLevel) {
					this.dialog_robot_data = app.getInstance().clone(this.robot[i])
				}
			}
			this.dialog_robot = true
		},
		handleEditRobot() {
			var input = document.getElementById('kcopy')
			if (!input) {
				input = document.createElement('input')
				input.setAttribute('id', 'kcopy')
				input.setAttribute('readonly', 'readonly')
				document.body.appendChild(input)
			}
			input.setAttribute('value', this.dialog_robot_data.config)
			input.select()
			input.setSelectionRange(0, 9999)
			document.execCommand('Copy')
			window.open('https://www.bejson.com/jsoneditoronline/')
		},
	},
}
</script>
