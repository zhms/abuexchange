<template>
	<div class="container">
		<!--筛选-->
		<div class="handle-box">
			<el-button type="primary" icon="el-icon-plus" class="mr10" v-show="auth('增')" @click="handleAdd">添加</el-button>
			<el-button type="primary" icon="el-icon-refresh" class="mr10" @click="handleQuery">刷新</el-button>
		</div>
		<!--表-->
		<div>
			<el-table :data="table_data" border max-height="620px" class="table" :cell-style="{ padding: '0' }">
				<el-table-column align="center" prop="RoleName" label="角色名" width="150"></el-table-column>
				<el-table-column align="center" label="状态" width="100">
					<template slot-scope="scope">
						<el-link :underline="false" type="danger" v-if="scope.row.IsDisabled == 1">禁用</el-link>
						<el-link :underline="false" type="primary" v-if="scope.row.IsDisabled == 0">启用</el-link>
					</template>
				</el-table-column>
				<el-table-column align="center" prop="Note" label="备注" width="500"></el-table-column>
				<el-table-column label="操作">
					<template slot-scope="scope">
						<el-button type="text" icon="el-icon-edit" v-show="auth('改')" @click="handleModify(scope.$index)">编辑</el-button>
						<el-button type="text" icon="el-icon-delete" class="red" v-show="auth('删')" @click="handleDel(scope.$index)">删除</el-button>
					</template>
				</el-table-column>
			</el-table>
		</div>
		<!-- 对话框 -->
		<div>
			<el-dialog :title="dialog_title" :visible.sync="dialog" width="600px" center>
				<el-form :inline="true" label-width="100px">
					<el-form-item label="角色:">
						<el-input :disabled="dialog_title == '修改角色'" v-model="dialog_data.RoleName"></el-input>
					</el-form-item>
					<el-form-item>
						<el-checkbox border label="禁用" v-model="dialog_data.IsDisabled" :true-label="1" :false-label="0"></el-checkbox>
					</el-form-item>
				</el-form>
				<el-form :inline="true" label-width="100px" style="margin-top:-17px">
					<el-form-item label="玩家管理:"></el-form-item>
				</el-form>
				<el-form>
					<el-form-item label="账号管理:" label-width="150px" style="margin-top:-17px">
						<el-checkbox v-model="authdata.玩家管理.账号管理.查" :true-label="1" :false-label="0">查看</el-checkbox>
						<el-checkbox v-model="authdata.玩家管理.账号管理.冻结解冻" :true-label="1" :false-label="0">冻结解冻</el-checkbox>
						<el-checkbox v-model="authdata.玩家管理.账号管理.设置取消超管" :true-label="1" :false-label="0">设置取消超管</el-checkbox>
						<el-checkbox v-model="authdata.玩家管理.账号管理.赠送金币" :true-label="1" :false-label="0">赠送金币</el-checkbox>
						<el-checkbox v-model="authdata.玩家管理.账号管理.修改备注" :true-label="1" :false-label="0">修改备注</el-checkbox>
						<el-checkbox v-model="authdata.玩家管理.账号管理.设置取消测试" :true-label="1" :false-label="0">设置取消测试</el-checkbox>
					</el-form-item>
				</el-form>
				<!-- 系统管理 -->
				<el-form :inline="true" label-width="100px" style="margin-top:-17px">
					<el-form-item label="系统管理:"></el-form-item>
				</el-form>
				<el-form>
					<el-form-item label="系统设置:" label-width="150px" style="margin-top:-17px">
						<el-checkbox v-model="authdata.系统管理.系统设置.查" :true-label="1" :false-label="0">查看</el-checkbox>
						<el-checkbox v-model="authdata.系统管理.系统设置.增" :true-label="1" :false-label="0">新增</el-checkbox>
						<el-checkbox v-model="authdata.系统管理.系统设置.改" :true-label="1" :false-label="0">修改</el-checkbox>
						<el-checkbox v-model="authdata.系统管理.系统设置.删" :true-label="1" :false-label="0">删除</el-checkbox>
					</el-form-item>
					<el-form-item label="账号管理:" label-width="150px" style="margin-top:-17px">
						<el-checkbox v-model="authdata.系统管理.账号管理.查" :true-label="1" :false-label="0">查看</el-checkbox>
						<el-checkbox v-model="authdata.系统管理.账号管理.增" :true-label="1" :false-label="0">新增</el-checkbox>
						<el-checkbox v-model="authdata.系统管理.账号管理.改" :true-label="1" :false-label="0">修改</el-checkbox>
						<el-checkbox v-model="authdata.系统管理.账号管理.删" :true-label="1" :false-label="0">删除</el-checkbox>
					</el-form-item>
					<el-form-item label="角色管理:" label-width="150px" style="margin-top:-17px">
						<el-checkbox v-model="authdata.系统管理.角色管理.查" :true-label="1" :false-label="0">查看</el-checkbox>
						<el-checkbox v-model="authdata.系统管理.角色管理.增" :true-label="1" :false-label="0">新增</el-checkbox>
						<el-checkbox v-model="authdata.系统管理.角色管理.改" :true-label="1" :false-label="0">修改</el-checkbox>
						<el-checkbox v-model="authdata.系统管理.角色管理.删" :true-label="1" :false-label="0">删除</el-checkbox>
					</el-form-item>
					<el-form-item label="系统日志:" label-width="150px" style="margin-top:-17px">
						<el-checkbox v-model="authdata.系统管理.系统日志.查" :true-label="1" :false-label="0">查看</el-checkbox>
					</el-form-item>
				</el-form>
				<!-- 备注 -->
				<el-form :inline="true" label-width="100px">
					<el-form-item label="备注:">
						<el-input type="textarea" v-model="dialog_data.Note" :rows="4"></el-input>
					</el-form-item>
				</el-form>
				<!-- 确定按钮 -->
				<span slot="footer" class="dialog-footer">
					<el-button type="primary" @click="handleConfirm">确 定</el-button>
				</span>
			</el-dialog>
		</div>
	</div>
