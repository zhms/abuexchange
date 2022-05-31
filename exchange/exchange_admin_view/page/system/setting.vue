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
				<el-table-column align="center" prop="SettingName" label="设置名" width="200"></el-table-column>
				<el-table-column align="center" prop="SettingValue" label="设置值" width="400"></el-table-column>
				<el-table-column align="center" prop="Note" label="说明" width="500"></el-table-column>
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
					<el-form-item label="设置名:">
						<el-input v-model="dialog_data.SettingName" :disabled="dialog_title == '修改配置'"></el-input>
					</el-form-item>
					<el-form-item label="设置值:">
						<el-input v-model="dialog_data.SettingValue"></el-input>
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
			table_data: [],
			dialog: false,
			dialog_title: null,
			dialog_data: { SettingName: null, SettingValue: null, Note: null },
			current_row: null,
		}
	},
	components: {},
	computed: {},
	created() {
		this.handleQuery()
	},
	methods: {
		auth(o) {
			return app.auth2('系统管理', '系统设置', o)
		},
		handleQuery() {
			app.post('/system/config/query', {}, (data) => {
				this.table_data = data
			})
		},
		handleAdd() {
			this.dialog_data = {}
			this.dialog_title = '添加配置'
			this.dialog = true
		},
		handleModify(index) {
			this.current_row = index
			this.dialog_title = '修改配置'
			this.dialog_data = app.clone(this.table_data[index])
			this.dialog = true
		},
		handleDel(index) {
			if (confirm('确定删除该项配置?')) {
				app.post('/system/config/delete', { SettingName: this.table_data[index].SettingName }, (data) => {
					this.table_data.splice(index, 1)
					this.$message.success('操作成功')
				})
			}
		},
		handleConfirm() {
			if (this.dialog_title == '修改配置') {
				app.post('/system/config/modify', this.dialog_data, () => {
					this.table_data[this.current_row] = app.clone(this.dialog_data)
					this.dialog = false
					this.$message.success('操作成功')
				})
			}
			if (this.dialog_title == '添加配置') {
				if (!this.dialog_data.SettingName) {
					this.$message.error('请填写设置名')
					return
				}
				if (!this.dialog_data.SettingValue) {
					this.$message.error('请填写设置值')
					return
				}
				if (!this.dialog_data.Note) {
					this.$message.error('请填写备注')
					return
				}
				app.post('/system/config/add', this.dialog_data, () => {
					this.table_data.push(this.dialog_data)
					this.dialog = false
					this.$message.success('操作成功')
				})
			}
		},
	},
}
</script>
