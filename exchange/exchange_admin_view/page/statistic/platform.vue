<template>
	<div class="container">
		<div>
			<el-form :inline="true" :model="filters">
				<el-form-item label="运营商:">
					<el-select v-model="filters.SellerId" placeholder="请选择" style="width: 150px; margin-right: 10px">
						<el-option v-for="item in sellers" :key="item.SellerName" :label="item.SellerName" :value="item.SellerId"> </el-option>
					</el-select>
				</el-form-item>
				<el-form-item label="查询日期:">
					<el-date-picker v-model="filters.QueryDate" align="right" type="date" placeholder="选择日期" style="width: 150px"> </el-date-picker>
				</el-form-item>
				<el-form-item>
					<el-button type="primary" class="mr10" @click="handleQuery()">查询</el-button>
				</el-form-item>
			</el-form>
		</div>
		<div>
			<el-table :data="table_data" border class="table" max-height="700px" :cell-style="{ padding: '0' }">
				<el-table-column align="center" prop="id" label="编号" width="50"></el-table-column>
				<el-table-column align="center" prop="SellerName" label="运营商" width="100"></el-table-column>
				<el-table-column align="center" prop="NewUser" label="新增人数" width="50"></el-table-column>
				<el-table-column align="center" prop="NewBindUser" label="新增正式人数" width="50"></el-table-column>
				<el-table-column align="center" prop="NewRechargeUser" label="新增充值人数" width="50"></el-table-column>
				<el-table-column align="center" prop="PresentScore" label="赠送彩金" width="80"></el-table-column>
				<el-table-column align="center" prop="DrawCommission" label="领取佣金" width="60"></el-table-column>
				<el-table-column align="center" prop="CalcCommission" label="结算佣金" width="60"></el-table-column>
				<el-table-column align="center" prop="BetScore" label="总投注" width="100"></el-table-column>
				<el-table-column align="center" prop="WinLostScore" label="总输赢" width="100"></el-table-column>
				<el-table-column align="center" prop="TaxScore" label="总税收" width="100"></el-table-column>
				<el-table-column align="center" prop="FlowScore" label="总流水" width="100"></el-table-column>
				<el-table-column align="center" prop="RechareScore" label="总充值" width="100"></el-table-column>
				<el-table-column align="center" prop="CashScore" label="总提款" width="100"></el-table-column>
				<el-table-column align="center" prop="ProfitScore" label="充提差" width="100"></el-table-column>
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
				SellerId: null,
				QueryDate: null,
			},
			sellers: [],
			table_data: [],
			dialog_data: [],
			dialog_title: null,
			dialog: false,
		}
	},
	components: {},
	computed: {},
	created() {
		app.post('/seller/list/query', {}, (result) => {
			this.sellers = [{ SellerId: null, SellerName: '全部' }]
			for (var i = 0; i < result.length; i++) {
				this.sellers.push({ SellerId: result[i].SellerId, SellerName: result[i].SellerName })
			}
			this.handleQuery()
		})
	},
	methods: {
		handleQuery() {
			app.post('/statistic/platform/query', this.filters, (result) => {
				for (var i = 0; i < result.length; i++) {
					result[i].id = i + 1
					for (var j = 0; j < this.sellers.length; j++) {
						result[i].WinLostScore = -result[i].WinLostScore
						if (result[i].SellerId == this.sellers[j].SellerId) {
							result[i].SellerName = this.sellers[j].SellerName
						}
						if (result[i].SellerId == 0) {
							result[i].SellerName = '全部运营商'
						}
					}
				}
				this.table_data = result
			})
		},
	},
}
</script>