</template>
<script>
import { app } from '@/api/app.js'
import '@/assets/css/k.css'
export default {
	data() {
		return {
			authdata: {
				玩家管理: {
					账号管理: { 查: 1, 冻结解冻: 1, 设置取消超管: 1, 赠送金币: 1, 修改备注: 1, 设置取消测试: 1 },
					在线玩家: { 查: 1 },
				},
				游戏系统: {
					游戏管理: { 查: 1, 改: 1, 增: 1, 删: 1 },
					房间管理: { 查: 1, 改: 1, 增: 1, 删: 1 },
					游戏记录: { 查: 1 },
					牌局记录: { 查: 1 },
				},
				系统管理: {
					系统设置: { 查: 1, 改: 1, 增: 1, 删: 1 },
					账号管理: { 查: 1, 改: 1, 增: 1, 删: 1 },
					角色管理: { 查: 1, 改: 1, 增: 1, 删: 1 },
					系统日志: { 查: 1 },
				},
				统计报表: {
					平台统计: { 查: 1 },
					游戏统计: { 查: 1 },
					房间统计: { 查: 1 },
				},
			},
			table_data: [],
			dialog_title: '',
			dialog: false,
			current_row: null,
			dialog_data: {},
			authtemplete: null,
		}
	},
	components: {},
	computed: {},
	created() {
		function initauthdata(data) {
			for (var i in data) {
				if (data[i] === 1) {
					data[i] = 0
				} else {
					initauthdata(data[i])
				}
			}
		}
		initauthdata(this.authdata)
		this.authtemplete = app.getInstance().clone(this.authdata)
		this.handleQuery()
	},
	methods: {
		auth(o) {
			return app.getInstance().auth('系统管理', '角色管理', o)
		},
		handleQuery() {
			app.getInstance().post('/system/role/query', {}, (data) => {
				this.table_data = data
			})
		},
		handleAdd() {
			this.dialog_title = '添加角色'
			this.dialog_data = {
				IsDisabled: 0,
			}
			this.authdata = app.getInstance().clone(this.authtemplete)
			this.dialog = true
		},
		handleModify(index) {
			this.current_row = index
			this.dialog_title = '修改角色'
			this.dialog_data = {
				RoleName: this.table_data[index].RoleName,
				IsDisabled: this.table_data[index].IsDisabled,
				Note: this.table_data[index].Note,
			}
			this.authdata = app.getInstance().clone(this.authtemplete)
			function checkdata(data, ad) {
				for (var i in data) {
					if (data[i] === 1) {
						if ((ad != null || ad != undefined) && ad[i] != null && ad[i] != undefined) {
							ad[i] = 1
						}
					} else if (data[i] === 0) {
						if ((ad != null || ad != undefined) && ad[i] != null && ad[i] != undefined) {
							ad[i] = 0
						}
					} else {
						checkdata(data[i], ad[i])
					}
				}
			}
			checkdata(JSON.parse(this.table_data[index].AuthData), this.authdata)
			this.dialog = true
		},
		handleDel(index) {
			if (confirm('确定删除该角色?')) {
				app.getInstance().post('/system/role/delete', { RoleName: this.table_data[index].RoleName }, (data) => {
					this.table_data.splice(index, 1)
					this.$message.success('操作成功')
				})
			}
		},
		handleConfirm() {
			if (this.dialog_title == '添加角色') {
				this.dialog_data.AuthData = this.authdata
				app.getInstance().post('/system/role/add', this.dialog_data, () => {
					var data = app.getInstance().clone(this.dialog_data)
					data.AuthData = JSON.stringify(this.authdata)
					this.table_data.push(data)
					this.$message.success('操作成功')
					this.dialog = false
				})
			}
			if (this.dialog_title == '修改角色') {
				this.dialog_data.AuthData = this.authdata
				app.getInstance().post('/system/role/modify', this.dialog_data, () => {
					var data = app.getInstance().clone(this.dialog_data)
					data.AuthData = JSON.stringify(this.authdata)
					this.table_data[this.current_row] = data
					this.$message.success('操作成功')
					this.dialog = false
				})
			}
		},
	},
}
</script>
