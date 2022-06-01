<template>
	<div class="container">
		<!-- 筛选 -->
		<div>
			<el-form :inline="true" :model="filters">
				<el-form-item label="账号:">
					<el-input v-model="filters.Account" style="width: 150px" :clearable="true"></el-input>
				</el-form-item>
				<el-form-item label="运营商:" v-show="zong">
					<el-select v-model="filters.SellerId" placeholder="运营商" style="width: 130px">
						<el-option v-for="item in seller" :key="item.SellerId" :label="item.SellerName" :value="item.SellerId"> </el-option>
					</el-select>
				</el-form-item>
				<el-form-item>
					<el-button type="primary" icon="el-icon-refresh" class="mr10" @click="handleQuery">查询</el-button>
					<el-button type="primary" icon="el-icon-plus" class="mr10" v-show="auth('增')" @click="handleAdd">添加</el-button>
				</el-form-item>
			</el-form>
		</div>
		<!--表-->
		<div>
			<el-table :data="table_data" border max-height="620px" class="table" :cell-style="{ padding: '0' }">
				<el-table-column align="center" prop="Id" label="id" width="80"></el-table-column>
				<el-table-column align="center" prop="Account" label="账号" width="100"></el-table-column>
				<el-table-column align="center" prop="SellerName" label="运营商" width="130"></el-table-column>
				<el-table-column align="center" prop="RoleSellerName" label="角色运营商" width="130"></el-table-column>
				<el-table-column align="center" prop="RoleName" label="角色" width="100"></el-table-column>
				<el-table-column align="center" label="状态" width="100">
					<template slot-scope="scope">
						<el-link :underline="false" type="primary" v-if="scope.row.State == 1">启用</el-link>
						<el-link :underline="false" type="danger" v-if="scope.row.State == 2">禁用</el-link>
					</template>
				</el-table-column>
				<el-table-column align="center" prop="LoginTime" label="登录时间" width="160"></el-table-column>
				<el-table-column align="center" prop="LoginIp" label="ip" width="120"></el-table-column>
				<el-table-column align="center" prop="LoginCount" label="登录次数" width="100"></el-table-column>
				<el-table-column align="center" prop="Remark" label="备注" width="300"></el-table-column>
				<el-table-column label="操作">
					<template slot-scope="scope">
						<el-button type="text" icon="el-icon-edit" v-show="auth('改')" @click="handleModify(scope.$index)">编辑</el-button>
						<el-button type="text" icon="el-icon-delete" class="red" v-show="auth('删')" @click="handleDel(scope.$index)">删除</el-button>
					</template>
				</el-table-column>
			</el-table>
			<div class="pagination">
				<el-pagination style="margin-top: 5px" background layout="total, prev, pager, next, jumper" :hide-on-single-page="true" :total="total" @current-change="handleQuery" :page-size="pagesize"></el-pagination>
			</div>
		</div>
		<!--对话框-->
		<div>
			<el-dialog :title="dialog.title" :visible.sync="dialog.show" width="415px" center>
				<el-form :inline="true" label-width="100px">
					<el-form-item label="账号:">
						<el-input v-model="dialog.data.Account" :disabled="dialog.type == 'modify'"></el-input>
					</el-form-item>
					<el-form-item label="运营商:">
						<el-input v-model="dialog.data.SellerName" style="width: 200px" :disabled="dialog.type == 'modify'"></el-input>
					</el-form-item>
					<el-form-item label="密码:">
						<el-input v-model="dialog.data.Password" show-password style="width: 200px"></el-input>
					</el-form-item>
					<el-form-item label="角色运营商:" v-show="zong">
						<el-select v-model="dialog.data.RoleSellerId" placeholder="运营商" style="width: 130px" @change="handleSelectRoleSeller">
							<el-option v-for="item in dialog.options.RoleSellers" :key="item.SellerId" :label="item.SellerName" :value="item.SellerId"> </el-option>
						</el-select>
					</el-form-item>
					<el-form-item label="角色:">
						<el-select v-model="dialog.data.RoleName" placeholder="请选择" style="width: 200px">
							<el-option v-for="item in dialog.options.RoleNames" :key="item.RoleName" :label="item.RoleName" :value="item.RoleName"></el-option>
						</el-select>
					</el-form-item>
				</el-form>
				<el-form :inline="true" label-width="100px">
					<el-form-item label="选项:">
						<el-checkbox border label="禁用" v-model="dialog.data.State" :true-label="2" :false-label="1"></el-checkbox>
					</el-form-item>
				</el-form>
				<el-form :inline="true" label-width="100px">
					<el-form-item label="备注:">
						<el-input type="textarea" v-model="dialog.data.Remark" :rows="4"></el-input>
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
import base from '@/api/base.js'
import '@/assets/css/k.css'
export default {
	extends: base,
	data() {
		return {
			filters: {
				Account: null,
			},
			dialog: {
				options: {
					RoleNames: [],
					RoleSellers: [],
				},
			},
		}
	},
	created() {
		this.handleQuery(1)
	},
	methods: {
		auth(o) {
			return app.auth2('系统管理', '账号管理', o)
		},
		handleSelectRoleSeller() {
			this.dialog.data.RoleName = null
			app.post('/admin/role/listall', { SellerId: this.dialog.data.RoleSellerId, IgnoreSeller: true }, (result) => {
				this.dialog.options.RoleNames = []
				for (let i = 0; i < result.data.length; i++) {
					this.dialog.options.RoleNames.push({ RoleName: result.data[i] })
				}
			})
		},
		handleQuery(page) {
			if (typeof page == 'object') page = 1
			var data = {
				page: page,
				pagesize: this.pagesize,
				Account: this.filters.Account || '',
				SellerId: parseInt(this.filters.SellerId || 0),
			}
			app.post('/user/list', data, (result) => {
				this.total = result.data.total
				this.table_data = result.data.data
				for (var i = 0; i < this.table_data.length; i++) {
					this.table_data[i].CreateTime = this.$moment(this.table_data[i].CreateTime).format('YYYY-MM-DD hh:mm:ss')
					this.table_data[i].LoginTime = this.$moment(this.table_data[i].LoginTime).format('YYYY-MM-DD hh:mm:ss')
					for (let j = 0; j < this.seller.length; j++) {
						if (this.seller[j].SellerId == this.table_data[i].SellerId) this.table_data[i].SellerName = this.seller[j].SellerName
						if (this.seller[j].SellerId == this.table_data[i].RoleSellerId) this.table_data[i].RoleSellerName = this.seller[j].SellerName
					}
				}
			})
		},
		handleAdd() {
			this.dialog_data = { IsDisabled: 0 }
			this.dialog_title = '添加账号'
			this.dialog = true
		},
		handleModify(index) {
			this.current_row = index
			this.dialog.title = '修改账号'
			this.dialog.data = app.clone(this.table_data[index])
			this.dialog.type = 'modify'
			this.dialog.options.RoleSellers = app.clone(this.seller)
			let v = this.dialog.data.RoleName
			this.handleSelectRoleSeller()
			this.dialog.data.RoleName = v
			for (let i = 0; i < this.dialog.options.RoleSellers.length; i++) {
				if (this.dialog.options.RoleSellers[i].SellerId == 0) {
					this.dialog.options.RoleSellers.splice(i, 1)
				}
			}
			this.dialog.show = true
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
			//if(this.dialog_data.)
			// if (this.dialog_title == '修改账号') {
			// 	app.post('/system/account/modify', this.dialog_data, () => {
			// 		this.table_data[this.current_row] = app.clone(this.dialog_data)
			// 		this.dialog = false
			// 		this.$message.success('操作成功')
			// 	})
			// }
			// if (this.dialog_title == '添加账号') {
			// 	app.post('/system/account/add', this.dialog_data, (data) => {
			// 		this.dialog_data.AdminId = data.AdminId
			// 		this.table_data.push(this.dialog_data)
			// 		this.dialog = false
			// 		this.$message.success('操作成功')
			// 	})
			// }
		},
	},
}
</script>
