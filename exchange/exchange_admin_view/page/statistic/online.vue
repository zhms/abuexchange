<template>
	<div class="container">
		<div class="handle-box">
			<el-form :inline="true" :model="filters">
				<el-form-item label="运营商:">
					<el-select v-model="filters.SellerId" placeholder="请选择" style="width:150px;margin-right: 10px;">
						<el-option v-for="item in sellers" :key="item.SellerName" :label="item.SellerName" :value="item.SellerId"> </el-option>
					</el-select>
					<el-form-item>
						<el-button type="primary" class="mr10" @click="handleQuery()">刷新</el-button>
					</el-form-item>
				</el-form-item>
			</el-form>
		</div>
		<div>
			<el-table :data="table_data" border class="table" max-height="700px" :cell-style="{ padding: '0' }">
				<el-table-column align="center" prop="id" label="编号" width="100"></el-table-column>
				<el-table-column align="center" prop="GameId" label="游戏Id" width="230"></el-table-column>
				<el-table-column align="center" prop="GameName" label="游戏名称" width="150"></el-table-column>
				<el-table-column align="center" prop="OnlineCount" label="在线人数" width="150"></el-table-column>
				<el-table-column label="房间在线">
					<template slot-scope="scope">
						<el-button type="text" @click="handleLook(scope.$index)">查看</el-button>
					</template>
				</el-table-column>
			</el-table>
		</div>
		<div>
			<el-dialog :title="dialog_title" :visible.sync="dialog" width="530px" center>
				<el-table :data="dialog_data" border class="table" max-height="700px" :cell-style="{ padding: '0' }">
					<el-table-column align="center" prop="id" label="编号" width="100"></el-table-column>
					<el-table-column align="center" prop="RoomName" label="房间名称" width="230"></el-table-column>
					<el-table-column align="center" prop="OnlineCount" label="在线人数"></el-table-column>
				</el-table>
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
				SellerId: null,
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
		this.handleQuery()
		app.getInstance().post('/seller/list/query', {}, (result) => {
			this.sellers = [{ SellerId: null, SellerName: '全部' }]
			for (var i = 0; i < result.length; i++) {
				this.sellers.push({ SellerId: result[i].SellerId, SellerName: result[i].SellerName })
			}
		})
	},
	methods: {
		handleQuery() {
			app.getInstance().post('/statistic/gameonline/query', { SellerId: this.filters.SellerId }, (result) => {
				for (var i = 0; i < result.length; i++) {
					result[i].id = i + 1
				}
				this.table_data = result
			})
		},
		handleLook(index) {
			this.dialog_title = this.table_data[index].GameName + ' 房间在线'
			var data = { GameId: this.table_data[index].GameId, SellerId: this.filters.SellerId }
			app.getInstance().post('/statistic/roomonline/query', data, (result) => {
				for (var i = 0; i < result.length; i++) {
					result[i].id = i + 1
				}
				this.dialog_data = result
				this.dialog = true
			})
		},
		handleConfirm(index) {
			this.dialog = false
		},
	},
}
</script>
