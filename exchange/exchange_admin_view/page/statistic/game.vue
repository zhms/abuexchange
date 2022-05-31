<template>
	<div class="container">
		<div>
			<el-form :inline="true" :model="filters">
				<el-form-item label="运营商:">
					<el-select v-model="filters.SellerId" placeholder="请选择" style="width: 150px; margin-right: 10px">
						<el-option v-for="item in sellers" :key="item.SellerName" :label="item.SellerName" :value="item.SellerId"> </el-option>
					</el-select>
				</el-form-item>
				<el-form-item label="游戏:">
					<el-select v-model="filters.GameId" placeholder="请选择" style="width: 150px; margin-right: 10px">
						<el-option v-for="item in games" :key="item.GameName" :label="item.GameName" :value="item.GameId"> </el-option>
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
				<el-table-column align="center" prop="GameId" label="游戏" width="150"></el-table-column>
				<el-table-column align="center" prop="PlayCount" label="参与" width="100"></el-table-column>
				<el-table-column align="center" prop="WinPlayCount" label="赢钱" width="80"></el-table-column>
				<el-table-column align="center" prop="LostPlayCount" label="输钱" width="80"></el-table-column>
				<el-table-column align="center" prop="BetScore" label="总投注" width="100"></el-table-column>
				<el-table-column align="center" prop="WinLostScore" label="总输赢" width="100"></el-table-column>
				<el-table-column align="center" prop="TaxScore" label="总税收" width="100"></el-table-column>
				<el-table-column align="center" prop="FlowScore" label="总流水" width="100"></el-table-column>
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
				GameId: null,
				QueryDate: null,
			},
			sellers: [],
			games: [],
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
			app.post('/game/game/query', {}, (games) => {
				this.games = [{ GameName: '全部', GameId: null }]
				for (var i = 0; i < games.length; i++) {
					var v = { GameName: games[i].GameName, GameId: games[i].GameId }
					this.games.push(v)
				}
				this.handleQuery()
			})
		})
	},
	methods: {
		handleQuery() {
			app.post('/statistic/game/query', this.filters, (result) => {
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
