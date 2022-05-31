<template>
	<div class="container">
		<!-- 筛选 -->
		<div>
			<el-form :inline="true" :model="filters">
				<el-form-item label="筛选:">
					<el-input v-model="filters.userid" style="width: 200px" :clearable="true" placeholder="玩家id"></el-input>
				</el-form-item>
				<el-form-item label="时间:">
					<el-date-picker v-model="query_time" type="datetimerange" range-separator="至" start-placeholder="开始日期" end-placeholder="结束日期"> </el-date-picker>
				</el-form-item>
				<el-form-item>
					<el-button type="primary" v-on:click="handleQuery(1)">查询</el-button>
				</el-form-item>
			</el-form>
		</div>
		<div>
			<el-table :data="score_data" border max-height="620px" class="table" :cell-style="{ padding: '0' }">
				<el-table-column align="center" prop="RecordTime" label="记录时间" width="160"></el-table-column>
				<el-table-column align="center" prop="ChangeScore" label="变化金币" width="120">
					<template slot-scope="scope">
						<span v-if="score_data[scope.$index].ChangeScore > 0" style="color: rgb(255, 0, 0)">{{ score_data[scope.$index].ChangeScore }}</span>
						<span v-else style="color: rgb(52, 180, 83)">{{ score_data[scope.$index].ChangeScore }}</span>
					</template>
				</el-table-column>
				<el-table-column align="center" prop="BeforeScore" label="变化前金币" width="120"></el-table-column>
				<el-table-column align="center" prop="AfterScore" label="变化后金币" width="120"></el-table-column>
				<el-table-column align="center" prop="ChangeBankScore" label="变化银行金币" width="120">
					<template slot-scope="scope">
						<span v-if="score_data[scope.$index].ChangeBankScore > 0" style="color: rgb(255, 0, 0)">{{ score_data[scope.$index].ChangeBankScore }}</span>
						<span v-else style="color: rgb(52, 180, 83)">{{ score_data[scope.$index].ChangeBankScore }}</span>
					</template>
				</el-table-column>
				<el-table-column align="center" prop="BeforeBankScore" label="变化前银行金币" width="120"></el-table-column>
				<el-table-column align="center" prop="AfterBankScore" label="变化后银行金币" width="120"></el-table-column>
				<el-table-column align="center" prop="ChangeReason" label="变化原因" width="130"></el-table-column>
				<el-table-column align="center" prop="Note" label="备注"></el-table-column>
			</el-table>
		</div>
		<div class="pagination">
			<el-pagination style="margin-top: 5px" background layout="total, prev, pager, next, jumper" :hide-on-single-page="true" :total="score_total" @current-change="handleQuery" :page-size="score_pagesize"></el-pagination>
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
				userid: null,
				starttime: null,
				endtime: null,
			},
			query_time: null,
			score_data: [],
			score_total: 0,
			score_pagesize: 20,
		}
	},
	components: {},
	computed: {},
	created() {},
	methods: {
		auth(o) {
			return app.auth2('玩家管理', '金流记录', o)
		},
		handleQuery(page) {
			var reqdata = {
				UserId: this.filters.userid,
				page: page,
				pagesize: this.score_pagesize,
			}
			if (this.query_time) {
				reqdata.starttime = this.$moment(this.query_time[0]).format('YYYY-MM-DD hh:mm:ss')
				reqdata.endtime = this.$moment(this.query_time[1]).format('YYYY-MM-DD hh:mm:ss')
			}
			console.log(reqdata)
			app.post('/user/scorelist/query', reqdata, (result) => {
				this.score_data = result.data
				this.score_total = result.total
				for (var i = 0; i < this.score_data.length; i++) {
					this.score_data[i].RecordTime = this.$moment(this.score_data[i].RecordTime).format('YYYY-MM-DD hh:mm:ss')
				}
			})
		},
	},
}
</script>
