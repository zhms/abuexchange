<template>
	<div class="container">
		<!-- 筛选 -->
		<div>
			<el-form :inline="true" :model="filters">
				<el-form-item label="玩家:">
					<el-input v-model="filters.UserId" style="width:120px" :clearable="true" placeholder=""></el-input>
				</el-form-item>
				<el-form-item label="游戏:">
					<el-select v-model="filters.GameId" placeholder="请选择" style="width:130px" @change="handleSelectGame()">
						<el-option v-for="item in games" :key="item.GameName" :label="item.GameName" :value="item.GameId"> </el-option>
					</el-select>
				</el-form-item>
				<el-form-item label="房间:" v-show="this.filters.GameId != null && this.filters.GameId != undefined">
					<el-select v-model="filters.RoomLevel" placeholder="请选择" style="width:150px">
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
				<el-table-column align="center" prop="id" label="编号" width="100"></el-table-column>
				<el-table-column align="center" prop="UserId" label="玩家" width="100"></el-table-column>
				<el-table-column align="center" prop="RoomName" label="房间名称" width="140"></el-table-column>
				<el-table-column align="center" prop="GameSerial" label="牌局号" width="160"></el-table-column>
				<el-table-column align="center" prop="BetScore" label="本局下注" width="160"></el-table-column>
				<el-table-column align="center" prop="WinLostScore" label="本局输赢" width="100">
					<template slot-scope="scope">
						<span v-if="table_data[scope.$index].WinLostScore > 0" style="color: rgb(255,0,0)">{{ table_data[scope.$index].WinLostScore }}</span>
						<span v-else style="color: rgb(52,180,83)">{{ table_data[scope.$index].WinLostScore }}</span>
					</template>
				</el-table-column>
				<el-table-column align="center" prop="TotalTax" label="本局税收" width="100"></el-table-column>
				<el-table-column align="center" prop="RecordTime" label="记录时间" width="150"></el-table-column>
				<el-table-column label="详细记录">
					<template slot-scope="scope">
						<el-button type="text" icon="el-icon-zoom-in" @click="handleDetail(scope.$index)">详情</el-button>
					</template>
				</el-table-column>
			</el-table>
			<div class="pagination">
				<el-pagination style="margin-top:5px" background layout="total, prev, pager, next, jumper" :hide-on-single-page="true" :total="total" @current-change="handleQuery" :page-size="pagesize"></el-pagination>
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
import '@/assets/css/k.css'
import moment from 'moment'
export default {
	data() {
		return {
			filters: {
				UserId: null,
				GameId: null,
				RoomLevel: null,
				starttime: null,
				endtime: null,
				onlywin: null,
				onlyfaild: null,
			},
			pagesize: 17,
			total: 0,
			games: [],
			rooms: [],
			option_room: [],
			table_data: [],
			dialog_data: {},
			query_time: null,
			dialog: false,
			detail: null,
			detail_title: '',
		}
	},
	components: {},
	computed: {},
	created() {
		this.handleQuery(1)
		app.getInstance().post(
			'/game/game/query',
			{},
			(games) => {
				this.games = [{ GameName: '全部', GameId: null }]
				for (var i = 0; i < games.length; i++) {
					this.games.push({ GameName: games[i].GameName, GameId: games[i].GameId })
				}
				app.getInstance().post(
					'/game/room/query',
					{},
					(rooms) => {
						for (var i = 0; i < rooms.rdata.length; i++) {
							this.rooms.push({ RoomName: rooms.rdata[i].RoomName, RoomLevel: rooms.rdata[i].RoomLevel, GameId: rooms.rdata[i].GameId })
						}
					},
					true
				)
			},
			true
		)
	},
	methods: {
		auth(o) {
			return app.getInstance().auth('游戏系统', '游戏记录', o)
		},
		handleDetail(index) {
			var data = {
				page: 1,
				pagesize: 1,
				GameSerial: this.table_data[index].GameSerial,
			}
			this.detail_title = '详细:' + this.table_data[index].id
			app.getInstance().post('/game/detail/query', data, (result) => {
				if (result.data.length > 0) {
					this.detail = JSON.stringify(JSON.parse(result.data[0].RecordData), null, '\t')
				}
				this.dialog = true
			})
		},
		handleSelectGame() {
			this.filters.level = null
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
				UserId: this.filters.UserId,
				GameId: this.filters.GameId,
				RoomLevel: this.filters.RoomLevel,
				starttime: starttime,
				endtime: endtime,
				onlywin: this.filters.onlywin,
				onlyfaild: this.filters.onlyfaild,
			}
			app.getInstance().post('/game/record/query', data, (result) => {
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
