<template>
	<div class="container">
		<!-- 筛选 -->
		<div class="handle-box">
			<el-button type="primary" icon="el-icon-plus" class="mr10" v-show="auth('增')" @click="handleAdd">添加</el-button>
			<el-button type="primary" icon="el-icon-refresh" class="mr10" @click="handleQuery">刷新</el-button>
		</div>
		<!--表-->
		<div>
			<el-table :data="table_data" border max-height="620px" class="table" :cell-style="{ padding: '0' }">
				<el-table-column align="center" prop="AdminId" label="id" width="100"></el-table-column>
				<el-table-column align="center" prop="Account" label="账号" width="100"></el-table-column>
				<el-table-column align="center" prop="NickName" label="昵称" width="100"></el-table-column>
				<el-table-column align="center" prop="RoleName" label="角色" width="100"></el-table-column>
				<el-table-column align="center" label="状态" width="100">
					<template slot-scope="scope">
						<el-link :underline="false" type="danger" v-if="scope.row.IsDisabled == 1">禁用</el-link>
						<el-link :underline="false" type="primary" v-if="scope.row.IsDisabled == 0">启用</el-link>
					</template>
				</el-table-column>
				<el-table-column align="center" prop="Note" label="备注" width="300"></el-table-column>
				<el-table-column align="center" prop="LastLoginTime" label="登录时间" width="160"></el-table-column>
				<el-table-column align="center" prop="LastLoginIp" label="ip" width="120"></el-table-column>
				<el-table-column align="center" prop="LoginCount" label="登录次数" width="100"></el-table-column>
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
			<el-dialog :title="dialog_title" :visible.sync="dialog" width="415px" center>
				<el-form :inline="true" label-width="100px">
					<el-form-item label="账号:">
						<el-input v-model="dialog_data.Account" :disabled="dialog_title == '修改账号'"></el-input>
					</el-form-item>
				</el-form>
				<el-form :inline="true" label-width="100px">
					<el-form-item label="密码:">
						<el-input v-model="dialog_data.Password" show-password style="width: 200px"></el-input>
					</el-form-item>
					<el-form-item label="昵称:">
						<el-input v-model="dialog_data.NickName"></el-input>
					</el-form-item>
					<el-form-item label="角色:">
						<el-select v-model="dialog_data.RoleName" placeholder="请选择" style="width: 200px">
							<el-option v-for="item in options" :key="item.value" :label="item.value" :value="item.value"></el-option>
						</el-select>
					</el-form-item>
				</el-form>
				<el-form :inline="true" label-width="100px">
					<el-form-item label="选项:">
						<el-checkbox border label="禁用" v-model="dialog_data.IsDisabled" :true-label="1" :false-label="0"></el-checkbox>
					</el-form-item>
				</el-form>
				<el-form :inline="true" label-width="100px">
					<el-form-item label="备注:">
						<el-input type="textarea" v-model="dialog_data.Note" :rows="4"></el-input>
					</el-form-item>
				</el-form>
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
			options: [],
			table_data: [],
			dialog: false,
			dialog_title: '',
			dialog_data: {},
		}
	},
	components: {},
	computed: {},
	created() {
		this.handleQuery()
		app.post('/system/role/query', {}, (data) => {
			this.options = []
			for (var i = 0; i < data.length; i++) {
				this.options.push({ value: data[i].RoleName })
			}
		})
	},
	methods: {
		auth(o) {
			return app.auth2('系统管理', '账号管理', o)
		},
		handleQuery() {
			app.post('/system/account/query', {}, (data) => {
				this.table_data = data
			})
		},
		handleAdd() {
			this.dialog_data = { IsDisabled: 0 }
			this.dialog_title = '添加账号'
			this.dialog = true
		},
		handleModify(index) {
			this.current_row = index
			this.dialog_title = '修改账号'
			this.dialog_data = app.clone(this.table_data[index])
			this.dialog = true
		},
		handleDel(index) {
			if (confirm('确定删除该配置?')) {
				app.post('/system/account/delete', { Account: this.table_data[index].Account }, (data) => {
					this.table_data.splice(index, 1)
					this.$message.success('操作成功')
				})
			}
		},
		handleConfirm() {
			if (this.dialog_title == '修改账号') {
				app.post('/system/account/modify', this.dialog_data, () => {
					this.table_data[this.current_row] = app.clone(this.dialog_data)
					this.dialog = false
					this.$message.success('操作成功')
				})
			}
			if (this.dialog_title == '添加账号') {
				app.post('/system/account/add', this.dialog_data, (data) => {
					this.dialog_data.AdminId = data.AdminId
					this.table_data.push(this.dialog_data)
					this.dialog = false
					this.$message.success('操作成功')
				})
			}
		},
	},
}
</script>
