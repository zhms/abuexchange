<template>
	<div class="container">
		<!-- 筛选 -->
		<div>
			<el-form :inline="true" :model="filters">
				<el-form-item label="牌局号:">
					<el-input v-model="filters.GameSerial" style="width: 120px" :clearable="true" placeholder=""></el-input>
				</el-form-item>
				<el-form-item label="游戏:">
					<el-select v-model="filters.GameId" placeholder="请选择" style="width: 130px" @change="handleSelectGame()">
						<el-option v-for="item in games" :key="item.GameName" :label="item.GameName" :value="item.GameId"> </el-option>
					</el-select>
				</el-form-item>
				<el-form-item label="房间:" v-show="this.filters.GameId != null && this.filters.GameId != undefined">
					<el-select v-model="filters.RoomLevel" placeholder="请选择" style="width: 150px">
						<el-option v-for="item in option_room" :key="item.RoomName" :label="item.RoomName" :value="item.RoomLevel"> </el-option>
					</el-select>
				</el-form-item>
				<el-form-item label="时间:">
					<el-date-picker v-model="query_time" type="datetimerange" range-separator="至" start-placeholder="开始日期" end-placeholder="结束日期"> </el-date-picker>
				</el-form-item>
				<el-form-item>
					<el-checkbox v-model="filters.onlywin">只看赢</el-checkbox>
				</el-form-item>
				<el-form-item>
					<el-checkbox v-model="filters.onlyfaild">只看输</el-checkbox>
				</el-form-item>
				<el-form-item>
					<el-button type="primary" v-on:click="handleQuery">查询</el-button>
				</el-form-item>
			</el-form>
		</div>
		<!-- 表 -->
		<div>
			<el-table :data="table_data" border max-height="620px" class="table" :cell-style="{ padding: '0' }">
				<el-table-column align="center" prop="GameSerial" label="牌局号" width="160"></el-table-column>
				<el-table-column align="center" prop="RoomName" label="房间名称" width="140"></el-table-column>
				<el-table-column align="center" prop="WinLostScore" label="系统输赢" width="100">
					<template slot-scope="scope">
						<span v-if="table_data[scope.$index].WinLostScore > 0" style="color: rgb(255, 0, 0)">{{ table_data[scope.$index].WinLostScore }}</span>
						<span v-else style="color: rgb(52, 180, 83)">{{ table_data[scope.$index].WinLostScore }}</span>
					</template>
				</el-table-column>
				<el-table-column align="center" prop="TotalTaxScore" label="系统税收" width="100"></el-table-column>
				<el-table-column align="center" prop="RecordTime" label="记录时间" width="150"></el-table-column>
				<el-table-column label="详细记录">
					<template slot-scope="scope">
						<el-button type="text" icon="el-icon-zoom-in" @click="handleDetail(scope.$index)">详情</el-button>
					</template>
				</el-table-column>
			</el-table>
			<div class="pagination">
				<el-pagination style="margin-top: 5px" background layout="total, prev, pager, next, jumper" :hide-on-single-page="true" :total="total" @current-change="handleQuery" :page-size="pagesize"></el-pagination>
			</div>
		</div>
		<!--对话框-->
		<div>
			<el-dialog :title="detail_title" :visible.sync="dialog" width="700px" center> <el-input type="textarea" :rows="30" v-model="detail" :disabled="false"> </el-input></el-dialog>
		</div>
	</div>
</template>

<script>
import { app } from '@/api/app.js'
import moment from 'moment'
import '@/assets/css/k.css'
export default {
	data() {
		return {
			filters: {
				GameSerial: null,
				GameId: null,
				RoomLevel: null,
				onlywin: null,
				onlyfaild: null,
			},
			pagesize: 17,
			total: 0,
			games: [],
			rooms: [],
			table_data: [],
			dialog_data: {},
			dialog: false,
			detail: null,
			option_room: [],
			gameid: null,
			query_time: null,
			detail_title: '',
		}
	},
	components: {},
	computed: {},
	created() {
		app.post('/game/game/query', {}, (games) => {
			this.games = [{ GameName: '全部', GameId: null }]
			for (var i = 0; i < games.length; i++) {
				this.games.push({ GameName: games[i].GameName, GameId: games[i].GameId })
			}
			app.post('/game/room/query', {}, (rooms) => {
				for (var i = 0; i < rooms.rdata.length; i++) {
					this.rooms.push({ RoomName: rooms.rdata[i].RoomName, RoomLevel: rooms.rdata[i].RoomLevel, GameId: rooms.rdata[i].GameId })
				}
			})
		})
	},
	methods: {
		auth(o) {
			return app.auth2('游戏系统', '牌局记录', o)
		},
		handleSelectGame() {
			this.filters.RoomLevel = null
			if (this.filters.GameId) {
				var opt = [{ RoomName: '全部', RoomLevel: null }]
				for (var i = 0; i < this.rooms.length; i++) {
					if (this.rooms[i].GameId == this.filters.GameId) {
						opt.push({ RoomName: this.rooms[i].RoomName, RoomLevel: this.rooms[i].RoomLevel })
					}
				}
				this.option_room = opt
			}
		},
		handleDetail(index) {
			this.detail = JSON.stringify(JSON.parse(this.table_data[index].RecordData), null, '\t')
			this.detail_title = '详细:' + this.table_data[index].GameSerial
			this.dialog = true
		},
		handleQuery(page) {
			if (typeof page == 'object') page = 1
			var starttime = null
			var endtime = null
			if (this.query_time) {
				starttime = moment(this.query_time[0].getTime()).format('YYYY-MM-DD HH:mm:ss')
				endtime = moment(this.query_time[1].getTime()).format('YYYY-MM-DD HH:mm:ss')
			}
			var data = {
				page: page,
				pagesize: this.pagesize,
				GameSerial: this.filters.GameSerial,
				GameId: this.filters.GameId,
				RoomLevel: this.filters.RoomLevel,
				starttime: starttime,
				endtime: endtime,
				onlywin: this.filters.onlywin,
				onlyfaild: this.filters.onlyfaild,
			}
			app.post('/game/detail/query', data, (result) => {
				this.table_data = result.data
				this.total = result.total
				for (var i = 0; i < this.table_data.length; i++) {
					for (var j = 0; j < this.rooms.length; j++) {
						if (this.table_data[i].RoomLevel == this.rooms[j].RoomLevel && this.table_data[i].GameId == this.rooms[j].GameId) {
							this.table_data[i].RoomName = this.rooms[j].RoomName
						}
					}
					this.table_data[i].RecordTime = moment(this.table_data[i].RecordTime).format('YYYY-MM-DD HH:mm:ss')
				}
			})
		},
	},
}
</script>
