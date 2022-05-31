<template>
	<div class="container">
		<div class="handle-box">
			<el-button type="primary" icon="el-icon-plus" class="mr10" v-show="auth('增')" @click="handleAdd">添加</el-button>
			<el-button type="primary" class="mr10" v-show="auth('增')" @click="handleQuery()">刷新</el-button>
		</div>
		<div>
			<el-table :data="table_data" border max-height="620px" class="table" :cell-style="{ padding: '0' }">
				<el-table-column align="center" prop="id" label="编号" width="100"></el-table-column>
				<el-table-column align="center" prop="MinScore" label="最小金额" width="100"></el-table-column>
				<el-table-column align="center" prop="MaxScore" label="最大金额" width="100"></el-table-column>
				<el-table-column align="center" prop="WinLostScorePercent" label="盈亏比" width="100"></el-table-column>
				<el-table-column align="center" prop="MinFinishScorePercent" label="最小杀放" width="100"></el-table-column>
				<el-table-column align="center" prop="MaxFinishScorePercent" label="最大杀放" width="100"></el-table-column>
				<el-table-column align="center" prop="IsDisabled" label="禁用" width="100">
					<template slot-scope="scope">
						<el-link :underline="false" type="danger" v-if="scope.row.IsDisabled == 1">是</el-link>
						<el-link :underline="false" type="primary" v-if="scope.row.IsDisabled == 0">否</el-link>
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
		<div class="pagination">
			<el-pagination style="margin-top: 5px" background layout="total, prev, pager, next, jumper" :hide-on-single-page="true" :total="total" @current-change="handleQuery" :page-size="pagesize"></el-pagination>
		</div>
		<div>
			<el-dialog :title="dialog_title" :visible.sync="dialog" width="400px" center>
				<el-form :inline="true" label-width="100px">
					<el-form-item label="最小金额:">
						<el-input v-model="dialog_data.MinScore"></el-input>
					</el-form-item>
					<el-form-item label="最大金额:">
						<el-input v-model="dialog_data.MaxScore"></el-input>
					</el-form-item>
					<el-form-item label="盈亏比:">
						<el-input v-model="dialog_data.WinLostScorePercent"></el-input>
					</el-form-item>
					<el-form-item label="最小杀放:">
						<el-input v-model="dialog_data.MinFinishScorePercent"></el-input>
					</el-form-item>
					<el-form-item label="最大杀放:">
						<el-input v-model="dialog_data.MaxFinishScorePercent"></el-input>
					</el-form-item>
					<el-form-item label="选项:">
						<el-checkbox border label="禁用" v-model="dialog_data.IsDisabled" :true-label="1" :false-label="0"></el-checkbox>
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
			total: 0,
			pagesize: 15,
			dialog: false,
			dialog_title: '',
			step_dialog: false,
			dialog_data: { IsDisabled: 0 },
			current_row: 0,
			current_page: null,
		}
	},
	components: {},
	computed: {},
	created() {
		this.handleQuery()
	},
	methods: {
		auth(o) {
			return app.auth2('个人操盘', '个控配置', o)
		},
		handleAdd() {
			this.dialog_title = '添加体验'
			this.dialog_data = { IsRecharged: 0, IsDisabled: 0 }
			this.dialog = true
		},
		handleQuery(page) {
			if (page && page != this.current_page) this.current_page = page
			if (!this.current_page) this.current_page = 1
			var reqdata = {
				page: this.current_page,
				pagesize: this.pagesize,
			}
			app.post('/control/list/query', reqdata, (result) => {
				this.table_data = result.data
				this.total = result.total
				console.log(result)
			})
		},
		handleModify(index) {
			this.dialog_title = '修改体验'
			this.dialog_data = app.clone(this.table_data[index])
			this.current_row = index
			this.dialog = true
		},
		handleDel(index) {
			if (confirm('确定删除该体验?')) {
				app.post('/control/list/delete', { id: this.table_data[index].id }, (result) => {
					this.table_data.splice(index, 1)
					this.$message.success('操作成功')
				})
			}
		},
		handleConfirm() {
			this.dialog = false
			if (this.dialog_title == '添加体验') {
				app.post('/control/list/add', this.dialog_data, (result) => {
					var data = app.clone(this.dialog_data)
					data.id = -1
					this.table_data.push(data)
					this.$message.success('操作成功')
				})
			}
			if (this.dialog_title == '修改体验') {
				app.post('/control/list/modify', this.dialog_data, (result) => {
					this.table_data[this.current_row] = app.clone(this.dialog_data)
					this.table_data = app.clone(this.table_data)
					this.$message.success('操作成功')
				})
			}
		},
		handleLook(index) {
			this.step_dialog_title = this.table_data[index].id.toString()
			this.step_data = []
			var count = 0
			for (var i = 1; i <= 20; i++) {
				var data = this.table_data[index]['Step' + i]
				if (data) {
					data = JSON.parse(data)
					data.id = count
					this.step_data.push(data)
					count++
				}
			}
			this.step_dialog = true
		},
		handleSubmitStep() {
			for (var i = 0; i < this.step_data.length; i++) {
				var str = JSON.stringify({ minscore_rate: Number(this.step_data[i].minscore_rate), maxscore_rate: Number(this.step_data[i].maxscore_rate) })
				this.table_data[this.current_row]['Step' + (i + 1)] = str
			}
			app.post('/control/list/modify', this.table_data[this.current_row], (result) => {
				this.step_dialog = false
				this.$message.success('操作成功')
			})
		},
		handleModifyStep(index) {
			this.step_modify_dialog = true
			this.current_step_data = this.step_data[index]
			this.current_step_row = index
			this.step_modify_dialog_title = '修改步骤:' + this.current_step_data.id.toString()
		},
		handleDelStep(index) {
			if (confirm('确定删除?')) {
				this.step_data.splice(index, 1)
				for (var i = 0; i < this.step_data.length; i++) {
					this.step_data[i].id = i
				}
			}
		},
		handleAddStep() {
			this.step_modify_dialog_title = '添加步骤'
			this.current_step_data = []
			this.step_modify_dialog = true
		},
		handleStepOpt() {
			this.step_modify_dialog = false
			if (this.step_modify_dialog_title == '添加步骤') {
				this.step_data.push(this.current_step_data)
				for (var i = 0; i < this.step_data.length; i++) {
					this.step_data[i].id = i
				}
			}
		},
	},
}
</script>
