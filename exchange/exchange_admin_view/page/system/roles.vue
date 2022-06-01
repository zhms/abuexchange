<template>
	<div class="container">
		<!-- 筛选 -->
		<div>
			<el-form :inline="true" :model="filters">
				<el-form-item label="运营商:" v-show="zong">
					<el-select v-model="filters.SellerId" placeholder="运营商" style="width: 130px">
						<el-option v-for="item in seller" :key="item.SellerId" :label="item.SellerName" :value="item.SellerId"> </el-option>
					</el-select>
				</el-form-item>
				<el-form-item>
					<el-button type="primary" icon="el-icon-plus" class="mr10" v-show="auth2('增')" @click="handleAdd">添加</el-button>
					<el-button type="primary" icon="el-icon-refresh" v-on:click="handleQuery">查询</el-button>
				</el-form-item>
			</el-form>
		</div>
		<!--表-->
		<div>
			<el-table :data="table_data" border max-height="620px" class="table" :cell-style="{ padding: '0' }">
				<el-table-column align="center" prop="SellerName" label="运营商" width="200"></el-table-column>
				<el-table-column align="center" prop="RoleName" label="角色名" width="200"></el-table-column>
				<el-table-column align="center" prop="ParentSellerName" label="上级角色运营商" width="200"></el-table-column>
				<el-table-column align="center" prop="Parent" label="上级角色" width="200"></el-table-column>
				<el-table-column label="操作">
					<template slot-scope="scope">
						<el-button type="text" icon="el-icon-edit" v-show="auth2('改')" @click="handleModify(scope.$index)">编辑</el-button>
						<el-button type="text" icon="el-icon-delete" class="red" v-show="auth2('删')" @click="handleDel(scope.$index)">删除</el-button>
					</template>
				</el-table-column>
			</el-table>
		</div>
		<!-- 对话框 -->
		<div>
			<el-dialog :title="dialog_title" :visible.sync="dialog" width="500px" center>
				<div>
					<el-form :inline="true" :model="filters">
						<el-form-item label="角色名:" v-show="zong" style="margin-left: 25px">
							<el-input v-model="dialog_select.RoleName" :disabled="dialog_type == 'alter'"></el-input>
						</el-form-item>
						<el-form-item label="上级运营商:" v-show="zong">
							<el-select v-model="dialog_select.ParentSellerId" placeholder="请选择" style="width: 130px" :disabled="dialog_type == 'alter'" @change="dialogSellerChange">
								<el-option v-for="item in dialog_select.seller" :key="item.SellerId" :label="item.SellerName" :value="item.SellerId"> </el-option>
							</el-select>
						</el-form-item>
						<el-form-item label="上级角色:" v-show="zong">
							<el-select v-model="dialog_select.Parent" placeholder="请选择" style="width: 130px" :disabled="dialog_type == 'alter'" @change="dialogRoleChange">
								<el-option v-for="item in dialog_select.parents" :key="item.RoleName" :label="item.RoleName" :value="item.RoleName"> </el-option>
							</el-select>
						</el-form-item>
						<el-form-item label="运营商:" v-show="zong" style="margin-left: 30px">
							<el-select v-model="dialog_select.SellerId" placeholder="请选择" style="width: 130px" :disabled="dialog_type == 'alter'">
								<el-option v-for="item in dialog_select.seller" :key="item.SellerId" :label="item.SellerName" :value="item.SellerId"> </el-option>
							</el-select>
						</el-form-item>
					</el-form>
				</div>
				<el-tree :default-checked-keys="dialog_select.ids" node-key="path" ref="authtree" :props="props" show-checkbox v-show="dialog_tree"> </el-tree>
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
			filters: {
				SellerId: app.currentSeller(),
			},
			dialog_select: {
				seller: [],
				SellerId: null,
				parents: [],
				ParentSellerId: null,
				Parent: null,
				ids: [],
				RoleName: null,
			},
			props: {
				label: 'name',
				children: 'children',
			},
			zong: app.zong(),
			seller: app.getSeller(),
			table_data: null,
			pagesize: 15,
			total: 0,
			dialog: false,
			dialog_type: null,
			dialog_data: {},
			dialog_title: null,
			dialog_tree: false,
		}
	},
	components: {},
	computed: {},
	created() {
		this.handleQuery()
	},
	methods: {
		auth2(o) {
			return app.auth2('系统管理', '角色管理', o)
		},
		dialogSellerChange() {
			this.dialog_tree = false
			this.dialog_select.Parent = null
			app.post('/admin/role/listall', { SellerId: this.dialog_select.ParentSellerId, IgnoreSeller: true }, (result) => {
				this.dialog_select.parents = []
				for (let i = 0; i < result.data.length; i++) {
					this.dialog_select.parents.push({ RoleName: result.data[i] })
				}
			})
		},
		dialogRoleChange() {
			if (this.dialog_select.Parent) {
				app.post('/admin/role/roledata', { SellerId: this.dialog_select.ParentSellerId, IgnoreSeller: true, RoleName: this.dialog_select.Parent }, (result) => {
					this.dialog_select.parentroledata = JSON.parse(result.data.RoleData)
					this.dialog_select.superroledata = JSON.parse(result.data.SuperRoleData)
					this.dialog_tree = true
					this.dialog_select.roledata = {}
					this.dialog_tree = true
					let treedata = this.getTreeData()
					this.$refs.authtree.root.setData(treedata.menu)
					this.dialog_select.ids = treedata.ids
				})
			}
		},
		handleQuery(page) {
			if (typeof page == 'object') page = 1
			var data = {
				SellerId: parseInt(this.filters.SellerId || 0),
				page: page,
				pagesize: this.pagesize,
			}
			app.post('/admin/role/list', data, (result) => {
				this.table_data = result.data.data
				this.total = result.data.total
				for (var i = 0; i < this.table_data.length; i++) {
					for (let j = 0; j < this.seller.length; j++) {
						if (this.seller[j].SellerId == this.table_data[i].SellerId) this.table_data[i].SellerName = this.seller[j].SellerName
						if (this.seller[j].SellerId == this.table_data[i].ParentSellerId) this.table_data[i].ParentSellerName = this.seller[j].SellerName
					}
				}
			})
		},
		handleAdd() {
			this.dialog_title = `添加角色`
			this.dialog_type = 'add'
			this.dialog_data = {}
			this.dialog_select.SellerId = null
			this.dialog_select.ParentSellerId = null
			this.dialog_select.Parent = null
			this.dialog_select.RoleName = null
			this.dialog_select.seller = app.clone(this.seller)
			for (let i = 0; i < this.dialog_select.seller.length; i++) {
				if (this.dialog_select.seller[i].SellerId == 0) {
					this.dialog_select.seller.splice(i, 1)
				}
			}
			this.dialog_tree = false
			this.dialog = true
		},
		handleModify(index) {
			this.current_row = index
			if (this.table_data[this.current_row].Parent == 'god') {
				this.$message.error('该角色不可修改')
				return
			}
			this.dialog_data = app.clone(this.table_data[this.current_row])
			this.dialog_select.SellerId = this.dialog_data.SellerId
			this.dialog_select.ParentSellerId = this.dialog_data.ParentSellerId
			this.dialog_select.Parent = this.dialog_data.Parent
			this.dialog_select.RoleName = this.dialog_data.RoleName
			this.dialog_select.seller = app.clone(this.seller)
			this.dialog_title = `修改角色`
			this.dialog_type = 'alter'
			this.dialog = true
			setTimeout(() => {
				this.$refs.authtree.root.setData([])
			}, 10)
			app.post('/admin/role/listall', { SellerId: this.dialog_select.SellerId, IgnoreSeller: true }, (result) => {
				this.dialog_select.parents = []
				for (let i = 0; i < result.data.length; i++) {
					this.dialog_select.parents.push({ RoleName: result.data[i] })
				}
				for (let i = 0; i < this.dialog_select.seller.length; i++) {
					if (this.dialog_select.seller[i].SellerId == 0) {
						this.dialog_select.seller.splice(i, 1)
					}
				}
				app.post('/admin/role/roledata', { SellerId: this.dialog_data.SellerId, IgnoreSeller: true, RoleName: this.dialog_data.Parent }, (resulta) => {
					this.dialog_select.parentroledata = JSON.parse(resulta.data.RoleData)
					this.dialog_select.superroledata = JSON.parse(resulta.data.SuperRoleData)
					app.post('/admin/role/roledata', { SellerId: this.dialog_data.SellerId, IgnoreSeller: true, RoleName: this.dialog_data.RoleName }, (resultb) => {
						this.dialog_select.roledata = JSON.parse(resultb.data.RoleData)
						this.dialog_tree = true
						let treedata = this.getTreeData()
						this.$refs.authtree.root.setData(treedata.menu)
						this.dialog_select.ids = treedata.ids
					})
				})
			})
		},
		handleDel(index) {
			this.current_row = index
			this.dialog_data = this.table_data[this.current_row]
			if (this.dialog_data.Parent == 'god') {
				this.$message.error('该角色不可删除')
				return
			}
			if (confirm('确定删除该角色?')) {
				let data = {
					SellerId: this.dialog_data.SellerId,
					RoleName: this.dialog_data.RoleName,
					IgnoreSeller: true,
				}
				app.post('/admin/role/delete', data, (result) => {
					this.dialog = false
					this.$message.success('操作成功')
					this.handleQuery()
				})
			}
		},
		handleConfirm() {
			if (!this.dialog_select.RoleName) return this.$message.error('请填写角色名')
			if (!this.dialog_select.ParentSellerId) return this.$message.error('请选择上角色运营商')
			if (!this.dialog_select.Parent) return this.$message.error('请选择上级角色')
			if (!this.dialog_select.SellerId) return this.$message.error('请选择运营商')
			let setdisable = (node) => {
				for (let n in node) {
					if (typeof node[n] == 'object') {
						setdisable(node[n])
					} else {
						node[n] = 0
					}
				}
			}
			let newroledata = app.clone(this.dialog_select.superroledata)
			setdisable(newroledata)
			let selected = this.$refs.authtree.getCheckedNodes()
			for (let i = 0; i < selected.length; i++) {
				if (!selected[i].leaf) continue
				let path = selected[i].path.split('.')
				let pn = newroledata
				for (let i = 0; i < path.length - 1; i++) {
					pn = pn[path[i]]
				}
				pn[path[path.length - 1]] = 1
			}
			if (this.dialog_type == 'alter') {
				let data = {
					SellerId: this.dialog_select.SellerId,
					RoleName: this.dialog_select.RoleName,
					RoleData: JSON.stringify(newroledata),
				}
				app.post('/admin/role/update', data, (result) => {
					this.dialog = false
					this.$message.success('操作成功')
				})
			}
			if (this.dialog_type == 'add') {
				let data = {
					ParentSellerId: this.dialog_select.ParentSellerId,
					Parent: this.dialog_select.Parent,
					SellerId: this.dialog_select.SellerId,
					RoleName: this.dialog_select.RoleName,
					RoleData: JSON.stringify(newroledata),
					IgnoreSeller: true,
				}
				app.post('/admin/role/add', data, (result) => {
					this.dialog = false
					this.$message.success('操作成功')
					this.handleQuery()
				})
			}
		},
		getTreeData() {
			let setdisable = (node) => {
				for (let n in node) {
					if (typeof node[n] == 'object') {
						setdisable(node[n])
					} else {
						node[n] = 0
					}
				}
			}
			setdisable(this.dialog_select.superroledata)
			let setenable = (parent, node) => {
				for (let n in node) {
					if (typeof node[n] == 'object') {
						let p = parent + `.${n}`
						setenable(p, node[n])
					} else {
						if (node[n] == 1) {
							let p = parent.split('.')
							let pn = this.dialog_select.superroledata
							for (let j = 0; j < p.length; j++) {
								pn = pn[p[j]]
							}
							pn[n] = 1
						}
					}
				}
			}
			for (let n in this.dialog_select.parentroledata) {
				let parent = `${n}`
				setenable(parent, this.dialog_select.parentroledata[n])
			}
			let menu = []
			let submenu = (node, root) => {
				for (let n in root) {
					if (typeof root[n] == 'object') {
						let subnode = {
							path: node.path + '.' + n,
							name: n,
							children: [],
						}
						node.children.push(subnode)
						submenu(subnode, root[n])
					} else {
						let path = node.path + '.' + n
						let p = path.split('.')
						let pr = this.dialog_select.parentroledata
						for (let i = 0; i < p.length; i++) {
							pr = pr[p[i]]
						}
						if (pr == 1) {
							let subnode = {
								path: path,
								name: n,
								leaf: true,
							}
							node.children.push(subnode)
						}
					}
				}
			}
			for (let n in this.dialog_select.superroledata) {
				let node = {
					path: n,
					name: n,
					children: [],
				}
				menu.push(node)
				submenu(node, this.dialog_select.superroledata[n])
			}
			let ids = []
			let getselected = (parent, node) => {
				for (let n in node) {
					if (typeof node[n] == 'object') {
						let p = parent + `.${n}`
						getselected(p, node[n])
					} else {
						if (node[n] == 1) {
							ids.push(`${parent}.${n}`)
						}
					}
				}
			}
			for (let n in this.dialog_select.roledata) {
				let parent = `${n}`
				getselected(parent, this.dialog_select.roledata[n])
			}
			for (let i = 0; i < menu.length; i++) {
				if (!menu[i].children) continue
				for (let j = 0; j < menu[i].children.length; j++) {
					if (!menu[i].children[j].children) continue
					for (let k = 0; k < menu[i].children[j].children.length; k++) {
						if (!menu[i].children[j].children[k].children) continue
						if (menu[i].children[j].children[k].children.length == 0) {
							menu[i].children[j].children.splice(k, 1)
							k--
						}
					}
				}
			}
			for (let i = 0; i < menu.length; i++) {
				if (!menu[i].children) continue
				for (let j = 0; j < menu[i].children.length; j++) {
					if (!menu[i].children[j].children) continue
					if (menu[i].children[j].children.length == 0) {
						menu[i].children.splice(j, 1)
						j--
					}
				}
			}
			for (let i = 0; i < menu.length; i++) {
				if (!menu[i].children) continue
				if (menu[i].children.length == 0) {
					menu.splice(i, 1)
					i--
				}
			}
			return { menu, ids }
		},
	},
}
</script>
