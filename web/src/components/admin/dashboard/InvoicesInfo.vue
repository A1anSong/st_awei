<template>
  <a-col :span="12">
    <a-spin :spinning="loading">
      <a-statistic title="订单：Success/Process/Failed" :value="accountsInvoicesText"/>
    </a-spin>
  </a-col>
</template>

<script>
import {notification} from "ant-design-vue";

const axios = require('axios').default

export default {
  name: "InvoicesInfo",
  data() {
    return {
      loading: true,
      success: null,
      inprocess: null,
      failed: null,
    }
  },
  computed: {
    accountsInvoicesText() {
      return this.success + '/' + this.inprocess + '/' + this.failed
    },
  },
  methods: {
    getInvoicesInfo() {
      axios.get('/api/dashInvoices')
          .then(response => {
            this.success = response.data.success
            this.inprocess = response.data.inprocess
            this.failed = response.data.failed
            this.loading = false
          })
          .catch(error => {
            notification.error({
              message: '获取订单信息错误',
              description: error.message,
              duration: 0,
            })
          })
    },
  },
  mounted() {
    this.getInvoicesInfo()
  },
}
</script>

<style scoped>

</style>