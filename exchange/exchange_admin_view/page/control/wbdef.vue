<template>
	<div class="container">
		<div>
			<el-table :data="table_data" border max-height="620px" class="table" :cell-style="{ padding: '0' }">
				<el-table-column align="center" prop="id" label="编号" width="100"></el-table-column>
				<el-table-column align="center" prop="DefName" label="名称" width="100"></el-table-column>
				<el-table-column align="center" prop="DefSlot" label="拉霸类" width="100"></el-table-column>
				<el-table-column align="center" prop="DefDuiZhan" label="对战类" width="100"></el-table-column>
				<el-table-column align="center" prop="DefFish" label="捕鱼类" width="100"></el-table-column>
				<el-table-column label="操作">
					<template slot-scope="scope">
						<el-button type="text" icon="el-icon-edit" v-show="auth('改')" @click="handleModify(scope.$index)">编辑</el-button>
					</template>
				</el-table-column>
			</el-table>
		</div>
		<div>
			<el-dialog :title="dialog_title" :visible.sync="dialog" width="400px" center>
				<el-form :inline="true" label-width="100px">
					<el-form-item label="拉霸类:">
						<el-input v-model="dialog_data.DefSlot"></el-input>
					</el-form-item>
					<el-form-item label="对战类:">
						<el-input v-model="dialog_data.DefDuiZhan"></el-input>
					</el-form-item>
					<el-form-item label="捕鱼类:">
						<el-input v-model="dialog_data.DefFish"></el-input>
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
			dialog_data: {},
			dialog_title: '',
			dialog: false,
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
			return app.getInstance().auth('个人操盘', '黑白定义', o)
		},
		handleQuery() {
			app.getInstance().post('/control/def/query', {}, (result) => {
				this.table_data = result
			})
		},
		handleModify(index) {
			this.current_row = index
			this.dialog_data = app.getInstance().clone(this.table_data[index])
			this.dialog_title = '修改-' + this.table_data[index].DefName
			this.dialog = true
		},
		handleConfirm() {
			app.getInstance().post('/control/def/modify', this.dialog_data, (result) => {
				this.table_data[this.current_row] = this.dialog_data
				this.table_data = app.getInstance().clone(this.table_data)
				this.dialog = false
			})
		},
	},
}
</script>
