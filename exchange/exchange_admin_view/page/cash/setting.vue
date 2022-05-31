<template>
	<div>
		<div class="ms-login">
			<el-form label-width="120px" class="ms-content">
				<el-form-item label="充值流水倍数:">
					<el-input v-model="form_data.RechargeScoreMultipe" placeholder="充值流水倍数"> </el-input>
				</el-form-item>
				<el-form-item label="彩金流水倍数:">
					<el-input placeholder="彩金流水倍数" v-model="form_data.PresentScoreMultipe"> </el-input>
				</el-form-item>
				<div class="login-btn">
					<el-button type="primary" @click="SaveSetting()">保存</el-button>
				</div>
			</el-form>
		</div>
	</div>
</template>

<script>
import { app } from '@/api/app.js'
export default {
	data: function () {
		return {
			form_data: {
				RechargeScoreMultipe: 0,
				PresentScoreMultipe: 0,
			},
		}
	},
	created() {
		app.post('/cash/setting/query', {}, (result) => {
			if (result.length > 0) this.form_data = result[0]
		})
	},
	methods: {
		SaveSetting() {
			app.post('/cash/setting/modify', this.form_data, (result) => {
				this.$message.success('保存成功')
			})
		},
	},
}
</script>

<style scoped>
.ms-login {
	position: absolute;
	left: 30%;
	top: 30%;
	width: 450px;
	margin: -190px 0 0 -175px;
	border-radius: 5px;
	background: rgba(255, 255, 255, 0.3);
	overflow: hidden;
}
.ms-content {
	padding: 30px 30px;
}
.login-btn {
	text-align: center;
}
.login-btn button {
	width: 100%;
	height: 36px;
	margin-bottom: 10px;
}
</style>
