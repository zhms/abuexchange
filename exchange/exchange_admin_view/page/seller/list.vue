<template>
	<div class="container">
		<div class="handle-box">
			<el-button type="primary" icon="el-icon-plus" class="mr10" @click="handleAdd()">添加</el-button>
			<el-button type="primary" class="mr10" @click="handleQuery()">刷新</el-button>
		</div>
		<div>
			<el-table :data="table_data" border class="table" max-height="700px" :cell-style="{ padding: '0' }">
				<el-table-column align="center" prop="id" label="编号" width="100"></el-table-column>
				<el-table-column align="center" prop="SellerId" label="运营商id" width="230"></el-table-column>
				<el-table-column align="center" prop="SellerName" label="运营商称" width="150"></el-table-column>
				<el-table-column align="center" prop="Note" label="状态" width="100">
					<template slot-scope="scope">
						<el-link :underline="false" type="danger" v-if="scope.row.IsDisabled == 1">禁用</el-link>
						<el-link :underline="false" type="primary" v-if="scope.row.IsDisabled == 0">启用</el-link>
					</template>
				</el-table-column>
				<el-table-column align="center" prop="Note" label="备注" width="300"></el-table-column>
				<el-table-column label="房间在线">
					<template slot-scope="scope">
						<el-button type="text" icon="el-icon-edit" @click="handleModify(scope.$index)">修改</el-button>
						<el-button type="text" icon="el-icon-delete" class="red" @click="handleDelete(scope.$index)">删除</el-button>
					</template>
				</el-table-column>
			</el-table>
		</div>
		<div>
			<el-dialog :title="dialog_title" :visible.sync="dialog" width="400px" center>
				<el-form :inline="true" label-width="100px">
					<el-form-item label="运营商Id:">
						<el-input v-model="dialog_data.SellerId" :disabled="dialog_title == '修改运营商'"></el-input>
					</el-form-item>
					<el-form-item label="运营商名称:">
						<el-input v-model="dialog_data.SellerName"></el-input>
					</el-form-item>
					<el-form-item label="选项:">
						<el-checkbox border label="禁用" v-model="dialog_data.IsDisabled" :true-label="1" :false-label="0"></el-checkbox>
					</el-form-item>
					<el-form-item label="注释:">
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
			dialog_data: [],
			dialog_title: null,
			dialog: false,
			current_row: 0,
		}
	},
	components: {},
	computed: {},
	created() {
		this.handleQuery()
	},
	methods: {
		handleQuery() {
			app.getInstance().post('/seller/list/query', {}, (result) => {
				this.table_data = result
			})
		},
		handleAdd() {
			this.dialog_title = '添加运营商'
			this.dialog_data = { IsDisabled: 0 }
			this.dialog = true
		},
		handleModify(index) {
			this.dialog_title = '修改运营商'
			this.current_row = index
			this.dialog_data = app.getInstance().clone(this.table_data[index])
			this.dialog = true
		},
		handleConfirm() {
			if (this.dialog_title == '添加运营商') {
				app.getInstance().post('/seller/list/add', this.dialog_data, (result) => {
					this.table_data.push(this.dialog_data)
					this.table_data = app.getInstance().clone(this.table_data)
					this.dialog = false
				})
			}
			if (this.dialog_title == '修改运营商') {
				app.getInstance().post('/seller/list/modify', this.dialog_data, (result) => {
					this.table_data[this.current_row] = this.dialog_data
					this.table_data = app.getInstance().clone(this.table_data)
					this.dialog = false
				})
			}
		},
		handleDelete(index) {
			if (confirm('确定删除该运营商?')) {
				app.getInstance().post('/seller/list/delete', { id: this.table_data[index].id }, (result) => {
					this.table_data.splice(index, 1)
				})
			}
		},
	},
}
</script>
