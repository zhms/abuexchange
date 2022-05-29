<template>
	<div class="container">
		<!-- 筛选 -->
		<div>
			<el-form :inline="true" :model="filters">
				<el-form-item label="管理员:">
					<el-input v-model="filters.NickName" style="width:150px" :clearable="true"></el-input>
				</el-form-item>
				<el-form-item label="操作:">
					<el-input v-model="filters.OptType" style="width:150px" :clearable="true"></el-input>
				</el-form-item>
				<el-form-item>
					<el-button type="primary" v-on:click="handleQuery">查询</el-button>
				</el-form-item>
			</el-form>
		</div>
		<!-- 表 -->
		<div>
			<el-table :data="table_data" border class="table" max-height="700px" :cell-style="{ padding: '0' }">
				<el-table-column align="center" prop="id" label="序号" width="80"></el-table-column>
				<el-table-column align="center" prop="NickName" label="管理员" width="100"></el-table-column>
				<el-table-column align="center" prop="RecordDate" label="时间" width="160"></el-table-column>
				<el-table-column align="center" prop="OptIp" label="ip" width="130"></el-table-column>
				<el-table-column align="center" prop="OptType" label="操作类型" width="200"></el-table-column>
				<el-table-column label="内容">
					<template slot-scope="scope">
						<el-button type="text" icon="el-icon-document-copy" @click="handleCopy(scope.$index)">复制</el-button>
					</template>
				</el-table-column>
			</el-table>
			<div class="pagination">
				<el-pagination style="margin-top:5px" background layout="total, prev, pager, next, jumper" :hide-on-single-page="true" :total="total" @current-change="handleQuery" :page-size="pagesize"></el-pagination>
			</div>
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
				NickName: null,
				OptType: null,
			},
			table_data: null,
			pagesize: 15,
			total: 0,
		}
	},
	components: {},
	computed: {},
	created() {
		this.handleQuery(1)
	},
	methods: {
		handleCopy(index) {
			var oInput = document.createElement('input')
			oInput.value = this.table_data[index].RequestData
			document.body.appendChild(oInput)
			oInput.select()
			document.execCommand('Copy')
			oInput.remove()
			this.$message.success('复制成功')
		},
		handleQuery(page) {
			if (typeof page == 'object') page = 1
			var data = {
				NickName: this.filters.NickName,
				OptType: this.filters.OptType,
				page: page,
				pagesize: this.pagesize,
			}
			app.getInstance().post('/system/log/query', data, (result) => {
				this.table_data = result.data
				this.total = result.total
				for (var i = 0; i < this.table_data.length; i++) {
					this.table_data[i].RecordDate = this.$moment(this.table_data[i].RecordDate).format('YYYY-MM-DD hh:mm:ss')
				}
			})
		},
	},
}
</script>
