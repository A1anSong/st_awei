<template>
  <a-col :span="12">
    <a-spin :spinning="loading">
      <a-statistic title="账号：Alive/Total" :value="accountsInfoText"/>
    </a-spin>
  </a-col>
</template>

<script>
const axios = require('axios').default
import {notification} from 'ant-design-vue'

export default {
  name: "AccountsInfo",
  data() {
    return {
      loading: true,
      active: null,
      total: null,
    }
  },
  computed: {
    accountsInfoText() {
      return this.active + '/' + this.total
    },
  },
  methods: {
    getAccountsInfo() {
      axios.get('/api/dashAccounts')
          .then(response => {
            this.active = response.data.active
            this.total = response.data.total
            this.loading = false
          })
          .catch(error => {
            notification.error({
              message: '获取账户信息错误',
              description: error.message,
              duration: 0,
            })
          })
    },
  },
  mounted() {
    this.getAccountsInfo()
  },
}
</script>

<style scoped>

</style>