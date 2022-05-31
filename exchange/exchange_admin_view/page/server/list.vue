<template>
	<div class="container">
		<!-- 筛选 -->
		<div>
			<el-form :inline="true">
				<el-form-item>
					<el-button type="primary" v-on:click="handleAddServer">添加服务器</el-button>
					<el-button type="primary" v-on:click="handleGetCode">更新最新代码</el-button>
					<el-button type="danger" v-on:click="handleStopAll">停止所有服务器</el-button>
					<el-button type="primary" v-on:click="handleStatus">刷新状态</el-button>
				</el-form-item>
			</el-form>
		</div>
		<!-- 表 -->
		<div>
			<el-table :data="table_data" border class="table" max-height="700px" :cell-style="{ padding: '5' }">
				<el-table-column align="center" prop="ServerName" label="服务器名称" width="230"></el-table-column>
				<el-table-column align="center" prop="ServerType" label="服务器类型" width="150"></el-table-column>
				<el-table-column align="center" prop="ServerHost" label="服务地址" width="150"></el-table-column>
				<el-table-column align="center" prop="ServerPort" label="端口" width="120"></el-table-column>
				<el-table-column align="center" prop="RunServer" label="运行服务器" width="150"></el-table-column>
				<el-table-column align="center" prop="Note" label="备注" width="250"></el-table-column>
				<el-table-column align="center" prop="Sort" label="排序" width="50"></el-table-column>
				<el-table-column align="center" label="状态" width="80">
					<template slot-scope="scope">
						<el-tag type="success" effect="plain" v-show="scope.row.running">{{ '运行中' }}</el-tag>
						<el-tag type="danger" effect="plain" v-show="!scope.row.running">{{ '未运行' }}</el-tag>
					</template>
				</el-table-column>
				<el-table-column label="操作">
					<template slot-scope="scope">
						<el-button type="primary" size="mini" v-show="!scope.row.running" @click="handleStartServer(scope.$index)">启动</el-button>
						<el-button type="danger" size="mini" v-show="scope.row.running" @click="handleStopServer(scope.$index)">停止</el-button>
						<el-button type="primary" size="mini" v-show="!scope.row.running" @click="handleModifyServer(scope.$index)">修改</el-button>
						<el-button type="danger" size="mini" v-show="!scope.row.running" @click="handleDeleteServer(scope.$index)">删除</el-button>
					</template>
				</el-table-column>
			</el-table>
		</div>
		<!--对话框-->
		<div>
			<el-dialog :title="dialog_title" :visible.sync="dialog" width="415px" center>
				<el-form :inline="true" label-width="100px">
					<el-form-item label="服务名称:">
						<el-input v-model="dialog_data.ServerName" :disabled="dialog_title == '修改服务器'"></el-input>
					</el-form-item>
					<el-form-item label="服务类型:">
						<el-input v-model="dialog_data.ServerType" :disabled="dialog_title == '修改服务器'"></el-input>
					</el-form-item>
					<el-form-item label="服务地址:">
						<el-input v-model="dialog_data.ServerHost"></el-input>
					</el-form-item>
					<el-form-item label="服务端口:">
						<el-input v-model="dialog_data.ServerPort"></el-input>
					</el-form-item>
					<el-form-item label="运行服务器:">
						<el-input v-model="dialog_data.RunServer"></el-input>
					</el-form-item>
					<el-form-item label="排序:">
						<el-input v-model="dialog_data.Sort"></el-input>
					</el-form-item>
					<el-form :inline="true" label-width="100px">
						<el-form-item label="备注:">
							<el-input type="textarea" v-model="dialog_data.Note" :rows="4"></el-input>
						</el-form-item>
					</el-form>
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
			dialog_title: '',
			dialog_data: {},
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
			return app.auth2('服务管理', '服务列表', o)
		},
		handleQuery() {
			app.post('/server/list/query', {}, (serverdata) => {
				this.table_data = serverdata
				this.handleStatus()
			})
		},
		handleModifyServer(index) {
			this.dialog_title = '修改服务'
			this.dialog_data = app.clone(this.table_data[index])
			this.current_row = index
			this.dialog = true
		},
		handleAddServer(index) {
			this.dialog_title = '添加服务'
			this.dialog_data = {}
			this.current_row = index
			this.dialog = true
		},
		handleConfirm() {
			if (this.dialog_title == '添加服务') {
				app.post('/server/list/add', this.dialog_data, (data) => {
					this.table_data.push(this.dialog_data)
					this.dialog = false
					this.$message.success('操作成功')
				})
			}
			if (this.dialog_title == '修改服务') {
				app.post('/server/list/modify', this.dialog_data, (data) => {
					this.table_data[this.current_row] = app.clone(this.dialog_data)
					this.dialog = false
					this.$message.success('操作成功')
				})
			}
		},
		handleDeleteServer(index) {
			if (confirm('确定删除该服务?')) {
				app.post('/server/list/delete', this.table_data[index], (data) => {
					this.table_data.splice(index, 1)
					this.$message.success('操作成功')
				})
			}
		},
		handleRestartAll() {
			app.post('/server/list/startall', {}, (data) => {
				this.handleStatus()
			})
		},
		handleStopAll() {
			app.post('/server/list/stopall', {}, (data) => {
				this.handleStatus()
			})
		},
		handleStartServer(index) {
			app.post('/server/list/start', { ServerName: this.table_data[index].ServerName }, (data) => {
				this.handleStatus()
			})
		},
		handleStopServer(index) {
			app.post('/server/list/stop', { ServerName: this.table_data[index].ServerName }, (data) => {
				this.handleStatus()
			})
		},
		handleStatus() {
			app.post('/server/list/status', {}, (runningdata) => {
				for (var i = 0; i < this.table_data.length; i++) {
					this.table_data[i].running = runningdata[this.table_data[i].ServerName]
				}
				this.table_data = app.clone(this.table_data)
			})
		},
		handleGetCode() {
			app.post('/server/list/getcode', {}, () => {})
		},
	},
}
</script>
