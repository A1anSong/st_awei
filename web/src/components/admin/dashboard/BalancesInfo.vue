<template>
  <a-col :span="12">
    <a-spin :spinning="loading">
      <a-statistic title="余额：Available/Pending" :value="balancesInfoText"/>
    </a-spin>
  </a-col>
</template>

<script>
import {notification} from "ant-design-vue";

const axios = require('axios').default

export default {
  name: "BalancesInfo",
  data() {
    return {
      loading: true,
      available: null,
      pending: null,
    }
  },
  computed: {
    balancesInfoText() {
      return '€' + this.available + '/€' + this.pending
    },
  },
  methods: {
    getBalancesInfo() {
      axios.get('/api/dashBalances')
          .then(response => {
            this.available = response.data.available / 100
            this.pending = response.data.pending / 100
            this.loading = false
          })
          .catch(error => {
            notification.error({
              message: '获取余额信息错误',
              description: error.message,
              duration: 0,
            })
          })
    },
  },
  mounted() {
    this.getBalancesInfo()
  },
}
</script>

<style scoped>

</style>