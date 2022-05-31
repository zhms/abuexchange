<template>
	<div class="container">
		<!-- 筛选 -->
		<div>
			<el-form :inline="true" :model="filters">
				<el-form-item label="玩家Id:">
					<el-input v-model="filters.UserId" style="width: 150px" :clearable="true"></el-input>
				</el-form-item>
				<el-form-item>
					<el-button type="primary" v-on:click="handleQuery">查询</el-button>
				</el-form-item>
			</el-form>
		</div>
		<div>
			<el-table :data="table_data" border max-height="620px" class="table" :cell-style="{ padding: '0' }">
				<el-table-column align="center" prop="UserId" label="玩家Id" width="100"></el-table-column>
				<el-table-column align="center" prop="Status" label="状态" width="100">
					<template slot-scope="scope">
						<el-link :underline="false" type="danger" v-if="scope.row.Status == 0">未执行</el-link>
						<el-link :underline="false" type="success" v-if="scope.row.Status == 1">待执行</el-link>
						<el-link :underline="false" type="primary" v-if="scope.row.Status == 2">正在执行</el-link>
					</template>
				</el-table-column>
				<el-table-column align="center" prop="StartScore" label="开始金币" width="100"></el-table-column>
				<el-table-column align="center" prop="FinishScore" label="结束金币" width="100"></el-table-column>
				<el-table-column label="操作">
					<template slot-scope="scope">
						<el-button type="text" @click="handleReset(scope.$index)">重置</el-button>
					</template>
				</el-table-column>
			</el-table>
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
				UserId: null,
			},
			table_data: [],
		}
	},
	components: {},
	computed: {},
	created() {},
	methods: {
		auth(o) {
			return app.auth2('个人操盘', '查看个控', o)
		},
		handleQuery() {
			if (this.filters.UserId == 0) this.filters.UserId = null
			app.post('/control/personal/query', { UserId: this.filters.UserId }, (result) => {
				this.table_data = result
			})
		},
		handleReset(index) {
			if (confirm('确定重置?')) {
				app.post('/control/personal/modify', { UserId: this.table_data[index].UserId, opt: 'reset' }, (result) => {
					this.table_data[index].Status = 0
					this.table_data[index].ExpId = 0
					this.table_data[index].StartScore = 0
					this.table_data[index].DestScore = 0
					this.table_data[index].MaxScore = 0
					this.table_data[index].Step = 0
					this.table_data = app.clone(this.table_data)
				})
			}
		},
		handleEnd(index) {
			if (confirm('确定结束?')) {
				app.post('/control/personal/modify', { UserId: this.table_data[index].UserId, opt: 'end' }, (result) => {
					this.table_data[index].Status = 1
					this.table_data[index].ExpId = 0
					this.table_data[index].StartScore = 0
					this.table_data[index].DestScore = 0
					this.table_data[index].MaxScore = 0
					this.table_data[index].Step = 0
					this.table_data = app.clone(this.table_data)
				})
			}
		},
	},
}
</script>
