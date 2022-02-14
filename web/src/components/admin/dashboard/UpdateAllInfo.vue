<template>
  <a-tooltip :title="taskTooltip">
    <a-progress :percent="process" :success="{ percent: success }" status="active"
                :format="()=>success+'%'" v-if="taskProcessing"/>
  </a-tooltip>
  <a-space>
    <a-button type="primary" @click="patchUpdateInfo" :loading="taskProcessing">更新信息</a-button>
    <a-typography-text>上次更新时间：
    </a-typography-text>
    <a-typography-text code>2006/01/02 15:04</a-typography-text>
  </a-space>
  <a-divider/>
  <a-typography-text>
    1.上次更新时间是假的，还没写完；
    <br/>
    2.有必要的话就催更一个更新状态详细信息；
    <br/>
    3.信息更新完成后5秒会自动刷新页面获取新数据；
    <br/>
    4.目前更新完成后至少5分钟服务器才会初始化更新状态。
  </a-typography-text>
</template>

<script>
import {notification} from "ant-design-vue";

const axios = require('axios').default

export default {
  name: "UpdateAllInfo",
  data() {
    return {
      taskProcessing: false,
      completeAccounts: 0,
      completeBalances: 0,
      completeInvoices: 0,
      completePayouts: 0,
      totalAccounts: 0,
      totalBalances: 0,
      totalInvoices: 0,
      totalPayouts: 0,
    }
  },
  computed: {
    process() {
      if (this.taskProcessing) {
        if (this.totalBalances === 0 || this.totalInvoices === 0) {
          return 50
        } else {
          return 100
        }
      } else {
        return 0
      }
    },
    success() {
      if (this.taskProcessing) {
        if (this.totalBalances === 0 || this.totalInvoices === 0) {
          return parseInt(this.completeAccounts / this.totalAccounts * 10)
        } else {
          return parseInt((this.completeAccounts + this.completeBalances + this.completeInvoices + this.completePayouts)
              / (this.totalAccounts + this.totalBalances + this.totalInvoices + this.totalPayouts)
              * 100)
        }
      } else {
        return 0
      }
    },
    taskTooltip() {
      if (this.taskProcessing) {
        if (this.totalBalances === 0 || this.totalInvoices === 0) {
          return '0 done / 账号，出款 in progress / 余额，订单 to do'
        } else {
          return '账号，出款 done / 余额，订单 in progress / 0 to do'
        }
      } else {
        return '没有任务进行'
      }
    }
  },
  methods: {
    getUpdateInfos() {
      axios.get('/api/updateInfos')
          .then(response => {
            if (response.data.taskProcessing === true) {
              this.completeAccounts = response.data.completeAccounts
              this.completeBalances = response.data.completeBalances
              this.completeInvoices = response.data.completeInvoices
              this.completePayouts = response.data.completePayouts
              this.totalAccounts = response.data.totalAccounts
              this.totalBalances = response.data.totalBalances
              this.totalInvoices = response.data.totalInvoices
              this.totalPayouts = response.data.totalPayouts
              if (this.taskProcessing === true &&
                  this.completeAccounts === this.totalAccounts &&
                  this.completeBalances === this.totalBalances &&
                  this.completeInvoices === this.totalInvoices &&
                  this.totalPayouts === this.totalPayouts) {
                setTimeout(5000)
                // this.$router.go(0)
                window.location.reload()
              }
            } else {
              this.loading = false
            }
          })
          .catch(error => {
            notification.error({
              message: '获取更新信息错误',
              description: error.message,
              duration: 0,
            })
          })
    },
    patchUpdateInfo() {
      axios.patch('/api/updateInfos')
          .then(response => {
            this.taskProcessing = true
            this.getUpdateInfos()
            setInterval(this.getUpdateInfos, 5000)
            notification.success({
              message: '更新操作成功',
              description: response.data,
            })
          })
          .catch(error => {
            if (error.response.status === 412) {
              notification.info({
                message: '更新操作失败',
                description: error.response.data,
              })
              this.taskProcessing = true
              this.getUpdateInfos()
              setInterval(this.getUpdateInfos, 5000)
            } else if (error.status === 501) {
              notification.warning({
                message: '更新操作失败',
                description: error.response.data,
                duration: 0,
              })
              this.taskProcessing = true
              this.getUpdateInfos()
              setInterval(this.getUpdateInfos, 5000)
            } else {
              notification.error({
                message: '更新操作失败',
                description: error.message,
                duration: 0,
              })
            }
          })
    }
  },
}
</script>

<style scoped>

</style>